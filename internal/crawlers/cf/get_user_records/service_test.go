package getuserrecords

import (
	"io"
	"testing"
	"time"

	"spider/internal/constants"
	"spider/internal/models"

	cf "github.com/laoin114514/codeforcesClient"
	easecrawler "spider/pkg/crawler"
)

func newTestLogger() *easecrawler.EaseLogger {
	return easecrawler.NewLogger(io.Discard, "", 0)
}

func TestBuildSubmissionTable(t *testing.T) {
	g := &GetUserRecords{log: newTestLogger()}
	record := &cf.Submission{
		ID:                  12345,
		CreationTimeSeconds: 1700000000,
		Verdict:             "OK",
		Problem:             &cf.Problem{ContestID: 1900, Index: "A", Name: "A. Test", Rating: 1200},
	}

	got := g.buildSubmissionTable(record, &models.CfUserData{Account: "alice"})

	if got.Sub_id != 12345 || got.Account != "alice" || got.Problem_id != "1900A" ||
		got.Problem_name != "A. Test" || got.Verdict != "OK" || got.Rating != 1200 {
		t.Fatalf("buildSubmissionTable() = %+v", got)
	}
	want := time.Unix(1700000000, 0).Add(constants.TimeZoneOffsetHours * time.Hour)
	if !got.Creation_time.Equal(want) {
		t.Fatalf("Creation_time = %v, want %v", got.Creation_time, want)
	}
}

func TestBuildSubmissionTableUnrated(t *testing.T) {
	g := &GetUserRecords{log: newTestLogger()}
	got := g.buildSubmissionTable(
		&cf.Submission{Problem: &cf.Problem{ContestID: 1, Index: "B"}},
		&models.CfUserData{},
	)
	if got.Rating != -1 {
		t.Fatalf("未评级题目 Rating = %d, want -1", got.Rating)
	}
}

func TestShouldSkipRecord(t *testing.T) {
	g := &GetUserRecords{log: newTestLogger()}
	user := &models.CfUserData{OldDataSet: map[int]bool{7: true}}

	cases := []struct {
		name   string
		record *cf.Submission
		skip   bool
	}{
		{"已入库", &cf.Submission{ID: 7, Problem: &cf.Problem{}}, true},
		{"评测中", &cf.Submission{ID: 8, Verdict: "TESTING", Problem: &cf.Problem{}}, true},
		{"新记录", &cf.Submission{ID: 9, Verdict: "OK", Problem: &cf.Problem{}}, false},
		{"缺少题目信息", &cf.Submission{ID: 10}, true},
	}
	for _, tc := range cases {
		if got := g.shouldSkipRecord(tc.record, user); got != tc.skip {
			t.Errorf("%s: shouldSkipRecord() = %v, want %v", tc.name, got, tc.skip)
		}
	}
}
