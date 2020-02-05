package backlog

import (
	"encoding/json"
	"fmt"
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

func (i *Issue) params(property customFieldProperties) (url.Values, error) {
	params := url.Values{
		"summary": {i.Summary},
		// "parentIssueId": {strconv.Itoa(i.ParentIssueID)},
		"description": {i.Description},
	}

	for fieldName, value := range i.CustomFields {
		//+key
		fieldID, err := property.findFieldID(fieldName)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("find field id of field name '%s' failed", fieldName))
		}
		key := fmt.Sprintf("customField_%d", fieldID)
		//-key

		fieldType, err := property.findFieldType(fieldName)
		if err != nil {
			return nil, err
		}

		//+add custom fields to params
		switch fieldType {
		// valueがstringのケース
		case CustomFieldTypeText, CustomFieldTypeSentence:
			v, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("custom field '%s' value is invalid", value)
			}
			params.Add(key, v)

		// valueがlistItemのケース
		case CustomFieldTypeSingleList, CustomFieldTypeRadio:
			if value == nil {
				params.Add(key, "")
				continue
			}

			v, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("custom field '%s' value is invalid", fieldName)
			}

			if value == "" {
				params.Add(key, "")
				continue
			}

			itemID, err := property.findItemID(fieldName, v)
			if err != nil {
				return nil, fmt.Errorf("find item '%s' in custom field '%s' failed", v, fieldName)
			}

			params.Add(key, strconv.Itoa(itemID))

		// valueがlistItemsのケース
		case CustomFieldTypeMultipleList, CustomFieldTypeCheckbox:
			vs, ok := value.([]string)
			if !ok {
				return nil, fmt.Errorf("custom field '%s' value is invalid", fieldName)
			}

			itemIDs := make([]string, len(vs))
			for i, v := range vs {
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
type CustomFields map[string]interface{}

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

		switch r.CustomFieldType {
		case CustomFieldTypeText, CustomFieldTypeSentence:
			var v string
			err = json.Unmarshal(*r.Value, &v)
			if err != nil {
				return err
			}
			_fs[r.Name] = v
		case CustomFieldTypeSingleList, CustomFieldTypeRadio:
			var v struct {
				ID           int    `json:"id"`
				Name         string `json:"name"`
				DisplayOrder int    `json:"displayOrder"`
			}
			err = json.Unmarshal(*r.Value, &v)
			if err != nil {
				return err
			}
			_fs[r.Name] = v.Name
		case CustomFieldTypeMultipleList, CustomFieldTypeCheckbox:
			var items []*customFieldListItem

			err = json.Unmarshal(*r.Value, &items)
			if err != nil {
				return err
			}

			vs := make([]string, len(items))
			for i, item := range items {
				vs[i] = item.Name
			}
			_fs[r.Name] = vs
		}
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
