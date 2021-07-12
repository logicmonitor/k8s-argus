// Package aerrors has all predefined argus errors on which errors.Is or errors.As can be used
package aerrors

import "errors"

// ErrNoChangeInUpdateOptions to ignore update
var ErrNoChangeInUpdateOptions = errors.New("update: no change in additional options")

// ErrResourceExists when resource exists
var ErrResourceExists = errors.New("resource exists")

// ErrInvalidCache when cache entry is invalid
var ErrInvalidCache = errors.New("invalid cache error")

// ErrCacheMiss when cache entry is not present
var ErrCacheMiss = errors.New("cache miss error")

// ErrResourceGroupIsNotEmpty when resource group has resources in it
var ErrResourceGroupIsNotEmpty = errors.New("resource group is not empty, it has resources in it")

// ErrResourceGroupParentIsNotValid parent resource group is not valid
var ErrResourceGroupParentIsNotValid = errors.New("parent of resource group is not valid")
