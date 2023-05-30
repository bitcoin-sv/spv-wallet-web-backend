package util

// WriterFunc wrapper type for function that is implementing io.Writer interface.
type WriterFunc func(p []byte) (n int, err error)

// Write proxy to implement io.Writer interface.
func (f WriterFunc) Write(p []byte) (n int, err error) {
	return f(p)
}
