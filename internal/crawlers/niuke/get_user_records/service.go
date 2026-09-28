package getuserrecords

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"spider/internal/models"
	"spider/internal/repository"

	easecrawler "spider/pkg/crawler"
)

// APIRecord 单条提交记录（完整映射API返回的所有字段）
type APIRecord struct {
	Problem struct {
		QuestionNum string `json:"questionNum"`
		Title       string `json:"title"`
	} `json:"problem"`
	Submission struct {
		TimeConsumption   int   `json:"timeConsumption"`
		CreatedDate       int64 `json:"createdDate"`
		MemoryConsumption int   `json:"memoryConsumption"`
	} `json:"submission"`
	Language string `json:"language"`
	Status   struct {
		Desc string `json:"desc"`
	} `json:"status"`
}

// APIResponse 牛客API响应结构
type APIResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    struct {
		Current   int         `json:"current"`
		Size      int         `json:"size"`
		Total     int         `json:"total"`
		TotalPage int         `json:"totalPage"`
		Records   []APIRecord `json:"records"`
	} `json:"data"`
}

// Service 业务逻辑服务
type Service struct {
	client      *http.Client
	concurrency int
	log         *easecrawler.EaseLogger
	repo        *repository.NiukeRepository
}

// NewService 创建服务实例
func NewService(log *easecrawler.EaseLogger, repo *repository.NiukeRepository) *Service {
	return &Service{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		concurrency: 10,
		log:         log,
		repo:        repo,
	}
}

// GetAndStore 获取并存储所有用户记录
func (s *Service) GetAndStore() error {
	// 从数据库获取用户列表
	users, err := s.repo.GetUserNameMap()
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %v", err)
	}

	s.log.Printf("共找到 %d 个用户需要爬取", len(users))

	// 并发处理所有用户
	var wg sync.WaitGroup
	userChan := make(chan models.NiukeUserDeliver, len(users))

	// 将用户放入通道
	for _, user := range users {
		userChan <- user
	}
	close(userChan)

	// 启动工作协程
	for i := 0; i < s.concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for user := range userChan {
				s.log.Printf("[Worker %d] 正在处理用户: %s (UID: %s)", workerID, user.RealName, user.Uid)
				if err := s.processUser(user); err != nil {
					s.log.Printf("[Worker %d] 处理用户 %s 失败: %v", workerID, user.RealName, err)
				}
			}
		}(i)
	}

	wg.Wait()
	s.log.Printf("所有用户处理完成")
	return nil
}

// processUser 处理单个用户
func (s *Service) processUser(user models.NiukeUserDeliver) error {
	// 获取第1页，得到总页数
	firstPage, err := s.fetchPage(user.Uid, 1)
	if err != nil {
		return fmt.Errorf("获取第1页失败: %v", err)
	}

	if !firstPage.Success {
		return fmt.Errorf("API返回错误: %s", firstPage.Msg)
	}

	totalPage := firstPage.Data.TotalPage
	total := firstPage.Data.Total

	s.log.Printf("用户 %s: 总记录数 %d, 共 %d 页", user.RealName, total, totalPage)

	// 处理第1页数据
	newRecords := s.filterNewRecords(user.Uid, user.RealName, firstPage.Data.Records, user.OldDataSet)
	if len(newRecords) > 0 {
		if err := s.repo.InsertUserRecords(user.Uid, newRecords); err != nil {
			return fmt.Errorf("插入记录失败: %v", err)
		}
		s.log.Printf("用户 %s: 第1页插入 %d 条新记录", user.RealName, len(newRecords))
	}

	// 如果只有1页，直接返回
	if totalPage <= 1 {
		return nil
	}

	// 并发获取剩余页面
	pageChan := make(chan int, totalPage-1)
	resultChan := make(chan struct {
		page    int
		records []models.NiukeSubmission
		err     error
	}, totalPage-1)

	// 将页码放入通道
	for page := 2; page <= totalPage; page++ {
		pageChan <- page
	}
	close(pageChan)

	// 启动页面爬取协程
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ { // 每个用户用3个并发
		wg.Add(1)
		go func() {
			defer wg.Done()
			for page := range pageChan {
				data, err := s.fetchPage(user.Uid, page)
				if err != nil {
					resultChan <- struct {
						page    int
						records []models.NiukeSubmission
						err     error
					}{page, nil, err}
					continue
				}

				records := s.filterNewRecords(user.Uid, user.RealName, data.Data.Records, user.OldDataSet)
				resultChan <- struct {
					page    int
					records []models.NiukeSubmission
					err     error
				}{page, records, nil}

				time.Sleep(100 * time.Millisecond)
			}
		}()
	}

	// 等待所有页面爬取完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	totalNew := len(newRecords)
	for result := range resultChan {
		if result.err != nil {
			s.log.Printf("用户 %s 第 %d 页失败: %v", user.RealName, result.page, result.err)
		} else if len(result.records) > 0 {
			if err := s.repo.InsertUserRecords(user.Uid, result.records); err != nil {
				s.log.Printf("用户 %s 插入记录失败: %v", user.RealName, err)
			} else {
				totalNew += len(result.records)
			}
		}
	}

	s.log.Printf("用户 %s: 共插入 %d 条新记录", user.RealName, totalNew)
	return nil
}

// filterNewRecords 过滤出新记录（去重）
func (s *Service) filterNewRecords(uid string, userName string, records []APIRecord, oldDataSet map[string]bool) []models.NiukeSubmission {
	var newRecords []models.NiukeSubmission
	for _, record := range records {
		submitTime := time.Unix(record.Submission.CreatedDate/1000, 0).Format("2006-01-02")
		subId := uid + "_" + record.Problem.QuestionNum + "_" + submitTime

		// 如果已存在，跳过
		if oldDataSet[subId] {
			continue
		}

		newRecords = append(newRecords, models.NiukeSubmission{
			SubId:       subId,
			Uid:         uid,
			UserName:    userName,
			ProblemId:   record.Problem.QuestionNum,
			ProblemName: record.Problem.Title,
			SubmitTime:  submitTime,
			Status:      record.Status.Desc,
			Language:    record.Language,
		})
	}
	return newRecords
}

// fetchPage 获取指定页的数据
func (s *Service) fetchPage(uid string, page int) (*APIResponse, error) {
	timestamp := time.Now().UnixMilli()
	apiURL := fmt.Sprintf("https://gw-c.nowcoder.com/api/sparta/user/question-training/submission-history?_=%d", timestamp)

	requestBody := fmt.Sprintf(`{"pageNo":%d,"pageSize":20,"userId":%s}`, page, uid)
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://www.nowcoder.com/")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://www.nowcoder.com")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result APIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
