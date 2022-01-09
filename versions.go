package shotgun_api

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type VersionData struct {
	ID           int64       `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Notes        []LinkField `json:"notes"`
	SubmittedAt  string      `json:"submitted_at"`
	ReviewStatus string      `json:"review_status"`
	Status       string      `json:"status"`
	Number       int64       `json:"number"`
	Task         LinkField   `json:"task"`
	Entity       LinkField   `json:"entity"`
	Project      LinkField   `json:"project"`
	DownloadURL  string      `json:"download_url"`
}

var VersionFields = []string{
	"id", "code", "open_notes",
	"created_at", "sg_review_status", "sg_status_list",
	"sg_version_number", "project", "sg_task", "entity",
	"sg_download_uri", "description",
}

func (v *VersionData) SetField(fieldName string, fieldValue interface{}) error {
	reqBody := map[string]interface{}{
		fieldName: fieldValue,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		logrus.WithError(err).Error("failed to create request body")
		return err
	}

	fields := make([]string, 0)
	req, err := NewUpdateRequest("Version", v.ID, fields, data)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"field_name":  fieldName,
			"field_value": fmt.Sprintf("%v", fieldValue),
		}).Error("failed to create request to set Version field")
		return err
	}

	var handler VersionRecordResponse
	if err = DoUpdateRequest(req, &handler); err != nil {
		logrus.Error("do not complete update Version request")
		return err
	}

	return nil
}

func (v *VersionData) GetResourcePublishPath() (*string, bool, error) {
	projectRoot := GetProjectPath(v.Project.Name)

	var entityGroup string
	var entityGroupSubfolder string
	switch v.Entity.Type {
	case "Shot":
		shotData, err := GetShotForID(v.Entity.ID)
		if err != nil {
			logrus.WithError(err).Errorf("failed to retrieve Shot for Version (%v)", v.ID)
			return nil, false, err
		}
		entityGroupSubfolder = "sequences"
		entityGroup = shotData.Sequence
	case "Asset":
		assetData, err := GetAssetForID(v.Entity.ID)
		if err != nil {
			logrus.WithError(err).Errorf("failed to retrieve Asset for Version (%v)", v.ID)
			return nil, false, err
		}
		entityGroupSubfolder = "assets"
		entityGroup = assetData.Group
	}

	taskData, err := GetTaskFromID(v.Task.ID)
	if err != nil {
		logrus.WithError(err).Errorf("failed to retrieve Task for Version (%v)", v.ID)
		return nil, false, err
	}

	stepData, err := GetStepForID(taskData.Step.ID)
	if err != nil {
		logrus.WithError(err).Errorf("failed to retrieve Step for Version (%v)", v.ID)
		return nil, false, err
	}

	versionPath := filepath.Join(
		projectRoot,
		"production",
		entityGroupSubfolder,
		entityGroup,
		v.Entity.Name,
		stepData.ShortName,
		"publish",
		fmt.Sprintf("v%03d", v.Number),
	)

	_, err = os.Stat(versionPath)
	versionPathExists := !os.IsNotExist(err)

	return &versionPath, versionPathExists, nil
}

type VersionRecord struct {
	ID         int64 `json:"id"`
	Attributes struct {
		Code          string `json:"code"`
		CreatedAt     string `json:"created_at"`
		ReviewStatus  string `json:"sg_review_status"`
		Status        string `json:"sg_status_list"`
		VersionNumber int64  `json:"sg_version_number"`
		DownloadURI   string `json:"sg_download_uri"`
		Description   string `json:"description"`
	} `json:"attributes"`
	Relationships struct {
		OpenNotes struct {
			Data []LinkField `json:"data"`
		} `json:"open_notes"`
		Task struct {
			Data LinkField `json:"data"`
		} `json:"sg_task"`
		Entity struct {
			Data LinkField `json:"data"`
		} `json:"entity"`
		Project struct {
			Data LinkField `json:"data"`
		} `json:"project"`
	} `json:"relationships"`
}

type VersionMultiRecordResponse struct {
	Data []VersionRecord `json:"data"`
}

type VersionRecordResponse struct {
	Data VersionRecord `json:"data"`
}

func (e *VersionRecordResponse) ReadRecord(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		logrus.Error("failed to unmarshal Version response")
		return err
	}
	return nil
}

func (t *VersionMultiRecordResponse) ReadRecord(data []byte) error {
	err := json.Unmarshal(data, &t)
	if err != nil {
		logrus.Error("failed to unmarshal data to VersionMultiRecord")
		return err
	}
	return nil
}

func GetVersionsForTask(taskID int64) ([]VersionData, error) {
	filters := ShotgunFilters{
		Expressions: []ShotgunFilterExpression{
			{"sg_task.Task.id", "is", taskID},
		},
	}
	sort := []SortParam{
		{
			FieldName: "sg_version_number",
			Direction: Descending,
		},
		{
			FieldName: "created_at",
			Direction: Descending,
		},
	}
	req, err := NewSearchRequest("Version", filters, VersionFields, nil, sort)
	if err != nil {
		logrus.Error("failed to create Version search request")
		return nil, err
	}

	var resp VersionMultiRecordResponse
	if err = DoSearchRequest(req, &resp); err != nil {
		logrus.Error("failed to make Version search request")
		return nil, err
	}

	var result []VersionData
	for _, record := range resp.Data {
		version := VersionData{
			ID:           record.ID,
			Name:         record.Attributes.Code,
			SubmittedAt:  record.Attributes.CreatedAt,
			ReviewStatus: record.Attributes.ReviewStatus,
			Status:       record.Attributes.Status,
			Number:       record.Attributes.VersionNumber,
			Task:         record.Relationships.Task.Data,
			Project:      record.Relationships.Project.Data,
			Entity:       record.Relationships.Entity.Data,
			DownloadURL:  record.Attributes.DownloadURI,
			Description:  record.Attributes.Description,
		}

		result = append(result, version)
	}

	return result, nil
}

func GetVersionForID(versionID int64) (*VersionData, error) {
	req, err := NewFindRequest("Version", versionID, VersionFields)
	if err != nil {
		logrus.Error("failed to create Version find request")
		return nil, err
	}

	var resp VersionRecordResponse
	if err = DoFindRequest(req, &resp); err != nil {
		logrus.Error("failed to make Version find request")
		return nil, err
	}

	result := &VersionData{
		ID:           resp.Data.ID,
		Name:         resp.Data.Attributes.Code,
		SubmittedAt:  resp.Data.Attributes.CreatedAt,
		ReviewStatus: resp.Data.Attributes.ReviewStatus,
		Status:       resp.Data.Attributes.Status,
		Number:       resp.Data.Attributes.VersionNumber,
		Task:         resp.Data.Relationships.Task.Data,
		Project:      resp.Data.Relationships.Project.Data,
		Entity:       resp.Data.Relationships.Entity.Data,
		DownloadURL:  resp.Data.Attributes.DownloadURI,
		Description:  resp.Data.Attributes.Description,
	}

	return result, nil
}

func FindOneVersion(filters ShotgunFilters, sortParams []SortParam) (*VersionData, error) {
	pageParam := PageParam{
		Size: 1,
	}
	req, err := NewSearchRequest("Version", filters, VersionFields, &pageParam, sortParams)
	if err != nil {
		logrus.Error("failed to create Version search request")
		return nil, err
	}

	var resp VersionMultiRecordResponse
	if err = DoSearchRequest(req, &resp); err != nil {
		logrus.Error("failed to make Version search request")
		return nil, err
	}

	if len(resp.Data) == 0 {
		logrus.Info("search results contain zero Version items")
		return nil, nil
	}

	version := VersionData{
		ID:           resp.Data[0].ID,
		Name:         resp.Data[0].Attributes.Code,
		SubmittedAt:  resp.Data[0].Attributes.CreatedAt,
		ReviewStatus: resp.Data[0].Attributes.ReviewStatus,
		Status:       resp.Data[0].Attributes.Status,
		Number:       resp.Data[0].Attributes.VersionNumber,
		Task:         resp.Data[0].Relationships.Task.Data,
		Project:      resp.Data[0].Relationships.Project.Data,
		Entity:       resp.Data[0].Relationships.Entity.Data,
		DownloadURL:  resp.Data[0].Attributes.DownloadURI,
		Description:  resp.Data[0].Attributes.Description,
	}

	return &version, nil
}
