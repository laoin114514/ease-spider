package models

// LuoguTeamResponse 洛谷团队响应
type LuoguTeamResponse struct {
	Members          LuoguTeamMembers      `json:"members"`
	Groups           []interface{}         `json:"groups"` // 空数组，暂时用interface{}
	GroupMemberCount LuoguGroupMemberCount `json:"groupMemberCount"`
}

// LuoguTeamMembers 团队成员信息
type LuoguTeamMembers struct {
	Result  []LuoguTeamMember `json:"result"`
	PerPage int               `json:"perPage"`
	Count   int               `json:"count"`
}

// LuoguTeamMember 单个团队成员
type LuoguTeamMember struct {
	Group      *interface{} `json:"group"`      // 通常为null
	User       LuoguUser    `json:"user"`       // 复用现有的LuoguUser结构体
	Type       int          `json:"type"`       // 成员类型
	Permission int          `json:"permission"` // 权限
	RealName   string       `json:"realName"`   // 真实姓名
}

// LuoguGroupMemberCount 组成员统计
type LuoguGroupMemberCount struct {
	Null int `json:"null"` // 无分组成员数量
}
