package shotgun_api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type ShotgunAuth struct {
	ExpiresIn    int    `json:"expires_in"`    // Time delta in seconds when this token expires.
	RefreshToken string `json:"refresh_token"` // Token used to refresh credentials.
	AccessToken  string `json:"access_token"`  // Token to access Shotgun.
	TokenType    string `json:"token_type"`    // Type of token.
}

func AuthenticateShotgunScript() (*ShotgunAuth, error) {
	data := url.Values{}
	data.Set("client_id", os.Getenv("SHOTGUN_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("SHOTGUN_SECRET"))
	data.Set("grant_type", "client_credentials")

	authURL := ShotgunURL + "/auth/access_token"
	req, _ := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("failed to make auth request with Shotgun")
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		logrus.Error("authentication to Shotgun failed")
		var errorResp ShotgunError
		err = json.Unmarshal(bodyBytes, &errorResp)
		if err != nil {
			logrus.Error("failed to unmarshal Shotgun response")
		}
		return nil, errorResp.FormatError()
	}

	a := ShotgunAuth{}
	if err = json.Unmarshal(bodyBytes, &a); err != nil {
		logrus.Error("failed to unmarshal Shotgun auth response")
		return nil, err
	}
	return &a, nil
}

func AuthenticateShotgunUser(username, password string) (*ShotgunAuth, error) {
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("grant_type", "password")

	authURL := ShotgunURL + "/auth/access_token"
	req, _ := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("failed to make user auth request with Shotgun.")
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		logrus.Error("failed to authenticate user with Shotgun")
		var errorResp ShotgunError
		err = json.Unmarshal(bodyBytes, &errorResp)
		if err != nil {
			logrus.Error("failed to unmarshal Shotgun user auth response")
			return nil, err
		}
		return nil, errorResp.FormatError()
	}

	a := ShotgunAuth{}
	if err = json.Unmarshal(bodyBytes, &a); err != nil {
		logrus.Error("failed to unmarshal Shotgun user auth response")
		return nil, err
	}

	return &a, nil
}
