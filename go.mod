module spider

go 1.26.1

require (
	github.com/go-resty/resty/v2 v2.16.5
	github.com/go-sql-driver/mysql v1.9.3
	github.com/joho/godotenv v1.5.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

// 内化（in-tree）的 Codeforces API 客户端：真实代码在 git submodule
// pkg/codeforcesAPIClient，用 replace 重定向到本地目录，不再从网络拉取。
require github.com/laoin114514/codeforcesClient v0.2.2

replace github.com/laoin114514/codeforcesClient => ./pkg/codeforcesAPIClient
