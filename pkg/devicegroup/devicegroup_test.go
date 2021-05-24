package devicegroup_test

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
)

var (
	testAndString = `hasCategory("foo") && auto.bar == "baz"`
	testOrString  = `hasCategory("foo") || auto.bar == "baz"`
)

func TestAppliesToBuilder(t *testing.T) {
	t.Parallel()
	builder := devicegroup.NewAppliesToBuilder().HasCategory("foo").And().Auto("bar").Equals("baz")
	if builder.String() != testAndString {
		t.Errorf("appliesTo string is invalid: %s", builder.String())
	}

	builder = devicegroup.NewAppliesToBuilder().HasCategory("foo").Or().Auto("bar").Equals("baz")
	if builder.String() != testOrString {
		t.Errorf("appliesTo string is invalid: %s", builder.String())
	}
}
