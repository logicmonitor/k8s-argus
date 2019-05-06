#!/bin/bash

set -e

GOPACKAGES=$(go list ./... | grep -v /vendor/)
GOFILES=$(find . -type f -name '*.go' -not -path "./vendor/*")

test_packages() {
  echo "Testing packages"
  local coverage_report="coverage.txt"
  local profile="profile.out"
  if [[ -f ${coverage_report} ]]; then
    rm ${coverage_report}
  fi
  touch ${coverage_report}
  for package in ${GOPACKAGES[@]}; do
    go test -v -race -coverprofile=${profile} -covermode=atomic $package
    if [ -f ${profile} ]; then
      cat ${profile} >> ${coverage_report}
      rm ${profile}
    fi
  done
}

lint_packages() {
  echo "Linting packages"
  gometalinter --aggregate --vendor --exclude="zz_generated" --exclude="api" --disable-all --deadline=1200s --enable=deadcode --enable=dupl --enable=errcheck --enable=goconst --enable=gocyclo --enable=gofmt --enable=goimports --enable=golint --enable=gosec --enable=gotypex --enable=ineffassign --enable=interfacer --enable=maligned --enable=misspell --enable=nakedret --enable=structcheck --enable=test --enable=testify --enable=unconvert --enable=unparam --enable=varcheck --enable=vet --enable=vetshadow ./...
}

format_files() {
  echo "Formatting files"
  local gofmtfiles="$(gofmt -l -d -s ${GOFILES})"
  if [ ! -z "${gofmtfiles}" ]; then
    echo -e "Failed gofmt files:\n${gofmtfiles}"
    exit 1
  fi
}

lint_packages
format_files
test_packages

exit 0
