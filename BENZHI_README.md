# URL 短链接本地目录

本项目提供一个不依赖外部网络或数据库的 URL 短链接目录：从本地 JSON 读取链接规则，解析短码并记录本次进程内的访问次数。它也可以按所属人列出可管理的链接。

在仓库根目录执行：

```sh
go build ./...
go run ./cmd/shortlink --input examples/links.json --code docs --at 2026-08-15T12:00:00Z
go test ./...
```

项目固定使用 `go.mod` 中的 Go 1.26.5 和 `toolchain go1.26.5`。`benzhi.Dockerfile` 使用固定的 `golang:1.26.5-alpine` 基础镜像，并设置 `GOTOOLCHAIN=local`；容器始终从源码执行依赖下载和 `go build ./...`，不会复制宿主机编译产物。

## 双架构标准验收

在具备 Docker Buildx 的环境中，以下命令分别构建 `linux/amd64` 与 `linux/arm64` 镜像；脚本随后启动容器，先在容器内运行 `go build ./...`，再执行实际短码解析入口：

```sh
./build_benzhi_docker.sh linux/amd64
./build_benzhi_docker.sh linux/arm64
```

也可以逐步手工执行。以下示例为 amd64；将两个 `linux/amd64` 替换为 `linux/arm64` 即可验收另一架构：

```sh
docker buildx build --load --platform linux/amd64 -f benzhi.Dockerfile -t url-short-link:benzhi .
docker run --rm --platform linux/amd64 --entrypoint go url-short-link:benzhi build ./...
docker run --rm --platform linux/amd64 url-short-link:benzhi --input examples/links.json --code docs --at 2026-08-15T12:00:00Z
docker run --rm --platform linux/amd64 --entrypoint go url-short-link:benzhi test ./...
```

镜像构建、容器内编译和测试均以退出码 0 为通过。短码解析命令成功时会输出 `code`、`target`、`owner`、`tags` 与 `visits`，其中 `target` 应为 JSON 目录中该短码对应的目标地址；无效或过期短码会以非零退出码结束。
