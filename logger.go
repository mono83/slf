package slf

// Logger interface describes logger structure
type Logger interface {
	Trace(message string, p ...Param)
	Debug(message string, p ...Param)
	Info(message string, p ...Param)
	Warning(message string, p ...Param)
	Error(message string, p ...Param)
	Alert(message string, p ...Param)
	Emergency(message string, p ...Param)
}
