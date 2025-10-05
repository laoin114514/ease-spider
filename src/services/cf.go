package services

import (
	"spider/src/models"
	"spider/src/repository"
	"spider/src/utils"
)

type CfService struct {
	repo *repository.CfRepository
	log  *utils.LogContainer
}

func NewCfService() *CfService {
	return &CfService{
		repo: repository.NewCfRepository(),
		log:  utils.NewLogContainer(),
	}
}

func (s *CfService) GetCfRecords(concurrency int) error {
	//构建工具类
	conCurrenter := utils.NewConCurrenter[models.CfUserData](concurrency)

	//获取cf用户数据
	cfUserDatas, err := s.repo.GetCfAccountData()
	if err != nil {
		return err
	}
	conCurrenter.Run(cfUserDatas, func(cfUserData models.CfUserData) error {
		return nil
	})
	return nil
}
