package models

// LuoguRecordsResponse 洛谷记录列表响应
type LuoguRecordsResponse struct {
	Code            int    `json:"code"`
	CurrentTemplate string `json:"currentTemplate"`
	CurrentData     struct {
		Records struct {
			Result []struct {
				Time    int `json:"time"`
				Memory  int `json:"memory"`
				Problem struct {
					PID        string `json:"pid"`
					Title      string `json:"title"`
					Difficulty int    `json:"difficulty"`
					FullScore  int    `json:"fullScore"`
					Type       string `json:"type"`
				} `json:"problem"`
				Contest *struct {
					ID          int64  `json:"id"`
					Name        string `json:"name"`
					StartTime   int64  `json:"startTime"`
					EndTime     int64  `json:"endTime"`
					Type        int    `json:"type"`
					Rated       bool   `json:"rated"`
					Description string `json:"description"`
				} `json:"contest"`
				SourceCodeLength int   `json:"sourceCodeLength"`
				SubmitTime       int64 `json:"submitTime"`
				Language         int   `json:"language"`
				User             struct {
					UID    int64   `json:"uid"`
					Name   string  `json:"name"`
					Avatar string  `json:"avatar"`
					Slogan string  `json:"slogan"`
					Badge  *string `json:"badge"`
				} `json:"user"`
				ID       int64 `json:"id"`
				Status   int   `json:"status"`
				EnableO2 bool  `json:"enableO2"`
				Score    int   `json:"score"`
			} `json:"result"`
			PerPage int `json:"perPage"`
			Count   int `json:"count"`
		} `json:"records"`
	} `json:"currentData"` //主要数据
	CurrentTitle string `json:"currentTitle"`
	CurrentTheme any    `json:"currentTheme"`
	CurrentTime  int64  `json:"currentTime"`
	CurrentUser  any    `json:"currentUser"`
}

// 用户提交记录
type UserRecord struct {
	Time    int `json:"time"`
	Memory  int `json:"memory"`
	Problem struct {
		PID        string `json:"pid"`
		Title      string `json:"title"`
		Difficulty int    `json:"difficulty"`
		FullScore  int    `json:"fullScore"`
		Type       string `json:"type"`
	} `json:"problem"`
	Contest *struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		StartTime   int64  `json:"startTime"`
		EndTime     int64  `json:"endTime"`
		Type        int    `json:"type"`
		Rated       bool   `json:"rated"`
		Description string `json:"description"`
	} `json:"contest"`
	SourceCodeLength int   `json:"sourceCodeLength"`
	SubmitTime       int64 `json:"submitTime"`
	Language         int   `json:"language"`
	User             struct {
		UID    int64   `json:"uid"`
		Name   string  `json:"name"`
		Avatar string  `json:"avatar"`
		Slogan string  `json:"slogan"`
		Badge  *string `json:"badge"`
	} `json:"user"`
	ID       int64 `json:"id"`
	Status   int   `json:"status"`
	EnableO2 bool  `json:"enableO2"`
	Score    int   `json:"score"`
}
