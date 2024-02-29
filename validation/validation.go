package validation

// interface for validation
type Validator interface {
	Validate(obj interface{}) error
}
