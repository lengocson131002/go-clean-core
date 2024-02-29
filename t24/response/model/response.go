package model

import (
	"strings"
)

type Response struct {
	ResponseCommon ResponseCommon
	ResponseRecord ResponseOfsRecord
}

type OfsDataField struct {
	OfsName    string
	MvPosition int
	SvPosition int
	Value      string
}

type ResponseOfsRecord struct {
	ResponseName   string
	ResponseId     string
	ResponseFields []OfsDataField
}

func NewResponseOfsRecord(responseName string, responseId string, responseFields []OfsDataField) ResponseOfsRecord {
	return ResponseOfsRecord{
		// TODO: handle remove record name
		ResponseName:   responseName,
		ResponseId:     responseId,
		ResponseFields: responseFields,
	}
}

// ResponseCommon represents a common response structure.
type ResponseCommon struct {
	MessageID        string
	TransactionID    string
	Status           string
	ErrorMessage     string
	ErrorMessageList []string
}

// NewTxnFailureResponse creates a new instance of ResponseCommon for a transaction failure.
func NewTxnFailureResponse(messageID, txnID string, status ResponseStatus, failureMessage string) ResponseCommon {
	return ResponseCommon{
		MessageID:        messageID,
		TransactionID:    txnID,
		Status:           string(status),
		ErrorMessage:     getErrorResponse(failureMessage),
		ErrorMessageList: getListOfErrorMessage(failureMessage),
	}
}

// NewTxnSuccessResponse creates a new instance of ResponseCommon for a transaction success.
func NewTxnSuccessResponse(messageID, txnID string, status ResponseStatus) ResponseCommon {
	return ResponseCommon{
		MessageID:        messageID,
		TransactionID:    txnID,
		Status:           string(status),
		ErrorMessage:     "",
		ErrorMessageList: nil,
	}
}

// NewEnqSuccessResponse creates a new instance of ResponseCommon for an enquiry success.
func NewEnqSuccessResponse() ResponseCommon {
	return ResponseCommon{
		MessageID:        "",
		TransactionID:    "",
		Status:           string(ResponseStatusSuccess),
		ErrorMessage:     "",
		ErrorMessageList: nil,
	}
}

// NewEnqSuccessEmptyResponse creates a new instance of ResponseCommon for an enquiry success with a response message.
func NewEnqSuccessEmptyResponse(responseMessage string) ResponseCommon {
	return ResponseCommon{
		MessageID:        "",
		TransactionID:    "",
		Status:           string(ResponseStatusSuccess),
		ErrorMessage:     responseMessage,
		ErrorMessageList: nil,
	}
}

// NewFailureResponse creates a new instance of ResponseCommon for a failure response.
func NewFailureResponse(errorMessage string) ResponseCommon {
	if errorMessage == "" {
		panic("Error message is null for an error response")
	}

	return ResponseCommon{
		MessageID:        "",
		TransactionID:    "",
		Status:           string(ResponseStatusError),
		ErrorMessage:     errorMessage,
		ErrorMessageList: nil,
	}
}

// GetListOfErrorMessage returns a list of error messages from the error response string.
func getListOfErrorMessage(errorResponse string) []string {
	errorMessageList := make([]string, 0)

	if errorResponse != "" {
		arrayErrorResponse := strings.Split(errorResponse, "$-$")

		for _, message := range arrayErrorResponse {
			errorMessageList = append(errorMessageList, message)
		}
	}

	return errorMessageList
}

// GetErrorResponse returns the error response string based on specific conditions.
func getErrorResponse(errorResponse string) string {
	if errorResponse != "" && strings.HasPrefix(errorResponse, "AAA") && strings.Contains(errorResponse, "DEPT.CODE:") {
		return ""
	}
	return errorResponse
}
