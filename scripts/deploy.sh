#!/bin/bash

set -e

NAMESPACE="logicmonitor"
NAME="argus"

docker login -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD}

if [[ ${TRAVIS_BRANCH} == "master" -a ${TRAVIS_TAG} != ""  ]]; then
  docker tag ${NAMESPACE}/${NAME}:latest ${NAMESPACE}/${NAME}:${TRAVIS_TAG}
  docker tag ${NAMESPACE}/${NAME}:latest ${NAMESPACE}/${NAME}:${TRAVIS_COMMIT}
  docker push ${NAMESPACE}/${NAME}:latest
  docker push ${NAMESPACE}/${NAME}:${TRAVIS_TAG}
  docker push ${NAMESPACE}/${NAME}:${TRAVIS_COMMIT}
elif [[ ${TRAVIS_BRANCH} == "master" ]]; then
  docker tag ${NAMESPACE}/${NAME}:latest ${NAMESPACE}/${NAME}:${TRAVIS_COMMIT}
  docker push ${NAMESPACE}/${NAME}:${TRAVIS_COMMIT}
fi
