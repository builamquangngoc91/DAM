package enums

// Error is a custom type for error messages
type Error int

const (
	BindJSONError                    Error = 200000
	InvalidRequestError              Error = 200001
	HashPasswordError                Error = 200002
	InternalError                    Error = 200003
	UserNotFoundError                Error = 200004
	WrongPasswordError               Error = 200005
	RedisError                       Error = 200006
	InvalidAuthorizationHeaderError  Error = 200007
	AuthorizationHeaderRequiredError Error = 200008
	UnexpectedError                  Error = 200009
	InvalidTokenError                Error = 200010
)
