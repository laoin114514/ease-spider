package services

import (
	"spider/src/repository"
	"spider/src/utils"
)

type Dingding struct {
	repo *repository.DingdingRepository
	// req  *utils.Request[models.DingdingResponse]
	log *utils.LogContainer
}

func NewDingdingService() *Dingding {
	return &Dingding{
		repo: repository.NewDingdingRepository(),
		// req:  utils.NewRequest[models.DingdingResponse](),
		log: utils.NewLogContainer(),
	}
}
