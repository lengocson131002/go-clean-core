package txn

import (
	"regexp"

	"github.com/lengocson131002/go-clean/pkg/t24/util"
)

func IsTxnResponseSuccess(ofsResponse string) bool {
	if ofsResponse == "" {
		return false
	}

	pattern := regexp.MustCompile(`\S*/\S*/(1)\S*,(((\S+:\d+:\d+=[\S\s]*)|(\S+=[\S\s]*:\d+:\d+))(,*))+`)
	if !pattern.MatchString(ofsResponse) {
		return false
	}

	vRes, err := util.ValidateTransactionResponse(ofsResponse)
	if err != nil {
		return false
	}

	return vRes
}

func IsTxnResponseFailure(ofsResponse string) bool {
	return ofsResponse != "" && regexp.MustCompile(`\S*/\S*/(-1|-2|-3)/[\S\s]*`).MatchString(ofsResponse)
}
