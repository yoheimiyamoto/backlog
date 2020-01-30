package backlog

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
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
	Status         IssueStatus   `json:"status"`
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

// Status represents
type IssueStatus struct {
	ID           int    `json:"id,omitempty"`
	ProjectID    int    `json:"projectId,omitempty"`
	Name         string `json:"name,omitempty"`
	Color        string `json:"color,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
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

func (i *Issue) Params(property customFieldProperties) (url.Values, error) {
	params := url.Values{
		"summary": {i.Summary},
		// "parentIssueId": {strconv.Itoa(i.ParentIssueID)},
		"description": {i.Description},
	}

	for _, f := range i.CustomFields {
		if f.Value == nil {
			continue
		}

		//+field name
		fieldID, err := property.findFieldID(f.Name)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("find field id of field name '%s' failed", f.Name))
		}
		fieldName := fmt.Sprintf("customField_%d", fieldID)
		//-field name

		var value string

		switch f.FieldType {
		case CustomFieldTypeText, CustomFieldTypeSentence:
			v, ok := f.Value.(string)
			if !ok {
				return nil, fmt.Errorf("custom field '%s' value is invalid", f.Name)
			}
			value = v
		case CustomFieldTypeSingleList, CustomFieldTypeRadio:
			v, ok := f.Value.(string)
			if !ok {
				return nil, fmt.Errorf("custom field '%s' value is invalid", f.Name)
			}
			itemID, err := property.findItemID(f.Name, v)
			if err != nil {
				return nil, fmt.Errorf("find item '%s' in custom field '%s' failed", v, f.Name)
			}
			value = strconv.Itoa(itemID)
		case CustomFieldTypeMultipleList, CustomFieldTypeCheckbox:

		}

		params.Add(fieldName, value)
	}

	return params, nil
}

//+custom field
type CustomFields []*CustomField

type CustomField struct {
	FieldType CustomFieldType `json:"fieldTypeId"`
	Name      string          `json:"name"`
	Value     interface{}     // string or []string
}

func (f *CustomField) String() string {
	if f.Value == nil {
		return ""
	}
	return fmt.Sprint(f.Value)
}

func (fs CustomFields) Find(fieldName string) (*CustomField, error) {
	var f *CustomField
	var seenCount int
	for _, _f := range fs {
		if seenCount > 1 {
			return nil, errors.New("field name is ambiguous")
		}
		if _f.Name == fieldName {
			f = _f
			seenCount++
		}
	}

	if f == nil {
		return nil, errors.New("field name is invalid")
	}

	return f, nil
}

func (f *CustomField) UnmarshalJSON(data []byte) error {
	type Raw struct {
		// ID        int              `json:"id"`
		FieldType CustomFieldType  `json:"fieldTypeId"`
		Name      string           `json:"name"`
		Value     *json.RawMessage `json:"value"` // string or customFieldListItem
	}

	var raw Raw

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// f.ID = raw.ID
	f.FieldType = raw.FieldType
	f.Name = raw.Name

	if raw.Value == nil {
		return nil
	}

	switch raw.FieldType {
	case CustomFieldTypeText, CustomFieldTypeSentence:
		var s string
		err := json.Unmarshal(*raw.Value, &s)
		if err != nil {
			return err
		}
		f.Value = s
	case CustomFieldTypeSingleList, CustomFieldTypeRadio:
		var i customFieldListItem
		err := json.Unmarshal(*raw.Value, &i)
		if err != nil {
			return err
		}
		f.Value = i.Name
	case CustomFieldTypeMultipleList, CustomFieldTypeCheckbox:
		var is []*customFieldListItem
		err := json.Unmarshal(*raw.Value, &is)
		if err != nil {
			return err
		}

		args := make([]string, len(is))
		for i, it := range is {
			args[i] = it.Name
		}
		f.Value = args
	}

	return nil
}

type customFieldListItem struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DisplayOrder int    `json:"displayOrder"`
}

//-custom field

//+custom field property
type customFieldProperties []*customFieldProperty

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
