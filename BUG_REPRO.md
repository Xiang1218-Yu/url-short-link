# 修复前故障复现（Docker）

## 项目与标准命令

本项目是从本地 JSON 目录解析 URL 短码并记录访问次数的 Go 命令行程序。在仓库根目录可执行：

```sh
go build ./...
go run ./cmd/shortlink --input examples/links.json --code docs --at 2026-08-15T12:00:00Z
go test ./...
```

## 环境构建与编译

已实际执行下列命令分别构建 linux/amd64 与 linux/arm64 镜像，并在相应容器中执行 `go build ./...`；两种平台的镜像构建和容器内编译均成功。

```sh
docker buildx build --load --platform linux/amd64 -f benzhi.Dockerfile -t url-short-link:benzhi .
docker run --rm --platform linux/amd64 --entrypoint go url-short-link:benzhi build ./...
docker buildx build --load --platform linux/arm64 -f benzhi.Dockerfile -t url-short-link:benzhi .
docker run --rm --platform linux/arm64 --entrypoint go url-short-link:benzhi build ./...
```

## 故障触发步骤

在修复前的 base 状态仓库根目录，先构建 amd64 镜像，再在容器中运行完整测试：

```sh
docker buildx build --load --platform linux/amd64 -f benzhi.Dockerfile -t url-short-link:benzhi .
docker run --rm --platform linux/amd64 --entrypoint go url-short-link:benzhi test ./...
```

## 实际错误输出

```text
ok  	url-short-link/cmd/shortlink	0.030s
?   	url-short-link/internal/domain	[no test files]
--- FAIL: TestResolveCountsOnlyActiveLinkVisits (0.00s)
    resolver_test.go:36: error=<nil>, want expired
FAIL
FAIL	url-short-link/internal/service	0.027s
--- FAIL: TestCancelledRecordDoesNotBlockLaterVisit (0.15s)
    memory_catalog_test.go:59: a canceled record blocked the next visit
FAIL
FAIL	url-short-link/internal/store	0.179s
?   	url-short-link/internal/transport	[no test files]
FAIL
```

## 期望行为

取消一次访问后，后续正常短码访问应及时完成；解析时间等于或晚于短码过期时间时，命令应以非零状态返回，且不输出成功解析结果或增加该短码访问次数。未过期的正常短码仍应输出其目标地址和访问结果。
