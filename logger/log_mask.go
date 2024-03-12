package logger

import (
	"regexp"
	"strings"
)

var (
	MaskedPatterns = []string{
		`\"password\"\s*:\s*\"(.*?)\"`,
		`<password>\s*(.*?)\s*<\/password>`,
		`<credentials>\s*(.*?)\s*<\/credentials>`,
		`<!\[CDATA\[\s*<password>\s*(.*?)\s*</password>\s*\]\]>`,
		`<!\[CDATA\[\s*<credentials>\s*(.*?)\s*</credentials>\s*\]\]>`,
		`((?:[A-Za-z0-9+\/]{64,}){1,}(?:[A-Za-z0-9+\/]{2}==|[A-Za-z0-9+\/]{3}=)?)`,
	}
)

func MaskSensitiveData(input string) string {
	for _, pattern := range MaskedPatterns {
		builder := strings.Builder{}
		builder.WriteString(input)

		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatchIndex(input, -1)
		for _, match := range matches {
			for i := 2; i < len(match); i += 2 {
				if match[i] != -1 && match[i+1] != -1 {
					input = input[:match[i]] + mask(input, match[i], match[i+1]) + input[match[i+1]:]
				}
			}
		}
	}
	return input
}

func mask(input string, from int, to int) string {
	if from < 0 && to < 0 {
		return ""
	}

	sb := strings.Builder{}
	for range input[from:to] {
		sb.WriteString("*")
	}

	return sb.String()
}
