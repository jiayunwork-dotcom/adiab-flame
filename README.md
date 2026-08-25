# adiab-flame — 定压绝热火焰温度核算

adiab-flame 是一个命令行核算工具：给定燃料、当量比和进气温度（JSON 文件），它按元素守恒（C/H/O/N 进出相等）与焓守恒（产物总焓 = 反应物总焓）迭代求解定压绝热火焰温度，并输出产物摩尔分数与原子残差。能力边界为气体燃料与空气的完全燃烧产物集合：内置燃料 CH4、C2H4、C2H6，空气按 21% O2 / 79% N2（摩尔）处理，产物集合为 CO2、H2O、O2、N2（富燃时含未燃燃料）；可选开启 CO/H2 离解平衡细化。它不做多组分扩散火焰、不做污染物 NOx 化学、不做火焰传播速度，也非危化安全审批系统。

## 用法

```text
go run . tad example/ch4-stoich.json
```

打印化学计量甲烷/空气（φ=1、298.15 K 进气）的绝热火焰温度、产物摩尔分数与原子残差。该算例的 Tad 约为 2329 K（开启离解细化后约 2246 K）。

其他子命令：

```text
go run . tad example/ch4-lean.json            # 贫燃 φ=0.8，Tad 下降
go run . tad example/ch4-rich.json            # 富燃 φ=1.2，产物含未燃 CH4
go run . tad example/ch4-hot-inlet.json       # 预热进气 600 K，Tad 上升但升幅小于进气升幅
go run . tad example/ch4-dissoc.json          # 开启 CO/H2 离解平衡细化
go run . tad example/c2h6-stoich.json         # 乙烷 φ=1
go run . tad example/c2h4-stoich.json         # 乙烯 φ=1
go run . tad <file> --dissoc                  # 命令行开关，等价于 JSON 里 dissociation: true
go run . checks CH4                           # 交叉规则自检（φ 排序、进气趋势、守恒）
go run . fuels                                # 列出支持燃料与生成焓
go run . help
```

配置 JSON 字段：

```json
{
  "fuel": "CH4",
  "equivalence_ratio": 1.0,
  "inlet_temperature_k": 298.15,
  "dissociation": false
}
```

`fuel` 必填；`equivalence_ratio` 与 `inlet_temperature_k` 必填且必须为正；`pressure_atm`（默认 1.0）与 `dissociation`（默认 false）可省略。未知 JSON 字段会被拒绝。

## 关键约定

- **摩尔基准**：全部组分按「每摩尔燃料」表述，空气需求由燃料分子式（CₐHᵦ 需 a+b/4 摩尔 O₂）一次导出，燃料与空气比不另设第二套基准。
- **焓模型**：h(T) = Δhf°(298.15 K) + ∫_{298.15}^T cp dτ，cp 为 NASA 多项式（系数表钉死，低温 200–1000 K、高温 1000–3500 K 两段），Δhf° 与 cp 多项式取自同一张系数表。焓守恒方程单调，用二分迭代，迭代上限 200。
- **产物组成**：φ≤1 时燃料完全燃烧，多余 O2 留出；φ>1 时供氧不足，仅 1/φ 的燃料完全燃烧，其余以燃料蒸气留在产物中。任意 φ 下 C/H/O/N 原子进出严格相等。
- **离解细化（可选）**：`dissociation: true` 时在同一温度迭代里求解 CO2⇌CO+½O2 与 H2O⇌H2+½O2 两个平衡，平衡常数由同一张 NASA 表经 Gibbs 能导出。富燃含未燃燃料时不适用该细化。
- **交叉规则**：同燃料同进气下 φ=1 的 Tad 高于 φ=0.8 与 φ=1.2；只提高进气温度，Tad 上升但升幅小于进气升幅；原子残差低于 1e-9；N2 不参与氧化，其摩尔数等于进气 N2。`checks` 子命令对全部内置燃料逐条验证。
- **非法输入**：φ≤0、进气温度≤0、未知燃料、未知 JSON 字段、缺失文件、温度迭代不收敛——一律 stderr 明文报错并非零退出。

## 构建与测试

```text
go build ./...
go test ./...
```

纯标准库，无第三方依赖，无 cgo。`go test -run TestMethaneStoichTadBand ./internal/flame` 可单独跑甲烷化学计量温度带测试。

## 许可

MIT，见 [LICENSE](./LICENSE)。
