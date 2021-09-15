package resourcegroup_test

import (
	"testing"

	"github.com/logicmonitor/k8s-argus/pkg/resourcegroup/dgbuilder"
)

var (
	testAndString = `hasCategory("foo") && auto.bar == "baz"`
	testOrString  = `hasCategory("foo") || auto.bar == "baz"`
)

func TestAppliesToBuilder(t *testing.T) {
	t.Parallel()
	builder := dgbuilder.NewAppliesToBuilder().HasCategory("foo").And().Auto("bar").Equals("baz")
	if builder.Build() != testAndString {
		t.Errorf("appliesTo string is invalid: %s", builder.Build())
	}

	builder = dgbuilder.NewAppliesToBuilder().HasCategory("foo").Or().Auto("bar").Equals("baz")
	if builder.Build() != testOrString {
		t.Errorf("appliesTo string is invalid: %s", builder.Build())
	}
}
