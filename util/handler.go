package util

// error handler util 만들기 -> Logger로 변경 필요
func ErrorHandler (err error) {
	if err != nil {
		panic(err)
	}
}