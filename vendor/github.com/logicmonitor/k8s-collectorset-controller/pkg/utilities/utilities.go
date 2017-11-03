package utilities

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/logicmonitor/k8s-collectorset-controller/pkg/metrics"

	"github.com/logicmonitor/lm-sdk-go"
)

// CheckAllErrors is a helper function to deal with the number of possible places that an API call can fail.
func CheckAllErrors(restResponse interface{}, apiResponse *logicmonitor.APIResponse, err error) error {
	var restResponseMessage string
	var restResponseStatus int64

	// Get the underlying concrete type.
	t := reflect.ValueOf(restResponse)

	// Check it the interface is a pointer and get the underlying value.
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Ensure that it is a struct, and get the necessary fields if they are available.
	if t.Kind() == reflect.Struct {
		field := t.FieldByName("Status")
		if field.IsValid() {
			restResponseStatus = field.Int()
		}
		field = t.FieldByName("Errmsg")
		if field.IsValid() {
			restResponseMessage = field.String()
		}
	}

	if restResponseStatus != http.StatusOK {
		metrics.RESTError()
		return fmt.Errorf("[REST] [%d] %s", restResponseStatus, restResponseMessage)
	}

	if apiResponse.StatusCode != http.StatusOK {
		metrics.APIError()
		return fmt.Errorf("[API] [%d] %s", apiResponse.StatusCode, restResponseMessage)
	}

	if err != nil {
		return fmt.Errorf("[ERROR] %v", err)
	}

	return nil
}
