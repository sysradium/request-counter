package common

type Logger interface {
	Printf(string, ...interface{})
}

var _ Logger = (*NullLogger)(nil)

type NullLogger struct{}

func (n *NullLogger) Printf(_ string, _ ...interface{}) {}
