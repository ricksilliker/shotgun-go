package shotgun_api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type EventData struct {
	ID          *int64                 `json:"id,omitempty"`
	EventType   string                 `json:"event_type"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"meta,omitempty"`
	Entity      LinkField              `json:"entity"`
	Project     LinkField              `json:"project"`
}

type EventRecord struct {
	ID         int64 `json:"id"`
	Attributes struct {
		EventType   string                 `json:"event_type"`
		Description string                 `json:"description"`
		Metadata    map[string]interface{} `json:"meta,omitempty"`
	} `json:"attributes"`
	Relationships struct {
		Entity struct {
			Data LinkField `json:"data"`
		} `json:"entity"`
		Project struct {
			Data LinkField `json:"data"`
		} `json:"project"`
	} `json:"relationships"`
}

type EventRecordResponse struct {
	Data EventRecord `json:"data"`
}

func (e *EventRecordResponse) ReadRecord(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		logrus.Error("failed to unmarshal EventRecord")
		return err
	}
	return nil
}

type EventMultiRecordResponse struct {
	Data []EventRecord `json:"data"`
}

func (t *EventMultiRecordResponse) ReadRecord(data []byte) error {
	err := json.Unmarshal(data, &t)
	if err != nil {
		logrus.Error("failed to unmarshal data to EventMultiRecord")
		return err
	}
	return nil
}

func NewEvent(data *EventData) error {
	reqBody, err := json.Marshal(data)
	if err != nil {
		logrus.Error("failed to marshal Event to JSON")
		return err
	}

	req, err := NewCreateRequest("EventLogEntry", reqBody)
	if err != nil {
		logrus.Error("failed to create new event request")
		return err
	}

	var resp EventRecordResponse
	if err = DoCreateRequest(req, &resp); err != nil {
		logrus.Error("failed to make new event request")
		return err
	}

	data.ID = &resp.Data.ID
	return nil
}

func GetNewEvents(lastEventID int64) ([]EventData, error) {
	var filters ShotgunFilters
	if lastEventID > 0 {
		filters.Expressions = append(filters.Expressions,
			ShotgunFilterExpression{"id", "greater_than", lastEventID},
		)
	}

	fields := []string{
		"id", "event_type", "project",
		"entity", "description", "meta",
	}

	var sort []SortParam
	if lastEventID > 0 {
		sort = []SortParam{
			{
				FieldName: "id",
				Direction: Ascending,
			},
		}
	} else {
		sort = []SortParam{
			{
				FieldName: "id",
				Direction: Descending,
			},
		}
	}

	var page PageParam
	if lastEventID > 0 {
		page = PageParam{
			Size:   25,
			Number: 0,
		}
	} else {
		page = PageParam{
			Size:   1,
			Number: 0,
		}
	}

	req, err := NewSearchRequest("EventLogEntry", filters, fields, &page, sort)
	if err != nil {
		logrus.Error("failed to create EventLogEntry search request")
		return nil, err
	}

	var resp EventMultiRecordResponse
	if err = DoSearchRequest(req, &resp); err != nil {
		logrus.Error("failed to make EventLogEntry search request")
		return nil, err
	}

	var result []EventData
	for _, record := range resp.Data {
		eventID := record.ID
		event := EventData{
			ID:          &eventID,
			EventType:   record.Attributes.EventType,
			Description: record.Attributes.Description,
			Project:     record.Relationships.Project.Data,
			Entity:      record.Relationships.Entity.Data,
			Metadata:    record.Attributes.Metadata,
		}

		result = append(result, event)
	}

	return result, nil
}
