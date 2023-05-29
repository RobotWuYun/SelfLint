package example

func myLog(format string, args ...interface{}) {
	leafFunc(format, args...)
}

func leafFunc(format string, args ...interface{}) {
	// do something
}
