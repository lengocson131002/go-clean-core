package txn

import (
	"fmt"
	"regexp"

	"github.com/lengocson131002/go-clean/pkg/t24/response/parser"
)

type ofsResponseParserFactory struct {
}

// GetOfsResponseParser implements parser.OfsResponseParserFactory.
func (*ofsResponseParserFactory) GetOfsResponseParser(responseData string) (parser.OfsResponseDataParser, error) {
	defaultPattern := regexp.MustCompile(`(\S+:\d+:\d+=[\S\s]*)`)
	aaPattern := regexp.MustCompile(`(\S+=[\S\s]*:\d+:\d+)`)

	if defaultPattern.MatchString(responseData) {
		return NewDefaultOfsResponseDataParser(responseData), nil
	}
	if aaPattern.MatchString(responseData) {
		return NewAaOfsResponseDataParser(responseData), nil
	}

	return nil, fmt.Errorf("Invalid ofs response data part encountered [" + responseData + "]")
}

func NewTxnOfsResponseParserFactory() parser.OfsResponseParserFactory {
	return &ofsResponseParserFactory{}
}
