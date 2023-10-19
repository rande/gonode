package fixtures

import (
	"strconv"

	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/blog"
	"github.com/rande/gonode/modules/feed"
	"github.com/rande/gonode/modules/media"
	"github.com/rande/gonode/modules/raw"
	"github.com/rande/gonode/modules/search"
	"github.com/rande/gonode/modules/user"
)

func GetFakeMediaNode(pos int) *base.Node {
	node := base.NewNode()

	node.Type = "media.image"
	node.Name = "The image " + strconv.Itoa(pos)
	node.Slug = "the-image-" + strconv.Itoa(pos)
	node.Data = &media.Image{
		Name:      "Go pic",
		Reference: "0x0",
	}
	node.Meta = &media.ImageMeta{
		SourceStatus: base.ProcessStatusInit,
	}

	return node
}

func GetFakePostNode(pos int) *base.Node {
	node := base.NewNode()

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

func GetFakeUserNode(pos int) *base.Node {
	node := base.NewNode()

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

func LoadFixtures(m *base.PgNodeManager, max int) error {

	var err error

	// create user
	admin := base.NewNode()

	admin.Nid = base.GetRootReference()
	admin.Type = "core.user"
	admin.Name = "The admin user"
	admin.Slug = "the-admin-user"
	admin.Data = &user.User{
		Username:    "admin",
		NewPassword: "admin",
		Roles:       []string{"ROLE_ADMIN", "ROLE_API", "node:api:master"},
	}
	admin.Meta = &user.UserMeta{
		PasswordCost: 12,
		PasswordAlgo: "bcrypt",
	}
	admin.Access = []string{"node:api:master", "node:api:read", "node:owner:XXXXXX"}

	m.Save(admin, false)

	for i := 1; i < max; i++ {
		node := GetFakeUserNode(i)
		node.UpdatedBy = admin.Nid
		node.CreatedBy = admin.Nid
		node.Access = []string{"node:api:master", "node:api:read", "node:owner:XXXXXX", "node:prism:render"}

		_, err = m.Save(node, false)

		helper.PanicOnError(err)
	}

	for i := 1; i < max; i++ {
		node := GetFakeMediaNode(i)
		node.UpdatedBy = admin.Nid
		node.CreatedBy = admin.Nid
		node.Access = []string{"node:api:master", "node:api:read", "node:owner:XXXXXX", "node:prism:render"}

		_, err = m.Save(node, false)

		helper.PanicOnError(err)
	}

	root := base.NewNode()
	root.Type = "core.root"
	root.Name = "Root path"
	root.Slug = "root-path"
	root.Meta = make(map[string]interface{})
	root.Data = make(map[string]interface{})
	root.Access = []string{"node:api:master", "node:api:read", "node:owner:XXXXXX", "node:prism:render"}

	_, err = m.Save(root, false)
	helper.PanicOnError(err)

	// create blog archives
	archive := base.NewNode()
	archive.Type = "search.index"
	archive.Name = "Blog Archive"
	archive.Slug = "blog"
	archive.Data = &search.Index{
		Type: search.NewParam([]string{"blog.post"}),
	}
	archive.Meta = &search.IndexMeta{}
	archive.Access = []string{
		"node:api:master",
		"node:api:read",
		"node:owner:XXXXXX",
		"node:prism:render",
	}

	_, err = m.Save(archive, false)
	helper.PanicOnError(err)

	_, err = m.Move(archive.Nid, root.Nid)
	helper.PanicOnError(err)

	for i := 1; i < max; i++ {
		node := GetFakePostNode(i)
		node.UpdatedBy = admin.Nid
		node.CreatedBy = admin.Nid
		node.Access = []string{"node:api:master", "node:api:read", "node:owner:XXXXXX", "node:prism:render"}

		_, err = m.Save(node, false)
		helper.PanicOnError(err)

		_, err = m.Move(node.Nid, archive.Nid)
		helper.PanicOnError(err)
	}

	// create feed archives
	f := base.NewNode()
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
	f.Access = []string{"node:api:master", "node:api:read", "node:owner:XXXXXX", "node:prism:render"}

	_, err = m.Save(f, false)
	helper.PanicOnError(err)

	_, err = m.Move(f.Nid, root.Nid)
	helper.PanicOnError(err)

	// create human.txt
	h := base.NewNode()
	h.Type = "core.raw"
	h.Name = "human.txt"
	h.Slug = "human.txt"
	h.Data = &raw.Raw{
		Name:        "human.txt",
		Content:     []byte("The human file"),
		ContentType: "text/plain",
	}
	h.Meta = &raw.RawMeta{}
	h.Access = []string{"node:api:master", "node:api:read", "node:owner:XXXXXX", "node:prism:render"}

	_, err = m.Save(h, false)
	helper.PanicOnError(err)

	_, err = m.Move(h.Nid, root.Nid)
	helper.PanicOnError(err)

	// create real image
	image := base.NewNode()

	image.Type = "media.image"
	image.Name = "The image for resize"
	image.Slug = "the-image-for-resize"
	image.Data = &media.Image{
		Name:      "Resize.jpg",
		Reference: "0x0",
		// from: https://github.com/nfnt/resize
		SourceUrl: "https://camo.githubusercontent.com/ef6fdc21c7c8e17354524f0982cdb52885335191/687474703a2f2f6e666e742e6769746875622e636f6d2f696d672f494d475f333639345f3732302e6a7067",
	}
	image.Meta = &media.ImageMeta{
		SourceStatus: base.ProcessStatusInit,
	}
	image.Access = []string{"node:api:master", "node:api:read", "node:owner:XXXXXX", "node:prism:render"}

	_, err = m.Save(image, false)

	helper.PanicOnError(err)

	return nil
}
