package err

import (
	"fmt"
	"testing"
)

func TestRecoverError(t *testing.T) {
	run(t)
	t.Logf("test success")

}

func run(t *testing.T) {
	defer RecoverError("test msg")
	t.Logf("test run start...")
	panic(fmt.Errorf("test panic"))
	t.Logf("test run end...")
}
