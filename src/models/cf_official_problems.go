package models

type CfOfficialProblemsResponse struct {
	CfResponse[cfOfficialResult]
}
type cfOfficialResult struct {
	Problems []CfProblem `json:"problems"`
}
