package shotgun_api

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

type ShotgunError struct {
	Errors []ShotgunErrorBody `json:"errors"`
}

type ShotgunErrorBody struct {
	Status    int    `json:"status"`
	ErrorCode int    `json:"code"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
}

func (e *ShotgunError) FormatError() error {
	var errStrings []string
	for _, err := range e.Errors {
		errStrings = append(errStrings, fmt.Sprintf("%v: %v", err.Title, err.Detail))
	}
	return fmt.Errorf(strings.Join(errStrings, "\n"))
}

func HandleError(response *http.Response) error {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	logrus.Info(string(bodyBytes))
	if err != nil {
		logrus.Error("failed to read Shotgun error response")
		return err
	}

	var errorResp ShotgunError
	if err = json.Unmarshal(bodyBytes, &errorResp); err != nil {
		logrus.Error("failed to unmarshal Shotgun error")
		return err
	}

	return errorResp.FormatError()
}
