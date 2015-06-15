package handlers

import (
	nc "github.com/rande/gonode/core"
	"github.com/stretchr/testify/assert"
	"testing"
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

	node.Data.(*User).Email = "invalid email"
	node.Data.(*User).Gender = "v"

	errors := nc.NewErrors()
	manager := &nc.MockedManager{}

	handler.Validate(node, manager, errors)

	a.Equal(3, len(errors))
	a.True(errors.HasErrors())

	a.True(errors.HasError("data.login"))
	a.Equal([]string{"Login cannot be empty"}, errors.GetError("data.login"))

	a.True(errors.HasError("data.email"))
	a.Equal([]string{"Email is not valid"}, errors.GetError("data.email"))

	a.True(errors.HasError("data.gender"))
	a.Equal([]string{"Invalid gender code"}, errors.GetError("data.gender"))
}

func GeneratePasswordTest(t *testing.T) {
	a := assert.New(t)
	handler, node := GetUserHandleNode()

	node.Data.(*User).NewPassword = "password"

	manager := &nc.MockedManager{}

	a.False(len(node.Data.(*User).Password) > 0)

	handler.PreInsert(node, manager)

	a.Equal(0, len(node.Data.(*User).NewPassword))
	a.True(len(node.Data.(*User).Password) > 0)
}

func Test_UserHandler_GeneratePassword_PreInsert(t *testing.T) {
	GeneratePasswordTest(t)
}

func Test_UserHandler_GeneratePassword_PreUpdate(t *testing.T) {
	GeneratePasswordTest(t)
}
