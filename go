#!/bin/bash

set -e

SCRIPT_DIR=$(
  cd $(dirname $0)
  pwd -P
)

TASK=$1
ARGS=${@:2}
## toolchain containerising helpers

account_id_for_name() {
  case $1 in
  'dev') echo "259510286099" ;;
  'qa') echo "259510286099" ;;
  'uat') echo "259510286099" ;;
  'prod') echo "978668668395" ;;
  'tools') echo "688318228301" ;;
  esac
}

account_for_env() {
  case $1 in
  'dev') echo "dev" ;;
  'qa') echo "dev" ;;
  'uat') echo "dev" ;;
  'prod') echo "prod" ;;
  esac
}

runs_inside_gocd() {
  test -n "${GO_JOB_NAME}"
}

docker_run() {
  local image_id
  if [ -z "${image_id}" ]; then
    echo -n "Building toolchain container; this might take a while..."
    image_id=$(docker build ${DOCKER_BUILD_ARGS} . -q)
    echo " Done."
  fi

  if runs_inside_gocd; then
    local args="-v godata:/godata -w $(pwd)"
  else
    local args="-it -v $(pwd):/workspace:cached -w /workspace"
  fi

  DOCKER_ARGS="${DOCKER_ARGS} -v ${HOME}/.aws:/root/.aws"
  docker run --rm \
    -u "$(id -u)" \
    --hostname $(hostname) \
    --env-file <(env | grep JET_) \
    --env-file <(env | grep AWS_) \
    --env-file <(env | grep TF_) \
    ${args} ${DOCKER_ARGS} ${image_id} "$@"
}

docker_ensure_volume() {
  docker volume inspect $1 >/dev/null 2>&1 || docker volume create $1 >/dev/null 2>&1
}

docker_ensure_network() {
  network_name=$1
  if [ ! "$(docker network ls | grep ${network_name})" ]; then
    echo "Creating ${network_name} network ..."
    docker network create ${network_name}
  fi
}

gradle() {
  docker_ensure_volume gradle-cache

  cd ${SCRIPT_DIR}

  DOCKER_BUILD_ARGS="-f ${SCRIPT_DIR}/toolchain-containers/Dockerfile.gradle"
  DOCKER_ARGS="${DOCKER_ARGS} -v gradle-cache:/home/gradle/.gradle"

  docker_run "$@"

  local exit=$?
  cd - >/dev/null
  return $exit
}

assume_role() {
  account_id="$1"
  role="$2"

  credentials=$(aws sts assume-role --role-arn "arn:aws:iam::${account_id}:role/${role}" \
    --role-session-name initial --duration-seconds 2700 | jq '.Credentials')
  export AWS_ACCESS_KEY_ID=$(echo "${credentials}" | jq -r .AccessKeyId)
  export AWS_SECRET_ACCESS_KEY=$(echo "${credentials}" | jq -r .SecretAccessKey)
  export AWS_SESSION_TOKEN=$(echo "${credentials}" | jq -r .SessionToken)
  unset AWS_SECURITY_TOKEN
}

push-container() {
  local repo_url="$1"
  local app_name=$2

  assume_role $(account_id_for_name "tools") "push-containers"

  aws ecr get-login --no-include-email --region ap-southeast-1 | bash

  local revision_tag=$(git rev-parse --short HEAD)

  for tag in "latest" ${revision_tag}; do
    docker tag ${app_name}:latest ${repo_url}:$tag
    docker push ${repo_url}:$tag
  done

  local var_name="$(echo "${app_name}" | tr "[:lower:]" "[:upper:]" | tr - _)_CONTAINER"

  echo "${var_name}_TAG=${revision_tag}" >${app_name}-container.info
  echo "${var_name}=${repo_url}:${revision_tag}" >>${app_name}-container.info
}

tf() {
  if runs_inside_gocd; then
    local docker_user_args="-u $(id -u)"
  else
    local docker_user_args=""
  fi

  DOCKER_ARGS="${DOCKER_ARGS} ${docker_user_args}"
  DOCKER_BUILD_ARGS="-f ${SCRIPT_DIR}/toolchain-containers/Dockerfile.terraform"
  docker_run "$@"
}

terraform() {
  cd ${SCRIPT_DIR}/infrastructure/app
  tf "$@"
  local exit=$?
  cd - >/dev/null
  return $exit
}

kubectl() {

 # cd ${SCRIPT_DIR}

  DOCKER_BUILD_ARGS="-f ${SCRIPT_DIR}/toolchain-containers/Dockerfile.kubernetes"

  docker_run "$@"

  local exit=$?
  return $exit
}

terraform_ecr() {
  cd ${SCRIPT_DIR}/infrastructure/ecr
  tf "$@"
  local exit=$?
  cd - >/dev/null
  return $exit
}

add_container_tag() {
  local repository_name=$1
  local image_tag=$2
  local new_image_tag=$3

  (
    assume_role $(account_id_for_name "tools") "push-containers"

    local image_manifest=$(aws ecr batch-get-image --region ap-southeast-1 \
      --repository-name ${repository_name} \
      --image-ids imageTag=${image_tag} \
      --query 'images[].imageManifest' \
      --output text)

    aws ecr put-image --region ap-southeast-1 \
      --repository-name ${repository_name} \
      --image-tag "${new_image_tag}" \
      --image-manifest "${image_manifest}"
  )
}

exec_psql() {

  pwd
#  cd ${SCRIPT_DIR}  

  DOCKER_BUILD_ARGS="-f ${SCRIPT_DIR}/toolchain-containers/Dockerfile.psql"

  docker_run "$@"

  local exit=$? 
  cd - >/dev/null 
  return $exit
}

## tasks
help__lint="checking code format"
task_lint() {
  gradle ktlint
}

help__fmt="format kotlin code style for cleanliness"
task_fmt() {
  gradle ktlintFormat
}

help__test="clean up, test and verify code coverage"
task_test() {
  gradle test
}

help__build="build jar"
task_build() {
  gradle build
}

help__dependency_check="check security on dependencies"
task_dependency_check() {
  gradle dependencyCheckAnalyze
}

help__static_check="check security on static codes"
task_static_check() {
  gradle check
}

help__running_app="running application"
task_running_app() {
  gradle clean bootRun
}

help__containerize="containerize application into docker image"
task_containerize() {
  docker build --pull -f ${SCRIPT_DIR}/Dockerfile.production -t loan-eligibility-service .
}

help__push_container="push image to ECR"
task_push_container() {
  push-container "688318228301.dkr.ecr.ap-southeast-1.amazonaws.com/coreplatform/loan-eligibility-service" loan-eligibility-service
}

help__apply="Provision backend microservices infrastructure"
task_apply() {
  local env=$1
  local account=$(account_for_env $env)

  if [ -z "${env}" ]; then
    echo "Needs environment"
    exit 1
  fi

  cd ${SCRIPT_DIR}/infrastructure
  if runs_inside_gocd; then
    args="-auto-approve"
  else
    args=""
  fi

  tf init
  tf workspace select $env || tf workspace new $env
  tf apply -var-file $env.tfvars $args

  cd - >/dev/null
}

help__destroy="Provision backend microservices infrastructure"
task_destroy() {
  local env=$1
  local account=$(account_for_env $env)

  if [ -z "${env}" ]; then
    echo "Needs environment"
    exit 1
  fi

  cd ${SCRIPT_DIR}/infrastructure
  if runs_inside_gocd; then
    args="-auto-approve"
  else
    args=""
  fi

  tf init
  tf workspace select $env || tf workspace new $env
  tf destroy -var-file $env.tfvars $args

  cd - >/dev/null
}

help__plan="Provision backend microservices infrastructure"
task_plan() {
  local env=$1
  local account=$(account_for_env $env)

  if [ -z "${env}" ]; then
    echo "Needs environment"
    exit 1
  fi

  cd ${SCRIPT_DIR}/infrastructure
  if runs_inside_gocd; then
    args="-auto-approve"
  else
    args=""
  fi

  tf init
  tf workspace select $env || tf workspace new $env
  tf plan -var-file $env.tfvars $args

  cd - >/dev/null
}

help__infrastructure_apply_ecr="provision ecr"
task_infrastructure_apply_ecr() {
  if runs_inside_gocd; then
    local args="-auto-approve"
  else
    local args=""
  fi

  terraform_ecr init
  terraform_ecr apply $args
}

help__kube_apply="kubectl apply deployment"
task_kube_apply() {
  local env=$1

  source loan-eligibility-service-container.info
  if [ -z "${LOAN_ELIGIBILITY_SERVICE_CONTAINER}" ]; then
    echo "expected LOAN_ELIGIBILITY_SERVICE_CONTAINER"
    exit 1
  fi

  if [ -z "${env}" ]; then
    echo "Needs environment"
    exit 1
  fi

  (
    assume_role $(account_id_for_name ${env}) "deploy-app"

    export LOAN_ELIGIBILITY_SERVICE_CONTAINER=${LOAN_ELIGIBILITY_SERVICE_CONTAINER}
    export DB_CONNECTION_STRING=$(aws rds describe-db-clusters --query '*[].{Endpoint:Endpoint}' --output=text | grep ${env}-global)
    envsubst <infrastructure/k8s/template/deployment.yaml > infrastructure/k8s/template/output.yaml
    cd ${SCRIPT_DIR}/infrastructure/k8s

    if runs_inside_gocd; then
      args="-auto-approve"
    else
      args=""
    fi

    aws eks --region ap-southeast-1 update-kubeconfig --name ${env}_eks_cluster
    cd - >/dev/null
    cp ~/.kube/config ./infrastructure/k8s/config
    kubectl kubectl apply -f infrastructure/k8s/template/output.yaml
    kubectl kubectl apply -f infrastructure/k8s/template/service.yaml
  )
}

help__init_db="docker to initial database"
task_init_db() {
  local env=$1

  (
    assume_role $(account_id_for_name ${env}) "deploy-app"
    export secret=$(aws secretsmanager get-secret-value --secret-id ${env}/coreplatform-db-secrets --query SecretString --output text --region ap-southeast-1)
    export rds_endpoint=$(aws rds describe-db-clusters --query '*[].{Endpoint:Endpoint}' --output=text | grep ${env}-global)
    export DB_USER=RDSUser
    #export loan_db_pass=$(openssl rand -base64 20)
    export loan_db_pass=$(aws secretsmanager get-secret-value --secret-id ${env}/loan-eligibility-db-secrets --query SecretString --output text --region ap-southeast-1)
    export loan_db_pass_encoded=$(echo -n "${loan_db_pass}" | base64)
    #export loan_db_pass_encoded=$(echo -n $loan_db_pass | base64)
    export connection_string=postgresql://${DB_USER}:${secret}@${rds_endpoint}/postgres
    export connection_string_rds=postgresql://${DB_USER}:${secret}@${rds_endpoint}/loan_eligibility
    export loan_user_connection_string=postgresql://loan_user:${loan_db_pass}@${rds_endpoint}/loan_eligibility

    envsubst <infrastructure/k8s/template/initdb.yaml > output.yaml

    envsubst '${loan_db_pass}' <toolchain-containers/init/init-loan-db.sql > output.sql
    aws eks --region ap-southeast-1 update-kubeconfig --name ${env}_eks_cluster
    cp ~/.kube/config ./infrastructure/k8s/config
    kubectl kubectl delete configmap loan-initdb-sql || true
#    kubectl kubectl create configmap loan-initdb-sql --from-file=output.sql --from-file=toolchain-containers/init/create-schema.sql
    kubectl kubectl create configmap loan-initdb-sql --from-file=toolchain-containers/init/create-schema.sql

    kubectl kubectl delete job loan-eligibility-service-init-db-job || true
    kubectl kubectl apply -f output.yaml

    kubectl kubectl delete secret loan-db-secret || true
    envsubst '${loan_db_pass_encoded}' <infrastructure/k8s/template/db-secret.yaml > db-secret.yaml
    kubectl kubectl apply -f db-secret.yaml
  )

}

## main
list_all_helps() {
  compgen -v | egrep "^help__.*"
}

NEW_LINE=$'\n'
if type -t "task_$TASK" &>/dev/null; then
  task_$TASK $ARGS
else
  echo "usage: $0 <task> [<..args>]"
  echo "task:"

  HELPS=""
  for help in $(list_all_helps); do

    HELPS="$HELPS    ${help/help__/} |-- ${!help}$NEW_LINE"
  done

  echo "$HELPS" | column -t -s "|"
  exit 1
fi
