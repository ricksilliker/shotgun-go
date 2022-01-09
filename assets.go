package shotgun_api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type AssetData struct {
	ID    int64  `json:"id"`
	Group string `json:"group"`
	Name  string `json:"name"`
}

type AssetRecord struct {
	ID         int64 `json:"id"`
	Attributes struct {
		Code string `json:"code"`
		Type string `json:"sg_asset_type"`
	} `json:"attributes"`
}

var assetFields = []string{
	"id", "code", "sg_asset_type",
}

type AssetMultiRecordResponse struct {
	Data []AssetRecord `json:"data"`
}

type AssetRecordResponse struct {
	Data AssetRecord `json:"data"`
}

func (t *AssetMultiRecordResponse) ReadRecord(data []byte) error {
	err := json.Unmarshal(data, &t)
	if err != nil {
		logrus.Error("failed to unmarshal data to AssetMultiRecord")
		return err
	}
	return nil
}

func (e *AssetRecordResponse) ReadRecord(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		logrus.Error("failed to unmarshal Asset response")
		return err
	}
	return nil
}

func GetAssetForID(assetID int64) (*AssetData, error) {
	req, err := NewFindRequest("Asset", assetID, assetFields)
	if err != nil {
		logrus.Error("failed to create Asset find request")
		return nil, err
	}

	var resp AssetRecordResponse
	if err = DoFindRequest(req, &resp); err != nil {
		logrus.Error("failed to make Asset find request")
		return nil, err
	}

	result := &AssetData{
		ID:    resp.Data.ID,
		Name:  resp.Data.Attributes.Code,
		Group: resp.Data.Attributes.Type,
	}

	return result, nil
}

func GetProjectAssets(projectID int64) ([]AssetData, error) {
	filters := ShotgunFilters{
		Expressions: []ShotgunFilterExpression{
			{"project.Project.id", "is", projectID},
		},
	}
	sort := []SortParam{
		{
			FieldName: "code",
			Direction: Ascending,
		},
	}
	req, err := NewSearchRequest("Asset", filters, assetFields, nil, sort)
	if err != nil {
		logrus.Error("failed to create Asset search request")
		return nil, err
	}

	var resp AssetMultiRecordResponse
	if err = DoSearchRequest(req, &resp); err != nil {
		logrus.Error("failed to make Asset search request")
		return nil, err
	}

	var result []AssetData
	for _, record := range resp.Data {
		asset := AssetData{
			ID:    record.ID,
			Name:  record.Attributes.Code,
			Group: record.Attributes.Type,
		}

		result = append(result, asset)
	}

	return result, nil
}
