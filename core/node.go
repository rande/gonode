package core

import (
	"time"
	"github.com/twinj/uuid"
	"encoding/json"
)

var (
	StatusNew       = 0
	StatusDraft     = 1
	StatusCompleted = 2
	StatusValidated = 3
)

type Reference struct {
	uuid.UUID
}

func (m *Reference) MarshalJSON() ([]byte, error) {
    // Manually calling Marshal for Contents
    cont, err := json.Marshal(uuid.Formatter(m.UUID, uuid.CleanHyphen))
    if err != nil {
        return nil, err
    }

    // Stitching it all together
    return cont, nil
}

func (m *Reference) UnmarshalJSON(data []byte) error {

	if len(data) < 32 {
		panic("invalid uuid size")
	}

	tmpUuid, err := uuid.ParseUUID(string(data[1:len(data)-1]))

	if err != nil {
		return err
	}

	m.UUID      = GetReference(tmpUuid)

	return nil
}

func GetReferenceFromString(reference string) Reference {
	v, err := uuid.ParseUUID(reference)

	if err != nil {
		panic(err)
	}

	return GetReference(v)
}

func GetReference(uuid uuid.UUID) Reference {
	return Reference{uuid}
}

type Node struct {
	id         int
	Uuid       Reference     `json:"uuid"`
	Type       string        `json:"type"`
	Name       string        `json:"name"`
	Slug       string        `json:"slug"`
	Data       interface {}  `json:"data"`
	Meta       interface {}  `json:"meta"`
	Status     int           `json:"status"`
	Weight     int           `json:"weight"`
	Revision   int           `json:"revision"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	Enabled    bool          `json:"enabled"`
	Deleted    bool          `json:"deleted"`
	Parents    []Reference   `json:"parents"`
	UpdatedBy  Reference     `json:"updated_by"`
	CreatedBy  Reference     `json:"created_by"`
	ParentUuid Reference     `json:"parent_uuid"`
	SetUuid    Reference     `json:"set_uuid"`
	Source     Reference     `json:"source"`
}

func (node *Node) Id() int {
	return node.id
}

func NewNode() *Node {
	return &Node{
		Uuid:       GetEmptyReference(),
		Source:     GetEmptyReference(),
		ParentUuid: GetEmptyReference(),
		UpdatedBy:  GetEmptyReference(),
		CreatedBy:  GetEmptyReference(),
		SetUuid:    GetEmptyReference(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Weight:     1,
		Revision:   1,
		Deleted:    false,
		Enabled:    true,
		Status:     StatusNew,
	}
}
