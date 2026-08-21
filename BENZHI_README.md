# adiab-flame — 定压绝热火焰温度核算

adiab-flame 是定压绝热火焰温度核算命令行工具：给定燃料、当量比和进气温度的 JSON 文件，按元素守恒与焓守恒迭代求解火焰温度，并打印产物摩尔分数与原子残差。纯标准库，无网络依赖，无 cgo。

## 构建 / 运行 / 测试

```text
go build ./...
go run . tad example/ch4-stoich.json   # CLI：核算甲烷 φ=1、298 K 进气并打印 Tad、产物摩尔分数、原子残差
go run . checks CH4                    # 交叉规则自检（φ 排序、进气趋势、守恒）
go test ./...                          # 单元测试（thermo / flame / cli）
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```

进容器后运行 `go build ./... && go test ./...`，再用 `go run . tad example/ch4-stoich.json` 验证 CLI 输出 Tad、产物摩尔分数与原子残差。
