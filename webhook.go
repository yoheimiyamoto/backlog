package backlog

import (
	"encoding/json"
	"time"
)

type Webhook struct {
	ID      int       `json:"id"`
	Project *Project  `json:"project"`
	Issue   *Issue    `json:"content"`
	Created time.Time `json:"created"`
}

func (w *Webhook) UnmarshalJSON(data []byte) error {

	type WebhookCustomField struct {
		FieldType CustomFieldType `json:"fieldTypeId"`
		Name      string          `json:"field"` // jsonのフィールド名がなぜかここだけ 'name' ではなくて、 'field' になっている糞仕様
		Value     interface{}     `json:"value"` // string or []string
	}

	type WebhookIssue struct {
		ID             int                   `json:"id"`
		ProjectID      int                   `json:"projectId"`
		IssueKey       string                `json:"issueKey"`
		KeyID          int                   `json:"keyId"`
		IssueType      IssueType             `json:"issueType"`
		Summary        string                `json:"summary"`
		Description    string                `json:"description"` // 詳細
		Resolution     Resolution            `json:"resolution"`
		Priority       Priority              `json:"priority"`
		Status         IssueStatusItem       `json:"status"`
		Assignee       Assignee              `json:"assignee"`
		Categories     []*Category           `json:"category"`
		Versions       []*Version            `json:"versions"`
		Milestones     []*Version            `json:"milestone"`
		StartDate      time.Time             `json:"startDate"`
		DueDate        time.Time             `json:"dueDate"`
		EstimatedHours int                   `json:"estimatedHours"`
		ActualHours    int                   `json:"actualHours"`
		ParentIssueID  int                   `json:"parentIssueId"`
		CreatedUser    *User                 `json:"createdUser"`
		Created        time.Time             `json:"created"`
		UpdatedUser    *User                 `json:"updatedUser"`
		Updated        time.Time             `json:"updated"`
		CustomFields   []*WebhookCustomField `json:"customFields"`
		Attachments    []*Attachment         `json:"attachments"`
		SharedFiles    []*SharedFile         `json:"sharedFiles"`
		Stars          []*Star               `json:"stars"`
		Comment        struct {
			ID      int    `json:"id"`
			Content string `json:"content"`
		} `json:"comment"`
	}

	raw := struct {
		ID      int           `json:"id"`
		Project *Project      `json:"project"`
		Issue   *WebhookIssue `json:"content"`
		Created time.Time     `json:"created"`
	}{}

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	customFields := make(CustomFields)
	for _, issue := range raw.Issue.CustomFields {
		customFields[issue.Name] = issue.Value
	}

	issue := Issue{
		ID:             raw.Issue.ID,
		ProjectID:      raw.Issue.ProjectID,
		IssueKey:       raw.Issue.IssueKey,
		KeyID:          raw.Issue.KeyID,
		IssueType:      raw.Issue.IssueType,
		Summary:        raw.Issue.Summary,
		Description:    raw.Issue.Description,
		Resolution:     raw.Issue.Resolution,
		Priority:       raw.Issue.Priority,
		Status:         raw.Issue.Status.Name,
		Assignee:       raw.Issue.Assignee,
		Categories:     raw.Issue.Categories,
		Versions:       raw.Issue.Versions,
		Milestones:     raw.Issue.Milestones,
		StartDate:      raw.Issue.StartDate,
		DueDate:        raw.Issue.DueDate,
		EstimatedHours: raw.Issue.EstimatedHours,
		ActualHours:    raw.Issue.ActualHours,
		ParentIssueID:  raw.Issue.ParentIssueID,
		CreatedUser:    raw.Issue.CreatedUser,
		Created:        raw.Issue.Created,
		UpdatedUser:    raw.Issue.UpdatedUser,
		Updated:        raw.Issue.Updated,
		CustomFields:   customFields,
		Attachments:    raw.Issue.Attachments,
		SharedFiles:    raw.Issue.SharedFiles,
		Stars:          raw.Issue.Stars,
		Comment:        raw.Issue.Comment,
	}

	w.ID = raw.ID
	w.Project = raw.Project
	w.Issue = &issue
	w.Created = raw.Created

	return nil
}

// func (r *Record) UnmarshalJSON(data []byte) error {
// 	var rawMsg json.RawMessage

// 	raw := struct {
// 		Type  string      `json:"type"`
// 		Value interface{} `json:"value"`
// 	}{Value: &rawMsg}

// 	err := json.Unmarshal(data, &raw)
// 	if err != nil {
// 		return err
// 	}
// 	var f Obj
// 	switch raw.Type {
// 	case "user":
// 		var _f User
// 		f = &_f
// 		// err := json.Unmarshal(rawMsg, &f)
// 		// if err != nil {
// 		// 	return err
// 		// }
// 		// r.Value = f
// 	case "organization":
// 		var _f Organization
// 		f = _f
// 		// err := json.Unmarshal(rawMsg, &f)
// 		// if err != nil {
// 		// 	return err
// 		// }
// 		// r.Value = f
// 	case "myString":
// 		var _f MyString2
// 		f = _f
// 	}

// 	err = json.Unmarshal(rawMsg, &f)
// 	if err != nil {
// 		return err
// 	}
// 	r.Value = f
// 	return nil
// }

type Project struct {
	ID                                int    `json:"id"`
	Key                               string `json:"projectKey"`
	Name                              string `json:"name"`
	ChartEnabled                      bool   `json:"chartEnabled"`
	SubtaskingEnabled                 bool   `json:"subtaskingEnabled"`
	ProjectLeaderCanEditProjectLeader bool   `json:"projectLeaderCanEditProjectLeader"`
	UseWikiTreeView                   bool   `json:"useWikiTreeView"`
	TextFormattingRule                string `json:"textFormattingRule"`
	Archived                          bool   `json:"archived"`
}
