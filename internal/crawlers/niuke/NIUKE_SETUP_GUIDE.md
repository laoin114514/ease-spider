# 牛客爬虫接入指南

## 概述

牛客爬虫已集成到 spider1 项目，可自动爬取牛客网用户提交记录。

## 需要创建的 MySQL 表

### 1. niuke_submissions 表（存储提交记录）

```sql
CREATE TABLE IF NOT EXISTS niuke_submissions (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增主键',
    sub_id VARCHAR(255) UNIQUE NOT NULL COMMENT '唯一提交ID：uid_problemId_submitTime',
    uid VARCHAR(50) NOT NULL COMMENT '牛客用户UID',
    user_name VARCHAR(100) NOT NULL COMMENT '用户真实姓名',
    problem_id VARCHAR(50) NOT NULL COMMENT '题目ID',
    problem_name VARCHAR(255) NOT NULL COMMENT '题目名称',
    submit_time DATE NOT NULL COMMENT '提交日期',
    status VARCHAR(50) NOT NULL COMMENT '提交状态（通过/失败等）',
    language VARCHAR(50) NOT NULL COMMENT '编程语言',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
    INDEX idx_uid (uid),
    INDEX idx_submit_time (submit_time),
    INDEX idx_problem_id (problem_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='牛客用户提交记录表';
```

### 2. user 表添加 niuke_uid 字段

```sql
-- 添加牛客UID字段
ALTER TABLE user ADD COLUMN niuke_uid VARCHAR(50) NULL COMMENT '牛客网UID' AFTER nowcoder_account;

-- 添加索引
ALTER TABLE user ADD INDEX idx_niuke (niuke_uid);
```

## 配置用户数据

在 `user` 表中配置需要爬取的用户的牛客 UID：

```sql
-- 查看现有用户
SELECT id, real_name, role_id FROM user WHERE role_id IN (1, 3);

-- 给用户添加牛客 UID（示例）
UPDATE user SET niuke_uid='12345678' WHERE id=1;
```


## 查看爬取结果

```sql
-- 查看所有牛客提交记录
SELECT * FROM niuke_submissions LIMIT 10;

-- 统计每个用户的提交数量
SELECT user_name, COUNT(*) as count FROM niuke_submissions GROUP BY user_name;

-- 查看某个用户的记录
SELECT * FROM niuke_submissions WHERE uid='12345678';
```

## 文件位置

牛客爬虫相关文件：

```
spider1/
├── internal/crawlers/niuke/get_user_records/
│   ├── base.go      # 爬虫入口
│   └── service.go   # 核心逻辑
├── internal/models/niuke.go
└── internal/repository/niuke.go
```

---

