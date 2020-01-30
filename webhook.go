package backlog

import "time"

type Webhook struct {
	ID      int       `json:"id"`
	Project *Project  `json:"project"`
	Issue   *Issue    `json:"content"`
	Created time.Time `json:"created"`
}

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
