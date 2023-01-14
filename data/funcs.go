package data

// throwError 抛出错误
func throwError(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}
