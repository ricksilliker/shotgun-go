package shotgun_api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"runtime"
)

func GetPlatformProjectsPath() string {
	switch runtime.GOOS {
	case "windows":
		return "O:\\projects"
	case "darwin":
		return "/Volumes/prod/projects"
	default:
		return "/prod/projects"
	}
}

var attachmentFields = []string{
	"id", "this_file", "name",
}

type AttachmentData struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	FileURL string `json:"file_url"`
}

type AttachmentRecord struct {
	ID         int64 `json:"id"`
	Attributes struct {
		File struct {
			URL string `json:"url"`
		} `json:"this_file"`
		Name string `json:"name"`
	} `json:"attributes"`
}

type AttachmentRecordResponse struct {
	Data AttachmentRecord `json:"data"`
}

func (e *AttachmentRecordResponse) ReadRecord(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		logrus.Error("failed to unmarshal Attachment response")
		return err
	}
	return nil
}

func GetAttachmentFromID(attachmentID int64) (*AttachmentData, error) {
	req, err := NewFindRequest("Attachment", attachmentID, attachmentFields)
	if err != nil {
		logrus.Error("failed to create Attachment find request")
		return nil, err
	}

	var resp AttachmentRecordResponse
	if err = DoFindRequest(req, &resp); err != nil {
		logrus.Error("failed to make Attachment find request")
		return nil, err
	}

	result := &AttachmentData{
		ID:      resp.Data.ID,
		Name:    resp.Data.Attributes.Name,
		FileURL: resp.Data.Attributes.File.URL,
	}

	return result, nil
}
