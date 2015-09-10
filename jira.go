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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Jira struct {
	Url string
	Jar *cookiejar.Jar
}

func New(url string) (*Jira, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &Jira{
		url,
		jar,
	}, nil
}

func NewAuth(uri, value string) (*Jira, error) {
	cookies := []*http.Cookie{
		&http.Cookie{
			Name:  "JSESSIONID",
			Value: value,
		},
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	cookieURL, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	jar.SetCookies(cookieURL, cookies)
	return &Jira{
		uri,
		jar,
	}, nil
}

type CreateSessionResponse struct {
	Session struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"session"`
	LoginInfo `json:"loginInfo"`
}

type LoginInfo struct {
	FailedLoginCount    int    `json:"failedLoginCount"`
	LoginCount          int    `json:"loginCount"`
	LastFailedLoginTime string `json:"lastFailedLoginTime"`
	PreviousLoginTime   string `json:"previousLoginTime"`
}

type GetSessionResponse struct {
	Self      string `json:"self"`
	Name      string `json:"name"`
	LoginInfo `json:"loginInfo"`
}

func (j *Jira) CreateSession(username, password string) (*CreateSessionResponse, error) {
	url := fmt.Sprintf("%srest/auth/1/session", j.Url)
	creds := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		username,
		password,
	}
	result, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}
	data := bytes.NewReader(result)
	res, err := j.request("POST", url, "application/json", data)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	createSessionResponse := &CreateSessionResponse{}
	return createSessionResponse, json.NewDecoder(res.Body).Decode(&createSessionResponse)
}

func (j *Jira) GetSession() (*GetSessionResponse, error) {
	url := fmt.Sprintf("%srest/auth/1/session", j.Url)
	res, err := j.request("GET", url, "", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	getSessionResponse := &GetSessionResponse{}
	return getSessionResponse, json.NewDecoder(res.Body).Decode(&getSessionResponse)
}

func (j *Jira) DeleteSession() error {
	url := fmt.Sprintf("%srest/auth/1/session", j.Url)
	res, err := j.request("DELETE", url, "", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusUnauthorized {
		return errors.New(http.StatusText(res.StatusCode))
	}
	return nil
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
	res, err := j.request("GET", url, "", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	rapidViews := &RapidViews{}
	return rapidViews, json.NewDecoder(res.Body).Decode(&rapidViews)
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
	res, err := j.request("GET", url, "", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	sprintQuery := &SprintQuery{}
	return sprintQuery, json.NewDecoder(res.Body).Decode(&sprintQuery)
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
	url := fmt.Sprintf("%srest/greenhopper/latest/rapid/charts/sprintreport?rapidViewId=%d&sprintId=%d", j.Url, rapidViewId, sprintId)
	res, err := j.request("GET", url, "", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	sprintReport := &SprintReport{}
	return sprintReport, json.NewDecoder(res.Body).Decode(&sprintReport)
}

func (j *Jira) request(method, url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	client := &http.Client{Jar: j.Jar}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	j.Jar.SetCookies(req.URL, res.Cookies())
	return res, nil
}
