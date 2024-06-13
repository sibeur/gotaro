package common

const (
	ErrValidationMsg  = "Validation error"
	ErrSlugInvalidMsg = "Slug is invalid"

	// Driver error messages
	ErrDriverAlreadyExistMsg   = "Driver already exist"
	ErrDriverNotFoundMsg       = "Driver not found"
	ErrDriverClientNotFoundMsg = "Driver client not found"

	// Rule error messages
	ErrRuleAlreadyExistMsg = "Rule already exist"
	ErrRuleNotFoundMsg     = "Rule not found"

	// Media error messages
	ErrMediaAlreadyExistMsg = "Media already exist"
	ErrFileSizeExceededMsg  = "File size exceeded"
	ErrFileMimeInvalidMsg   = "File mime invalid"

	// API Client error messages
	ErrAPIClientAlreadyExistMsg = "API client already exist"
	ErrAPIClientNotFoundMsg     = "API client not found"

	// Auth error messages
	ErrAuthenticationFailedMsg = "API Key or Secret Key invalid."
	ErrJWTSecretNotFoundMsg    = "Secret JWT not defined."
	ErrJWTTokenInvalidMsg      = "Token invalid."
	ErrUnauthorizedMsg         = "Unauthorized."

	TemporaryFolder = "tmp"

	APIClientSuperAdminScope = "super-admin"
	APIClientUploaderScope   = "uploader"

	// JWT Issuer
	JWTIssuerAccessToken  = "gotaro-access-token"
	JWTIssuerRefreshToken = "gotaro-refresh-token"
)
