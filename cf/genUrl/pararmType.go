package genUrl

type User_info struct {
	Handles              []string `required`
	CheckHistoricHandles bool
}
type User_status struct {
	Handle         string `required`
	From           int
	Count          int
	IncludeSources bool
}
type User_blogEntries struct {
	Handle string `required`
}
type User_friend struct {
	Handle     string `while apikey`
	OnlyOnline bool
}
type User_ratedList struct {
	Handle         string `while apikey`
	ActiveOnly     bool
	IncludeRetired bool
	ContestId      string
}
type User_rating struct {
	Handle string `required`
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////  Contest type
type Contest_hacks struct {
	Handle    string `while apikey`
	ContestId string
	AsManager bool
}
type Contest_list struct {
	Handle    string `while apikey`
	Gym       bool
	GroupCode int
}
type Contest_rating struct {
	Handle    string `while apikey`
	ContestId string
}
type Contest_standings struct {
	Handle           string `while apikey`
	ContestId        string
	AsManager        bool
	From             int
	Count            int
	Room             string
	ShowUnofficial   bool
	ParticipantTypes string
}
type Contest_status struct {
	ContestId      string
	Asmanager      bool
	Handle         string
	From           int
	Count          int
	includeSources bool
}
