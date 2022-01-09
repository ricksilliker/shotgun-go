package shotgun_api

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func NewThumbnailRequest(entityID int64, entityType, fieldName string) (*http.Request, error) {
	url := ShotgunURL + fmt.Sprintf("/entity/%v/%v/%v", entityType, entityID, fieldName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.WithError(err).Error("failed to create media download request")
		return nil, err
	}

	req.Header.Add("Accept", "application/json")

	return req, nil
}

type GetThumbnailResponse struct {
	Data string `json:"data"`
}

func DoGetThumbnail(req *http.Request) (string, error) {
	auth, err := AuthenticateShotgunScript()
	if err != nil {
		logrus.Error("authentication failed")
		return "", err
	}
	token := fmt.Sprintf("%v %v", auth.TokenType, auth.AccessToken)
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("failed to do media download request")
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("failed to read GET thumbnail response")
		return "", err
	}

	var thumbnailResp GetThumbnailResponse
	if err = json.Unmarshal(bodyBytes, &thumbnailResp); err != nil {
		logrus.WithError(err).Error("failed to unmarshal GET thumbnail response")
		return "", err
	}

	return thumbnailResp.Data, nil
}

func GetThumbnailURL(entityID int64, entityType, fieldName string) string {
	req, err := NewThumbnailRequest(entityID, entityType, fieldName)
	if err != nil {
		return ""
	}

	url, err := DoGetThumbnail(req)
	if err != nil {
		return ""
	}

	return url
}
