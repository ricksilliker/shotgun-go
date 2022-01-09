package shotgun_api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

var projectFields = []string{
	"id", "name", "image",
}

type ProjectData struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail_url"`
}

type ProjectRecord struct {
	ID         int64 `json:"id"`
	Attributes struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	} `json:"attributes"`
}

type ProjectRecordResponse struct {
	Data ProjectRecord `json:"data"`
}

func (e *ProjectRecordResponse) ReadRecord(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		logrus.Error("failed to unmarshal Project response")
		return err
	}
	return nil
}

type ProjectMultiRecordResponse struct {
	Data []ProjectRecord `json:"data"`
}

func (t *ProjectMultiRecordResponse) ReadRecord(data []byte) error {
	err := json.Unmarshal(data, &t)
	if err != nil {
		logrus.Error("failed to unmarshal data to ProjectMultiRecord")
		return err
	}
	return nil
}

func GetAllProjectsForUser(username string) ([]ProjectData, error) {
	filters := ShotgunFilters{
		Expressions: []ShotgunFilterExpression{
			{"users.HumanUser.login", "contains", username},
		},
	}
	sort := []SortParam{
		{
			FieldName: "name",
			Direction: Descending,
		},
	}
	req, err := NewSearchRequest("Project", filters, projectFields, nil, sort)
	if err != nil {
		logrus.Error("failed to create Project search request")
		return nil, err
	}

	var resp ProjectMultiRecordResponse
	if err = DoSearchRequest(req, &resp); err != nil {
		logrus.Error("failed to make Project search request")
		return nil, err
	}

	var result []ProjectData
	for _, record := range resp.Data {
		project := ProjectData{
			ID:        record.ID,
			Name:      record.Attributes.Name,
			Thumbnail: record.Attributes.Image,
		}

		result = append(result, project)
	}

	return result, nil
}

func GetProjectPath(projectName string) string {
	return filepath.Join(GetPlatformProjectsPath(), projectName)
}

func GetProjectFromID(projectID int64) (*ProjectData, error) {
	req, err := NewFindRequest("Project", projectID, projectFields)
	if err != nil {
		logrus.Error("failed to create Task find request")
		return nil, err
	}

	var resp ProjectRecordResponse
	if err = DoFindRequest(req, &resp); err != nil {
		logrus.Error("failed to make Task find request")
		return nil, err
	}

	result := &ProjectData{
		ID:        resp.Data.ID,
		Name:      resp.Data.Attributes.Name,
		Thumbnail: resp.Data.Attributes.Image,
	}

	return result, nil
}
