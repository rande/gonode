package base

import (
	"github.com/rande/gonode/core/vault"
)

func GetVaultMetadata(node *Node) (meta vault.VaultMetadata) {

	// @todo: see if we can add more information here.
	meta = vault.NewVaultMetadata()
	meta["nid"] = node.Nid
	meta["type"] = node.Type
	meta["name"] = node.Name
	meta["created_at"] = node.CreatedAt
	meta["updated_at"] = node.UpdatedAt
	meta["meta"] = node.Meta
	meta["data"] = node.Data
	meta["revision"] = node.Revision
	meta["updated_by"] = node.UpdatedBy
	meta["created_by"] = node.CreatedBy
	meta["source"] = node.Source
	meta["set_nid"] = node.SetNid
	meta["parent_nid"] = node.ParentNid

	return
}
