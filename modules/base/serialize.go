// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"bytes"
	"encoding/json"
	"github.com/rande/gonode/core/helper"
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
	Handlers      Handlers
}

func (s *Serializer) AddSerializer(name string, f NodeSerializer) {
	s.serializers[name] = f
}

func (s *Serializer) AddDeserializer(name string, f NodeDeserializer) {
	s.deserializers[name] = f
}

func (s *Serializer) Serialize(w io.Writer, data interface{}) error {
	switch d := data.(type) {
	case *Node:
		if _, ok := s.serializers[d.Type]; ok {
			return s.serializers[d.Type](w, d)
		}
	}

	return Serialize(w, data)
}

func (s *Serializer) Deserialize(r io.Reader, o interface{}) error {
	var buffer bytes.Buffer
	read, err := buffer.ReadFrom(r)

	reader := bytes.NewReader(buffer.Bytes())

	helper.PanicOnError(err)
	helper.PanicIf(read == 0, "no data read from the request")

	switch o.(type) {
	case *Node:
		node := o.(*Node)
		if node.Type == "" {
			// we need to deserialize twice to load the correct Meta/Data structure
			err := Deserialize(reader, node)

			helper.PanicOnError(err)

			reader.Seek(0, 0)
			node.Data, node.Meta = s.Handlers.Get(node).GetStruct()
		}

		if _, ok := s.deserializers[node.Type]; ok {
			return s.deserializers[node.Type](reader, node)
		}
	}

	return Deserialize(reader, o)
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
