package customError

import "fmt"

type ErrorCode = string

const EncodeJSONErrorCode ErrorCode = "encode_json_failed"
const OpenAIAPIErrorCode ErrorCode = "openai_api_error"
const OpenAIInvalidResponseErrorCode ErrorCode = "openai_invalid_response"
const BadRequestQueryParam ErrorCode = "bad_or_missing_query_param"
const UnknownErrorCode ErrorCode = "unknown_error"

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
