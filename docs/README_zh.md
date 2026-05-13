# reqrisk

**中文文档 | [English](../README.md)**

`reqrisk` 是一个轻量、可解释的 Go 核心库，用来把已经提取好的请求特征转换成风险结果。

它不负责解析 `*http.Request`，也不负责采集请求体或生成原始特征；它只关注一件事：

`feature -> signal -> score -> explanation`

## 它能做什么

`reqrisk` 接收已经计算好的特征，例如：

- 熵值
- 结构复杂度
- 字符集异常
- 指纹稀有度或波动性

然后输出：

- 单信号分析结果
- 总分
- 风险等级
- 带证据的 findings
- 轻量 tags
- 建议动作

## 它不做什么

`reqrisk` 有意保持聚焦，它不会：

- 解析原始 `*http.Request`
- 采集或修改请求体
- 从原始 payload 直接计算熵、指纹或结构特征
- 充当 WAF
- 自动阻断流量

这让它既可以单独使用，也可以作为 `web-profiler` 这类上游特征提取工具的下游风险评估核心。

## 安装

```bash
go get github.com/gofurry/reqrisk
```

## 快速开始

一个可直接运行的示例放在 [`../example/main.go`](../example/main.go)。

```go
assessor := reqrisk.New()

result := assessor.Evaluate(reqrisk.Features{
    Entropy: &reqrisk.EntropyFeature{
        Value:      7.9,
        SampleSize: 1024,
        Source:     "body",
    },
    Complexity: &reqrisk.ComplexityFeature{
        Depth:          7,
        FieldCount:     42,
        MaxArrayLength: 18,
        Source:         "json",
    },
})
```

当前 benchmark 基线记录在 [`benchmark_baseline.md`](benchmark_baseline.md)。

## V1 内置信号

- `entropy`
- `complexity`
- `charset`
- `fingerprint`

每个信号也都支持单独分析：

```go
report := reqrisk.AssessEntropy(reqrisk.EntropyFeature{
    Value:      7.7,
    SampleSize: 900,
    Source:     "body",
})
```

## 风险等级

默认有四个等级：

- `low`
- `medium`
- `high`
- `critical`

建议动作包括：

- `observe`
- `review`
- `challenge`
- `block`

---

## Presets

`reqrisk` 提供了几种轻量预设，方便从不同的风险灵敏度开始：

- `PresetBalanced`：保持默认基线
- `PresetSensitive`：降低阈值，让边界信号更早出现
- `PresetConservative`：提高阈值，只让更强的信号参与

示例：

```go
assessor := reqrisk.New(reqrisk.WithPreset(reqrisk.PresetSensitive))
```

预设只是起点，后续你仍然可以继续叠加显式的策略覆盖。

## 仓库目录

- 根包：对外稳定的公共 API
- `example/`：可直接运行的示例程序
- `internal/model`：输入输出数据结构
- `internal/policy`：策略、配置和归一化逻辑
- `internal/core`：评估引擎和内置信号实现

对外导入路径保持在根包，内部实现按职责拆分，方便后续维护和扩展。

## License

[MIT](../LICENSE)
