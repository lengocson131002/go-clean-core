package parser

import (
	"errors"
	"fmt"
	"html"
	"strings"

	"github.com/lengocson131002/go-clean/pkg/t24/response/model"
)

// OfsResponseDataParser is an interface for parsing OFS response data.
type OfsResponseDataParser interface {
	GetOfsName() string
	GetOfsValue() string
	GetMvPosition() int
	GetSvPosition() int
}

type ResponseDataParser interface {
	ParseResponseData(responseName, ofsResponse string) (model.Response, error)
}

type OfsResponseParserFactory interface {
	GetOfsResponseParser(responseData string) (OfsResponseDataParser, error)
}

type AbstractParser struct {
	Parserfactory OfsResponseParserFactory
}

func (a *AbstractParser) ParseData(ofsResponseData string) ([]model.OfsDataField, error) {
	if ofsResponseData == "" {
		return nil, errors.New("ofsResponseData is empty")
	}

	responseDataList := make([]model.OfsDataField, 0)

	splitedData := strings.Split(ofsResponseData, ",")
	for pos := len(splitedData) - 1; pos >= 0; pos-- {
		if strings.Index(splitedData[pos], "=") == -1 && pos > 0 {
			splitedData[pos-1] = splitedData[pos-1] + "," + splitedData[pos]
		} else {
			responseData := splitedData[pos]

			if strings.HasPrefix(responseData, "OVERRIDE:") {
				responseData = a.unescapeXML(responseData)
				override := strings.Split(responseData[9:], "=")
				fieldID := override[0]
				overString := a.parseOfsOverride(override[1])
				responseData = fieldID + "=" + overString
			}

			ofsField, err := a.buildOfsField(responseData)
			if err != nil {
				return nil, err
			}
			responseDataList = append(responseDataList, ofsField)
		}
	}

	return responseDataList, nil
}

func (a *AbstractParser) unescapeXML(input string) string {
	return html.UnescapeString(input)
}

func (a *AbstractParser) parseOfsOverride(ofsOverride string) string {
	overrideContent := ofsOverride

	msgParts := strings.Split(ofsOverride, "{")
	var msgFormat string
	var msgParameters []string

	if len(msgParts) > 0 {
		iPos1 := strings.Index(msgParts[0], "}")

		if iPos1 >= 0 {
			msgFormat = msgParts[0][iPos1+1:]
			overrideContent = msgFormat
		} else {
			msgFormat = msgParts[0]
		}

		if len(msgParts) > 1 {
			partAttrIndex := 1

			msgParameters = strings.Split(msgParts[1], "}")

			if len(msgParameters) > 0 {
				msgOverride := strings.Split(msgFormat, "&")

				var sb strings.Builder

				if len(msgOverride) > 0 {
					for i := 0; i < len(msgOverride); i++ {
						sb.WriteString(msgOverride[i])

						if i < len(msgParameters) {
							sb.WriteString(msgParameters[i])
						}
					}
				}

				if sb.Len() > 0 {
					overrideContent = sb.String()
				}

				partAttrIndex = 2
			}

			if partAttrIndex < len(msgParts) {
				var sb strings.Builder
				sb.WriteString(overrideContent)

				for partAttrIndex < len(msgParts) {
					if len(msgParts[partAttrIndex]) > 0 {
						switch attrIndex := partAttrIndex - 1; attrIndex {
						case 0:
							sb.WriteString(fmt.Sprintf("(Currency[%s]) ", msgParts[partAttrIndex]))
						case 1:
							sb.WriteString(fmt.Sprintf("(Amount[%s]) ", msgParts[partAttrIndex]))
						case 2:
							sb.WriteString(fmt.Sprintf("(Account[%s]) ", msgParts[partAttrIndex]))
						case 3:
							sb.WriteString(fmt.Sprintf("(Customer[%s]) ", msgParts[partAttrIndex]))
						case 4:
							sb.WriteString(fmt.Sprintf("(Transaction[%s]) ", msgParts[partAttrIndex]))
						case 5:
							sb.WriteString(fmt.Sprintf("(Limit[%s]) ", msgParts[partAttrIndex]))
						case 6:
							sb.WriteString(fmt.Sprintf("(Restriction[%s]) ", msgParts[partAttrIndex]))
						}
					}
					partAttrIndex++
				}
				overrideContent = sb.String()
			}
		}
	}
	return overrideContent
}

func (a *AbstractParser) buildOfsField(responseData string) (model.OfsDataField, error) {
	ofsParser, err := a.Parserfactory.GetOfsResponseParser(responseData)
	if err != nil {
		return model.OfsDataField{}, err
	}

	return model.OfsDataField{
		OfsName:    ofsParser.GetOfsName(),
		Value:      ofsParser.GetOfsValue(),
		MvPosition: ofsParser.GetMvPosition(),
		SvPosition: ofsParser.GetSvPosition(),
	}, nil
}
