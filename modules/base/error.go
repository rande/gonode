// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"errors"
	"net/http"

	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/security"
)

var (
	ValidationError             = errors.New("Unable to validate data")
	RevisionError               = errors.New("Wrong revision while saving")
	NotFoundError               = errors.New("Unable to find the node")
	InvalidReferenceFormatError = errors.New("Unable to parse the reference")
	AlreadyDeletedError         = errors.New("Unable to find the node")
	NoStreamHandler             = errors.New("No stream handler defined")
	AccessForbiddenError        = errors.New("Access forbidden")
	InvalidVersionError         = errors.New("Wrong node version")
	InvalidUuidLengthError      = errors.New("Invalid UUID length")
)

type validationError struct {
	message string
}

func (e *validationError) Error() string {
	return e.message
}

type alreadyDeletedError struct {
	message string
}

func (e *alreadyDeletedError) Error() string {
	return e.message
}

type noStreamHandlerError struct {
	message string
}

func (e *noStreamHandlerError) Error() string {
	return e.message
}

type revisionError struct {
	s string
}

func (e *revisionError) Error() string {
	return e.s
}

type notFoundError struct {
	message string
}

func (e *notFoundError) Error() string {
	return e.message
}

type invalidReferenceFormatError struct {
	message string
}

func (e *invalidReferenceFormatError) Error() string {
	return e.message
}

func NewRevisionError(message string) error {
	return &revisionError{message}
}

// use for model validation
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

func HandleError(req *http.Request, res http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	statusCode := http.StatusInternalServerError

	switch err {
	case NotFoundError:
		statusCode = http.StatusNotFound
	case AlreadyDeletedError:
		statusCode = http.StatusGone
	case AccessForbiddenError, security.AccessForbiddenError:
		statusCode = http.StatusForbidden
	case RevisionError:
		statusCode = http.StatusConflict
	case ValidationError:
		statusCode = http.StatusPreconditionFailed
	case InvalidVersionError:
		statusCode = http.StatusBadRequest
	}

	helper.SendWithHttpCode(res, statusCode, err.Error())
}
