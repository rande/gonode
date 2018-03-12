package user

import "io"
import "github.com/rande/gonode/modules/base"

type PublicUserMeta struct {
}

type PublicUserData struct {
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
	NewPassword string   `json:"newpassword,omitempty"`
}

func UserSerializer(w io.Writer, node *base.Node) error {
	privateMeta := node.Meta.(*UserMeta)
	privateData := node.Data.(*User)

	publicMeta := &PublicUserMeta{}
	publicData := &PublicUserData{
		FirstName:   privateData.FirstName,
		LastName:    privateData.LastName,
		Email:       privateData.Email,
		DateOfBirth: privateData.DateOfBirth,
		Locked:      privateData.Locked,
		Enabled:     privateData.Enabled,
		Expired:     privateData.Expired,
		Roles:       privateData.Roles,
		Gender:      privateData.Gender,
		Locale:      privateData.Locale,
		Timezone:    privateData.Timezone,
		Username:    privateData.Username,
		NewPassword: privateData.NewPassword,
	}

	node.Meta = publicMeta
	node.Data = publicData

	// serialize
	err := base.Serialize(w, node)

	node.Meta = privateMeta
	node.Data = privateData

	if err != nil {
		return err
	}

	return nil
}

func UserDeserializer(r io.Reader, node *base.Node) error {
	publicMeta := &PublicUserMeta{}
	publicData := &PublicUserData{}

	privateMeta := node.Meta.(*UserMeta)
	privateData := node.Data.(*User)

	node.Meta = publicMeta
	node.Data = publicData

	err := base.Deserialize(r, node)

	node.Meta = privateMeta
	node.Data = privateData

	if err != nil {
		return err
	}

	privateData.FirstName = publicData.FirstName
	privateData.LastName = publicData.LastName
	privateData.Email = publicData.Email
	privateData.DateOfBirth = publicData.DateOfBirth
	privateData.Locked = publicData.Locked
	privateData.Enabled = publicData.Enabled
	privateData.Expired = publicData.Expired
	privateData.Roles = publicData.Roles
	privateData.Gender = publicData.Gender
	privateData.Locale = publicData.Locale
	privateData.Timezone = publicData.Timezone
	privateData.Username = publicData.Username
	privateData.NewPassword = publicData.NewPassword

	return nil
}
