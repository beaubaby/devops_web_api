#!/bin/bash
set -e

aws ecr get-login --registry-ids 252817234305 --no-include-email --region ap-southeast-1 | bash
docker run --detach \
  --name demo \
  --restart always \
  --net=host \
  -e SERVER_PORT=80 \
  -e SPRING_PROFILES_ACTIVE=${active_spring_profiles} \
  ${container_url}