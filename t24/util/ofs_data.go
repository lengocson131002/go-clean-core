package util

import (
	"regexp"
	"strings"
	"time"
)

const (
	QUESTIONMARK_REPLACE       = "%?%"
	CARET_REPLACE              = "%^%"
	PIPE_REPLACE               = "%|%"
	QUOTES_REPLACE             = "|"
	COMMA_REPLACE              = "?"
	SLASH_REPLACE              = "^"
	UNDERSCORE_REPLACE         = "'_'"
	REGEX_VALID_RESPONSE_DATA1 = `(\\S+:\\d+:\\d+=[\\S\\s]*)`
	REGEX_VALID_RESPONSE_DATA2 = `(\\S+=\\S+:\\d+:\\d+)`
	REGEX_TXN_VALID_RESPONSE   = `\S*/\S*/(1)\S*,(((\S+:\\d+:\\d+=[\\S\\s]*)|(\S+=\\S+:\\d+:\\d+))(,*))+`
	REGEX_TXN_ERROR_RESPONSE   = `\S*/\S*/(-1|-2|-3)/[\\S\\s]*`
	REGEX_ENQ_VALID_RESPONSE   = `([\\S\\s]*,[\\S\\s]*,\\S+:\\d+:\\d+=[\\S\\s]*)`
	REGEX_ENQ_ERROR_RESPONSE   = `([\\S\\s]*,[\\S\\s]*,[\\S\\s]*)`
)

func CleanToXML(originalData string) string {
	return originalData // No direct escapeXML equivalent in Go, as HTML escape is usually not needed in Go.
}

func CleanFromXML(originalData string) string {
	return originalData // No direct unescapeXML equivalent in Go, as HTML unescape is usually not needed in Go.
}

func ConvertToOfs(originalValue string) string {
	if ValidateFieldValueForConversion(originalValue) {
		cleanValue := strings.ReplaceAll(originalValue, "?", QUESTIONMARK_REPLACE)
		cleanValue = strings.ReplaceAll(cleanValue, "^", CARET_REPLACE)
		cleanValue = strings.ReplaceAll(cleanValue, "|", PIPE_REPLACE)
		cleanValue = strings.ReplaceAll(cleanValue, "\"", QUOTES_REPLACE)
		cleanValue = strings.ReplaceAll(cleanValue, ",", COMMA_REPLACE)
		cleanValue = strings.ReplaceAll(cleanValue, "/", SLASH_REPLACE)
		return strings.ReplaceAll(cleanValue, "_", UNDERSCORE_REPLACE)
	}
	return originalValue
}

func ConvertXMLToOfsDateFormat(originalDate string) (string, error) {
	if originalDate != "null" {
		xmlFormat := "2006-01-02"
		date, err := time.Parse(xmlFormat, originalDate)
		if err != nil {
			return "", err
		}
		ofsFormat := "20060102"
		return date.Format(ofsFormat), nil
	}
	return "NULL", nil
}

func CleanRecordName(recordName string) string {
	if recordName == "" {
		return recordName
	}
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(recordName, ".", ""), ",", ""), "%", "")
}

func BuildRecordTypeName(recordName string) string {
	if recordName == "" {
		return recordName
	}
	recordName = CleanRecordName(recordName)
	return recordName + "Type"
}

func ValidateFieldValueForConversion(value string) bool {
	return value != "|-|"
}

func ValidateTransactionResponse(ofsResponse string) (bool, error) {
	responseCommonStr, err := ExtractResponseCommon(ofsResponse)
	if err != nil {
		return false, err
	}

	lastSlashIndex := strings.LastIndex(responseCommonStr, "/")
	secondLastSlashIndex := strings.LastIndex(responseCommonStr[:lastSlashIndex], "/")
	ofsResponseSubString := responseCommonStr[secondLastSlashIndex+1:]
	return ofsResponseSubString[0] == '1', nil
}

func ExtractResponseCommon(ofsResponse string) (string, error) {
	firstCommaIndex, err := checkAndGetFirstCommaIndex(ofsResponse)
	if err != nil {
		return "", err
	}

	responseCommonPart := ofsResponse[:firstCommaIndex]
	if !IsResponseCommonWithThreeSlashes(responseCommonPart) {
		return responseCommonPart + "/", nil
	}
	return responseCommonPart, nil
}

func IsResponseCommonWithThreeSlashes(responseCommonPart string) bool {
	regexThreeSlashes := `\S*/\S*/\S*/\S*`
	return regexp.MustCompile(regexThreeSlashes).MatchString(responseCommonPart)
}
