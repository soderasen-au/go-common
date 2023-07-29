package loggers

type NullWriter struct{}

func (*NullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
