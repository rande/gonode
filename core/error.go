package core

var (
	ValidationError = &validationError{"Unable to validate date"}
	RevisionError   = &revisionError{"Wrong revision while saving"}
)

type validationError struct {
	message string
}

func (e *validationError) Error() string {
	return e.message
}

type revisionError struct {
	s string
}

func (e *revisionError) Error() string {
	return e.s
}

func NewRevisionError(message string) error {
	return &revisionError{message}
}

func NewErrors() Errors {
	return Errors{}
}

type Errors map[string][]string

func (es Errors) AddError(field string, message string) {

	if _, ok := es[field]; !ok {
		es[field] = []string{}
	}

	es[field] = append(es[field], message)
}

func (es Errors) HasError(field string) bool {
	if _, ok := es[field]; !ok {
		return false
	}

	return len(es[field]) > 0
}

func (es Errors) GetError(field string) []string {
	if _, ok := es[field]; !ok {
		return nil
	}

	return es[field]
}

func (es Errors) HasErrors() bool {

	for _, errors := range es {
		if len(errors) > 0 {
			return true
		}
	}

	return false
}
