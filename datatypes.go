package shotgun_api

import (
	"fmt"
	"strings"
)

type PageParam struct {
	Size   int `json:"size,omitempty"`
	Number int `json:"number,omitempty"`
}

type LinkField struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
}

type ShotgunFilterExpression struct {
	Field    string
	Relation string
	Value    interface{}
}

type ShotgunFilters struct {
	Expressions []ShotgunFilterExpression
}

func (f *ShotgunFilters) SerializeFilters() [][]interface{} {
	result := make([][]interface{}, len(f.Expressions))
	for i, filter := range f.Expressions {
		result[i] = []interface{}{filter.Field, filter.Relation, filter.Value}
	}
	return result
}

type SortDirection int

const (
	Ascending SortDirection = iota
	Descending
)

type SortParam struct {
	FieldName string
	Direction SortDirection
}

func SerializeSortParameters(params []SortParam) string {
	var serializedParams []string
	for _, param := range params {
		switch param.Direction {
		case Ascending:
			serializedParams = append(serializedParams, fmt.Sprintf("%v", param.FieldName))
		case Descending:
			serializedParams = append(serializedParams, fmt.Sprintf("-%v", param.FieldName))
		}
	}
	return strings.Join(serializedParams, ",")
}

type EntityData struct {
	ID int64 `json:"id,omitempty"`
}

type Record struct {
	Data EntityData `json:"data"`
}

type RecordResponseHandler interface {
	ReadRecord(data []byte) error
}

type MultiRecord struct {
	Data []EntityData `json:"data"`
}

type MultiRecordHandler interface {
	ReadRecord(data []byte) error
}

type MediaAltType string

const (
	Original  MediaAltType = "original"
	Thumbnail MediaAltType = "thumbnail"
)
