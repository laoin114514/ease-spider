package getteamcontests

import (
	"testing"
	"time"

	"spider/config/db"
	"spider/internal/constants"

	cf "github.com/laoin114514/codeforcesClient"
)

func TestBuildTeamContestTable(t *testing.T) {
	g := &GetTeamContests{}
	got := g.buildTeamContestTable(&cf.Contest{
		ID:               1000,
		Name:             "团队训练赛",
		PreparedBy:       "233zhang",
		StartTimeSeconds: 1700000000,
	})

	want := db.Cf_team_contests{
		Contest_id:   1000,
		Contest_name: "团队训练赛",
		PrePare_by:   "233zhang",
		Start_time:   time.Unix(1700000000, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
	}
	if got.Contest_id != want.Contest_id || got.Contest_name != want.Contest_name || got.PrePare_by != want.PrePare_by {
		t.Fatalf("buildTeamContestTable() = %+v", got)
	}
	if !got.Start_time.Equal(want.Start_time) {
		t.Fatalf("Start_time = %v, want %v", got.Start_time, want.Start_time)
	}
}

func TestBuildTeamContestTableFallsBackToNow(t *testing.T) {
	g := &GetTeamContests{}
	if got := g.buildTeamContestTable(&cf.Contest{ID: 1, Name: "x"}); got.Start_time.IsZero() {
		t.Fatal("StartTimeSeconds 为 0 时应回退到当前时间，而不是零值时间")
	}
}
