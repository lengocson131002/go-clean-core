package txn

import (
	"strconv"
	"strings"
)

// AaOfsResponseDataParser is a struct implementing OfsResponseDataParser for the Aa case.
type AaOfsResponseDataParser struct {
	responseData string
}

// NewAaOfsResponseDataParser creates a new instance of AaOfsResponseDataParser.
func NewAaOfsResponseDataParser(responseData string) *AaOfsResponseDataParser {
	return &AaOfsResponseDataParser{responseData: responseData}
}

func (parser *AaOfsResponseDataParser) GetOfsName() string {
	firstEqualsPosition := strings.Index(parser.responseData, "=")
	return parser.responseData[:firstEqualsPosition]
}

func (parser *AaOfsResponseDataParser) GetMvPosition() int {
	firstColonPos := strings.Index(parser.responseData, ":")
	secondColonPos := strings.Index(parser.responseData[firstColonPos+1:], ":") + firstColonPos + 1
	mvPositionStr := parser.responseData[firstColonPos+1 : secondColonPos]

	mvPosition, err := strconv.Atoi(mvPositionStr)
	if err != nil {
		// TODO: handle error
		return -1
	}

	return mvPosition
}

func (parser *AaOfsResponseDataParser) GetSvPosition() int {
	firstColonPos := strings.Index(parser.responseData, ":")
	secondColonPos := strings.Index(parser.responseData[firstColonPos+1:], ":") + firstColonPos + 1
	svPositionStr := parser.responseData[secondColonPos+1:]

	svPosition, err := strconv.Atoi(svPositionStr)
	if err != nil {
		// TODO: handler error
		return -1
	}

	return svPosition
}

func (parser *AaOfsResponseDataParser) GetOfsValue() string {
	firstEqualsPos := strings.Index(parser.responseData, "=")
	firstColonPos := strings.Index(parser.responseData, ":")

	value := parser.responseData[firstEqualsPos+1 : firstColonPos]
	if len(value) > 1 && strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		return value[1 : len(value)-1]
	}
	return value
}
