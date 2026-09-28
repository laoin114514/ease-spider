//go:build tools

// Package tools 只用于把内化的依赖固定在 go.mod 里，不参与正常构建。
package tools

// codeforcesAPIClient 以 git submodule 的形式内化在 pkg/codeforcesAPIClient，
// go.mod 用 replace 把 github.com/laoin114514/codeforcesClient 重定向到该目录。
//
// 这个文件带 tools 构建标签：go mod tidy 会把它计入依赖图（因此 require/replace
// 不会被误删），而 go build / go vet 默认不编译它，客户端不会进入构建产物。
// 业务代码直接 import "github.com/laoin114514/codeforcesClient" 即可。
import _ "github.com/laoin114514/codeforcesClient"
