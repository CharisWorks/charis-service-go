package util

type Error struct {
	Message string `json:"message"`
}

func NewError(message string) *Error {
	return &Error{
		Message: message,
	}
}
func (e *Error) Error() string {
	return e.Message
}
