package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lengocson131002/go-clean/pkg/t24/response/model"
)

// IsNullOrEmpty checks if a string is null or empty
func IsNullOrEmpty(str string) bool {
	return str == "" || strings.TrimSpace(str) == ""
}

// GetHexString converts a byte array to a hexadecimal string
func GetHexString(value []byte) (string, error) {
	var result strings.Builder
	for _, v := range value {
		result.WriteString(fmt.Sprintf("%02x", v))
	}
	return result.String(), nil
}

// CreateEnqCondition creates an Enquiry condition string
func CreateEnqCondition(fieldName, operator, value string) string {
	if IsNullOrEmpty(fieldName) {
		return ""
	}
	return fmt.Sprintf(",%s:%s=%s", fieldName, operator, value)
}

// CreateOfsSingleValue creates an Ofs single value string
func CreateOfsSingleValue(t24FieldID string, value interface{}) string {
	if IsNullOrEmpty(t24FieldID) {
		return ""
	}
	if IsNullOrEmpty(fmt.Sprintf("%v", value)) {
		return fmt.Sprintf(",%s:1=", t24FieldID)
	}
	return fmt.Sprintf(",%s:1=%v", t24FieldID, value)
}

// CreateOfsMultiValue creates an Ofs multi-value string
func CreateOfsMultiValue(t24FieldID string, mvPosition, svPosition int, value string) string {
	if IsNullOrEmpty(t24FieldID) {
		return ""
	}
	if IsNullOrEmpty(value) {
		return fmt.Sprintf(",%s:%d:%d=", t24FieldID, mvPosition, svPosition)
	}
	return fmt.Sprintf(",%s:%d:%d=%s", t24FieldID, mvPosition, svPosition, value)
}

// CreateEnquiryRequest creates an Enquiry request string
func CreateEnquiryRequest(enquiryName, t24User, messageData string) string {
	t24App := "ENQUIRY.SELECT"
	return createEnquiryRequest(t24App, enquiryName, t24User, "", messageData)
}

// CreateTransactionRequest creates a Transaction request string
func CreateTransactionRequest(t24App, t24Vers, t24Func, t24Auth, t24Process, t24User, t24Company, t24MsgID, t24MsgData string) string {
	if IsNullOrEmpty(t24App) && !IsNullOrEmpty(t24Vers) {
		t24App = strings.ToUpper(t24Vers)
		t24Vers = ""
	}
	return createTransactionRequest(t24App, t24Vers, t24Func, t24Auth, t24Process, t24User, "", t24Company, "", "", t24MsgID, t24MsgData)
}

// ExtractResponseData extracts response data from the Ofs response
func ExtractResponseData(ofsResponse string) (string, error) {
	if IsNullOrEmpty(ofsResponse) {
		return "", nil
	}
	firstCommaIndex, err := checkAndGetFirstCommaIndex(ofsResponse)
	if err != nil {
		return "", err
	}
	return ofsResponse[firstCommaIndex+1:], nil
}

// SplitByNumber splits a string into slices of a specified size
func SplitByNumber(str string, size int) []string {
	if size < 1 || IsNullOrEmpty(str) {
		return nil
	}
	var result []string
	for i := 0; i < len(str); i += size {
		end := i + size
		if end > len(str) {
			end = len(str)
		}
		result = append(result, str[i:end])
	}
	return result
}

// ParseDouble parses a string to a float64
func ParseDouble(s string) float64 {
	if IsNullOrEmpty(s) {
		return 0
	}
	result, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return result
}

// ReplaceWhiteSpacesStartEnd replaces leading and trailing white spaces with '#'
func ReplaceWhiteSpacesStartEnd(value string) string {
	result := value
	if strings.HasPrefix(result, " ") {
		result = "#" + result[1:]
	}
	if strings.HasSuffix(result, " ") {
		result = result[:len(result)-1] + "#"
	}
	return result
}

// ComparatorBySVPositionASC is a comparator for sorting OfsDataField by svPosition in ascending order
var ComparatorBySVPositionASC = func(s1, s2 *model.OfsDataField) int {
	svPosition1 := s1.SvPosition
	svPosition2 := s2.SvPosition
	return svPosition1 - svPosition2
}

func createEnquiryRequest(t24App, enquiryName, t24User, pass, messageData string) string {
	formatEnquiry := "%s,,%s/%s,%s"
	formatEnquiryMessageDate := "%s,,%s/%s,%s%s"

	if IsNullOrEmpty(messageData) {
		return fmt.Sprintf(formatEnquiry, t24App, t24User, pass, enquiryName)
	}
	return fmt.Sprintf(formatEnquiryMessageDate, t24App, t24User, pass, enquiryName, messageData)
}

func createTransactionRequest(t24App, t24Vers, t24Func, t24Auth, t24Process, t24User, t24Pass, t24Company, t24Replace, t24GTS, t24MsgID, t24MsgData string) string {
	t24App = strings.ToUpper(t24App)
	t24Vers = strings.ToUpper(t24Vers)
	t24Pass = "******"

	msg := fmt.Sprintf("%s,", t24App)
	msg += fmt.Sprintf("%s/", t24Vers)
	if !IsNullOrEmpty(t24Process) {
		msg += fmt.Sprintf("%s/", t24Func)
		msg += fmt.Sprintf("%s", t24Process)
	}

	if !IsNullOrEmpty(t24GTS) || !IsNullOrEmpty(t24Auth) {
		msg += fmt.Sprintf("/%s", t24GTS)
		if !IsNullOrEmpty(t24Auth) {
			msg += fmt.Sprintf("/%s", t24Auth)
		}
	}

	msg += fmt.Sprintf(",%s", t24User)
	msg += fmt.Sprintf("/%s", t24Pass)
	if !IsNullOrEmpty(t24Company) || "YES" == strings.ToUpper(t24Replace) {
		msg += fmt.Sprintf("/%s", t24Company)
		if "YES" == strings.ToUpper(t24Replace) {
			msg += "///1"
		}
	}

	msg += ","
	if !IsNullOrEmpty(t24MsgID) {
		msg += t24MsgID
	}

	if !IsNullOrEmpty(t24MsgData) {
		msg += ","
		msg += t24MsgData
	}

	msg = strings.ReplaceAll(msg, "//,", "/,")
	msg = strings.ReplaceAll(msg, "/,", ",")
	msg = strings.TrimSpace(msg)

	return msg
}

func checkAndGetFirstCommaIndex(ofsResponse string) (int, error) {
	firstCommaIndex := strings.Index(ofsResponse, ",")
	if firstCommaIndex < 0 {
		return -1, errors.New("[checkAndGetFirstCommaIndexInvalid] ofsResponse returned from T24 [" + ofsResponse + "]")
	}
	return firstCommaIndex, nil
}
