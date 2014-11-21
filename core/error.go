package gonode

type revisionError struct {
	s string
}

func (e *revisionError) Error() string {
	return e.s
}

func NewRevisionError(message string) error {
	return &revisionError{message}
}
