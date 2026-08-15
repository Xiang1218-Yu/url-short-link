# 修复前故障复现（Docker）

## 项目与标准命令
这是一个从本地 JSON 目录读取短链接记录、按短码解析并输出链接信息的命令行工具。在仓库根目录可执行：

```sh
go build ./...
go run ./cmd/shortlink --input examples/links.json --code docs --at 2026-08-15T12:00:00Z
go test ./...
```

## 环境构建与编译
已在修复前的基线状态分别执行以下命令；linux/amd64 与 linux/arm64 的镜像构建和容器内 `go build ./...` 均成功：

```sh
docker buildx build --load --platform linux/amd64 -f benzhi.Dockerfile -t url-short-link-bug001-base-amd64:delivery .
docker run --rm --platform linux/amd64 --entrypoint go url-short-link-bug001-base-amd64:delivery build ./...
docker buildx build --load --platform linux/arm64 -f benzhi.Dockerfile -t url-short-link-bug001-base-arm64:delivery .
docker run --rm --platform linux/arm64 --entrypoint go url-short-link-bug001-base-arm64:delivery build ./...
```

## 故障触发步骤
在修复前的基线仓库根目录，完成上述 linux/amd64 镜像构建后执行：

```sh
docker run --rm --platform linux/amd64 --entrypoint go url-short-link-bug001-base-amd64:delivery test ./...
```

## 实际错误输出

```text
--- FAIL: TestRunRejectsCatalogLinkWithoutOwner (0.00s)
    main_test.go:40: catalog link without an owner was accepted
FAIL
FAIL	url-short-link/cmd/shortlink	0.029s
?   	url-short-link/internal/domain	[no test files]
ok  	url-short-link/internal/service	0.027s
ok  	url-short-link/internal/store	0.022s
?   	url-short-link/internal/transport	[no test files]
FAIL
test_exit_status=1
```

## 期望行为
当目录中的短链接负责人字段为空或仅包含空白字符时，加载或解析应返回可见错误，不能将该短码作为有效链接输出；负责人填写完整的目录仍应可以正常解析。
