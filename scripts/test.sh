#!/bin/bash

set -e

GOPACKAGES=$(go list ./... | grep -v /vendor/)
GOFILES=$(find . -type f -name '*.go' -not -path "./vendor/*")

COVERAGE_REPORT=coverage.txt
PROFILE=profile.out

echo "Running tests"
if [[ -f ${COVERAGE_REPORT} ]]; then
  rm ${COVERAGE_REPORT}
fi
touch ${COVERAGE_REPORT}
for package in ${GOPACKAGES[@]}; do
  go test -v -race -coverprofile=${PROFILE} -covermode=atomic $package
  if [ -f ${PROFILE} ]; then
    cat ${PROFILE} >> ${COVERAGE_REPORT}
    rm ${PROFILE}
  fi
done

echo "Linting packages"
gometalinter --vendor --enable-all --disable=gas --disable=gotype --disable=lll --deadline=600s ./...

echo "Formatting go files"
GOFMTFILES="$(gofmt -l -d -s ${GOFILES})"
if [ ! -z "${GOFMTFILES}" ]; then
 echo -e "Failed gofmt files:\n${GOFMTFILES}"
 exit 1
fi

exit 0