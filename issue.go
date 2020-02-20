package backlog

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type CustomFieldType int32

// CustomFieldType
const (
	CustomFieldTypeText CustomFieldType = iota + 1
	CustomFieldTypeSentence
	CustomFieldTypeNumber
	CustomFieldTypeDate
	CustomFieldTypeSingleList
	CustomFieldTypeMultipleList
	CustomFieldTypeCheckbox
	CustomFieldTypeRadio
)

type Issue struct {
	ID             int           `json:"id,"`
	ProjectID      int           `json:"projectId"`
	IssueKey       string        `json:"issueKey"`
	KeyID          int           `json:"keyId"`
	IssueType      IssueType     `json:"issueType"`
	Summary        string        `json:"summary"`
	Description    string        `json:"description"` // 詳細
	Resolution     Resolution    `json:"resolution"`
	Priority       Priority      `json:"priority"`
	Status         string        `json:"status"`
	Assignee       Assignee      `json:"assignee"`
	Categories     []*Category   `json:"category"`
	Versions       []*Version    `json:"versions"`
	Milestones     []*Version    `json:"milestone"`
	StartDate      time.Time     `json:"startDate"`
	DueDate        time.Time     `json:"dueDate"`
	EstimatedHours int           `json:"estimatedHours"`
	ActualHours    int           `json:"actualHours"`
	ParentIssueID  int           `json:"parentIssueId"`
	CreatedUser    *User         `json:"createdUser"`
	Created        time.Time     `json:"created"`
	UpdatedUser    *User         `json:"updatedUser"`
	Updated        time.Time     `json:"updated"`
	CustomFields   CustomFields  `json:"customFields"`
	Attachments    []*Attachment `json:"attachments"`
	SharedFiles    []*SharedFile `json:"sharedFiles"`
	Stars          []*Star       `json:"stars"`
	Comment        struct {
		ID      int    `json:"id"`
		Content string `json:"content"`
	} `json:"comment"`
}

type Date time.Time

func (d *Date) UnmarshalJSON(data []byte) error {
	var raw string

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	if raw == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return err
	}
	*d = Date(t)

	return nil
}

func (d *Date) String() string {
	return time.Time(*d).Format("2006-01-02")
}

func (i *Issue) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID             int             `json:"id,"`
		ProjectID      int             `json:"projectId"`
		IssueKey       string          `json:"issueKey"`
		KeyID          int             `json:"keyId"`
		IssueType      IssueType       `json:"issueType"`
		Summary        string          `json:"summary"`
		Description    string          `json:"description"` // 詳細
		Resolution     Resolution      `json:"resolution"`
		Priority       Priority        `json:"priority"`
		Status         issueStatusItem `json:"status"` // status
		Assignee       Assignee        `json:"assignee"`
		Categories     []*Category     `json:"category"`
		Versions       []*Version      `json:"versions"`
		Milestones     []*Version      `json:"milestone"`
		StartDate      time.Time       `json:"startDate"`
		DueDate        time.Time       `json:"dueDate"`
		EstimatedHours int             `json:"estimatedHours"`
		ActualHours    int             `json:"actualHours"`
		ParentIssueID  int             `json:"parentIssueId"`
		CreatedUser    *User           `json:"createdUser"`
		Created        time.Time       `json:"created"`
		UpdatedUser    *User           `json:"updatedUser"`
		Updated        time.Time       `json:"updated"`
		CustomFields   CustomFields    `json:"customFields"`
		Attachments    []*Attachment   `json:"attachments"`
		SharedFiles    []*SharedFile   `json:"sharedFiles"`
		Stars          []*Star         `json:"stars"`
		Comment        struct {
			ID      int    `json:"id"`
			Content string `json:"content"`
		} `json:"comment"`
	}

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	*i = Issue{
		ID:             raw.ID,
		ProjectID:      raw.ProjectID,
		IssueKey:       raw.IssueKey,
		KeyID:          raw.KeyID,
		IssueType:      raw.IssueType,
		Summary:        raw.Summary,
		Description:    raw.Description,
		Resolution:     raw.Resolution,
		Priority:       raw.Priority,
		Status:         raw.Status.Name, // status name
		Assignee:       raw.Assignee,
		Categories:     raw.Categories,
		Versions:       raw.Versions,
		Milestones:     raw.Milestones,
		StartDate:      raw.StartDate,
		DueDate:        raw.DueDate,
		EstimatedHours: raw.EstimatedHours,
		ActualHours:    raw.ActualHours,
		ParentIssueID:  raw.ParentIssueID,
		CreatedUser:    raw.CreatedUser,
		Created:        raw.Created,
		UpdatedUser:    raw.UpdatedUser,
		Updated:        raw.Updated,
		CustomFields:   raw.CustomFields,
		Attachments:    raw.Attachments,
		SharedFiles:    raw.SharedFiles,
		Stars:          raw.Stars,
		Comment:        raw.Comment,
	}

	return nil
}

type IssueType struct {
	ID           int    `json:"id,omitempty"`
	ProjectID    int    `json:"projectId,omitempty"`
	Name         string `json:"name,omitempty"`
	Color        string `json:"color,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

type Resolution struct {
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// Priority represents
type Priority struct {
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// Assignee represents
type Assignee struct {
	ID          int    `json:"id,omitempty"`
	UserID      string `json:"userId,omitempty"`
	Name        string `json:"name,omitempty"`
	RoleType    int    `json:"roleType,omitempty"`
	Lang        string `json:"lang,omitempty"`
	MailAddress string `json:"mailAddress,omitempty"`
}

type User struct {
	ID           int    `json:"id,omitempty"`
	UserID       string `json:"userId,omitempty"`
	Name         string `json:"name,omitempty"`
	RoleType     int    `json:"roleType,omitempty"`
	Lang         string `json:"lang,omitempty"`
	MailAddress  string `json:"mailAddress,omitempty"`
	NulabAccount string `json:"nulabAccount,omitempty"`
}

type Category struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

type Version struct {
	ID             int       `json:"id,omitempty"`
	ProjectID      int       `json:"projectId,omitempty"`
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
	StartDate      time.Time `json:"startDate,omitempty"`
	ReleaseDueDate time.Time `json:"releaseDueDate,omitempty"`
	Archived       bool      `json:"archived,omitempty"`
	DisplayOrder   int       `json:"displayOrder,omitempty"`
}

type Attachment struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Size        int64     `json:"size,omitempty"`
	CreatedUser User      `json:"createdUser,omitempty"`
	Created     time.Time `json:"created,omitempty"`
}

type SharedFile struct {
	ID          int       `json:"id,omitempty"`
	Type        string    `json:"type,omitempty"`
	Directory   string    `json:"dir,omitempty"`
	Name        string    `json:"name,omitempty"`
	Size        int64     `json:"size,omitempty"`
	CreatedUser *User     `json:"createdUser,omitempty"`
	Created     time.Time `json:"created,omitempty"`
	UpdatedUser *User     `json:"updatedUser,omitempty"`
	Updated     time.Time `json:"updated,omitempty"`
}

type Star struct {
	ID        int       `json:"id,omitempty"`
	Comment   string    `json:"comment,omitempty"`
	URL       string    `json:"url,omitempty"`
	Title     string    `json:"name,omitempty"`
	Created   time.Time `json:"created,omitempty"`
	Presenter *User     `json:"presenter,omitempty"`
}

type Issues []*Issue

// IssueStatusItem
type issueStatusItem struct {
	ID           int    `json:"id,omitempty"`
	ProjectID    int    `json:"projectId,omitempty"`
	Name         string `json:"name,omitempty"`
	Color        string `json:"color,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

//+sort
func (x Issues) Len() int {
	return len(x)
}

func (x Issues) Less(i, j int) bool {
	return x[i].ID < x[j].ID
}

func (x Issues) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

//-sort

func (i *Issue) params(property customFieldProperties, issueStatusItems []*issueStatusItem) (url.Values, error) {

	//+issueStatus
	var issueStatus string
	for _, item := range issueStatusItems {
		if item.Name == i.Status {
			issueStatus = strconv.Itoa(item.ID)
		}
	}
	//-issueStatus

	params := url.Values{
		"summary": {i.Summary},
		// "parentIssueId": {strconv.Itoa(i.ParentIssueID)},
		"description": {i.Description},
		"statusId":    {issueStatus},
	}

	for fieldName, customField := range i.CustomFields {
		//+key
		fieldID, err := property.findFieldID(fieldName)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("find field id of field name '%s' failed", fieldName))
		}
		key := fmt.Sprintf("customField_%d", fieldID)
		//-key

		if customField == nil {
			params.Add(key, "")
			continue
		}

		//+add custom fields to params
		switch f := customField.(type) {
		// valueがstringのケース
		case TextCustomField, SentenceCustomField:
			params.Add(key, fmt.Sprint(f))
		// valueがlistItemのケース
		case SingleListCustomField, RadioCustomField:
			itemID, err := property.findItemID(fieldName, fmt.Sprint(f))
			if err != nil {
				return nil, fmt.Errorf("find item '%s' in custom field '%s' failed", f, fieldName)
			}
			params.Add(key, strconv.Itoa(itemID))
		// valueがlistItemsのケース
		case MultipleListCustomField:
			itemIDs := make([]string, len(f))
			for i, v := range f {
				itemID, err := property.findItemID(fieldName, v)
				if err != nil {
					return nil, fmt.Errorf("find item '%s' in custom field '%s' failed", v, fieldName)
				}
				itemIDs[i] = strconv.Itoa(itemID)
			}
			params.Add(key, strings.Join(itemIDs, ","))
		case CheckboxCustomField:
			itemIDs := make([]string, len(f))
			for i, v := range f {
				itemID, err := property.findItemID(fieldName, v)
				if err != nil {
					return nil, fmt.Errorf("find item '%s' in custom field '%s' failed", v, fieldName)
				}
				itemIDs[i] = strconv.Itoa(itemID)
			}
			params.Add(key, strings.Join(itemIDs, ","))
		}
		//-add custom fields to params
	}

	return params, nil
}

//+custom field
// value = string or []string or nil
type CustomFields map[string]CustomField

type CustomField interface {
	String() string
}

type TextCustomField string

func (f TextCustomField) String() string {
	return string(f)
}

type SentenceCustomField string

func (f SentenceCustomField) String() string {
	return string(f)
}

type NumberCustomField int

func (f NumberCustomField) String() string {
	return string(f)
}

type DateCustomField time.Time

type SingleListCustomField string

func (f SingleListCustomField) String() string {
	return string(f)
}

type MultipleListCustomField []string

func (f MultipleListCustomField) String() string {
	return strings.Join(f, ",")
}

type CheckboxCustomField []string

func (f CheckboxCustomField) String() string {
	return strings.Join(f, ",")
}

type RadioCustomField string

func (f RadioCustomField) String() string {
	return string(f)
}

func (fs *CustomFields) UnmarshalJSON(data []byte) error {
	type Raw struct {
		Name            string           `json:"name"`
		CustomFieldType CustomFieldType  `json:"fieldTypeId"`
		Value           *json.RawMessage `json:"value"`
	}
	var raws []*Raw

	err := json.Unmarshal(data, &raws)
	if err != nil {
		return err
	}

	//+フィールド名重複確認
	seen := make(map[string]bool)

	for _, r := range raws {
		if seen[r.Name] {
			return fmt.Errorf("custom field name '%s' is ambiguous", r.Name)
		}
	}
	//-フィールド名重複確認

	_fs := make(CustomFields)

	for _, r := range raws {

		if r.Value == nil {
			_fs[r.Name] = nil
			continue
		}

		var customField CustomField

		switch r.CustomFieldType {
		case CustomFieldTypeText:
			var _f TextCustomField
			err = json.Unmarshal(*r.Value, &_f)
			if err != nil {
				return err
			}
			customField = _f
		case CustomFieldTypeSentence:
			var _f SentenceCustomField
			err = json.Unmarshal(*r.Value, &_f)
			if err != nil {
				return err
			}
			customField = _f
		case CustomFieldTypeSingleList:
			var item customFieldListItem
			err = json.Unmarshal(*r.Value, &item)
			if err != nil {
				log.Println(string(*r.Value))
				return err
			}
			customField = SingleListCustomField(item.Name)
		case CustomFieldTypeRadio:
			var item customFieldListItem
			err = json.Unmarshal(*r.Value, &item)
			if err != nil {
				return err
			}
			customField = RadioCustomField(item.Name)
		case CustomFieldTypeCheckbox:
			var items []*customFieldListItem

			err = json.Unmarshal(*r.Value, &items)
			if err != nil {
				return err
			}

			vs := make([]string, len(items))
			for i, item := range items {
				vs[i] = item.Name
			}
			customField = CheckboxCustomField(vs)
		case CustomFieldTypeMultipleList:
			var items []*customFieldListItem
			err = json.Unmarshal(*r.Value, &items)
			if err != nil {
				return err
			}

			vs := make([]string, len(items))
			for i, item := range items {
				vs[i] = item.Name
			}

			customField = MultipleListCustomField(vs)
		}

		_fs[r.Name] = customField
	}

	*fs = _fs
	return nil
}

// func (f *CustomField) UnmarshalJSON(data []byte) error {
// 	type Raw struct {
// 		FieldType CustomFieldType  `json:"fieldTypeId"`
// 		Name      string           `json:"name"`
// 		Value     *json.RawMessage `json:"value"` // string or customFieldListItem
// 	}

// 	var raw Raw

// 	if err := json.Unmarshal(data, &raw); err != nil {
// 		return err
// 	}

// 	f.Name = raw.Name

// 	if raw.Value == nil {
// 		return nil
// 	}

// 	switch raw.FieldType {
// 	case CustomFieldTypeText, CustomFieldTypeSentence:
// 		var s string
// 		err := json.Unmarshal(*raw.Value, &s)
// 		if err != nil {
// 			return err
// 		}
// 		f.Value = s
// 	case CustomFieldTypeSingleList, CustomFieldTypeRadio:
// 		var i customFieldListItem
// 		err := json.Unmarshal(*raw.Value, &i)
// 		if err != nil {
// 			return err
// 		}
// 		f.Value = i.Name
// 	case CustomFieldTypeMultipleList, CustomFieldTypeCheckbox:
// 		var is []*customFieldListItem
// 		err := json.Unmarshal(*raw.Value, &is)
// 		if err != nil {
// 			return err
// 		}

// 		args := make([]string, len(is))
// 		for i, it := range is {
// 			args[i] = it.Name
// 		}
// 		f.Value = args
// 	}

// 	return nil
// }

type customFieldListItem struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DisplayOrder int    `json:"displayOrder"`
}

//-custom field

//+custom field property
type customFieldProperties []*customFieldProperty

func (ps customFieldProperties) findFieldType(fieldName string) (CustomFieldType, error) {
	var property *customFieldProperty

	var seenCount int
	for _, p := range ps {
		if seenCount > 1 {
			return 0, fmt.Errorf("field name '%s' is ambiquous", fieldName)
		}
		if p.Name == fieldName {
			property = p
			seenCount++
		}
	}

	if property == nil {
		return 0, fmt.Errorf("field name '%s' is not found", fieldName)
	}

	return property.FieldType, nil
}

func (ps customFieldProperties) findFieldID(fieldName string) (int, error) {
	var property *customFieldProperty

	var seenCount int
	for _, p := range ps {
		if seenCount > 1 {
			return 0, fmt.Errorf("field name '%s' is ambiquous", fieldName)
		}
		if p.Name == fieldName {
			property = p
			seenCount++
		}
	}

	return property.ID, nil
}

func (ps customFieldProperties) findItemID(fieldName string, itemName string) (int, error) {
	var property *customFieldProperty
	var item *customFieldListItem

	//+property
	var seenCount int
	for _, p := range ps {
		if seenCount > 1 {
			return 0, fmt.Errorf("field name '%s' is ambiquous", fieldName)
		}
		if p.Name == fieldName {
			property = p
			seenCount++
		}
	}

	if property == nil {
		return 0, fmt.Errorf("field name is invalid")
	}
	//+property

	//+list item
	seenCount = 0
	for _, _item := range property.ListItems {
		if seenCount > 1 {
			return 0, fmt.Errorf("field list item name '%s' is ambiguous", item.Name)
		}
		if _item.Name == itemName {
			item = _item
			seenCount++
		}
	}

	if item == nil {
		return 0, fmt.Errorf("field list item name '%s' is invalid", itemName)
	}
	//-list item

	return item.ID, nil
}

type customFieldProperty struct {
	ID                   int                    `json:"id"`
	FieldType            CustomFieldType        `json:"typeID"`
	Name                 string                 `json:"name"`
	Description          string                 `json:"description"`
	Required             bool                   `json:"required"`
	ApplicableIssueTypes []interface{}          `json:"applicableIssueTypes"`
	AllowAddItem         bool                   `json:"allowAddItem"`
	ListItems            []*customFieldListItem `json:"items"` // fieldTypeがlistの時のみ
}

//-custom field property
