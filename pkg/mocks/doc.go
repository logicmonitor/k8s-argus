package mocks

import (

	// as go modules copies only those packages those are used in code, else it cannot
	_ "github.com/golang/mock/mockgen/model"
)
