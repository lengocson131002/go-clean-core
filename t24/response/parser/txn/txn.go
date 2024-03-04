package txn

import (
	"regexp"
	"strings"

	"github.com/lengocson131002/go-clean/pkg/t24/response/model"
	"github.com/lengocson131002/go-clean/pkg/t24/response/parser"
)

type TransactionResponseDataParser struct {
	parser.AbstractParser
}

func NewTransactionResponseDataParser() parser.ResponseDataParser {
	parserFac := NewTxnOfsResponseParserFactory()
	return &TransactionResponseDataParser{
		parser.AbstractParser{parserFac},
	}
}

func (parser *TransactionResponseDataParser) ParseResponseData(responseName, ofsResponse string) (model.Response, error) {
	if IsTxnResponseSuccess(ofsResponse) {
		return parser.buildSuccessResponse(responseName, ofsResponse)
	}
	if IsTxnResponseFailure(ofsResponse) {
		return parser.buildFailureResponse(ofsResponse)
	}
	return parser.buildErrorResponse(ofsResponse)
}

func (parser *TransactionResponseDataParser) ParseResponseDataWithoutResponseName(ofsResponse string) (model.Response, error) {
	if IsTxnResponseSuccess(ofsResponse) {
		return parser.buildSuccessResponseWithoutResponseName(ofsResponse)
	}
	if IsTxnResponseFailure(ofsResponse) {
		return parser.buildFailureResponse(ofsResponse)
	}
	return parser.buildErrorResponse(ofsResponse)
}

func (parser *TransactionResponseDataParser) buildSuccessResponse(responseName, ofsResponse string) (model.Response, error) {
	responseCommonStr := parser.extractResponseCommon(ofsResponse)
	responseCommon, err := parser.buildResponseCommon(responseCommonStr)
	if err != nil {
		return model.Response{}, err
	}

	appName := parser.extractApplicationName(responseCommonStr)
	if appName != "" {
		responseName = appName
	}

	responseRecord, err := parser.buildResponseRecord(responseCommon.TransactionID, ofsResponse, responseName)
	if err != nil {
		return model.Response{}, err
	}

	return model.Response{
			ResponseCommon: responseCommon,
			ResponseRecord: responseRecord,
		},
		nil
}

func (parser *TransactionResponseDataParser) buildSuccessResponseWithoutResponseName(ofsResponse string) (model.Response, error) {
	responseCommonStr := parser.extractResponseCommon(ofsResponse)
	responseCommon, err := parser.buildResponseCommon(responseCommonStr)
	if err != nil {
		return model.Response{}, err
	}

	return model.Response{
			ResponseCommon: responseCommon,
		},
		nil
}

func (parser *TransactionResponseDataParser) buildErrorResponse(ofsResponse string) (model.Response, error) {
	resComm := model.NewFailureResponse(ofsResponse)
	return model.Response{
		ResponseCommon: resComm,
	}, nil
}

func (parser *TransactionResponseDataParser) buildFailureResponse(ofsResponse string) (model.Response, error) {
	var failureMessage string
	var propertyName string
	var errorMessageProcessedStatus bool

	responseCommonStr := parser.extractResponseCommon(ofsResponse)

	if !strings.HasPrefix(responseCommonStr, "AAA") && strings.HasPrefix(responseCommonStr, "AA") {
		propertyName = parser.extractPropertyName(responseCommonStr)
		errorMessageProcessedStatus = true
	}

	if !strings.HasPrefix(responseCommonStr, "AAA") || (propertyName != "") {
		failureMessage = ofsResponse
		errorMessageProcessedStatus = true
	}

	if strings.HasPrefix(responseCommonStr, "AAA") && !errorMessageProcessedStatus {
		errorMessage := parser.extractResponseData(ofsResponse)
		if !strings.Contains(errorMessage, "DEPT.CODE:") {
			failureMessage = ofsResponse
			errorMessageProcessedStatus = true
		}
	}

	responseCommon, err := parser.buildFailureResponseCommon(responseCommonStr, failureMessage)
	if err != nil {
		return model.Response{}, err
	}

	responseRecord := model.ResponseOfsRecord{ResponseId: responseCommon.TransactionID}
	return model.Response{ResponseCommon: responseCommon, ResponseRecord: responseRecord}, nil
}

func (parser *TransactionResponseDataParser) buildFailureResponseCommon(responseCommon, failureMessage string) (model.ResponseCommon, error) {
	statusStr := parser.extractStatus(responseCommon)
	messageID := parser.extractMessageID(responseCommon)
	txnID := parser.extractTransactionID(responseCommon)
	responseStatus := model.From(statusStr)

	return model.NewTxnFailureResponse(messageID, txnID, responseStatus, failureMessage), nil
}

func (parser *TransactionResponseDataParser) extractStatus(responseCommon string) string {
	secondSlashIndex := parser.checkAndGetSecondSlashIndex(responseCommon)
	thirdSlashIndex := parser.checkAndGetThirdSlashIndex(responseCommon)
	return responseCommon[secondSlashIndex+1 : thirdSlashIndex]
}

func (parser *TransactionResponseDataParser) checkAndGetSecondSlashIndex(responseCommonPart string) int {
	thirdSlashIndex := parser.checkAndGetThirdSlashIndex(responseCommonPart)
	strippedResponseCommon := responseCommonPart[:thirdSlashIndex]

	secondSlashIndex := strings.LastIndex(strippedResponseCommon, "/")
	return secondSlashIndex
}

func (parser *TransactionResponseDataParser) extractPropertyName(responseCommonStr string) string {
	propertyName := responseCommonStr

	indexOfSlash := strings.Index(propertyName, "/")
	if indexOfSlash != -1 {
		propertyName = propertyName[:indexOfSlash]

		if strings.Contains(propertyName, "-") {
			indexOfFirstDash := strings.Index(propertyName, "-")
			indexOfLastDash := strings.LastIndex(propertyName, "-")

			if indexOfFirstDash != -1 && indexOfLastDash != -1 && indexOfFirstDash != indexOfLastDash {
				propertyName = propertyName[indexOfFirstDash+1 : indexOfLastDash]
			} else {
				return ""
			}
		} else {
			return ""
		}
	} else {
		return ""
	}

	return propertyName
}

func (parser *TransactionResponseDataParser) buildResponseRecord(txnID, ofsResponse, applicationName string) (model.ResponseOfsRecord, error) {
	responseData := parser.extractResponseData(ofsResponse)
	responseFields, err := parser.ParseData(responseData)
	if err != nil {
		return model.ResponseOfsRecord{}, err
	}

	return model.NewResponseOfsRecord(txnID, applicationName, responseFields), nil
}

func (parser *TransactionResponseDataParser) extractResponseData(ofsResponse string) string {
	firstCommaIndex := parser.checkAndGetFirstCommaIndex(ofsResponse)
	dataLength := len(ofsResponse)
	return ofsResponse[firstCommaIndex+1 : dataLength]
}

func (parser *TransactionResponseDataParser) extractApplicationName(responseCommonPart string) string {
	thirdSlashIndex := parser.checkAndGetThirdSlashIndex(responseCommonPart)
	if thirdSlashIndex < 0 {
		return ""
	}
	return responseCommonPart[thirdSlashIndex+1:]
}

func (parser *TransactionResponseDataParser) checkAndGetThirdSlashIndex(responseCommonPart string) int {
	thirdSlashIndex := strings.LastIndex(responseCommonPart, "/")
	return thirdSlashIndex
}

func (parser *TransactionResponseDataParser) extractResponseCommon(ofsResponse string) string {
	firstCommaIndex := parser.checkAndGetFirstCommaIndex(ofsResponse)
	responseCommonPart := ofsResponse[:firstCommaIndex]
	if !parser.isResponseCommonWithThreeSlashes(responseCommonPart) {
		return responseCommonPart + "/"
	}
	return responseCommonPart
}

func (parser *TransactionResponseDataParser) buildResponseCommon(responseCommon string) (model.ResponseCommon, error) {
	statusStr := parser.extractStatusMessage(responseCommon)
	messageID := parser.extractMessageID(responseCommon)
	txnID := parser.extractTransactionID(responseCommon)
	responseStatus := model.From(statusStr)

	return model.ResponseCommon{
		MessageID:     messageID,
		TransactionID: txnID,
		Status:        string(responseStatus),
	}, nil
}

func (parser *TransactionResponseDataParser) extractTransactionID(responseCommon string) string {
	responseCommon = parser.getTrimmedMessage(responseCommon)
	lastSlashIndex := strings.LastIndex(responseCommon, "/")
	return responseCommon[:lastSlashIndex]
}

func (parser *TransactionResponseDataParser) extractStatusMessage(responseCommon string) string {
	lastSlashIndex := strings.LastIndex(responseCommon, "/")
	secondLastSlashIndex := strings.LastIndex(responseCommon[:lastSlashIndex-1], "/")
	ofsResponseSubString := responseCommon[secondLastSlashIndex+1 : lastSlashIndex]
	return ofsResponseSubString
}

func (parser *TransactionResponseDataParser) extractMessageID(responseCommon string) string {
	responseCommon = parser.getTrimmedMessage(responseCommon)
	lastSlashIndex := strings.LastIndex(responseCommon, "/")
	return responseCommon[lastSlashIndex+1:]
}

func (parser *TransactionResponseDataParser) getTrimmedMessage(responseCommon string) string {
	lastSlashIndex := strings.LastIndex(responseCommon, "/")
	secondLastSlashIndex := strings.LastIndex(responseCommon[:lastSlashIndex-1], "/")
	return responseCommon[:secondLastSlashIndex]
}

func (parser *TransactionResponseDataParser) checkAndGetFirstCommaIndex(ofsResponse string) int {
	firstCommaIndex := strings.Index(ofsResponse, ",")
	return firstCommaIndex
}

func (parser *TransactionResponseDataParser) isResponseCommonWithThreeSlashes(responseCommonPart string) bool {
	regexThreeSlashes := "\\S*/\\S*/\\S*/\\S*"
	return regexp.MustCompile(regexThreeSlashes).MatchString(responseCommonPart)
}
