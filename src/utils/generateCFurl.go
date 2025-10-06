package utils

import (
	"fmt"
	"math/rand"
	"spider/src/models"
	"spider/src/repository"
	"time"
)

const BaseUrl = "https://codeforces.com/api/"

type GenerateCFurl struct {
	baseUrl    string
	repository *repository.CfRepository
	User       *user
	Contest    *contest
	ProblemSet *problemSet
	*models.CfUserData
}
type user struct{}
type contest struct{}
type problemSet struct{}

func NewGenerateCFurl() *GenerateCFurl {
	return &GenerateCFurl{
		baseUrl:    BaseUrl,
		repository: repository.NewCfRepository(),
		User:       &user{},
		Contest:    &contest{},
		ProblemSet: &problemSet{},
	}
}

// 组合url和apikey
func combineUrlWithApikey[T any](handle string, method string, pararms T) (string, error) {
	baseUrl := BaseUrl
	repository := repository.NewCfRepository()
	//获取apikey
	apikey, secret, err := repository.GetCfApikey(handle)
	fmt.Printf("参数：%v,账号：%v\n", pararms, handle)
	if err != nil {
		return combineUrlWithNoApikey(method, pararms)
	}
	now := time.Now()
	time := now.Unix()
	rand.Seed(now.UnixNano())
	randomKey := rand.Intn(999999)
	//将参数转换为字符串
	pararmStr, err := NewStructTransfer(pararms).
		AddParam("apiKey", apikey).
		AddParam("time", time).
		ToOrderStr()
	if err != nil {
		return "", err
	}
	tail := fmt.Sprintf("%v?%v", method, pararmStr)
	hashCode := NewHashEncoder().Hash512(fmt.Sprintf("%v/%v?%v#%v", randomKey, method, tail, secret))
	url := fmt.Sprintf("%v%v?%v&apiSig=%v%v", baseUrl, method, tail, randomKey, hashCode)
	return url, nil
}
func combineUrlWithNoApikey[T any](method string, pararms T) (string, error) {
	baseUrl := BaseUrl
	pararmStr, err := NewStructTransfer(pararms).
		ToOrderStr()
	if err != nil {
		return "", err
	}
	tail := fmt.Sprintf("%v?%v", method, pararmStr)
	return fmt.Sprintf("%v%v?%v", baseUrl, method, tail), nil
}

// ============================================User============================================//
func (u *user) Status(query *models.UserStatusParams) (string, error) {
	return combineUrlWithApikey(query.Handle, "user.status", query)
}
func (u *user) Rating(query *models.UserRatingParams) (string, error) {
	return combineUrlWithApikey(query.Handle, "user.rating", query)
}
func (u *user) RatedList(query *models.UserRatedListParams) (string, error) {
	return combineUrlWithApikey("", "user.ratedList", query)
}

func (u *user) Info(query *models.UserInfoParams) (string, error) {
	return combineUrlWithApikey("", "user.info", query)
}

func (u *user) Friends(query *models.UserFriendsParams) (string, error) {
	return combineUrlWithApikey("", "user.friends", query)
}

func (u *user) BlogEntries(query *models.UserBlogEntriesParams) (string, error) {
	return combineUrlWithApikey("", "user.blogEntries", query)
}

func (u *user) RecentActions(query *models.RecentActionsParams) (string, error) {
	return combineUrlWithApikey("", "user.recentActions", query)
}

// ============================================Contest============================================//
func (c *contest) List(query *models.ContestListParams) (string, error) {
	return combineUrlWithApikey("", "contest.list", query)
}
func (c *contest) Standings(query *models.ContestStandingsParams) (string, error) {
	return combineUrlWithApikey("", "contest.standings", query)
}
func (c *contest) Status(query *models.ContestStatusParams) (string, error) {
	return combineUrlWithApikey("", "contest.status", query)
}

// ============================================ProblemSet============================================//
func (p *problemSet) Problems(query *models.ProblemsetProblemsParams) (string, error) {
	return combineUrlWithApikey("", "problemset.problems", query)
}
func (p *problemSet) RecentStatus(query *models.ProblemsetRecentStatusParams) (string, error) {
	return combineUrlWithApikey("", "problemset.recentStatus", query)
}
