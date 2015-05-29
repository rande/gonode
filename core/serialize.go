package core

import (
	"encoding/json"
	"io"
)

type NodeSerializer func(w io.Writer, node *Node) error
type NodeDeserializer func(r io.Reader, node *Node) error

func NewSerializer() *Serializer {
	return &Serializer{
		serializers:   make(map[string]NodeSerializer),
		deserializers: make(map[string]NodeDeserializer),
	}
}

type Serializer struct {
	serializers   map[string]NodeSerializer
	deserializers map[string]NodeDeserializer
}

func (s *Serializer) AddSerializer(name string, f NodeSerializer) {
	s.serializers[name] = f
}

func (s *Serializer) AddDeserializer(name string, f NodeDeserializer) {
	s.deserializers[name] = f
}

func (s *Serializer) Serialize(w io.Writer, node *Node) error {
	if _, ok := s.serializers[node.Type]; ok {
		return s.serializers[node.Type](w, node)
	}

	return Serialize(w, node)
}

func (s *Serializer) Deserialize(r io.Reader, node *Node) error {
	if _, ok := s.deserializers[node.Type]; ok {
		return s.deserializers[node.Type](r, node)
	}

	return Deserialize(r, node)
}

func Serialize(w io.Writer, data interface{}) error {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)

	return err
}

func Deserialize(r io.Reader, data interface{}) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(data)

	return err
}
