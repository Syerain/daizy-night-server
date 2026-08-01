# Conf: Go Production-Ready Configuration

[![Go Report Card](https://goreportcard.com/badge/github.com/oy3o/conf)](https://goreportcard.com/report/github.com/oy3o/conf)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`conf` 是一个为 Go 语言打造的**生产级**配置管理库。它集成了 `Viper` 的强大加载能力和 `Validator` 的精准校验，并针对**微服务**、**云原生**和**高频 API** 场景进行了深度优化。

它的核心设计哲学是：**默认即最佳（Convention over Configuration）** 与 **极致性能（Performance at Scale）**。

## 核心特性 (Features)

*   **泛型支持 (Generics)**: 使用 `Load[T]()` 直接返回强类型结构体，告别繁琐的类型断言和变量声明。
*   **生产环境安全 (Strict Env)**: 独创 `env:"strict"` 标签。在生产环境 (`GO_ENV=production`) 下，强制检查敏感字段（如密码）是否直接来自环境变量，防止密钥意外写入配置文件。
*   **混合验证模式 (Hybrid Validation)**:
    *   **反射模式**: 使用 Tag 开发，简单快捷。
    *   **极速模式**: 实现 `SelfValidatable` 接口，**性能提升 40x - 100x**，专为热点路径设计。
*   **国际化校验 (I18n)**: 内置中文 (`zh`) 和英文 (`en`) 错误提示，自动解析 `mapstructure` 标签，报错信息准确友好。
*   **智能默认值**: 支持 `default` 标签设置默认值。
*   **严格解码**: 防止配置文件中出现未定义的字段（避免拼写错误被忽略）。

## 安装 (Installation)

```bash
go get github.com/oy3o/conf
```

## 快速开始 (Quick Start)

### 1. 定义配置结构体

```go
package main

import (
    "fmt"
    "github.com/oy3o/conf"
)

type DBConfig struct {
    // 支持 mapstructure 映射，validate 校验，default 默认值
    Host string `mapstructure:"host" validate:"required"`
    Port int    `mapstructure:"port" validate:"min=1024" default:"3306"`
    
    // 【核心特性】生产环境强制要求从环境变量获取，不允许写在文件里
    Password string `mapstructure:"password" validate:"required" env:"strict"`
}

type Config struct {
    AppName string   `mapstructure:"app_name" default:"MyApp"`
    Debug   bool     `mapstructure:"debug"`
    DB      DBConfig `mapstructure:"db"`
}

func main() {
    // 一行代码加载 + 解析 + 验证 + 默认值
    // 默认搜索当前目录下的 config.yaml
    cfg := conf.MustLoad[Config]("myapp", 
        conf.WithLocale("zh"), // 开启中文报错
    )

    fmt.Printf("App: %s, DB Port: %d\n", cfg.AppName, cfg.DB.Port)
}
```

### 2. 配置文件 (config.yaml)

```yaml
app_name: "SuperService"
debug: true
db:
  host: "127.0.0.1"
  # port 使用默认值 3306
  # password 不写在这里，通过环境变量传递
```

### 3. 运行

**开发环境 (GO_ENV=dev):**
```bash
# 即使没有密码，开发环境也会跳过 strict 检查
export GO_ENV=dev
go run main.go
```

**生产环境 (GO_ENV=production):**
```bash
export GO_ENV=production
export MYAPP_DB_PASSWORD="secret_password" # 必须设置，否则启动报错
go run main.go
```

---

## 高级特性

### 1. 混合验证模式 (Performance)

对于 QPS 极高的场景，反射带来的开销虽然很小，但依然存在。`conf` 支持通过实现接口来**绕过反射**。

```go
type FastRequest struct {
    APIKey string
}

// 实现 SelfValidatable 接口
// Conf 会自动检测并直接调用此方法，耗时仅需 ~4.5ns
func (r *FastRequest) Validate() error {
    if len(r.APIKey) != 32 {
        return fmt.Errorf("invalid api_key")
    }
    return nil
}
```

### 2. 生产环境强制 Env 检查

在 `GO_ENV=production` 或 `prod` 时，标记了 `env:"strict"` 的字段**必须**存在于系统环境变量中。Viper 从配置文件读取的值将被视为无效。这从代码层面杜绝了"误将生产密码提交到 Git 仓库"的风险。

## 配置选项 (Options)

加载配置时支持以下 Option：

| Option | 说明 | 默认值 |
| :--- | :--- | :--- |
| `WithSearchPaths(paths...)` | 配置文件搜索路径 | `.` 和 `./config` |
| `WithFileName(name)` | 配置文件名 | `config` |
| `WithFileType(type)` | 文件类型 (yaml, json, toml...) | `yaml` |
| `WithLocale(lang)` | 验证错误语言 (`zh`, `en`, `""`) | `zh` |

## 性能基准测试 (Benchmarks)

基于 Intel i9-13900HX 的测试数据：

| 验证模式 | 耗时 (ns/op) | 内存分配 (B/op) | 性能提升 |
| :--- | :--- | :--- | :--- |
| **Reflect (Tag)** | 152.3 ns | 24 B | 基准 |
| **Interface (Fast)** | **4.5 ns** | **0 B** | **~33x** |
| **Direct Call** | 0.45 ns | 0 B | 理论极限 |

> **建议**: 配置加载场景使用 **Tag 模式**（开发效率高）；在极度敏感的热点代码中使用 **Interface/原生 模式**; 。
