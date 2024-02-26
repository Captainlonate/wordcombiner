package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	ce "captainlonate/wordcombiner/internal/customError"
)

// All API responses must have this structure, regardless of success or failure,
// regardless of which endpoint or what data. Always 200. Always JSON. Always this structure.
// Any struct that will ever be placed in the `Data` field of the ApiResponse
// must begin with DTO_
type ApiResponse struct {
	Success bool        `json:"success"`
	Error   error       `json:"error"`
	Data    interface{} `json:"data"`
}

func apiResponseSuccess(data interface{}) *ApiResponse {
	return &ApiResponse{
		Success: true,
		Data:    data,
		Error:   nil,
	}
}

func apiResponseFailure(code ce.ErrorCode, message string) *ApiResponse {
	return &ApiResponse{
		Success: false,
		Data:    nil,
		Error: &ce.ApiError{
			Code:            code,
			FriendlyMessage: message,
		},
	}
}

// When a route handler wants to respond with JSON, we can use
// this convenience utility. It will handle encoding any struct
// to a JSON string before sending it as the response.
func sendJSON(w http.ResponseWriter, data interface{}) {
	jsonResponse, err := json.Marshal(data)

	if err != nil {
		fmt.Println("Failed to encode JSON before response:", err)
		jsonResponse, _ = json.Marshal(apiResponseFailure(ce.EncodeJSONErrorCode, "Failed to encode JSON before response."))
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
