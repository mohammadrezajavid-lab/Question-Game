package richerror

type Kind int

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
)

type Operation string

type RichError struct {
	wrappedError error
	operation    Operation
	message      string
	kind         Kind
	meta         map[string]interface{}
}

func NewRichError(operation Operation) *RichError {

	return &RichError{
		wrappedError: nil,
		operation:    operation,
		message:      "",
		kind:         0,
		meta:         nil,
	}
}

/*func NewRichError(args ...interface{}) RichError {
	r := RichError{}

	for _, val := range args {
		switch val.(type) {
		case error:
			r.wrappedError = val.(error)
		case Operation:
			r.operation = val.(Operation)
		case string:
			r.message = val.(string)
		case Kind:
			r.kind = val.(Kind)
		case map[string]interface{}:
			r.meta = val.(map[string]interface{})
		}
	}

	return r
}*/

/*func NewRichError(err error, operation, message string, kind Kind) *RichError {

	return &RichError{
		wrappedError: err,
		operation:    operation,
		message:      message,
		kind:         kind,
	}
}*/

func (re RichError) GetError() error {
	return re.wrappedError
}

func (re RichError) GetOperation() Operation {
	return re.operation
}

func (re RichError) GetMessage() string {
	return re.message
}

func (re RichError) GetKind() Kind {
	return re.kind
}

func (re RichError) GetMeta() map[string]interface{} {
	return re.meta
}

func (re RichError) WithError(err error) RichError {
	re.wrappedError = err

	return re
}

func (re RichError) WithOperation(operation Operation) RichError {
	re.operation = operation

	return re
}

func (re RichError) WithMessage(message string) RichError {
	re.message = message

	return re
}

func (re RichError) WithKind(kind Kind) RichError {
	re.kind = kind

	return re
}

func (re RichError) WithMeta(meta map[string]interface{}) RichError {
	re.meta = meta

	return re
}
func (re RichError) Error() string {

	return re.message
}
