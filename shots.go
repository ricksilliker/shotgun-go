package shotgun_api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type ShotData struct {
	ID       int64   `json:"id"`
	Sequence string  `json:"sequence"`
	Name     string  `json:"name"`
	Assets   []int64 `json:"assets"`
	Status   string  `json:"status"`
}

type ShotRecord struct {
	ID         int64 `json:"id"`
	Attributes struct {
		Code   string `json:"code"`
		Status string `json:"sg_status_list"`
	} `json:"attributes"`
	Relationships struct {
		Sequence struct {
			Data LinkField `json:"data"`
		} `json:"sg_sequence"`
		Assets struct {
			Data []LinkField `json:"data"`
		} `json:"assets"`
	} `json:"relationships"`
}

var shotFields = []string{
	"id", "code", "sg_sequence", "assets",
	"sg_status_list",
}

type ShotRecordResponse struct {
	Data ShotRecord `json:"data"`
}

func (e *ShotRecordResponse) ReadRecord(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		logrus.Error("failed to unmarshal Shot response")
		return err
	}
	return nil
}

type ShotMultiRecordResponse struct {
	Data []ShotRecord `json:"data"`
}

func (t *ShotMultiRecordResponse) ReadRecord(data []byte) error {
	err := json.Unmarshal(data, &t)
	if err != nil {
		logrus.Error("failed to unmarshal data to ShotMultiRecord")
		return err
	}
	return nil
}

func GetShotForID(shotID int64) (*ShotData, error) {
	req, err := NewFindRequest("Shot", shotID, shotFields)
	if err != nil {
		logrus.Error("failed to create Shot find request")
		return nil, err
	}

	var resp ShotRecordResponse
	if err = DoFindRequest(req, &resp); err != nil {
		logrus.Error("failed to make Shot find request")
		return nil, err
	}

	result := &ShotData{
		ID:       resp.Data.ID,
		Name:     resp.Data.Attributes.Code,
		Status:   resp.Data.Attributes.Status,
		Sequence: resp.Data.Relationships.Sequence.Data.Name,
	}

	var assets []int64
	for _, ass := range resp.Data.Relationships.Assets.Data {
		assets = append(assets, ass.ID)
	}
	result.Assets = assets

	return result, nil
}

func GetShots(sequenceID int64, sortBy []SortParam) ([]ShotData, error) {
	filters := ShotgunFilters{
		Expressions: []ShotgunFilterExpression{
			{"sg_sequence.Sequence.id", "is", sequenceID},
		},
	}

	req, err := NewSearchRequest("Shot", filters, shotFields, nil, sortBy)
	if err != nil {
		logrus.Error("failed to create Shot search request")
		return nil, err
	}

	var resp ShotMultiRecordResponse
	if err = DoSearchRequest(req, &resp); err != nil {
		logrus.Error("failed to make Shot search request")
		return nil, err
	}

	var result []ShotData
	for _, record := range resp.Data {
		sh := ShotData{
			ID:       record.ID,
			Name:     record.Attributes.Code,
			Status:   record.Attributes.Status,
			Sequence: record.Relationships.Sequence.Data.Name,
		}

		var assets []int64
		for _, ass := range record.Relationships.Assets.Data {
			assets = append(assets, ass.ID)
		}
		sh.Assets = assets

		result = append(result, sh)
	}

	return result, nil
}
