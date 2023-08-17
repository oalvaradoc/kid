package fault

// define the Authorization related error code
var (
	Success         = 0
	ErrorCodeOffset = 0

	//CommonError       = ErrorCodeOffset + 100  // General error start code, such as parameter error, etc.
	AuthorityError = ErrorCodeOffset + 200 // Authorization related error start code
	//AccessError       = ErrorCodeOffset + 500  // Menu and access control

	// Authorization related error start code
	InvalidToken          = AuthorityError + 1 // Invalid token
	TokenExpired          = AuthorityError + 2 // token expired
	NoTokenParameter      = AuthorityError + 3 // Request is missing token
	PermissionDenied      = AuthorityError + 4 // Insufficient user rights
	UndefinedRing         = AuthorityError + 5 // Undefined ring
	InvalidTopic          = AuthorityError + 6 // Invalid Topic
	PermissionOutOfDesign = AuthorityError + 7 // Authorized operations beyond design
)
