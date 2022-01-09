package shotgun_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

func NewUpdateRequest(entityType string, entityID int64, fields []string, body []byte) (*http.Request, error) {
	url := ShotgunURL + fmt.Sprintf("/entity/%v/%v", entityType, entityID)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		logrus.WithError(err).Error("failed to create update request")
		return nil, err
	}

	q := req.URL.Query()
	q.Add("fields", strings.Join(fields, ","))
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func DoUpdateRequest(req *http.Request, handler RecordResponseHandler) error {
	auth, err := AuthenticateShotgunScript()
	if err != nil {
		logrus.Error("authentication failed")
		return err
	}
	req.Header.Add("Accept", "application/json")
	token := fmt.Sprintf("%v %v", auth.TokenType, auth.AccessToken)
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("failed to do update request")
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
