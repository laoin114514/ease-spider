package utils

import "spider/src/repository"

type GenerateCFurl struct {
	repository *repository.CfRepository
	User       *user
	Contest    *contest
	ProblemSet *problemSet
}
type user struct{}
type contest struct{}
type problemSet struct{}

func NewGenerateCFurl() *GenerateCFurl {
	return &GenerateCFurl{
		repository: repository.NewCfRepository(),
		User:       &user{},
		Contest:    &contest{},
		ProblemSet: &problemSet{},
	}
}
func (r *GenerateCFurl) combineUrlWithApikey[T any](handle string, method string, pararms T) string {
	apikey, secret, err := r.repository.GetCfApikey(handle)
	if err != nil {
		return ""
	}
	
}
