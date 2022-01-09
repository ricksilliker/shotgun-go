package shotgun_api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type SequenceData struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Shots  []int64 `json:"shots"`
	Status string  `json:"status"`
}

var sequenceFields = []string{
	"id", "code", "shots", "sg_status_list",
}

type SequenceRecord struct {
	ID         int64 `json:"id"`
	Attributes struct {
		Code   string `json:"code"`
		Status string `json:"sg_status_list"`
	} `json:"attributes"`
	Relationships struct {
		Shots struct {
			Data []LinkField `json:"data"`
		} `json:"shots,omitempty"`
	} `json:"relationships"`
}

type SequenceRecordResponse struct {
	Data SequenceRecord `json:"data"`
}

func (e *SequenceRecordResponse) ReadRecord(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		logrus.Error("failed to unmarshal Sequence response")
		return err
	}
	return nil
}

type SequenceMultiRecordResponse struct {
	Data []SequenceRecord `json:"data"`
}

func (t *SequenceMultiRecordResponse) ReadRecord(data []byte) error {
	err := json.Unmarshal(data, &t)
	if err != nil {
		logrus.Error("failed to unmarshal data to SequenceMultiRecord")
		return err
	}
	return nil
}

func GetSequences(projectID int64, sortBy []SortParam) ([]SequenceData, error) {
	filters := ShotgunFilters{
		Expressions: []ShotgunFilterExpression{
			{"project.Project.id", "is", projectID},
		},
	}

	req, err := NewSearchRequest("Sequence", filters, sequenceFields, nil, sortBy)
	if err != nil {
		logrus.Error("failed to create Sequence search request")
		return nil, err
	}

	var resp SequenceMultiRecordResponse
	if err = DoSearchRequest(req, &resp); err != nil {
		logrus.Error("failed to make Sequence search request")
		return nil, err
	}

	var result []SequenceData
	for _, record := range resp.Data {
		seq := SequenceData{
			ID:     record.ID,
			Name:   record.Attributes.Code,
			Status: record.Attributes.Status,
		}

		var shots []int64
		for _, item := range record.Relationships.Shots.Data {
			shots = append(shots, item.ID)
		}
		seq.Shots = shots

		result = append(result, seq)
	}

	return result, nil
}
