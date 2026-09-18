package getuserrecords

import (
	"context"
	"fmt"
	"spider/config"
	"spider/config/db"
	"spider/internal/constants"
	"spider/internal/models"
	"spider/pkg/luogu2api"
	"strconv"
	"time"

	easecrawler "github.com/laoin114514/ease-crawler"
)

// GetAndStore 获取并入库洛谷用户提交记录
func (g *GetUserRecords) GetAndStore() error {
	client, err := g.getClient()
	if err != nil {
		return err
	}

	luoguUserDelivers, err := g.repo.GetUserNameMap()
	if err != nil {
		return err
	}
	g.log.Printf("总共 %d个用户需要获取提交记录", len(luoguUserDelivers))

	// 通过并发器来获取洛谷用户提交记录
	conCurrenter := easecrawler.NewConCurrenter[models.LuoguUserDeliver](config.AppConfig.Luogu.LuoguRecordsConcurrency)
	conCurrenter.SetLogger(g.log)
	// 单个用户的失败由并发器收集并打印进插件日志，这里不重复返回
	_ = conCurrenter.Run(luoguUserDelivers, func(luoguUser models.LuoguUserDeliver) error {
		return g.processUserRecords(client, luoguUser)
	})

	return nil
}

// ChangePrivateProblem 修改私有题目难度为unknown - 对外提供的接口
func (g *GetUserRecords) ChangePrivateProblem() error {
	count, err := g.repo.ChangePrivateProblem()
	if err != nil {
		g.log.Errorf("修改私有题目难度为unknow失败 %s", err.Error())
		return err
	}
	if count == 0 {
		return nil
	}
	g.log.Printf("修改私有题目难度为unknow完成 %d", count)
	return nil
}

// getClient 构造 luogu2api SDK 客户端：懒加载，构造一次后复用。
//
// 记录接口需要洛谷登录态，服务端会从号池选号、cookie 失效时自动换号重试，
// 所以这里只需要服务地址与访问令牌，不需要账号密码。
func (g *GetUserRecords) getClient() (*luogu2api.Client, error) {
	if g.sdk != nil {
		return g.sdk, nil
	}
	cfg := config.AppConfig.Luogu2Api
	client, err := luogu2api.NewSDK(cfg.BaseURL, cfg.AdminToken)
	if err != nil {
		return nil, fmt.Errorf("初始化luogu2api SDK失败(检查配置 luogu2api.baseUrl / luogu2api.adminToken) %s", err.Error())
	}
	g.sdk = client
	return g.sdk, nil
}

// processUserRecords 处理单个用户的提交记录 - 私有方法
func (g *GetUserRecords) processUserRecords(client *luogu2api.Client, luoguUser models.LuoguUserDeliver) error {
	uid, err := luoguUID(&luoguUser)
	if err != nil {
		return err
	}

	// 获取初始化数据：总数和总页数（第 1 页顺带交给增量爬取，不重复请求）
	initPage, err := g.fetchRecordPage(client, uid, 1)
	if err != nil {
		if luogu2api.IsNotFound(err) {
			return fmt.Errorf("%s的uid不存在", luoguUser.RealName)
		}
		return fmt.Errorf("%s获取提交记录失败 %s", luoguUser.RealName, err.Error())
	}
	if initPage.Count == 0 {
		g.log.Printf("%s在洛谷没有提交记录", luoguUser.RealName)
		return nil
	}

	// 增量爬取
	err = g.fetchRecordsIncrementally(client, &luoguUser, initPage)
	if err != nil {
		return err
	}

	// 如果已爬取数据与数据库已爬取数据数量一致，则不进行全量爬取
	if g.isDataConsistent(luoguUser, initPage.Count) {
		//如果增量爬取没有数据，则不进行全量爬取
		if luoguUser.Count == 0 {
			g.log.Printf("%s无新增过题数据", luoguUser.RealName)
			return nil
		}
		g.log.Printf("%s获取提交记录完成 %d", luoguUser.RealName, luoguUser.Count)
		return nil
	}

	// 爬取数据与数据库已爬取数据不一致，进行全量爬取
	g.logDataInconsistency(luoguUser, initPage.Count)
	err = g.fetchRecordsFully(client, &luoguUser, initPage.TotalPages)
	if err != nil {
		return fmt.Errorf("%s重新获取提交记录失败 %s", luoguUser.RealName, err.Error())
	}

	if luoguUser.Count == 0 {
		return fmt.Errorf("%s重新获取提交记录失败，无新增过题数据", luoguUser.RealName)
	}
	g.log.Printf("%s重新获取提交记录完成 %d", luoguUser.RealName, luoguUser.Count)
	return nil
}

// fetchRecordPage 通过 luogu2api SDK 获取指定页的提交记录 - 私有方法
func (g *GetUserRecords) fetchRecordPage(client *luogu2api.Client, uid int, page int) (*luogu2api.RecordPage, error) {
	// SDK 自带 30s 单次请求超时，这里的 ctx 只负责取消传播，不再叠加 deadline
	return client.Record.ListByUser(context.Background(), uid, luogu2api.RecordListParams{Page: page})
}

// fetchRecordsIncrementally 增量爬取不重复数据 - 私有方法
func (g *GetUserRecords) fetchRecordsIncrementally(client *luogu2api.Client, luoguUser *models.LuoguUserDeliver, firstPage *luogu2api.RecordPage) error {
	uid, err := luoguUID(luoguUser)
	if err != nil {
		return err
	}
	luoguUser.Count = 0

	// 第 1 页初始化时已经取回，直接处理
	if g.processPageRecords(firstPage.Records, luoguUser, 1, true) {
		return nil
	}

	for page := 2; page <= firstPage.TotalPages; page++ {
		records, err := g.fetchRecordPage(client, uid, page)
		if err != nil {
			return fmt.Errorf("%s获取第%d页提交记录失败 %s", luoguUser.RealName, page, err.Error())
		}
		if g.processPageRecords(records.Records, luoguUser, page, true) {
			return nil
		}
	}
	return nil
}

// fetchRecordsFully 全量爬取所有页数的数据
func (g *GetUserRecords) fetchRecordsFully(client *luogu2api.Client, luoguUser *models.LuoguUserDeliver, totalPages int) error {
	uid, err := luoguUID(luoguUser)
	if err != nil {
		return err
	}
	luoguUser.Count = 0

	for page := 1; page <= totalPages; page++ {
		records, err := g.fetchRecordPage(client, uid, page)
		if err != nil {
			// 全量爬取是数据不一致时的兜底修复，单页失败不中断，继续翻后面的页
			g.log.Errorf("%s获取第%d页提交记录失败 %s", luoguUser.RealName, page, err.Error())
			continue
		}
		g.processPageRecords(records.Records, luoguUser, page, false)
	}
	return nil
}

// processPageRecords 处理单页记录，返回是否应当停止继续翻页。
//
// incremental 为 true（增量爬取）时：更新的记录只会出现在最前面几页，遇到已入库的记录
// 或插入失败就说明增量到头了，返回 true 交由调用方结束本次增量；
// 为 false（全量爬取）时：跳过已入库的记录，单条插入失败只记日志，继续处理本页剩余记录。
func (g *GetUserRecords) processPageRecords(records []luogu2api.RecordSummary, luoguUser *models.LuoguUserDeliver, page int, incremental bool) bool {
	for i := range records {
		record := &records[i]
		if g.shouldSkipRecord(record, luoguUser, page) {
			continue
		}

		if luoguUser.OldDataSet[strconv.Itoa(record.ID)] {
			if incremental {
				return true
			}
			continue
		}

		err := db.Insert_luogu_sub(g.buildSubmissionTable(record))
		if err != nil {
			g.log.Errorf("%s处理第%d页提交记录失败 %s", luoguUser.RealName, page, err.Error())
			if incremental {
				return true
			}
			continue
		}
		luoguUser.Count++
	}
	return false
}

// shouldSkipRecord 判断是否应该跳过该提交记录
func (g *GetUserRecords) shouldSkipRecord(record *luogu2api.RecordSummary, luoguUser *models.LuoguUserDeliver, page int) bool {
	// 如果提交记录状态为还在测评，则跳过
	nowTime := time.Now().Unix()
	if record.Status == 0 && nowTime-record.SubmitTime < 60*5 {
		g.log.Warnf("%s第%d页提交记录还在测评 %s||https://www.luogu.com.cn/record/%d", luoguUser.RealName, page, record.Problem.PID, record.ID)
		return true
	}
	// 如果洛谷UID不匹配，则更新洛谷UID
	if record.User.UID == 0 {
		record.User.UID, _ = strconv.Atoi(luoguUser.Uid)
	}
	return false
}

// buildSubmissionTable 构建提交记录表
func (g *GetUserRecords) buildSubmissionTable(record *luogu2api.RecordSummary) db.Luogu_all_submissions {
	return db.Luogu_all_submissions{
		Sub_id:        strconv.Itoa(record.ID),
		Uid:           strconv.Itoa(record.User.UID),
		Problem_id:    record.Problem.PID,
		Problem_name:  record.Problem.Title,
		Difficulty:    luoguDifficulty(record.Problem.Difficulty),
		Is_pass:       record.Status == constants.LuoguStatusAccepted,
		Creation_time: time.Unix(record.SubmitTime, 0).Add(constants.TimeZoneOffsetHours * time.Hour),
	}
}

// luoguUID 解析洛谷UID；数据库里的值不是数字时返回带真实姓名的错误 - 私有方法
func luoguUID(luoguUser *models.LuoguUserDeliver) (int, error) {
	uid, err := strconv.Atoi(luoguUser.Uid)
	if err != nil {
		return 0, fmt.Errorf("%s的洛谷UID %q 不是合法的数字 %s", luoguUser.RealName, luoguUser.Uid, err.Error())
	}
	return uid, nil
}

// luoguDifficulty 难度值转难度名；越界（洛谷新增难度等级）时回退为 unknown - 私有方法
func luoguDifficulty(difficulty int) string {
	if difficulty < 0 || difficulty >= len(constants.LuoguDifficultyMap) {
		return "unknown"
	}
	return constants.LuoguDifficultyMap[difficulty]
}

// isDataConsistent 检查数据一致性
func (g *GetUserRecords) isDataConsistent(luoguUser models.LuoguUserDeliver, totalCount int) bool {
	return len(luoguUser.OldDataSet)+luoguUser.Count == totalCount
}

// logDataInconsistency 记录数据不一致日志
func (g *GetUserRecords) logDataInconsistency(luoguUser models.LuoguUserDeliver, actualCount int) {
	msg := fmt.Sprintf("%s爬取实际数量%d，数据库已爬取数量%d", luoguUser.RealName, actualCount, len(luoguUser.OldDataSet)+luoguUser.Count)
	g.log.Warnf(msg)
}
