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

	if re.wrappedError != nil {

		return re.wrappedError
	}

	re, ok := re.wrappedError.(RichError)
	if !ok {

		return re.wrappedError
	}

	return re.GetError()
}

func (re RichError) GetOperation() Operation {
	return re.operation
}

func (re RichError) GetMessage() string {

	if re.message != "" {

		return re.message
	}

	re, ok := re.wrappedError.(RichError)
	if !ok {

		return re.wrappedError.Error()
	}

	return re.GetMessage()
}

func (re RichError) GetKind() Kind {

	if re.kind != 0 {

		return re.kind
	}

	re, ok := re.wrappedError.(RichError)
	if !ok {

		return re.kind
	}

	return re.GetKind()
}

func (re RichError) GetMeta() map[string]interface{} {

	if re.meta != nil {

		return re.meta
	}

	re, ok := re.wrappedError.(RichError)
	if !ok {
		return nil
	}

	return re.GetMeta()
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

	if re.message == "" {

		re.wrappedError.Error()
	}

	return re.message
}
