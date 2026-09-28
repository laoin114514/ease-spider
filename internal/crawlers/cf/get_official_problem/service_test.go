package getofficialproblem

import (
	"testing"

	"spider/config/db"

	cf "github.com/laoin114514/codeforcesClient"
)

func TestBuildProblemTable(t *testing.T) {
	g := &GetOfifcialProblems{}
	got := g.buildProblemTable(&cf.Problem{
		ContestID: 1900,
		Index:     "A",
		Name:      "A. Test",
		Points:    500,
		Rating:    1200,
		Tags:      []string{"greedy", "math"},
	})

	want := db.Cf_official_problems{
		Problem_id: "1900A",
		Title:      "A. Test",
		Points:     500,
		Rating:     1200,
		Tags:       `["greedy","math"]`,
	}
	if got != want {
		t.Fatalf("buildProblemTable() = %+v, want %+v", got, want)
	}
}

func TestBuildProblemTableUnratedAndNoTags(t *testing.T) {
	g := &GetOfifcialProblems{}
	got := g.buildProblemTable(&cf.Problem{ContestID: 1, Index: "B"})
	if got.Rating != -1 {
		t.Fatalf("Rating = %d, want -1", got.Rating)
	}
	// 空标签在迁移前会拼成 [""]，这里显式钉住该既有行为
	if got.Tags != `[""]` {
		t.Fatalf("Tags = %q, want %q", got.Tags, `[""]`)
	}
}
