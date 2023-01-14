package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	var errs = new(Collection)
	errs.Add(errors.New("error1"))
	fmt.Println(errs.AsError())
	errs.Add(errors.New("error2"))
	fmt.Println(errs.AsError())
}
