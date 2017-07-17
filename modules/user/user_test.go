// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package user

import (
	"testing"

	"github.com/rande/gonode/modules/base"
	"github.com/stretchr/testify/assert"
)

func GetUserHandleNode() (base.Handler, *base.Node) {
	node := base.NewNode()
	handler := &UserHandler{}

	node.Data, node.Meta = handler.GetStruct()

	node.Meta.(*UserMeta).PasswordCost = 1 // speed up test

	return handler, node
}

func Test_UserHandler_Validate_EmptyData(t *testing.T) {
	a := assert.New(t)

	handler, node := GetUserHandleNode()
	a.IsType(&UserMeta{}, node.Meta)
	a.IsType(&User{}, node.Data)

	node.Data.(*User).Email = "invalid email"
	node.Data.(*User).Gender = "v"

	errors := base.NewErrors()
	manager := &base.MockedManager{}

	if h, ok := handler.(base.ValidateNodeHandler); ok {
		h.Validate(node, manager, errors)
	} else {
		a.Fail("handler does not implement base.ValidateNodeHandler")
	}

	a.Equal(3, len(errors))
	a.True(errors.HasErrors())

	a.True(errors.HasError("data.username"))
	a.Equal([]string{"Username cannot be empty"}, errors.GetError("data.username"))

	a.True(errors.HasError("data.email"))
	a.Equal([]string{"Email is not valid"}, errors.GetError("data.email"))

	a.True(errors.HasError("data.gender"))
	a.Equal([]string{"Invalid gender code"}, errors.GetError("data.gender"))
}

func GeneratePasswordTest(t *testing.T) {
	a := assert.New(t)
	handler, node := GetUserHandleNode()

	node.Data.(*User).NewPassword = "password"

	manager := &base.MockedManager{}

	a.False(len(node.Data.(*User).Password) > 0)

	if h, ok := handler.(base.DatabaseNodeHandler); ok {
		h.PreInsert(node, manager)
	} else {
		a.Fail("handler does not implement base.DatabaseNodeHandler")
	}

	a.Equal(0, len(node.Data.(*User).NewPassword))
	a.True(len(node.Data.(*User).Password) > 0)
}

func Test_UserHandler_GeneratePassword_PreInsert(t *testing.T) {
	GeneratePasswordTest(t)
}

func Test_UserHandler_GeneratePassword_PreUpdate(t *testing.T) {
	GeneratePasswordTest(t)
}
