package fixtures

import (
	"github.com/rande/gonode/core"
	"github.com/rande/gonode/modules/blog"
	"github.com/rande/gonode/modules/feed"
	"github.com/rande/gonode/modules/media"
	"github.com/rande/gonode/modules/search"
	"github.com/rande/gonode/modules/user"
	"strconv"
)

func GetFakeMediaNode(pos int) *core.Node {
	node := core.NewNode()

	node.Type = "media.image"
	node.Name = "The image " + strconv.Itoa(pos)
	node.Slug = "the-image-" + strconv.Itoa(pos)
	node.Data = &media.Image{
		Name:      "Go pic",
		Reference: "0x0",
	}
	node.Meta = &media.ImageMeta{}

	return node
}

func GetFakePostNode(pos int) *core.Node {
	node := core.NewNode()

	node.Type = "blog.post"
	node.Name = "The blog post " + strconv.Itoa(pos)
	node.Slug = "the-blog-post-" + strconv.Itoa(pos)
	node.Data = &blog.Post{
		Title:   "Go pic",
		Content: "The Content of my blog post",
		Tags:    []string{"sport", "tennis", "soccer"},
	}
	node.Meta = &blog.PostMeta{
		Format: "markdown",
	}

	return node
}

func GetFakeUserNode(pos int) *core.Node {
	node := core.NewNode()

	node.Type = "core.user"
	node.Name = "The user " + strconv.Itoa(pos)
	node.Slug = "the-user-" + strconv.Itoa(pos)
	node.Data = &user.User{
		Username:    "user" + strconv.Itoa(pos),
		NewPassword: "user" + strconv.Itoa(pos),
	}
	node.Meta = &user.UserMeta{
		PasswordCost: 1,
		PasswordAlgo: "bcrypt",
	}

	return node
}

func LoadFixtures(m *core.PgNodeManager, max int) error {

	var err error

	// create user
	admin := core.NewNode()

	admin.Uuid = core.GetRootReference()
	admin.Type = "core.user"
	admin.Name = "The admin user"
	admin.Slug = "the-admin-user"
	admin.Data = &user.User{
		Username:    "admin",
		NewPassword: "admin",
	}
	admin.Meta = &user.UserMeta{
		PasswordCost: 12,
		PasswordAlgo: "bcrypt",
	}

	m.Save(admin, false)

	for i := 1; i < max; i++ {
		node := GetFakeUserNode(i)
		node.UpdatedBy = admin.Uuid
		node.CreatedBy = admin.Uuid

		_, err = m.Save(node, false)

		core.PanicOnError(err)
	}

	for i := 1; i < max; i++ {
		node := GetFakeMediaNode(i)
		node.UpdatedBy = admin.Uuid
		node.CreatedBy = admin.Uuid

		_, err = m.Save(node, false)

		core.PanicOnError(err)
	}

	for i := 1; i < max; i++ {
		node := GetFakePostNode(i)
		node.UpdatedBy = admin.Uuid
		node.CreatedBy = admin.Uuid

		_, err = m.Save(node, false)

		core.PanicOnError(err)
	}

	// create blog archives
	archive := core.NewNode()
	archive.Type = "core.index"
	archive.Name = "Blog Archive"
	archive.Slug = "blog"
	archive.Data = &search.Index{
		Type: search.NewParam([]string{"blog.post"}),
	}
	archive.Meta = &search.IndexMeta{}

	_, err = m.Save(archive, false)

	core.PanicOnError(err)

	// create feed archives
	f := core.NewNode()
	f.Type = "feed.index"
	f.Name = "Feed Archive"
	f.Slug = "feed"
	f.Data = &feed.Feed{
		Title:       "Archive blog",
		Description: "This is a description.",
		Index: &search.Index{
			Type: search.NewParam([]string{"blog.post"}),
		},
	}
	f.Meta = &search.IndexMeta{}

	_, err = m.Save(f, false)

	core.PanicOnError(err)

	return nil
}
