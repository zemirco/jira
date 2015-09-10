package jira

import (
	"testing"
)

var jira, _ = New("https://jira.atlassian.com/")

func TestRapidViews(t *testing.T) {
	rapidViews, err := jira.RapidViews()
	if err != nil {
		t.Fatal(err)
	}
	if len(rapidViews.Views) == 0 {
		t.Error("no rapid views or error")
	}
}

func TestSprintQuery(t *testing.T) {
	query, err := jira.SprintQuery(1159)
	if err != nil {
		t.Fatal(err)
	}
	if len(query.Sprints) == 0 {
		t.Error("no sprints or error")
	}
	if query.RapidViewId != 1159 {
		t.Error("wrong rapid view id")
	}
}

func TestSprintReport(t *testing.T) {
	report, err := jira.SprintReport(1159, 922)
	if err != nil {
		t.Fatal(err)
	}
	if report.Sprint.Id != 922 {
		t.Error("wrong sprint id")
	}
}

// func TestSession(t *testing.T) {
// 	// create session
// 	createSessionResponse, err := jira.CreateSession("username", "password")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if createSessionResponse.Session.Name != "JSESSIONID" {
// 		t.Error("wrong session name")
// 	}
//
// 	// get session
// 	res, err := jira.GetSession()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if res.Name == "" {
// 		t.Error("wrong session name")
// 	}
//
// 	// delete session
// 	err = jira.DeleteSession()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	// get session
// 	res, err = jira.GetSession()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if res.Name != "" {
// 		t.Error("session should not have a name")
// 	}
// }
