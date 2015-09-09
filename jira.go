// Package jira creates a client for the JIRA REST API.
//
//  jira := New("https://jira.atlassian.com/")
//  rapidViews, err := jira.RapidViews()
//  if err != nil {
//    panic(err)
//  }
//  fmt.Printf("%+v", rapidViews)
package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Jira struct {
	Url      string
	Username string
	Password string
}

func New(url string) *Jira {
	return NewAuth(url, "", "")
}

func NewAuth(url, username, password string) *Jira {
	return &Jira{
		url,
		username,
		password,
	}
}

type RapidViews struct {
	Views []View `json:"views"`
}

type View struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	CanEdit              bool   `json:"canEdit"`
	SprintSupportEnabled bool   `json:"sprintSupportEnabled"`
	Filter               `json:"filter"`
}

type Filter struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Query           string `json:"query"`
	Owner           `json:"owner"`
	CanEdit         bool `json:"canEdit"`
	IsOrderedByRank bool `json:"isOrederedByRank"`
}

type Owner struct {
	UserName     string `json:"userName"`
	DisplayName  string `json:"displayName"`
	RenderedLink string `json:"renderedLink"`
}

// RapidViews gets all available Rapid Views.
func (j *Jira) RapidViews() (*RapidViews, error) {
	url := fmt.Sprintf("%srest/greenhopper/latest/rapidviews/list", j.Url)
	body, err := j.request(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	rapidViews := &RapidViews{}
	return rapidViews, json.NewDecoder(body).Decode(&rapidViews)
}

type SprintQuery struct {
	Sprints     []Sprint `json:"sprints"`
	RapidViewId int      `json:"rapidViewId"`
}

type Sprint struct {
	Id               int    `json:"id"`
	Sequence         int    `json:"sequence"`
	Name             string `json:"name"`
	State            string `json:"state"`
	LinkedPagesCount int    `json:"linkedPagesCount"`
}

// SprintQuery get all sprints for given Rapid View Id.
func (j *Jira) SprintQuery(rapidViewId int) (*SprintQuery, error) {
	url := fmt.Sprintf("%srest/greenhopper/latest/sprintquery/%d", j.Url, rapidViewId)
	body, err := j.request(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	sprintQuery := &SprintQuery{}
	return sprintQuery, json.NewDecoder(body).Decode(&sprintQuery)
}

type SprintReport struct {
	Contents struct {
		IncompletedIssues []Issue `json:"incompletedIssues"`
	} `json:"contents"`
	Sprint `json:"sprint"`
}

type Issue struct {
	Id                int    `json:"id"`
	Key               string `json:"key"`
	Summary           string `json:"summary"`
	EpicField         `json:"epicField"`
	EstimateStatistic `json:"estimateStatistic"`
}

type EpicField struct {
	EpicKey   string `json:"epicKey"`
	EpicColor string `json:"epicColor"`
	Text      string `json:"text"`
}

type EstimateStatistic struct {
	StatFieldValue struct {
		Value float64 `json:"value"`
	} `json:"statFieldValue"`
}

// SprintReport get all issues for given Rapid View Id and Sprint Id.
func (j *Jira) SprintReport(rapidViewId, sprintId int) (*SprintReport, error) {
	url := fmt.Sprintf("%s/rest/greenhopper/latest/rapid/charts/sprintreport?rapidViewId=%d&sprintId=%d", j.Url, rapidViewId, sprintId)
	body, err := j.request(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	sprintReport := &SprintReport{}
	return sprintReport, json.NewDecoder(body).Decode(&sprintReport)
}

func (j *Jira) request(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if j.Username != "" && j.Password != "" {
		req.SetBasicAuth(j.Username, j.Password)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}
