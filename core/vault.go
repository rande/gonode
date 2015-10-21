package core

import (
	"github.com/rande/gonode/vault"
)

func GetVaultMetadata(node *Node) (meta vault.VaultMetadata) {

	// @todo: see if we can add more information here.
	meta = vault.NewVaultMetadata()
	meta["uuid"] = node.Uuid.CleanString()
	meta["type"] = node.Type
	meta["name"] = node.Name
	meta["created_at"] = node.CreatedAt
	meta["updated_at"] = node.UpdatedAt
	meta["meta"] = node.Meta
	meta["data"] = node.Data
	meta["revision"] = node.Revision
	meta["updated_by"] = node.UpdatedBy.CleanString()
	meta["created_by"] = node.CreatedBy.CleanString()
	meta["source"] = node.Source.CleanString()
	meta["set_uuid"] = node.SetUuid.CleanString()
	meta["parent_uuid"] = node.ParentUuid.CleanString()

	return
}
