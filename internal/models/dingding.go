package models

// 钉钉打卡数据响应结构体
type DingdingCheckUpData struct {
	Errcode int64                 `json:"errcode"`
	Errmsg  string                `json:"errmsg"`
	Records []DingdingCheckRecord `json:"recordresult"`
}

// 钉钉打卡记录结构体
type DingdingCheckRecord struct {
	BaseCheckTime  int64  `json:"baseCheckTime"`  // 基准打卡时间
	BaseMacAddr    string `json:"baseMacAddr"`    // 基准MAC地址
	BizId          string `json:"bizId"`          // 业务ID
	CheckType      string `json:"checkType"`      // 打卡类型：OnDuty(上班) OffDuty(下班)
	CorpId         string `json:"corpId"`         // 企业ID
	DeviceSN       string `json:"deviceSN"`       // 设备序列号
	GmtCreate      int64  `json:"gmtCreate"`      // 创建时间
	GmtModified    int64  `json:"gmtModified"`    // 修改时间
	GroupId        int64  `json:"groupId"`        // 组ID
	Id             int64  `json:"id"`             // 记录ID
	IsLegal        string `json:"isLegal"`        // 是否合法：Y/N
	LocationMethod string `json:"locationMethod"` // 定位方式：ATM
	LocationResult string `json:"locationResult"` // 定位结果：Normal
	PlanCheckTime  int64  `json:"planCheckTime"`  // 计划打卡时间
	PlanId         int64  `json:"planId"`         // 计划ID
	SourceType     string `json:"sourceType"`     // 来源类型：ATM
	TimeResult     string `json:"timeResult"`     // 时间结果：Normal
	UserAddress    string `json:"userAddress"`    // 用户地址
	UserCheckTime  int64  `json:"userCheckTime"`  // 用户打卡时间
	UserId         string `json:"userId"`         // 用户ID
	WorkDate       int64  `json:"workDate"`       // 工作日期
}

//钉钉token响应结构体
type DingdingTokenResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
}
