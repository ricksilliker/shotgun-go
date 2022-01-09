package shotgun_api

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

var ShotgunURL = os.Getenv("SHOTGUN_URL") + "/api/v1"

var userFields = []string{
	"firstname", "lastname", "groups", "login", "id", "sg_status_list",
}

type UserData struct {
	ID        int64       `json:"id"`
	Firstname string      `json:"first_name"`
	Lastname  string      `json:"last_name"`
	Login     string      `json:"login"`
	Status    string      `json:"status"`
	Groups    []LinkField `json:"groups"`
}

type UserRecord struct {
	ID         int64 `json:"id"`
	Attributes struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Login     string `json:"login"`
		Status    string `json:"sg_status_list"`
	} `json:"attributes"`
	Relationships struct {
		Groups struct {
			Data []LinkField `json:"data"`
		} `json:"groups"`
	} `json:"relationships"`
}

type UserMultiRecord struct {
	Data []UserData `json:"data"`
}

type UserMultiRecordResponse struct {
	Data []UserRecord `json:"data"`
}

type UserRecordResponse struct {
	Data UserRecord `json:"data"`
}

func (e *UserRecordResponse) ReadRecord(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		logrus.Error("failed to unmarshal UserRecord response")
		return err
	}
	return nil
}

func (t *UserMultiRecordResponse) ReadRecord(data []byte) error {
	err := json.Unmarshal(data, &t)
	if err != nil {
		logrus.Error("failed to unmarshal data to UserMultiRecord")
		return err
	}
	return nil
}

func GetShotgunUserByLogin(login string) (*UserData, error) {
	filters := ShotgunFilters{
		Expressions: []ShotgunFilterExpression{
			{"login", "is", login},
		},
	}
	page := PageParam{
		Size: 1,
	}
	req, err := NewSearchRequest("HumanUser", filters, userFields, &page, nil)
	if err != nil {
		logrus.Error("failed to create Shotgun User search request")
		return nil, err
	}

	var resp UserMultiRecordResponse
	err = DoSearchRequest(req, &resp)
	if err != nil {
		logrus.Error("failed to make Shotgun User search request")
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no Users found with login: %v", login)
	}

	return &UserData{
		ID:        resp.Data[0].ID,
		Firstname: resp.Data[0].Attributes.Firstname,
		Lastname:  resp.Data[0].Attributes.Lastname,
		Login:     resp.Data[0].Attributes.Login,
		Status:    resp.Data[0].Attributes.Status,
		Groups:    resp.Data[0].Relationships.Groups.Data,
	}, nil
}

func GetUserForID(userID int64) (*UserData, error) {
	req, err := NewFindRequest("HumanUser", userID, userFields)
	if err != nil {
		logrus.Error("failed to create Shotgun User search request")
		return nil, err
	}

	var resp UserRecordResponse
	err = DoFindRequest(req, &resp)
	if err != nil {
		logrus.Error("failed to make Shotgun User search request")
		return nil, err
	}

	return &UserData{
		ID:        resp.Data.ID,
		Firstname: resp.Data.Attributes.Firstname,
		Lastname:  resp.Data.Attributes.Lastname,
		Login:     resp.Data.Attributes.Login,
		Status:    resp.Data.Attributes.Status,
		Groups:    resp.Data.Relationships.Groups.Data,
	}, nil
}
