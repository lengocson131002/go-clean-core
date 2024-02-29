package util

type OfsOperatorConstant string

const (
	EQ   OfsOperatorConstant = "EQ"
	GE   OfsOperatorConstant = "GE"
	GT   OfsOperatorConstant = "GT"
	LE   OfsOperatorConstant = "LE"
	LK   OfsOperatorConstant = "LK"
	LT   OfsOperatorConstant = "LT"
	NE   OfsOperatorConstant = "NE"
	UL   OfsOperatorConstant = "UL"
	RG   OfsOperatorConstant = "RG"
	NR   OfsOperatorConstant = "NR"
	CT   OfsOperatorConstant = "CT"
	NC   OfsOperatorConstant = "NC"
	BW   OfsOperatorConstant = "BW"
	EW   OfsOperatorConstant = "EW"
	DNBW OfsOperatorConstant = "DNBW"
	DNEW OfsOperatorConstant = "DNEW"
	SAID OfsOperatorConstant = "SAID"
)

type OfsConstant string

const (
	T24_SUCCESS            OfsConstant = "SUCCESS"
	T24_ERROR              OfsConstant = "T24 ERROR"
	T24_OVERRIDE           OfsConstant = "T24 OVERRIDE"
	T24_OFFLINE            OfsConstant = "T24 OFFLINE"
	T24_FUNCTION_INPUT     OfsConstant = "I"
	T24_FUNCTION_SHOW      OfsConstant = "S"
	T24_FUNCTION_AUTHORISE OfsConstant = "A"
	T24_FUNCTION_REVERSE   OfsConstant = "R"
	T24_FUNCTION_DELETE    OfsConstant = "D"
	T24_OPERATION_PROCESS  OfsConstant = "PROCESS"
	T24_OPERATION_VALIDATE OfsConstant = "VALIDATE"
)
