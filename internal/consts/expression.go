package consts

type expression string

const (
	ExprNull         expression = "$NULL"
	ExprBlank        expression = "$BLANK"
	ExprIndetermined expression = "$INDETERMINED"

	ExprUnreachableCase    expression = "unreachable case contacted"
	ExprUnsupportedFeature expression = "feature unsupported"

	ExprAtomid       expression = "Attr-AtomID"
	ExprUsername     expression = "Attr-Username"
	ExprNickname     expression = "Attr-Nickname"
	ExprPassword     expression = "Attr-Password"
	ExprEntrycode    expression = "Attr-Entrycode"
	ExprRegistercode expression = "Attr-Registercode"

	ExprRegister expression = "Process-Register"
	ExprLogin    expression = "Process-Login"

	ExprHttpInternalServerError expression = "internal server error"
	ExprHttpOk                  expression = "ok"
)

// Biz
const (
	BizExprErrorUnknown    expression = "unknown error"
	BizExprErrorBadRequest expression = "bad request"
)
