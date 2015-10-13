package fixtures

import (
	nc "github.com/rande/gonode/core"
	nh "github.com/rande/gonode/handlers"
	"strconv"
)

func GetFakeMediaNode(pos int) *nc.Node {
	node := nc.NewNode()

	node.Type = "media.image"
	node.Name = "The image " + strconv.Itoa(pos)
	node.Slug = "the-image-" + strconv.Itoa(pos)
	node.Data = &nh.Image{
		Name:      "Go pic",
		Reference: "0x0",
	}
	node.Meta = &nh.ImageMeta{}

	return node
}

func GetFakePostNode(pos int) *nc.Node {
	node := nc.NewNode()

	node.Type = "blog.post"
	node.Name = "The blog post " + strconv.Itoa(pos)
	node.Slug = "the-blog-post-" + strconv.Itoa(pos)
	node.Data = &nh.Post{
		Title:   "Go pic",
		Content: "The Content of my blog post",
		Tags:    []string{"sport", "tennis", "soccer"},
	}
	node.Meta = &nh.PostMeta{
		Format: "markdown",
	}

	return node
}

func GetFakeUserNode(pos int) *nc.Node {
	node := nc.NewNode()

	node.Type = "core.user"
	node.Name = "The user " + strconv.Itoa(pos)
	node.Slug = "the-user-" + strconv.Itoa(pos)
	node.Data = &nh.User{
		Login:       "user" + strconv.Itoa(pos),
		NewPassword: "user" + strconv.Itoa(pos),
	}
	node.Meta = &nh.UserMeta{
		PasswordCost: 12,
		PasswordAlgo: "bcrypt",
	}

	return node
}

func LoadFixtures(m *nc.PgNodeManager, max int) error {

	var err error

	// create user
	admin := nc.NewNode()

	admin.Uuid = nc.GetRootReference()
	admin.Type = "core.user"
	admin.Name = "The admin user"
	admin.Slug = "the-admin-user"
	admin.Data = &nh.User{
		Login:       "admin",
		NewPassword: "admin",
	}
	admin.Meta = &nh.UserMeta{
		PasswordCost: 12,
		PasswordAlgo: "bcrypt",
	}

	m.Save(admin)

	for i := 1; i < max; i++ {
		node := GetFakeUserNode(i)
		node.UpdatedBy = admin.Uuid
		node.CreatedBy = admin.Uuid

		_, err = m.Save(node)

		nc.PanicOnError(err)
	}

	for i := 1; i < max; i++ {
		node := GetFakeMediaNode(i)
		node.UpdatedBy = admin.Uuid
		node.CreatedBy = admin.Uuid

		_, err = m.Save(node)

		nc.PanicOnError(err)
	}

	for i := 1; i < max; i++ {
		node := GetFakePostNode(i)
		node.UpdatedBy = admin.Uuid
		node.CreatedBy = admin.Uuid

		_, err = m.Save(node)

		nc.PanicOnError(err)
	}

	return nil
}
