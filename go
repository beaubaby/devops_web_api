#!/usr/bin/env bash

set -e
set -o pipefail

TASK=$1
ARGS=${@:2}

APP_NAME="demo"
JAVA_DOCKER_IMAGE="adoptopenjdk/openjdk11:jdk-11.0.3_7"
#JAVA_DOCKER_IMAGE="gradle:6.3-jdk11"
TERRAFORM_DOCKER_IMAGE="hashicorp/terraform:0.12.8"

docker_ensure_volume() {
  docker volume inspect $1 >/dev/null 2>&1 || docker volume create $1 >/dev/null 2>&1
}

docker_ensure_network() {
  docker network inspect $1 >/dev/null 2>&1 || docker network create $1 >/dev/null 2>&1
}

docker_run() {
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

  DOCKER_IMAGE="${JAVA_DOCKER_IMAGE}"
  DOCKER_ARGS="${DOCKER_ARGS} -v demo-gradle-cache:/root/.gradle"

  docker_run ./gradlew $@
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

help__test="test"
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
