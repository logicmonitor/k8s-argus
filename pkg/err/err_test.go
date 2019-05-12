package err

import (
	"fmt"
	"testing"
)

func TestRecoverError(t *testing.T) {
	run()
	fmt.Println("Complete test")

}

func run() {
	defer RecoverError("Test")
	fmt.Println("Start to run")
	panic("Test panic")
}
