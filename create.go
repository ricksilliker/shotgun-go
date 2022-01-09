package shotgun_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type CreateResponseHandler interface {
	ReadRecord(data []byte) error
}

func NewCreateRequest(entityType string, data []byte) (*http.Request, error) {
	createURL := ShotgunURL + fmt.Sprintf("/entity/%v", entityType)
	body := bytes.NewBuffer(data)
	req, _ := http.NewRequest("POST", createURL, body)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	auth, err := AuthenticateShotgunScript()
	if err != nil {
		logrus.WithError(err).Error("failed to authenticate with Shotgun")
		return nil, err
	}
	token := fmt.Sprintf("%v %v", auth.TokenType, auth.AccessToken)
	req.Header.Add("Authorization", token)
	return req, nil
}

func DoCreateRequest(req *http.Request, handler RecordResponseHandler) error {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("failed to do create request")
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
		logrus.Error("failed to read MultiRecord")
		return err
	}

	return nil
}
