// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
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
	ErrValidation             = errors.New("unable to validate data")
	ErrRevision               = errors.New("wrong revision while saving")
	ErrNotFound               = errors.New("unable to find the node")
	ErrInvalidReferenceFormat = errors.New("unable to parse the reference")
	ErrAlreadyDeleted         = errors.New("unable to find the node")
	ErrNoStreamHandler        = errors.New("no stream handler defined")
	ErrAccessForbidden        = errors.New("access forbidden")
	ErrInvalidVersion         = errors.New("wrong node version")
	ErrInvalidUuidLength      = errors.New("invalid UUID length")
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
	case ErrNotFound:
		statusCode = http.StatusNotFound
	case ErrAlreadyDeleted:
		statusCode = http.StatusGone
	case ErrAccessForbidden, security.ErrAccessForbidden:
		statusCode = http.StatusForbidden
	case ErrRevision:
		statusCode = http.StatusConflict
	case ErrValidation:
		statusCode = http.StatusPreconditionFailed
	case ErrInvalidVersion:
		statusCode = http.StatusBadRequest
	}

	helper.SendWithHttpCode(res, statusCode, err.Error())
}
