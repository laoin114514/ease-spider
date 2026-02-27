package crawler

import (
	"database/sql"
	"spider/config"
	"time"
)

type Crawler struct {
	Name     string
	Interval time.Duration
	Enable   bool
}
type Context struct {
	DB     *sql.DB
	Config *config.Config
}
type Task func(*Context)
