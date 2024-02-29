package model

type ResponseStatus string

const (
	ResponseStatusSuccess  ResponseStatus = "SUCCESS"
	ResponseStatusError    ResponseStatus = "T24 ERROR"
	ResponseStatusOverride ResponseStatus = "T24 OVERRIDE"
	ResponseStatusOffline  ResponseStatus = "T24 OFFLINE"
)

func FromString(value string) ResponseStatus {
	switch value {
	case "SUCCESS":
		return ResponseStatusSuccess
	case "T24 ERROR":
		return ResponseStatusError
	case "T24 OVERRIDE":
		return ResponseStatusOverride
	case "T24 OFFLINE":
		return ResponseStatusOffline
	default:
		return ""
	}
}

func From(statusStr string) ResponseStatus {
	switch statusStr {
	case "1":
		return ResponseStatusSuccess
	case "-1":
		return ResponseStatusError
	case "-2":
		return ResponseStatusOverride
	case "-3":
		return ResponseStatusOffline
	default:
		return ""
	}
}
