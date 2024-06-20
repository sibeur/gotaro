package common

import "time"

const (
	ErrValidationMsg  = "Validation error"
	ErrSlugInvalidMsg = "Slug is invalid"

	// Driver error messages
	ErrDriverAlreadyExistMsg   = "Driver already exist"
	ErrDriverNotFoundMsg       = "Driver not found"
	ErrDriverClientNotFoundMsg = "Driver client not found"
	ErrDriverNotInitiate       = "Driver not initiate"

	// GCP Driver error message
	ErrBucketNotExistMsg               = "Bucket not exist"
	ErrNotHaveStorageAdminPrivilageMsg = "Not have storage admin privilage"

	// Rule error messages
	ErrRuleAlreadyExistMsg = "Rule already exist"
	ErrRuleNotFoundMsg     = "Rule not found"

	// Media error messages
	ErrMediaAlreadyExistMsg = "Media already exist"
	ErrMediaNotFoundMsg     = "Media not found"
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

	// Media default config
	TemporaryFolder     = "tmp"
	DefaultSignedURLTTL = time.Minute * 10

	// API Client default scope
	APIClientSuperAdminScope = "super-admin"
	APIClientUploaderScope   = "uploader"

	// JWT Issuer
	JWTIssuerAccessToken  = "gotaro-access-token"
	JWTIssuerRefreshToken = "gotaro-refresh-token"

	// Cache Keys
	CacheGetMediaKey       = "gotaro:media:%s:%s"
	CacheMediaSignedUrlKey = "gotaro:media:signedUrl:%s:%s"

	// Cache TTL
	DefaultGetMediaCacheTTL = 60 * 10
)
