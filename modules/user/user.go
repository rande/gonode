// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package user

import (
	"regexp"

	v "github.com/asaskevich/govalidator"
	"github.com/rande/gonode/modules/base"
	"golang.org/x/crypto/bcrypt"
)

var (
	validPassword, _  = regexp.Compile("{([a-zA-Z0-9]*)}(.*)")
	validPasswordAlgo = []string{"plain", "bcrypt"}
)

const (
	USER_GENDER_MALE   = "m"
	USER_GENDER_FEMALE = "f"
)

type UserMeta struct {
	PasswordCost int    `json:"password_cost"`
	PasswordAlgo string `json:"password_algo"`
}

type User struct {
	FirstName   string   `json:"firstname"`
	LastName    string   `json:"lastname"`
	Email       string   `json:"email"`
	DateOfBirth string   `json:"dateofbirth"`
	Locked      bool     `json:"locked"`
	Enabled     bool     `json:"enabled"`
	Expired     bool     `json:"expired"`
	Roles       []string `json:"roles"`
	Gender      string   `json:"gender"`
	Locale      string   `json:"locale"`
	Timezone    string   `json:"timezone"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	NewPassword string   `json:"newpassword,omitempty"`
}

func (u *User) GetRoles() []string {
	return u.Roles
}

func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) GetUsername() string {
	return u.Username
}

type UserHandler struct {
}

func (h *UserHandler) GetStruct() (base.NodeData, base.NodeMeta) {
	return &User{}, &UserMeta{
		PasswordCost: 12,
		PasswordAlgo: "bcrypt",
	}
}

func (h *UserHandler) PreInsert(node *base.Node, m base.NodeManager) error {
	updatePassword(node)

	return nil
}

func (h *UserHandler) PreUpdate(node *base.Node, m base.NodeManager) error {
	updatePassword(node)

	return nil
}

func (h *UserHandler) PostInsert(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *UserHandler) PostUpdate(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *UserHandler) Validate(node *base.Node, m base.NodeManager, errors base.Errors) {
	data := node.Data.(*User)

	if data.Username == "" {
		errors.AddError("data.username", "Username cannot be empty")
	}

	if data.Email != "" && !v.IsEmail(data.Email) {
		errors.AddError("data.email", "Email is not valid")
	}

	if data.Gender != "" && (data.Gender != USER_GENDER_FEMALE && data.Gender != USER_GENDER_MALE) {
		errors.AddError("data.gender", "Invalid gender code")
	}
}

func updatePassword(node *base.Node) error {
	data := node.Data.(*User)
	meta := node.Meta.(*UserMeta)

	if data.NewPassword == "" {
		return nil
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), meta.PasswordCost)

	if err != nil {
		return err
	}

	data.Password = string(password)
	data.NewPassword = ""

	return nil
}
