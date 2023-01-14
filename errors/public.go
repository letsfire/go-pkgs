package errors

var (
	ClientError = func(s string, args ...interface{}) Error {
		return NewErrorf("client.error", s, args...)
	}

	Unauthorized = func(s string, args ...interface{}) Error {
		return NewErrorf("unauthorized", s, args...)
	}

	Forbidden = func(s string, args ...interface{}) Error {
		return NewErrorf("forbidden", s, args...)
	}

	NotFound = func(s string, args ...interface{}) Error {
		return NewErrorf("not.found", s, args...)
	}

	TooManyRequests = func(s string, args ...interface{}) Error {
		return NewErrorf("to.many.requests", s, args...)
	}

	TooManyErrReqs = func(s string, args ...interface{}) Error {
		return NewErrorf("to.many.error.requests", s, args...)
	}

	ServerError = func(s string, args ...interface{}) Error {
		return NewErrorf("server.error", s, args...)
	}

	CallRPCError = func(err Error, s string, args ...interface{}) Error {
		return NewErrorf("call.rpc.error", s, args...).AddPrev(err)
	}
)
