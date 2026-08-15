# 修复前故障复现（Docker）

## 项目与标准命令
本项目是一个读取本地 JSON 目录并解析短链接的命令行工具。仓库根目录的标准命令如下：

```sh
go build ./...
go run ./cmd/shortlink --input examples/links.json --code docs --at 2026-08-15T12:00:00Z
go test ./...
```

`go build ./...` 与示例短码解析命令可以正常执行；本故障由目录快照泄漏测试稳定触发，`go test ./...` 在修复前预期失败。

## 环境构建与编译
已依次在两个平台执行以下命令：

```sh
docker buildx build --load --platform linux/amd64 -f benzhi.Dockerfile -t url-short-link:base-amd64 .
docker run --rm --platform linux/amd64 --entrypoint go url-short-link:base-amd64 build ./...
docker buildx build --load --platform linux/arm64 -f benzhi.Dockerfile -t url-short-link:base-arm64 .
docker run --rm --platform linux/arm64 --entrypoint go url-short-link:base-arm64 build ./...
```

linux/amd64 与 linux/arm64 的镜像构建、容器内 `go build ./...` 均成功；目标故障在下节的测试命令中触发。

## 故障触发步骤
在仓库根目录执行：

```sh
docker buildx build --load --platform linux/amd64 -f benzhi.Dockerfile -t url-short-link:base-amd64 .
docker run --rm --platform linux/amd64 --entrypoint go url-short-link:base-amd64 test ./...
```

调用方先取得短码 `docs` 的查询结果并改写其中的标签元素，再次解析相同短码时，目录中的标签也已被改写；对应测试稳定失败。

## 实际错误输出
```text
ok  	url-short-link/cmd/shortlink	0.055s
?   	url-short-link/internal/domain	[no test files]
ok  	url-short-link/internal/service	0.050s
--- FAIL: TestMemoryCatalogReturnsIndependentSnapshots (0.01s)
    memory_catalog_test.go:26: catalog leaked caller mutation: []string{"changed"}
FAIL
FAIL	url-short-link/internal/store	0.059s
?   	url-short-link/internal/transport	[no test files]
FAIL
```

## 期望行为
每次查询短链接都应向调用方返回独立的标签快照。调用方为展示而修改返回结果中的标签后，再次解析同一短码仍应保留目录原始标签，`TestMemoryCatalogReturnsIndependentSnapshots` 应通过。
