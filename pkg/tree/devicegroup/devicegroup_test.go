package devicegroup

import (
	"testing"
)

var testAndString = `hasCategory("foo") && auto.bar == "baz"`
var testOrString = `hasCategory("foo") || auto.bar == "baz"`

func TestAppliesToBuilder(t *testing.T) {
	builder := NewAppliesToBuilder().HasCategory("foo").And().Auto("bar").Equals("baz")
	if builder.String() != testAndString {
		t.Errorf("appliesTo string is invalid: %s", builder.String())
	}

	builder = NewAppliesToBuilder().HasCategory("foo").Or().Auto("bar").Equals("baz")
	if builder.String() != testOrString {
		t.Errorf("appliesTo string is invalid: %s", builder.String())
	}
}
