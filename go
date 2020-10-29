#!/usr/bin/env bash

set -e
set -o pipefail
SCRIPT_DIR=$(cd $(dirname $0) ; pwd -P)

TASK=$1
ARGS=${@:2}

APP_NAME="demo"
TERRAFORM_DOCKER_IMAGE="hashicorp/terraform:0.12.8"
CONTAINER_URL="252817234305.dkr.ecr.ap-southeast-1.amazonaws.com/demo"

docker_ensure_volume() {
  docker volume inspect $1 >/dev/null 2>&1 || docker volume create $1 >/dev/null 2>&1
}

docker_ensure_network() {
  docker network inspect $1 >/dev/null 2>&1 || docker network create $1 >/dev/null 2>&1
}

docker_run() {
#  local args="-it -v $(pwd):/workspace:cached -w /workspace"
  local image_id

  if [ -z "${DOCKER_IMAGE}" ]; then
    echo -n "Building toolchain container; this might take a while..."
    DOCKER_IMAGE=$(docker build ${DOCKER_BUILD_ARGS} . -q)
    echo " Done."
  fi

  DOCKER_ARGS="${DOCKER_ARGS} -v ${HOME}/.aws:/root/.aws"
  docker run --rm \
    --hostname $(hostname) \
    --env-file <(env | grep AWS_) \
    --env-file <(env | grep TF_) \
    ${DOCKER_ARGS} ${DOCKER_IMAGE} "$@"
}

gradle() {
  docker_ensure_volume demo-gradle-cache

  DOCKER_BUILD_ARGS="-f ${SCRIPT_DIR}/toolchain-containers/Dockerfile.gradle"
  DOCKER_ARGS="${DOCKER_ARGS} -v demo-gradle-cache:/root/.gradle"

  docker_run gradle $@

  docker build -t gradle -f ${SCRIPT_DIR}/toolchain-containers/Dockerfile.gradle . -q
  docker run -v demo-gradle-cache:/root/.gradle gradle gradle

  local exit=$?
  return $exit
}

tf() {
  DOCKER_IMAGE="${TERRAFORM_DOCKER_IMAGE}"
  DOCKER_ARGS="${DOCKER_ARGS}"

  docker_run "$@"
}

terraform_ecr() {
  cd ${SCRIPT_DIR}/infrastructure/ecr
  tf "$@"
  local exit=$?
  cd - >/dev/null
  return $exit
}

terraform_app() {
  cd ${SCRIPT_DIR}/infrastructure/
  tf "$@"
  local exit=$?
  cd - >/dev/null
  return $exit
}

# task gradle

help__build="build to jar"
task_build() {
  gradle build
}

help__lint="analyzes code for stylistic errors and suspicious constructs"
task_lint() {
  gradle ktlintCheck
}

help__fmt="format kotlin code style for cleanliness"
task_fmt() {
  gradle ktlintFormat
}

help__test="clean up, test and verify code coverage"
task_test() {
  gradle test
}

# task local database
help__startDb="Start the database locally and bind port to port 3306"
task_startDb() {
  task_stopDb

  mkdir -p data
  docker_ensure_network postgres_container

  cd scripts/
  docker-compose up -d
}

task_stopDb() {
  if docker ps | grep "postgres_container" > /dev/null; then
    docker stop postgres_container
  fi
}

# task docker build and push

help__containerize="containerize application into docker image"
task_containerize() {
    docker build --pull --label org.label-schema.vcs-ref=$(git rev-parse HEAD) -f ${SCRIPT_DIR}/Dockerfile.production -t demo .
}

help__infrastructure_apply_ecr="provision ecr"
task_infrastructure_apply_ecr() {

  terraform_ecr init
  terraform_ecr apply $args
}

push_container() {
  local repo_url="$1"
  local app_name=$2

  aws ecr get-login --no-include-email --region ap-southeast-1 | bash

  local revision_tag=$(git rev-parse --short HEAD)

  for tag in "latest" ${revision_tag}; do
    docker tag ${app_name}:latest ${repo_url}:$tag
    docker push ${repo_url}:$tag
  done

  local var_name="$(echo "${app_name}" | tr "[:lower:]" "[:upper:]" | tr - _)_CONTAINER"

  echo "${var_name}_TAG=${revision_tag}" > ${app_name}-container.info
  echo "${var_name}=${repo_url}:${revision_tag}" >> ${app_name}-container.info
}

help__push_container="push image to ECR"
task_push_container() {
  push_container "252817234305.dkr.ecr.ap-southeast-1.amazonaws.com/demo" demo
}

help__infrastructure_apply_app="provision app infra"
task_infrastructure_apply_app() {
  local env=$1

  if [ -z "${env}" ] ; then
    echo "Needs environment"
    exit 1
  fi

  source topup-backend-container.info

  if [ -z "${TOPUP_BACKEND_CONTAINER}" ]; then
    echo "expected TOPUP_BACKEND_CONTAINER"
    exit 1
  fi

  terraform_app init
  terraform_app workspace select $env || terraform_app workspace new $env

  terraform_app apply -var-file $env.tfvars \
                -var application_image_url=${CONTAINER_URL} \
                $restore_args $args

  cd - >/dev/null
}

# task helper
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
  for help in $(list_all_helps)
  do

    HELPS="$HELPS    ${help/help__/} |-- ${!help}$NEW_LINE"
  done

  echo "$HELPS" | column -t -s "|"
  exit 1
fi
