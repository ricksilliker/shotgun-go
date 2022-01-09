package shotgun_api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type NoteData struct {
	ID        int64  `json:"id"`
	Author    string `json:"author"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type NoteRecord struct {
	ID         int64 `json:"id"`
	Attributes struct {
		Subject   string `json:"subject"`
		Body      string `json:"content"`
		CreatedAt string `json:"created_at"`
	} `json:"attributes"`
	Relationships struct {
		Author struct {
			Data LinkField `json:"data"`
		} `json:"user"`
	} `json:"relationships"`
}

type NoteMultiRecordResponse struct {
	Data []NoteRecord `json:"data"`
}

type NoteRecordResponse struct {
	Data NoteRecord `json:"data"`
}

func (e *NoteRecordResponse) ReadRecord(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		logrus.Error("failed to unmarshal Note response")
		return err
	}
	return nil
}

func (t *NoteMultiRecordResponse) ReadRecord(data []byte) error {
	err := json.Unmarshal(data, &t)
	if err != nil {
		logrus.Error("failed to unmarshal data to NoteMultiRecord")
		return err
	}
	return nil
}

func GetAllNotesForTask(taskID int64) ([]NoteData, error) {
	filters := ShotgunFilters{
		Expressions: []ShotgunFilterExpression{
			{"tasks.Task.id", "is", taskID},
		},
	}
	fields := []string{
		"id", "user", "subject", "content", "created_at",
	}
	sort := []SortParam{
		{
			FieldName: "created_at",
			Direction: Ascending,
		},
	}
	req, err := NewSearchRequest("Note", filters, fields, nil, sort)
	if err != nil {
		logrus.Error("failed to create Note search request")
		return nil, err
	}

	var resp NoteMultiRecordResponse
	if err = DoSearchRequest(req, &resp); err != nil {
		logrus.Error("failed to make Note search request")
		return nil, err
	}

	var result []NoteData
	for _, record := range resp.Data {
		task := NoteData{
			ID:        record.ID,
			Subject:   record.Attributes.Subject,
			Body:      record.Attributes.Body,
			CreatedAt: record.Attributes.CreatedAt,
			Author:    record.Relationships.Author.Data.Name,
		}

		result = append(result, task)
	}

	return result, nil
}
