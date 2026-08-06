package consts

type expression string

const (
	ExprNull         expression = "$NULL"
	ExprBlank        expression = "$BLANK"
	ExprIndetermined expression = "$INDETERMINED"

	ExprUnreachableCase    expression = "unreachable case contacted"
	ExprUnsupportedFeature expression = "feature unsupported"

	ExprFailedRegistercodeWithdraw expression = "failed to withdraw registercode"

	//ExprAtomid       expression = "Attr-AtomID"
	ExprUsername     expression = "Attr-Username"
	ExprNickname     expression = "Attr-Nickname"
	ExprPassword     expression = "Attr-Password"
	ExprEntrycode    expression = "Attr-Entrycode"
	ExprRegistercode expression = "Attr-Registercode"

	ExprUser expression = "User"

	ExprEnckey expression = "Crypto-PrivateKey"
	ExprDeckey expression = "Crypto-PublicKey"

	ExprRegister expression = "Process-Register"
	ExprLogin    expression = "Process-Login"

	ExprHttpInternalServerError expression = "internal server error"
	ExprHttpOk                  expression = "ok"
	ExprHttpTooManyRequests     expression = "too many requests"
)

// Biz
const (
	BizExprErrorUnknown    expression = "unknown error"
	BizExprErrorBadRequest expression = "bad request"
)

// program
const (
	ExprContextKeyJWT expression = "jwt_user"
	ExprUserID        expression = "uid"
)

// json
const (
	JsonExprAtomid       expression = "atomid"
	JsonExprUsername     expression = "username"
	JsonExprNickname     expression = "nickname"
	JsonExprPassword     expression = "password"
	JsonExprRegistercode expression = "registercode"
	JsonExprEntrycode    expression = "entrycode"
	JsonExprAccessToken  expression = "access_token"
	JsonExprRefreshToken expression = "refresh_token"
	JsonExprRole         expression = "role"
)
