#!/usr/bin/env bash

set -e
set -o pipefail
SCRIPT_DIR=$(cd $(dirname $0) ; pwd -P)

TASK=$1
ARGS=${@:2}

APP_NAME="demo"
TERRAFORM_DOCKER_IMAGE="hashicorp/terraform:0.12.8"

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

# task gradle

help__build="build to jar"
task_build() {
  gradle assemble
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

help__dockerBuild="build docker image"
task_dockerBuild() {
  docker build \
    --pull \
    --label org.label-schema.vcs-ref=$(git rev-parse HEAD) \
    -f Dockerfile.production \
    -t ${APP_NAME} .
}

help__dockerPush="push docker image to ECR"
task_dockerPush() {

  local repo_url="252817234305.dkr.ecr.ap-southeast-1.amazonaws.com/api/demo/"
  local latest_tag="${repo_url}:latest"
  local revision=$(git rev-parse --short HEAD)
  local revision_tag="${repo_url}:${revision}"

  for tag in ${latest_tag} ${revision_tag};
  do
    docker tag ${APP_NAME} ${tag}
  done
  (
    aws ecr get-login --no-include-email --region ap-southeast-1 | bash
    docker push ${latest_tag}
    docker push ${revision_tag}
  )

  local var_name="$(echo "${APP_NAME}" | tr "[:lower:]" "[:upper:]" | tr - _)_CONTAINER"
  echo "${var_name}_TAG=${revision}" > ${APP_NAME}-container.info
  echo "${var_name}=${revision_tag}" >> ${APP_NAME}-container.info
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
