package customerror

import "fmt"

type ErrorCode = string

const EncodeJSONErrorCode ErrorCode = "encode_json_failed"
const OpenAIAPIErrorCode ErrorCode = "openai_api_error"
const OpenAIInvalidResponseErrorCode ErrorCode = "openai_invalid_response"
const BadRequestQueryParam ErrorCode = "bad_or_missing_query_param"
const UnknownErrorCode ErrorCode = "unknown_error"

// When handling errors and responding to the client, we won't pass
// error status codes directly to the client. Instead, we'll always
// return a 200 status code, and one of these hand tailored custom error objects.
// Every error should have a human readable message, and a machine readable error code.
// It will make it easier for the client to "switch" on the error code.
type ApiError struct {
	Code            ErrorCode `json:"error_code"`
	FriendlyMessage string    `json:"human_msg"`
	Err             error
}

// In order for ApiError to satisfy the error interface,
// it needs an Error() method.
func (r *ApiError) Error() string {
	return fmt.Sprintf("error_code %s: err %v", r.Code, r.Err)
}
