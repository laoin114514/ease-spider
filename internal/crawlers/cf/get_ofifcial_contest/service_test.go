package getofifcialcontest

import (
	"testing"
	"time"

	"spider/config/db"
	"spider/internal/constants"

	cf "github.com/laoin114514/codeforcesClient"
)

func TestBuildOfficialContestTable(t *testing.T) {
	g := &GetOfifcialContest{}
	got := g.buildOfficialContestTable(&cf.Contest{
		ID:               2000,
		Name:             "Codeforces Round 1000",
		Phase:            "FINISHED",
		StartTimeSeconds: 1700000000,
	})

	want := db.Cf_official_contests{
		Official_contest_id:   2000,
		Official_contest_name: "Codeforces Round 1000",
		Phase:                 "FINISHED",
		Start_time:            time.Unix(1700000000, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
	}
	if got.Official_contest_id != want.Official_contest_id ||
		got.Official_contest_name != want.Official_contest_name ||
		got.Phase != want.Phase {
		t.Fatalf("buildOfficialContestTable() = %+v", got)
	}
	if !got.Start_time.Equal(want.Start_time) {
		t.Fatalf("Start_time = %v, want %v", got.Start_time, want.Start_time)
	}
}
