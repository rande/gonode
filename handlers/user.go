package handlers

import (
	nc "github.com/rande/gonode/core"
	"regexp"
)

var (
	validPassword, _ = regexp.Compile("{([a-zA-Z0-9]*)}(.*)");
	validPasswordAlgo = []string{"plain", "md5", "bcrypt"}
)

type UserMeta struct {

}

type User struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserHandler struct {

}

func (h *UserHandler) GetStruct() (nc.NodeData, nc.NodeMeta) {
	return &User{}, &UserMeta{}
}

func (h *UserHandler) PreInsert(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *UserHandler) PreUpdate(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *UserHandler) PostInsert(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *UserHandler) PostUpdate(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *UserHandler) Validate(node *nc.Node, m nc.NodeManager, errors nc.Errors) {

	data := node.Data.(*User)

	if (data.Login == "") {
		errors.AddError("data.login", "Login cannot be empty")
	}

	if (data.Name == "") {
		errors.AddError("data.name", "Name cannot be empty")
	}

	if (data.Password == "") {
		errors.AddError("data.password", "Password cannot be empty")
	} else if (!validPassword.Match([]byte(data.Password))) {
		errors.AddError("data.password", "Invalid password format")
	}

	if (!errors.HasError("data.password")) {
		result := validPassword.FindStringSubmatch(data.Password)

		invalid := true
		for _, v := range validPasswordAlgo {
			if v == result[1] {
				invalid = false
			}
		}

		if invalid {
			errors.AddError("data.password", "Invalid algorithm selected")
		}
	}
}
