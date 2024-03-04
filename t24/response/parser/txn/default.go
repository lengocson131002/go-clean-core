package txn

import (
	"strconv"
	"strings"

	"github.com/lengocson131002/go-clean/pkg/t24/response/parser"
)

// DefaultOfsResponseDataParser is a struct implementing OfsResponseDataParser for the default case.
type DefaultOfsResponseDataParser struct {
	responseData string
}

// NewDefaultOfsResponseDataParser creates a new instance of DefaultOfsResponseDataParser.
func NewDefaultOfsResponseDataParser(responseData string) parser.OfsResponseDataParser {
	return &DefaultOfsResponseDataParser{responseData: responseData}
}

func (parser *DefaultOfsResponseDataParser) GetOfsName() string {
	firstColonPos := strings.Index(parser.responseData, ":")
	return parser.responseData[:firstColonPos]
}

func (parser *DefaultOfsResponseDataParser) GetSvPosition() int {
	firstColonPos := strings.Index(parser.responseData, ":")
	secondColonPos := strings.Index(parser.responseData[firstColonPos+1:], ":") + firstColonPos + 1
	equalToPos := strings.Index(parser.responseData, "=")
	svPositionStr := parser.responseData[secondColonPos+1 : equalToPos]

	svPosition, err := strconv.Atoi(svPositionStr)
	if err != nil {
		// TODO: handle error
		return -1
	}

	return svPosition
}

func (parser *DefaultOfsResponseDataParser) GetMvPosition() int {
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

func (parser *DefaultOfsResponseDataParser) GetOfsValue() string {
	firstEqualsPos := strings.Index(parser.responseData, "=")
	value := parser.responseData[firstEqualsPos+1:]
	if len(value) > 1 && strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		return value[1 : len(value)-1]
	}
	return value
}
