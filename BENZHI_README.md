Go 实现的定压绝热火焰温度计算服务，默认启动 HTTP 在 :8080 提供单点求解与当量比扫描接口，也可通过 tad/sweep/checks 子命令在命令行对 JSON 算例做离线计算。

## 构建与启动

```bash
go build -o adiab-flame .
./adiab-flame                              # 启动 HTTP 服务 :8080
./adiab-flame tad example/ch4-stoich.json  # 命令行计算
```

## 评测镜像

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh adiab-flame
```
