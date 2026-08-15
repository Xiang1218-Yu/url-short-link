# 修复前故障复现（Docker）

## 项目与标准命令

本项目从本地 JSON 目录读取短链接规则，支持解析短码、记录访问次数，以及按所属人列出链接。在仓库根目录执行以下标准命令：

```sh
go build ./...
go run ./cmd/shortlink --input examples/links.json --code docs --at 2026-08-15T12:00:00Z
go test ./...
```

其中 `go build ./...` 可以成功；包含未填写 `expires_at` 的长期链接时，运行入口和测试都会触发下述故障。

## 环境构建与编译

实际执行的镜像构建和容器内编译命令如下：

```sh
docker buildx build --load --platform linux/amd64 -f benzhi.Dockerfile -t url-short-link:bug005-base-amd64 .
docker run --rm --platform linux/amd64 --entrypoint go url-short-link:bug005-base-amd64 build ./...
docker buildx build --load --platform linux/arm64 -f benzhi.Dockerfile -t url-short-link:bug005-base-arm64 .
docker run --rm --platform linux/arm64 --entrypoint go url-short-link:bug005-base-arm64 build ./...
```

linux/amd64 与 linux/arm64 的镜像构建和容器内 `go build ./...` 均成功；目标故障发生在包含长期链接的目录被加载后执行测试或解析时。

## 故障触发步骤

在仓库根目录构建基线镜像后，执行：

```sh
docker buildx build --load --platform linux/arm64 -f benzhi.Dockerfile -t url-short-link:bug005-base-arm64 .
docker run --rm --platform linux/arm64 --entrypoint go url-short-link:bug005-base-arm64 test ./...
```

该命令稳定以退出码 1 失败。

## 实际错误输出

```text
--- FAIL: TestRunResolvesCatalogLink (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x156368]

goroutine 8 [running]:
testing.tRunner.func1.2({0x190160, 0x329190})
	/usr/local/go/src/testing/testing.go:1974 +0x1a0
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1977 +0x318
panic({0x190160?, 0x329190?})
	/usr/local/go/src/runtime/panic.go:860 +0x12c
url-short-link/internal/domain.Link.Clone(...)
	/workspace/internal/domain/link.go:51
url-short-link/internal/store.NewMemoryCatalog({0x7261f8872000, 0x3, 0x1c75a2?})
	/workspace/internal/store/memory_catalog.go:31 +0x1a8
url-short-link/cmd/shortlink.run({0x1d4568, 0x35ab20}, {0x7261f8822240, 0x6, 0x6}, {0x1d38c0, 0x7261f8830690})
	/workspace/cmd/shortlink/main.go:41 +0x220
url-short-link/cmd/shortlink.TestRunResolvesCatalogLink(0x7261f885c248)
	/workspace/cmd/shortlink/main_test.go:12 +0xbc
testing.tRunner(0x7261f885c248, 0x1d0ab0)
	/usr/local/go/src/testing/testing.go:2036 +0xc4
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:2101 +0x3a8
FAIL	url-short-link/cmd/shortlink	0.003s
?   	url-short-link/internal/domain	[no test files]
--- FAIL: TestResolveCountsOnlyActiveLinkVisits (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x12da68]

goroutine 20 [running]:
testing.tRunner.func1.2({0x15e2c0, 0x2d8e40})
	/usr/local/go/src/testing/testing.go:1974 +0x1a0
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1977 +0x318
panic({0x15e2c0?, 0x2d8e40?})
	/usr/local/go/src/runtime/panic.go:860 +0x12c
url-short-link/internal/domain.Link.Clone(...)
	/workspace/internal/domain/link.go:51
url-short-link/internal/store.NewMemoryCatalog({0x143acfea3ea0, 0x2, 0x0?})
	/workspace/internal/store/memory_catalog.go:31 +0x1a8
url-short-link/internal/service.newCatalog(0x143acfee8248, {0x143acfea3ea0, 0x2, 0x2})
	/workspace/internal/service/resolver_test.go:15 +0x40
url-short-link/internal/service.TestResolveCountsOnlyActiveLinkVisits(0x143acfee8248)
	/workspace/internal/service/resolver_test.go:24 +0x124
testing.tRunner(0x143acfee8248, 0x19a970)
	/usr/local/go/src/testing/testing.go:2036 +0xc4
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:2101 +0x3a8
FAIL	url-short-link/internal/service	0.003s
--- FAIL: TestMemoryCatalogReturnsIndependentSnapshots (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x12d628]

goroutine 20 [running]:
testing.tRunner.func1.2({0x15dfe0, 0x2d8e80})
	/usr/local/go/src/testing/testing.go:1974 +0x1a0
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1977 +0x318
panic({0x15dfe0?, 0x2d8e80?})
	/usr/local/go/src/runtime/panic.go:860 +0x12c
url-short-link/internal/domain.Link.Clone(...)
	/workspace/internal/domain/link.go:51
url-short-link/internal/store.NewMemoryCatalog({0x196f06aa3f00, 0x1, 0x2c01a0?})
	/workspace/internal/store/memory_catalog.go:31 +0x1a8
url-short-link/internal/store.TestMemoryCatalogReturnsIndependentSnapshots(0x196f06ae8248)
	/workspace/internal/store/memory_catalog_test.go:12 +0xa4
testing.tRunner(0x196f06ae8248, 0x19a498)
	/usr/local/go/src/testing/testing.go:2036 +0xc4
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:2101 +0x3a8
FAIL	url-short-link/internal/store	0.003s
?   	url-short-link/internal/transport	[no test files]
FAIL
```

退出结果：`exit status 1`。

## 期望行为

未填写 `expires_at` 的记录表示长期有效短链接。目录加载、查询、列表和解析这类链接时应正常完成且不发生 panic；解析长期链接应能够记录访问次数。已填写过期时间的链接仍应在当前时间达到或超过其过期时间时返回过期错误，并且不增加访问次数。
