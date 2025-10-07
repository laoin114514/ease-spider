package models

type CfTeamContestProblemsResponse struct {
	CfResponse[cfTeamContestProblemsResult]
}
type cfTeamContestProblemsResult struct {
	Contest  CfContest   `json:"contest"`
	Problems []CfProblem `json:"problems"`
}
