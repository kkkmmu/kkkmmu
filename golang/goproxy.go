package main

import (
	"net/http"
	"os"

	"github.com/goproxy/goproxy"
)

func main() {
	g := goproxy.New()
	g.GoBinEnv = append(
		os.Environ(),
		"GOPROXY=https://goproxy.cn,direct", // 使用 goproxy.cn 作为上游代理
		"GOPRIVATE=git.example.com",         // 解决私有模块的拉取问题（比如可以配置成公司内部的代码源）
	)
	http.ListenAndServe("localhost:8080", g)
}
