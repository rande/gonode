package core

import (
	"io"
)

type NodeData interface {}
type NodeMeta interface {}

type DownloadData struct {
	ContentType  string
	Filename     string
	CacheControl string
	Pragma       string
	Expires      string
	Stream       func(node *Node, w io.Writer)
}

type Handler interface {
	GetStruct() (NodeData, NodeMeta) // Data, Meta

	PreUpdate(node *Node, m NodeManager) error
	PostUpdate(node *Node, m NodeManager) error
	PreInsert(node *Node, m NodeManager) error
	PostInsert(node *Node, m NodeManager) error
	Validate(node *Node, m NodeManager, e Errors)
	GetDownloadData(node *Node) *DownloadData
}

func GetDownloadData() *DownloadData {
	return &DownloadData{
		ContentType: "application/octet-stream",
		Filename: "gonode-notype.bin",
		CacheControl: "private",
		Stream: func(node *Node, w io.Writer) {
			io.WriteString(w, "No content defined to be download for this node")
		},
	}
}

