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
	// 同一个洛谷 uid 可能绑在 user 表的多行上（例如"张天杰"和"张天杰1"）：
	// 同一个账号被并发爬两次既浪费请求，也会让两边同时插入同一条记录
	luoguUserDelivers = dedupeByUID(luoguUserDelivers, g.log.Warnf)
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

// crawlPlan 一轮增量爬取之后该做什么
type crawlPlan int

const (
	// planConsistent 库内已同步：洛谷返回的记录都在库里，且没有多出来的行
	planConsistent crawlPlan = iota
	// planMissing 库内比洛谷少：可能有缺口，值得全量核对一遍
	planMissing
	// planExtra 库内比洛谷多：洛谷删除/隐藏过记录，或 user.luogu_uid 变更留下的旧行
	planExtra
)

// syncedCount 本轮结束时"库里该用户已有 + 本轮新增 + 还在测评"的条数 - 私有函数
func syncedCount(luoguUser models.LuoguUserDeliver) int {
	return len(luoguUser.OldDataSet) + luoguUser.Count + luoguUser.Pending
}

// planAfterIncremental 用条数判断是否需要全量核对 - 私有函数
//
// 这里只能当启发式用：洛谷删除/隐藏的记录不会从库里消失，user.luogu_uid 被改过时旧行也
// 还留在别的 uid 下，所以"库内比洛谷多"是正常的数据状态，不是需要修复的错误——原来的实现
// 把它当成"不一致"而反复触发全量爬取，是日志刷屏的根源之一。
func planAfterIncremental(luoguUser models.LuoguUserDeliver, total int) crawlPlan {
	switch got := syncedCount(luoguUser); {
	case got == total:
		return planConsistent
	case got < total:
		return planMissing
	default:
		return planExtra
	}
}

// dedupeByUID 按洛谷 uid 去重，重复的行只记一条告警 - 私有函数
func dedupeByUID(users []models.LuoguUserDeliver, warnf func(format string, args ...interface{})) []models.LuoguUserDeliver {
	seen := make(map[string]string, len(users))
	unique := make([]models.LuoguUserDeliver, 0, len(users))
	for _, user := range users {
		if first, ok := seen[user.Uid]; ok {
			warnf("洛谷UID %s 同时对应 %q 和 %q，本轮只爬取一次；建议清理 user 表的重复绑定", user.Uid, first, user.RealName)
			continue
		}
		seen[user.Uid] = user.RealName
		unique = append(unique, user)
	}
	return unique
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
		// 洛谷对不存在的 uid 也返回 200 + 0 条，所以这里只能提示，不能断言"确实没有提交"
		g.log.Warnf("%s在洛谷没有提交记录（uid=%s，若与实际不符请核对 user.luogu_uid）", luoguUser.RealName, luoguUser.Uid)
		return nil
	}

	// 增量爬取
	err = g.fetchRecordsIncrementally(client, &luoguUser, initPage)
	if err != nil {
		return err
	}

	switch planAfterIncremental(luoguUser, initPage.Count) {
	case planConsistent:
		if luoguUser.Count == 0 {
			g.log.Printf("%s无新增提交记录（洛谷共%d条）", luoguUser.RealName, initPage.Count)
			return nil
		}
		g.log.Printf("%s获取提交记录完成 新增%d条（洛谷共%d条）", luoguUser.RealName, luoguUser.Count, initPage.Count)
		return nil

	case planExtra:
		// 库里的行比洛谷多：全量爬取既补不了也删不掉，只会每轮打印一堆主键冲突，直接跳过
		g.log.Warnf("%s库内已有%d条、还在测评%d条，多于洛谷的%d条（洛谷删过记录，或 user.luogu_uid 变更遗留的旧行），跳过全量核对",
			luoguUser.RealName, len(luoguUser.OldDataSet)+luoguUser.Count, luoguUser.Pending, initPage.Count)
		return nil
	}

	// 库内比洛谷少：可能有缺口，全量核对一遍。写入是幂等的，已存在的记录不会再报主键冲突
	g.log.Warnf("%s库内只有%d条、还在测评%d条，少于洛谷的%d条，开始全量核对",
		luoguUser.RealName, len(luoguUser.OldDataSet)+luoguUser.Count, luoguUser.Pending, initPage.Count)
	if err := g.fetchRecordsFully(client, &luoguUser, initPage.TotalPages); err != nil {
		return err
	}

	if luoguUser.Count == 0 {
		// 全量核对一条都没补进来：说明洛谷返回的记录库里其实都有，只是不全挂在当前 uid 下，
		// 属于历史 uid 变更遗留，不是爬取失败（原先这里返回错误，于是每轮都报"无新增过题数据"）
		g.log.Warnf("%s全量核对完成：洛谷%d条记录库中都已存在，但当前 uid 下只有%d条，建议核对 user.luogu_uid 与历史行",
			luoguUser.RealName, initPage.Count, len(luoguUser.OldDataSet))
		return nil
	}
	g.log.Printf("%s全量核对完成，补录%d条", luoguUser.RealName, luoguUser.Count)
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
	luoguUser.Pending = 0

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
	luoguUser.Pending = 0

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
// 就说明增量到头了，返回 true 交由调用方结束本次增量；
// 为 false（全量爬取）时：跳过已入库的记录，单条插入失败只记日志，继续处理本页剩余记录。
//
// "已入库"以数据库为准（upsert 的返回值），不再只看当前 uid 名下的行：记录挂在别的 uid 下
// 也是已经爬过，不会再被当成新记录去撞主键。
func (g *GetUserRecords) processPageRecords(records []luogu2api.RecordSummary, luoguUser *models.LuoguUserDeliver, page int, incremental bool) bool {
	for i := range records {
		record := &records[i]
		skip, ownerKnown := g.shouldSkipRecord(record, luoguUser, page)
		if skip {
			// 还在测评：没入库，但也不算缺失，单独计数避免被判成"数据不一致"
			luoguUser.Pending++
			continue
		}

		if luoguUser.OldDataSet[strconv.Itoa(record.ID)] {
			if incremental {
				return true
			}
			continue
		}

		inserted, err := db.Upsert_luogu_sub(g.buildSubmissionTable(record), ownerKnown)
		if err != nil {
			g.log.Errorf("%s处理第%d页提交记录失败 %s", luoguUser.RealName, page, err.Error())
			if incremental {
				return true
			}
			continue
		}
		if !inserted {
			// 记录已在库中（可能挂在别的 uid 下）：增量遇到即到头，全量继续翻后面的页
			if incremental {
				return true
			}
			continue
		}
		luoguUser.Count++
	}
	return false
}

// shouldSkipRecord 判断是否应该跳过该提交记录，并说明记录里的提交者 uid 是否可信 - 私有方法
func (g *GetUserRecords) shouldSkipRecord(record *luogu2api.RecordSummary, luoguUser *models.LuoguUserDeliver, page int) (skip bool, ownerKnown bool) {
	// 如果提交记录状态为还在测评，则跳过
	nowTime := time.Now().Unix()
	if record.Status == 0 && nowTime-record.SubmitTime < 60*5 {
		g.log.Warnf("%s第%d页提交记录还在测评 %s||https://www.luogu.com.cn/record/%d", luoguUser.RealName, page, record.Problem.PID, record.ID)
		return true, false
	}
	// 洛谷偶尔不回提交者 uid：用查询用的 uid 兜底；这种兜底值不能用来纠正库里已有的行
	if record.User.UID == 0 {
		record.User.UID, _ = strconv.Atoi(luoguUser.Uid)
		return false, false
	}
	return false, true
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
		return 0, fmt.Errorf("%s的洛谷UID %q 不是合法的数字，请把 user.luogu_uid 改成洛谷数字 UID %s", luoguUser.RealName, luoguUser.Uid, err.Error())
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
