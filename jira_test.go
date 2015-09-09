package jira

import (
	"testing"
)

var jira = New("https://jira.atlassian.com/")

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
	t.Logf("%+v", report)
}
