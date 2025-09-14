package genUrl

type User struct{}

func (u *User) Info(query User_info) (string, error) {
	handleArr := ""
	for _, handle := range query.Handles {
		handleArr += handle + ";"
	}
	handleArr = handleArr[:len(handleArr)-1]
	return defaultUrl("user.info", map[string]any{
		"handles":              handleArr,
		"checkHistoricHandles": query.CheckHistoricHandles,
	})
}

func (u *User) Status(useApikey bool, query User_status) (string, error) {
	pararms := map[string]any{
		"handle":         query.Handle,
		"from":           query.From,
		"count":          query.Count,
		"includeSources": query.IncludeSources,
	}
	if useApikey {
		back, err := apiKeyUrl("user.status", pararms)
		if err != nil {
			return defaultUrl("user.status", pararms)
		}
		return back, err
	}
	return defaultUrl("user.status", pararms)
}

type Contest struct{}

func (c Contest) List(query Contest_list) string {
	pararms := map[string]any{
		"handle":    query.Handle,
		"gym":       query.Gym,
		"groupCode": query.GroupCode,
	}
	url, err := apiKeyUrl("contest.list", pararms)
	if err != nil {
		return ""
	}
	return url
}

func (c Contest) Standings(query Contest_standings) string {
	pararms := map[string]any{
		"handles":   query.Handle,
		"contestId": query.ContestId,
		"asManager": query.AsManager,
		"from":      query.From,
		"count":     query.Count,
		// "room":             query.Room,
		// "showUnofficial":   query.ShowUnofficial,
		// "participantTypes": query.ParticipantTypes,
	}
	url, err := apiKeyUrl("contest.standings", pararms)
	if err != nil {
		return ""
	}
	return url
}

type ProblemSet struct{}

func (p ProblemSet) Problems() (string, error) {
	return defaultUrl("problemset.problems", nil)
}
