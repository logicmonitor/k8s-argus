package err

import (
	"fmt"
	"testing"
)

func TestRecoverError(t *testing.T) {
	run(t)
	t.Logf("Test success")

}

func run(t *testing.T) {
	defer RecoverError("Test msg")
	t.Logf("Test run start...")
	panic(fmt.Errorf("test panic"))
	t.Logf("Test run end...")
}
