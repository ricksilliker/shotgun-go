package shotgun_api

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

type PublishedFileData struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	CreatedAt   string    `json:"created_at"`
	Entity      LinkField `json:"entity"`
	Project     LinkField `json:"project"`
	Version     LinkField `json:"version"`
	Task        LinkField `json:"task"`
	DownloadURI string    `json:"download_uri"`
	WindowsFile string    `json:"windows_file"`
	MacFile     string    `json:"mac_file"`
	LinuxFile   string    `json:"linux_file"`
}

var publishedFileFields = []string{
	"id", "created_at", "name",
	"entity", "project",
	"path", "path_cache",
	"sg_file_size", "sg_download_uri",
	"version", "task",
}

func (p *PublishedFileData) SetField(fieldName string, fieldValue interface{}) error {
	reqBody := map[string]interface{}{
		fieldName: fieldValue,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		logrus.WithError(err).Error("failed to create request body")
		return err
	}

	req, err := NewUpdateRequest("PublishedFile", p.ID, publishedFileFields, data)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"field_name":  fieldName,
			"field_value": fmt.Sprintf("%v", fieldValue),
		}).Error("failed to create request to set PublishedFile field")
		return err
	}

	var handler PublishedFileRecordResponse
	if err = DoUpdateRequest(req, &handler); err != nil {
		logrus.Error("do not complete update PublishedFile request")
		return err
	}

	return nil
}

type Attachment struct {
	ID               int64     `json:"id"`
	LinkType         string    `json:"link_type"`
	Name             string    `json:"name"`
	LocalStorage     LinkField `json:"local_storage"`
	LocalPathMac     string    `json:"local_path_mac"`
	LocalPathLinux   string    `json:"local_path_linux"`
	LocalPathWindows string    `json:"local_path_windows"`
}

type PublishedFileRecord struct {
	ID         int64 `json:"id"`
	Attributes struct {
		Name        string     `json:"content"`
		DownloadURI string     `json:"sg_download_uri"`
		CreatedAt   string     `json:"created_at"`
		Path        Attachment `json:"path"`
		PathCache   string     `json:"path_cache"`
		FileSize    int64      `json:"sg_file_size"`
	} `json:"attributes"`
	Relationships struct {
		Version struct {
			Data LinkField `json:"data"`
		} `json:"version"`
		Task struct {
			Data LinkField `json:"data"`
		} `json:"task"`
		Entity struct {
			Data LinkField `json:"data"`
		} `json:"entity"`
		Project struct {
			Data LinkField `json:"data"`
		} `json:"project"`
	} `json:"relationships"`
}

type PublishedFileMultiRecordResponse struct {
	Data []PublishedFileRecord `json:"data"`
}

type PublishedFileRecordResponse struct {
	Data PublishedFileRecord `json:"data"`
}

func (e *PublishedFileRecordResponse) ReadRecord(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		logrus.Error("failed to unmarshal PublishedFile response")
		return err
	}
	return nil
}

func (t *PublishedFileMultiRecordResponse) ReadRecord(data []byte) error {
	err := json.Unmarshal(data, &t)
	if err != nil {
		logrus.Error("failed to unmarshal data to PublishedFileMultiRecord")
		return err
	}
	return nil
}

func GetPublishedFileForID(publishedFileID int64) (*PublishedFileData, error) {
	req, err := NewFindRequest("PublishedFile", publishedFileID, publishedFileFields)
	if err != nil {
		logrus.Error("failed to create PublishedFile find request")
		return nil, err
	}

	var resp PublishedFileRecordResponse
	if err = DoFindRequest(req, &resp); err != nil {
		logrus.Error("failed to make PublishedFile find request")
		return nil, err
	}

	result := &PublishedFileData{
		ID:          resp.Data.ID,
		Name:        resp.Data.Attributes.Name,
		Project:     resp.Data.Relationships.Project.Data,
		Entity:      resp.Data.Relationships.Entity.Data,
		Task:        resp.Data.Relationships.Task.Data,
		Version:     resp.Data.Relationships.Version.Data,
		CreatedAt:   resp.Data.Attributes.CreatedAt,
		DownloadURI: resp.Data.Attributes.DownloadURI,
		MacFile:     resp.Data.Attributes.Path.LocalPathMac,
		WindowsFile: resp.Data.Attributes.Path.LocalPathWindows,
		LinuxFile:   resp.Data.Attributes.Path.LocalPathLinux,
		Size:        resp.Data.Attributes.FileSize,
	}

	return result, nil
}
