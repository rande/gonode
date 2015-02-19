package handlers

import (
	"github.com/stretchr/testify/assert"
	"testing"
	nc "github.com/rande/gonode/core"
	"github.com/rande/gonode/test/mock"
)

func GetUserHandleNode() (nc.Handler, *nc.Node) {
	node := nc.NewNode()
	handler := &UserHandler{}

	node.Data, node.Meta = handler.GetStruct()

	return handler, node
}

func Test_UserHandler_Validate_EmptyData(t *testing.T) {
	a := assert.New(t)

	handler, node := GetUserHandleNode()
	a.IsType(&UserMeta{}, node.Meta)
	a.IsType(&User{}, node.Data)

	errors := nc.NewErrors()
	manager := &mock.MockedManager{}

	handler.Validate(node, manager, errors)

	a.Equal(3, len(errors))
	a.True(errors.HasErrors())

	a.True(errors.HasError("data.login"))
	a.Equal([]string{"Login cannot be empty"}, errors.GetError("data.login"))

	a.True(errors.HasError("data.name"))
	a.Equal([]string{"Name cannot be empty"}, errors.GetError("data.name"))

	a.True(errors.HasError("data.password"))
	a.Equal([]string{"Password cannot be empty"}, errors.GetError("data.password"))
}

func Test_UserHandler_Validate_InvalidPassword(t *testing.T) {

	a := assert.New(t)
	handler, node := GetUserHandleNode()

	node.Data.(*User).Password = "password"

	errors := nc.NewErrors()
	manager := &mock.MockedManager{}
	handler.Validate(node, manager, errors)

	a.True(errors.HasError("data.password"))
	a.Equal([]string{"Invalid password format"}, errors.GetError("data.password"))
}

func Test_UserHandler_Validate_ValidPassword(t *testing.T) {

	a := assert.New(t)
	handler, node := GetUserHandleNode()

	node.Data.(*User).Password = "{plain}password"

	errors := nc.NewErrors()
	manager := &mock.MockedManager{}
	handler.Validate(node, manager, errors)

	a.False(errors.HasError("data.password"))
}

func Test_UserHandler_Validate_InvalidPasswordAlgo(t *testing.T) {
	a := assert.New(t)
	handler, node := GetUserHandleNode()

	node.Data.(*User).Password = "{wrong}password"

	errors := nc.NewErrors()
	manager := &mock.MockedManager{}
	handler.Validate(node, manager, errors)

	a.True(errors.HasError("data.password"))
	a.Equal([]string{"Invalid algorithm selected"}, errors.GetError("data.password"))
}
