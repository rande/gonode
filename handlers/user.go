package handlers

import (
	gn "github.com/rande/gonode/core"
)

type UserMeta struct {

}

type User struct {
	Name     string
	Login    string
	Password string
}

type UserHandler struct {

}

func (h *UserHandler) GetStruct() (gn.NodeData, gn.NodeMeta) {
	return &User{}, &UserMeta{}
}
