# 修复前故障复现（Docker）

## 项目与标准命令
本项目是一个从本地 JSON 目录加载短链接并按短码解析或按负责人列出链接的 Go 命令行程序。在仓库根目录可执行以下标准命令：

```sh
go build ./...
go run ./cmd/shortlink --input examples/links.json --code docs --at 2026-08-15T12:00:00Z
go test ./...
```

## 环境构建与编译
在修复前源码状态下，以下命令均已实际执行。linux/amd64 与 linux/arm64 的镜像构建以及容器内 `go build ./...` 均成功：

```sh
docker buildx build --load --platform linux/amd64 -f benzhi.Dockerfile -t url-short-link:bug003-base-linux-amd64 .
docker run --rm --platform linux/amd64 --entrypoint go url-short-link:bug003-base-linux-amd64 build ./...
docker buildx build --load --platform linux/arm64 -f benzhi.Dockerfile -t url-short-link:bug003-base-linux-arm64 .
docker run --rm --platform linux/arm64 --entrypoint go url-short-link:bug003-base-linux-arm64 build ./...
```

## 故障触发步骤
在修复前源码的仓库根目录构建 linux/arm64 镜像后，执行以下命令：

```sh
docker buildx build --load --platform linux/arm64 -f benzhi.Dockerfile -t url-short-link:bug003-base-linux-arm64 .
docker run --rm --platform linux/arm64 --entrypoint go url-short-link:bug003-base-linux-arm64 test ./...
```

## 实际错误输出

```text
--- FAIL: TestRunReturnsCancellationBeforeReadingInput (0.00s)
    main_test.go:39: error=open catalog: open missing-catalog.json: no such file or directory, want context cancellation
FAIL
FAIL	url-short-link/cmd/shortlink	0.049s
?   	url-short-link/internal/domain	[no test files]
--- FAIL: TestResolveHonorsCancelledContext (0.00s)
    resolver_test.go:53: error=<nil>, want canceled
FAIL
FAIL	url-short-link/internal/service	0.050s
ok  	url-short-link/internal/store	0.051s
?   	url-short-link/internal/transport	[no test files]
FAIL
```

## 期望行为
当调用方在读取本地目录之前或短码解析过程中取消请求时，命令应立即返回取消错误，不应继续读取目录、返回正常解析结果或增加短链接访问次数。
