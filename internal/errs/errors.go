package errs

import "fmt"

// business logic
type ErrRegisterLogic struct {
	Type  errType
	Http  int
	Field string
}

func (e *ErrRegisterLogic) Error() string {
	return (fmt.Sprintf("error in register; field::%s; type::%s", e.Field, e.Type.Say()))
}
func (e *ErrRegisterLogic) StatusCode() int { return e.Http }
func (e *ErrRegisterLogic) Respond() string {
	return fmt.Sprintf("failure in register; wrong with %s", e.Field)
}

// validation
type ErrValidation struct {
	Type  errType
	Http  int
	Field string
	Value string
}

func (e *ErrValidation) Error() string {
	return fmt.Sprintf("error in validation; field::%s; details::%s", e.Field, e.Type.Say())
	//return ("validation error at " + e.Field + ";" + e.Type.Say() + ";found:" + e.Value)
}
func (e *ErrValidation) StatusCode() int { return e.Http }
func (e *ErrValidation) Respond() string {
	return fmt.Sprintf("failure in params validation; field::%s;", e.Field)
}

// user login
type ErrUserLogin struct {
	Type errType
	Http int
	User string
}

func (e *ErrUserLogin) Error() string {
	return fmt.Sprintf("error in user login; user::%s; details::%s", e.User, e.Type.Say())
}
func (e *ErrUserLogin) StatusCode() int { return e.Http }
func (e *ErrUserLogin) Respond() string {
	return fmt.Sprintf("failure in user login; details::%s", e.Type.Say())
}

// registercode
type ErrRegistercode struct {
	Type errType
	Http int
}

func (e *ErrRegistercode) Error() string {
	return ("error in registercode; details::" + e.Type.Say())
}
func (e *ErrRegistercode) StatusCode() int { return e.Http }
func (e *ErrRegistercode) Respond() string {
	return fmt.Sprintf("failure in checking registercode; details::%s", e.Type.Say())
}

// support/unsupported
type ErrSupport struct {
	Type errType
	Http int
}

func (e *ErrSupport) Error() string   { return e.Type.Say() }
func (e *ErrSupport) StatusCode() int { return e.Http }
func (e *ErrSupport) Respond() string { return e.Type.Say() }

// unknown
type ErrUnknown struct {
	Type errType
	Http int
}

func (e *ErrUnknown) Error() string   { return e.Type.Say() }
func (e *ErrUnknown) StatusCode() int { return e.Http }
func (e *ErrUnknown) Respond() string { return e.Type.Say() }

// db record
type ErrDbRecord struct {
	Type  errType
	Http  int
	Field string
}

func (e *ErrDbRecord) Error() string {
	return fmt.Sprintf("error in db record; field::%s; details::%s", e.Field, e.Type.Say())
}
func (e *ErrDbRecord) StatusCode() int { return e.Http }
func (e *ErrDbRecord) Respond() string {
	return fmt.Sprintf("error in db record; field::%s; details::%s", e.Field, e.Type.Say())
}
