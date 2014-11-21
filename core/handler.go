package gonode

type NodeData interface {}
type NodeMeta interface {}

type Handler interface {
	GetStruct() (NodeData, NodeMeta) // Data, Meta
}
