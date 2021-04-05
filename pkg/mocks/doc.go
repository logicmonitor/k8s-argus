package mocks

// marker to download package in vendor otherwise go generate mocks for unit testing fails with no package in vendor
// as go modules copies only those packages those are used in code, else it cannot
import _ "github.com/golang/mock/mockgen/model"
