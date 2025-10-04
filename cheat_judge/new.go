package cheat_judge

import "spider/config/db"

// NewCheatingDetector 创建新的检测器实例
// 创建检测器
func NewCheatingDetector() *CheatingDetector {
	return &CheatingDetector{
		db: db.Pool,
	}
}
