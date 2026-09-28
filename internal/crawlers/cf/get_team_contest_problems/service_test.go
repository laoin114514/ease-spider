package getteamcontestproblems

import (
	"testing"

	"spider/config/db"

	cf "github.com/laoin114514/codeforcesClient"
)

func TestBuildTeamProblemTable(t *testing.T) {
	g := &GetTeamContestProblems{}
	got := g.buildTeamProblemTable(&cf.Problem{ContestID: 1000, Index: "A", Name: "A. 题", Rating: 1500})

	// 字段映射与迁移前一致：Team_contest_name 取题目标题，Official_contest_ID 不赋值
	want := db.Cf_team_problems{
		Team_contest_id:   1000,
		Team_contest_name: "A. 题",
		Problem_name:      "A. 题",
		Rating:            1500,
	}
	if got != want {
		t.Fatalf("buildTeamProblemTable() = %+v, want %+v", got, want)
	}
}

func TestBuildTeamProblemTableUnrated(t *testing.T) {
	g := &GetTeamContestProblems{}
	if got := g.buildTeamProblemTable(&cf.Problem{Index: "B"}); got.Rating != -1 {
		t.Fatalf("Rating = %d, want -1", got.Rating)
	}
}
