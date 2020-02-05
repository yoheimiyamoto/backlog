package backlog

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// APIEndpoint constants
const (
	APIEndpointBase       = "https://%s.backlog.jp"
	APIEndpointRecord     = "/k/v1/record.json"
	APIEndpointRecords    = "/k/v1/records.json"
	APIEndpointApp        = "/k/v1/app.json"
	APIEndpointFormField  = "/k/v1/app/form/fields.json"
	APIEndpointFormLayout = "/k/v1/app/form/layout.json"
	APIEndpointFile       = "/k/v1/file.json"
)

type ParentChildType int32

const (
	ParentChildTypeAll ParentChildType = iota
	ParentChildTypeExceptingChildren
	ParentChildTypeChildren
	ParentChildTypeStandAlone
	ParentChildTypeParents
)

// Repository ...
type Repository struct {
	client Client
}

func NewRepository(subdomain, apiKey string) *Repository {
	c := newClient(subdomain, apiKey, nil)
	return &Repository{client: c}
}

func (repo *Repository) FindIssue(id int) (*Issue, error) {
	url := fmt.Sprintf("api/v2/issues/%d", id)

	data, err := repo.client.get(url, nil)
	if err != nil {
		return nil, err
	}

	var issue Issue
	err = json.Unmarshal(data, &issue)
	if err != nil {
		return nil, err
	}

	return &issue, nil
}

func (repo *Repository) SearchIssues(q SearchIssueQuery) ([]*Issue, error) {
	data, err := repo.client.get("api/v2/issues", url.Values(q))
	if err != nil {
		return nil, err
	}

	var issues []*Issue
	err = json.Unmarshal(data, &issues)
	if err != nil {
		return nil, err
	}

	return issues, nil
}

//+SearchIssueQuery
type SearchIssueQuery url.Values

func NewSearchIssueQuery() SearchIssueQuery {
	q := SearchIssueQuery{}
	return q
}

func (q SearchIssueQuery) SetProjectID(id int) {
	url.Values(q).Add("projectId[]", strconv.Itoa(id))
	return
}

func (q SearchIssueQuery) SetIssueID(id int) {
	url.Values(q).Add("id[]", strconv.Itoa(id))
	return
}

func (q SearchIssueQuery) SetParentIssueID(id int) {
	url.Values(q).Add("parentIssueId[]", strconv.Itoa(id))
	return
}

func (q SearchIssueQuery) SetSort(how string) SearchIssueQuery {
	url.Values(q).Add("sort", how)
	return q
}

//-SearchIssueQuery

func (repo *Repository) UpdateIssue(issue *Issue) error {
	url := fmt.Sprintf("api/v2/issues/%d", issue.ID)

	property, err := repo.getCustomFieldProperty(issue.ProjectID)
	if err != nil {
		return errors.Wrap(err, "get custom field property failed")
	}

	params, err := issue.params(property)
	if err != nil {
		return errors.Wrap(err, "create params failed")
	}
	log.Printf("params: %v", params)

	_, err = repo.client.put(url, params)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) getCustomFieldProperty(projectID int) (customFieldProperties, error) {
	url := fmt.Sprintf("api/v2/projects/%d/customFields", projectID)

	data, err := repo.client.get(url, nil)
	if err != nil {
		return nil, err
	}

	var ps customFieldProperties
	err = json.Unmarshal(data, &ps)
	if err != nil {
		return nil, err
	}

	return ps, nil
}
