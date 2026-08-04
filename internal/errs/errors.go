package errs

// business logic
type ErrRegisterLogic struct {
	Type  errType
	Http  int
	Field string
}

func (e *ErrRegisterLogic) Error() string {
	return ("Register error:" +
		e.Field +
		" - " +
		e.Type.Say())
}

func (e *ErrRegisterLogic) StatusCode() int { return e.Http }

// params validation; http400;
type ErrValidation struct {
	Type  errType
	Http  int
	Field string
	Value string
}

func (e *ErrValidation) Error() string {
	return ("validation error at " +
		e.Field +
		";" +
		e.Type.Say() +
		";found:" +
		e.Value)
}

func (e *ErrValidation) StatusCode() int { return e.Http }

type ErrUserLogin struct {
	Type errType
	Http int
	User string
}

func (e *ErrUserLogin) Error() string {
	return ("User login error: " +
		e.User +
		" - " +
		e.Type.Say())
}

func (e *ErrUserLogin) StatusCode() int { return e.Http }

type ErrRegistercode struct {
	Type errType
	Http int
}

func (e *ErrRegistercode) Error() string {
	return ("registercode failure: " +
		e.Type.Say())
}
func (e *ErrRegistercode) StatusCode() int { return e.Http }

type ErrSupport struct {
	Type errType
	Http int
}

func (e *ErrSupport) Error() string   { return e.Type.Say() }
func (e *ErrSupport) StatusCode() int { return e.Http }

type ErrUnknown struct {
	Type errType
	Http int
}

func (e *ErrUnknown) Error() string   { return e.Type.Say() }
func (e *ErrUnknown) StatusCode() int { return e.Http }

type ErrDbRecord struct {
	Type errType
	Http int
}

func (e *ErrDbRecord) Error() string   { return e.Type.Say() }
func (e *ErrDbRecord) StatusCode() int { return e.Http }
