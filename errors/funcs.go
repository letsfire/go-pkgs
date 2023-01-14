package errors

import (
	"fmt"
	"sync"
)

// Collection 错误收集
type Collection struct {
	sync.Mutex
	errors []error
}

func (c *Collection) Add(errs ...error) {
	c.Lock()
	defer c.Unlock()
	for _, err := range errs {
		if err != nil {
			c.errors = append(c.errors, err)
		}
	}
}

func (c *Collection) First() error {
	if len(c.errors) == 0 {
		return nil
	}
	return c.errors[0]
}

func (c *Collection) AsError() error {
	if len(c.errors) == 0 {
		return nil
	} else if len(c.errors) == 1 {
		return c.First()
	}
	return fmt.Errorf("%d errors, %#v", len(c.errors), c.errors)
}

// First 返回首个错误
func First(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// Panic 必须没有错误
func Panic(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}

// Wrap 包装错误
func Wrap(err error, msg string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s, error = %s", fmt.Sprintf(msg, args...), err)
}

// Retry 错误重试
func Retry(maxNum int, fn func(int) error) error {
	var err error
	for i := 0; i < maxNum; i++ {
		if err = fn(i); err == nil {
			break
		}
	}
	return err
}

// SerialUntil 按序处理直到发生错误
func SerialUntil(fns ...func() error) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

// CallIfNoErr 没有错误继续调用,直至发生错误
func CallIfNoErr(err error, fns ...func() error) error {
	if err != nil {
		return err
	}
	return SerialUntil(fns...)
}

// CallIfNoErr2 没有错误继续调用,直至发生错误
func CallIfNoErr2(err error, fns ...func()) error {
	if err == nil {
		for _, fn := range fns {
			fn()
		}
	}
	return err
}
