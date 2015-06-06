package handlers

import (
	v "github.com/asaskevich/govalidator"
	nc "github.com/rande/gonode/core"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

var (
	validPassword, _  = regexp.Compile("{([a-zA-Z0-9]*)}(.*)")
	validPasswordAlgo = []string{"plain", "md5", "bcrypt"}
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
	Login       string   `json:"login"`
	Password    string   `json:"password"`
	NewPassword string   `json:"newpassword,omitempty"`
}

type UserHandler struct {
}

func (h *UserHandler) GetStruct() (nc.NodeData, nc.NodeMeta) {
	return &User{}, &UserMeta{
		PasswordCost: 12,
		PasswordAlgo: "bcrypt",
	}
}

func (h *UserHandler) PreInsert(node *nc.Node, m nc.NodeManager) error {
	updatePassword(node)

	return nil
}

func (h *UserHandler) PreUpdate(node *nc.Node, m nc.NodeManager) error {
	updatePassword(node)

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

	if data.Login == "" {
		errors.AddError("data.login", "Login cannot be empty")
	}

	if data.Email != "" && !v.IsEmail(data.Email) {
		errors.AddError("data.email", "Email is not valid")
	}

	if data.Gender != "" && (data.Gender != USER_GENDER_FEMALE && data.Gender != USER_GENDER_MALE) {
		errors.AddError("data.gender", "Invalid gender code")
	}
}

func (h *UserHandler) GetDownloadData(node *nc.Node) *nc.DownloadData {
	return nc.GetDownloadData()
}

func (h *UserHandler) Load(data []byte, meta []byte, node *nc.Node) error {
	return nc.HandlerLoad(h, data, meta, node)
}

func updatePassword(node *nc.Node) error {
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
