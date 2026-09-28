// Package cfclient 把 user 表里的 CF 凭据装配成内化的 Codeforces 客户端
// （github.com/laoin114514/codeforcesClient，源码在 pkg/codeforcesAPIClient）。
//
// 客户端自带限流与 429/503/5xx 重试（见其 internal/http/transport.go），
// 所以各插件不再需要自己拼签名 URL，也不用手写重试。
package cfclient

import (
	cf "github.com/laoin114514/codeforcesClient"

	"spider/internal/repository"
)

// NewSigned 构造带号池签名器的客户端：调用时用 client.WithHandle(账号) 指定身份。
//
// rps <= 0 表示不限流；一般传对应任务的并发数，避免突发请求被 CF 判成 429。
func NewSigned(keys map[string]repository.CfApiKey, rps int) *cf.Client {
	pool := make(map[string]struct{ ApiKey, Secret string }, len(keys))
	for account, key := range keys {
		pool[account] = struct{ ApiKey, Secret string }{ApiKey: key.ApiKey, Secret: key.Secret}
	}
	return cf.NewClient(cf.WithSigner(cf.NewPoolSigner(pool)), cf.WithRateLimit(rps))
}

// NewPlain 构造不带签名器的客户端，用于不需要凭据的公共接口。
//
// 带签名器的客户端会无条件给每个请求签名（PoolSigner 找不到 handle 就返回 ErrAuth），
// 因此“退回匿名请求”必须换一个没有签名器的客户端，而不是复用上面那个。
func NewPlain(rps int) *cf.Client {
	return cf.NewClient(cf.WithRateLimit(rps))
}
