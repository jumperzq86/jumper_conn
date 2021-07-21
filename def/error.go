package def


const (
	ErrConnClosedCode           = 11011
	ErrConnUnexpectedClosedCode = 11012
	ErrInvalidConnParamCode     = 11013

	ErrGetExternalIpCode = 12011
	ErrGetMacAddrCode    = 12012
)

type Error struct{
	Code int32
	Message string
}

func (this *Error) Error() string{
	return this.Message
}

func New(code int32, message string) *Error{
	return &Error{
		Code:    code,
		Message: message,
	}
}

var (
	ErrConnClosed           = New(ErrConnClosedCode, "conn is closed.")
	ErrConnUnexpectedClosed = New(ErrConnUnexpectedClosedCode, "conn is unexpected closed.")
	ErrInvalidConnParam     = New(ErrInvalidConnParamCode, "create conn invalid param.")

	ErrGetExternalIp = New(ErrGetExternalIpCode, "get external ip failed.")
	ErrGetMacAddr    = New(ErrGetMacAddrCode, "get mac addr failed.")
)
