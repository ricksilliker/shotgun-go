package shotgun_api

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ActivityType string

const (
	NewActivityType    ActivityType = "create"
	ChangeActivityType              = "update"
	DeleteActivityType              = "delete"
)

type ActivityData struct {
	ID         int64        `json:"id"`
	Type       ActivityType `json:"type"`
	EntityType string       `json:"entity_type"`
	EntityID   int64        `json:"entity_id"`
	CreatedAt  string       `json:"created_at"`
	CreatedBy  struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"created_by"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Links       []string `json:"attachments"`
	Media       []string `json:"media"`
	UserGroups  []int64  `json:"user_groups"`
}

type ActivityRecord struct {
	Data struct {
		EntityType       string                 `json:"entity_type"`
		EntityID         int64                  `json:"entity_id"`
		LatestUpdateID   int64                  `json:"latest_update_id"`
		EarliestUpdateID int64                  `json:"earliest_update_id"`
		Updates          []ActivityUpdateRecord `json:"updates"`
	} `json:"data"`
}

type ActivityUpdateRecord struct {
	ID       int64        `json:"id"`
	Type     ActivityType `json:"update_type"`
	Metadata struct {
		Type       string `json:"type"`
		EntityType string `json:"entity_type"`
		EntityID   int64  `json:"entity_id"`
	} `json:"meta"`
	CreatedAt     string                 `json:"created_at"`
	Read          bool                   `json:"read"`
	PrimaryEntity map[string]interface{} `json:"primary_entity"`
	CreatedBy     struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"created_by"`
}

type VersionActivityFields struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	VersionNumber int    `json:"sg_version_number"`
	Description   string `json:"description"`
	DownloadURL   string `json:"sg_download_uri"`
	Movie         struct {
		URL string `json:"url"`
	} `json:"sg_uploaded_movie"`
	Entity     *LinkField  `json:"entity"`
	Task       *LinkField  `json:"sg_task"`
	UserGroups []LinkField `json:"user.HumanUser.groups"`
}

type NoteActivityFields struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Subject     string `json:"subject"`
	Body        string `json:"content"`
	Attachments []struct {
		ID int64 `json:"id"`
	} `json:"attachments"`
	Links      []LinkField `json:"note_links"`
	UserGroups []LinkField `json:"user.HumanUser.groups"`
}

func GetEntityActivity(entityType string, entityID int64, pageSize, latestActivityID int) ([]ActivityData, error) {
	activityStreamURL := ShotgunURL + fmt.Sprintf("/entity/%v/%v/activity_stream", entityType, entityID)
	req, err := http.NewRequest("GET", activityStreamURL, nil)
	if err != nil {
		logrus.Error("failed to create get activity_stream request")
		return nil, err
	}
	q := req.URL.Query()
	q.Add("limit", strconv.Itoa(pageSize))
	if latestActivityID > 0 {
		q.Add("max_id", strconv.Itoa(latestActivityID))
	}
	q.Add("entity_fields[Note]", "subject,content,attachments,note_links,user.HumanUser.groups")
	q.Add("entity_fields[Version]", "sg_version_number,description,sg_download_uri,sg_uploaded_movie,entity,sg_task,user.HumanUser.groups")
	req.URL.RawQuery = q.Encode()

	auth, err := AuthenticateShotgunScript()
	if err != nil {
		logrus.Error("failed to authorize script")
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%v %v", auth.TokenType, auth.AccessToken))

	resp, err := Client.Do(req)
	if err != nil {
		logrus.Error("failed to do get activity_stream request")
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, HandleError(resp)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("failed to read response body for get activity_stream request")
		return nil, err
	}

	var record ActivityRecord
	logrus.Debug("Decoding json response for activity stream.")
	if err = json.Unmarshal(bodyBytes, &record); err != nil {
		logrus.Error("failed to unmarshal JSON response from get activity_stream request")
		return nil, err
	}

	logrus.Debug("Start building ActivityData list..")
	var result []ActivityData
	for _, update := range record.Data.Updates {
		switch update.Metadata.EntityType {
		case "Version":
			item := formatVersion(update)
			if item != nil {
				result = append(result, *item)
				logrus.Debugf("Update added: %#v", item)
			}
		case "Note":
			item := formatNote(update)
			if item != nil {
				result = append(result, *item)
				logrus.Debugf("Update added: %#v", item)
			}
		default:
			logrus.Debugf("Update skipped because it doesnt meet entity type requirements: %#v", update)
		}
	}

	return result, nil
}

func formatVersion(record ActivityUpdateRecord) *ActivityData {
	jsonStr, _ := json.Marshal(record.PrimaryEntity)
	var keys []string
	for k := range record.PrimaryEntity {
		keys = append(keys, k)
	}
	var versionFields VersionActivityFields
	json.Unmarshal(jsonStr, &versionFields)

	if versionFields.Type != "Version" {
		logrus.Debug("Record is not a Version entity, skipping.")
		return nil
	}

	if versionFields.Name == "" {
		logrus.Debug("Version record has no name.")
		return nil
	}

	if versionFields.Movie.URL == "" && versionFields.DownloadURL == "" {
		logrus.WithFields(record.PrimaryEntity).Debugf("Version record has no interesting data.")
	}

	var item ActivityData
	item.EntityType = versionFields.Type
	item.EntityID = versionFields.ID
	for _, group := range versionFields.UserGroups {
		item.UserGroups = append(item.UserGroups, group.ID)
	}

	item.CreatedAt = record.CreatedAt
	item.CreatedBy.ID = record.CreatedBy.ID
	item.CreatedBy.Type = record.CreatedBy.Type
	item.CreatedBy.Name = record.CreatedBy.Name
	item.ID = record.ID
	item.Type = record.Type

	switch record.Type {
	case "create":
		item.Title = fmt.Sprintf("New Version: %v", versionFields.Name)
	case "update":
		item.Title = fmt.Sprintf("Updated Version: %v", versionFields.Name)
	case "delete":
		item.Title = fmt.Sprintf("Removed Version: %v", versionFields.Name)
	default:
		item.Title = fmt.Sprintf("Version from %v", record.CreatedBy.Name)
	}

	item.Description = fmt.Sprintf("Task: %v \n%v: %v", versionFields.Task.Name, versionFields.Entity.Type, versionFields.Entity.Name)

	if item.Description == "" {
		item.Description = fmt.Sprintf("Version update from %v", record.CreatedBy.Name)
	} else {
		item.Description += fmt.Sprintf("\n %v", versionFields.Description)
	}

	if versionFields.DownloadURL != "" {
		item.Links = append(item.Links, versionFields.DownloadURL)
	}

	if versionFields.Movie.URL != "" {
		item.Media = append(item.Media, versionFields.Movie.URL)
	}

	return &item
}

func formatNote(record ActivityUpdateRecord) *ActivityData {
	jsonStr, _ := json.Marshal(record.PrimaryEntity)
	var keys []string
	for k := range record.PrimaryEntity {
		keys = append(keys, k)
	}
	var noteFields NoteActivityFields
	json.Unmarshal(jsonStr, &noteFields)

	if noteFields.Type != "Note" {
		logrus.Debug("Record is not a Note entity, skipping.")
		return nil
	}

	if noteFields.Name == "" {
		logrus.Debug("Note record has no name.")
		return nil
	}

	if noteFields.Subject == "" && noteFields.Body == "" {
		logrus.WithFields(record.PrimaryEntity).Debugf("Note record has no interesting data.")
	}

	var item ActivityData
	item.EntityType = noteFields.Type
	item.EntityID = noteFields.ID
	for _, group := range noteFields.UserGroups {
		item.UserGroups = append(item.UserGroups, group.ID)
	}

	item.CreatedAt = record.CreatedAt
	item.CreatedBy.ID = record.CreatedBy.ID
	item.CreatedBy.Type = record.CreatedBy.Type
	item.CreatedBy.Name = record.CreatedBy.Name
	item.ID = record.ID
	item.Type = record.Type
	item.Title = noteFields.Subject
	if item.Title == "" {
		msg := "Note(s): "
		for _, link := range noteFields.Links {
			msg += link.Name + ", "
		}
		if len(noteFields.Links) == 0 {
			msg += fmt.Sprintf("from %v", item.CreatedBy.Name)
		}
		item.Title = msg
	}

	item.Description = noteFields.Body
	if item.Description == "" {
		item.Description = noteFields.Name
	}

	for _, a := range noteFields.Attachments {
		attachment, err := GetAttachmentFromID(a.ID)
		if err != nil {
			logrus.WithField("attachment_id", a).Error("failed to retrieve Attachment")
			continue
		}

		item.Links = append(item.Links, attachment.FileURL)
	}

	return &item
}
