package shotgun_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type SearchRequest struct {
	Filters [][]interface{} `json:"filters"`
	Fields  []string        `json:"fields"`
	Page    *PageParam      `json:"page,omitempty"`
	Sort    string          `json:"sort,omitempty"`
}

type SearchResponseHandler interface {
	ReadMultiRecord(data []byte) error
}

func NewSearchRequest(entityType string, filters ShotgunFilters, fields []string, page *PageParam, sort []SortParam) (*http.Request, error) {
	url := ShotgunURL + fmt.Sprintf("/entity/%v/_search", entityType)

	body := SearchRequest{
		filters.SerializeFilters(),
		fields,
		page,
		SerializeSortParameters(sort),
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal search request")
		return nil, err
	}

	data := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		logrus.WithError(err).Error("failed to create search request")
		return nil, err
	}

	return req, nil
}

func DoSearchRequest(req *http.Request, handler MultiRecordHandler) error {
	auth, err := AuthenticateShotgunScript()
	if err != nil {
		logrus.Error("authentication failed")
		return err
	}
	token := fmt.Sprintf("%v %v", auth.TokenType, auth.AccessToken)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/vnd+shotgun.api3_array+json")
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("failed to do search request")
		return err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		var errorResp ShotgunError
		err = json.Unmarshal(bodyBytes, &errorResp)
		if err != nil {
			logrus.WithField("shotgun_error", errorResp).Error("failed to unmarshal Shotgun error")
			return err
		}
		return errorResp.FormatError()
	}

	if err = handler.ReadRecord(bodyBytes); err != nil {
		logrus.Error("failed to read search MultiRecord")
		return err
	}

	return nil
}
