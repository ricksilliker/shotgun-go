package shotgun_api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type TaskData struct {
	ID         int64       `json:"id"`
	Name       string      `json:"name"`
	DueDate    string      `json:"due_date"`
	NoteCount  int64       `json:"open_notes_count"`
	Status     string      `json:"status"`
	Step       *LinkField  `json:"step"`
	Entity     LinkField   `json:"entity"`
	Project    LinkField   `json:"project"`
	AssignedTo []LinkField `json:"task_assignees"`
	Thumbnail  string      `json:"thumbnail"`
	UpdatedAt  string      `json:"updated_at"`
}

var taskFields = []string{
	"id", "content", "project",
	"entity", "task_assignees",
	"due_date", "open_notes_count",
	"sg_status_list", "image", "step",
	"updated_at",
}

type TaskRecord struct {
	ID         int64 `json:"id"`
	Attributes struct {
		Name      string `json:"content"`
		DueDate   string `json:"due_date"`
		NoteCount int64  `json:"open_notes_count"`
		Status    string `json:"sg_status_list"`
		Image     string `json:"image"`
		UpdatedAt string `json:"updated_at"`
	} `json:"attributes"`
	Relationships struct {
		Step struct {
			Data LinkField `json:"data"`
		} `json:"step"`
		Entity struct {
			Data LinkField `json:"data"`
		} `json:"entity"`
		Project struct {
			Data LinkField `json:"data"`
		} `json:"project"`
		AssignedTo struct {
			Data []LinkField `json:"data"`
		} `json:"task_assignees,omitempty"`
	} `json:"relationships"`
}

type TaskMultiRecordResponse struct {
	Data []TaskRecord `json:"data"`
}

type TaskRecordResponse struct {
	Data TaskRecord `json:"data"`
}

func (e *TaskRecordResponse) ReadRecord(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		logrus.Error("failed to unmarshal Task response")
		return err
	}
	return nil
}

func (t *TaskMultiRecordResponse) ReadRecord(data []byte) error {
	err := json.Unmarshal(data, &t)
	if err != nil {
		logrus.Error("failed to unmarshal data to TaskMultiRecord")
		return err
	}
	return nil
}

func GetAllTasksForUser(username string) ([]TaskData, error) {
	user, err := GetShotgunUserByLogin(username)
	if err != nil {
		logrus.WithField("login", username).Error("could not find User")
		return nil, err
	}

	return GetAllTasksForUserID(user.ID)
}

func GetAllTasksForUserID(userID int64) ([]TaskData, error) {
	filters := ShotgunFilters{
		Expressions: []ShotgunFilterExpression{
			{"task_assignees.HumanUser.id", "is", userID},
		},
	}

	sort := []SortParam{
		{
			FieldName: "due_date",
			Direction: Ascending,
		},
	}
	req, err := NewSearchRequest("Task", filters, taskFields, nil, sort)
	if err != nil {
		logrus.Error("failed to create Task search request")
		return nil, err
	}

	var resp TaskMultiRecordResponse
	if err = DoSearchRequest(req, &resp); err != nil {
		logrus.Error("failed to make Task search request")
		return nil, err
	}

	var result []TaskData
	for _, record := range resp.Data {
		task := TaskData{
			ID:         record.ID,
			Name:       record.Attributes.Name,
			Status:     record.Attributes.Status,
			DueDate:    record.Attributes.DueDate,
			NoteCount:  record.Attributes.NoteCount,
			Project:    record.Relationships.Project.Data,
			Entity:     record.Relationships.Entity.Data,
			Step:       &record.Relationships.Step.Data,
			AssignedTo: record.Relationships.AssignedTo.Data,
			Thumbnail:  record.Attributes.Image,
		}

		result = append(result, task)
	}

	return result, nil
}

func GetTaskFromID(taskID int64) (*TaskData, error) {
	req, err := NewFindRequest("Task", taskID, taskFields)
	if err != nil {
		logrus.Error("failed to create Task find request")
		return nil, err
	}

	var resp TaskRecordResponse
	if err = DoFindRequest(req, &resp); err != nil {
		logrus.Error("failed to make Task find request")
		return nil, err
	}

	result := &TaskData{
		ID:         resp.Data.ID,
		Name:       resp.Data.Attributes.Name,
		Status:     resp.Data.Attributes.Status,
		DueDate:    resp.Data.Attributes.DueDate,
		NoteCount:  resp.Data.Attributes.NoteCount,
		Project:    resp.Data.Relationships.Project.Data,
		Entity:     resp.Data.Relationships.Entity.Data,
		Step:       &resp.Data.Relationships.Step.Data,
		AssignedTo: resp.Data.Relationships.AssignedTo.Data,
		Thumbnail:  resp.Data.Attributes.Image,
	}

	return result, nil
}

func GetTasksForAsset(assetID int64) ([]TaskData, error) {
	filters := ShotgunFilters{
		Expressions: []ShotgunFilterExpression{
			{"entity.Asset.id", "is", assetID},
		},
	}
	sort := []SortParam{
		{
			FieldName: "due_date",
			Direction: Ascending,
		},
	}
	req, err := NewSearchRequest("Task", filters, taskFields, nil, sort)
	if err != nil {
		logrus.Error("failed to create Task search request")
		return nil, err
	}

	var resp TaskMultiRecordResponse
	if err = DoSearchRequest(req, &resp); err != nil {
		logrus.Error("failed to make Task search request")
		return nil, err
	}
	var result []TaskData
	for _, record := range resp.Data {
		task := TaskData{
			ID:         record.ID,
			Name:       record.Attributes.Name,
			Status:     record.Attributes.Status,
			DueDate:    record.Attributes.DueDate,
			NoteCount:  record.Attributes.NoteCount,
			Project:    record.Relationships.Project.Data,
			Entity:     record.Relationships.Entity.Data,
			Step:       &record.Relationships.Step.Data,
			AssignedTo: record.Relationships.AssignedTo.Data,
			Thumbnail:  record.Attributes.Image,
			UpdatedAt:  record.Attributes.UpdatedAt,
		}

		result = append(result, task)
	}

	return result, nil
}

func GetTasksForShot(shotID int64) ([]TaskData, error) {
	filters := ShotgunFilters{
		Expressions: []ShotgunFilterExpression{
			{"entity.Shot.id", "is", shotID},
		},
	}
	sort := []SortParam{
		{
			FieldName: "due_date",
			Direction: Ascending,
		},
	}
	req, err := NewSearchRequest("Task", filters, taskFields, nil, sort)
	if err != nil {
		logrus.Error("failed to create Task search request")
		return nil, err
	}

	var resp TaskMultiRecordResponse
	if err = DoSearchRequest(req, &resp); err != nil {
		logrus.Error("failed to make Task search request")
		return nil, err
	}

	var result []TaskData
	for _, record := range resp.Data {
		task := TaskData{
			ID:         record.ID,
			Name:       record.Attributes.Name,
			Status:     record.Attributes.Status,
			DueDate:    record.Attributes.DueDate,
			NoteCount:  record.Attributes.NoteCount,
			Project:    record.Relationships.Project.Data,
			Entity:     record.Relationships.Entity.Data,
			Step:       &record.Relationships.Step.Data,
			AssignedTo: record.Relationships.AssignedTo.Data,
			Thumbnail:  record.Attributes.Image,
			UpdatedAt:  record.Attributes.UpdatedAt,
		}

		result = append(result, task)
	}

	return result, nil
}

type StepData struct {
	ID        int64  `json:"id"`
	LongName  string `json:"long_name"`
	ShortName string `json:"short_name"`
}

type StepRecord struct {
	ID         int64 `json:"id"`
	Attributes struct {
		Code      string `json:"code"`
		ShortName string `json:"short_name"`
	} `json:"attributes"`
}

type StepRecordResponse struct {
	Data StepRecord `json:"data"`
}

func (e *StepRecordResponse) ReadRecord(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		logrus.Error("failed to unmarshal Step response")
		return err
	}
	return nil
}

func GetStepForID(stepID int64) (*StepData, error) {
	fields := []string{
		"id", "code", "short_name",
	}
	req, err := NewFindRequest("Step", stepID, fields)
	if err != nil {
		logrus.Error("failed to create Step find request")
		return nil, err
	}

	var resp StepRecordResponse
	if err = DoFindRequest(req, &resp); err != nil {
		logrus.Error("failed to make Step find request")
		return nil, err
	}

	result := &StepData{
		ID:        resp.Data.ID,
		LongName:  resp.Data.Attributes.Code,
		ShortName: resp.Data.Attributes.ShortName,
	}

	return result, nil
}
