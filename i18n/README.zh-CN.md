# Daizy Night

[English](README.md) | [中文](README.zh-CN.md)

**Daizy Night** 是面向 ATOM Reforge 社区的服务。

项目名 _Daizy Night_ 是 _The Hazy Night_ 的谐音。

> 说明：服务端仍处于早期阶段（`v0.5.4`），API 仍在演进中。

## 功能特性

- **注册码注册** — 每次注册都需要 *Registercode*（注册码）：由 Ed25519 签名、带过期校验，首次使用会被**原子地**占用（一次性），且每次占用均留痕供审计。
- **账号密码登录（legacy）** — 密码使用 **argon2id** 哈希处理，绝不存储明文。
- **双 JWT** — 登录同时签发 *access token* 与 *refresh token*，均使用 Ed25519（EdDSA）签名；refresh token 持久化保存在服务端，支持吊销。
- **认证与授权** — `echo-jwt` 中间件负责认证，另有角色控制（`user` / `admin`）。
- **速率限制** — 按客户端 IP 维度的内存令牌桶；`rate`、`burst`、`expiresIn` 均可配置。
- **GitHub OAuth（预留）** — OAuth 注册/登录方式已在常量层定义，但尚未实现（当前返回 `FeatureUnsupported`）。

## 技术栈

- [Go](https://go.dev) 1.26
- [Echo](https://echo.labstack.com) v5 — HTTP 框架
- [GORM](https://gorm.io) + SQLite — ORM 与存储
- [confx](https://github.com/atomreforge/confx)（基于 viper）— 配置加载
- [argon2id](https://github.com/alexedwards/argon2id) — 密码哈希
- [golang-jwt v5](https://github.com/golang-jwt/jwt) + [echo-jwt v5](https://github.com/labstack/echo-jwt) — JWT 签名与校验
- `crypto/ed25519` — JWT 与注册码使用的密钥对

## 相关仓库

- **服务端**（本仓库）：<https://github.com/atomreforge/daizy-night-server>
- **计划客户端**：<https://github.com/atomreforge/daizy-night-app>
- **测试用 CLI 客户端**：<https://github.com/Syerain/dnappcli>

## 快速开始

### 1. 获取代码

```bash
git clone https://github.com/atomreforge/daizy-night-server.git
cd daizy-night-server
```

### 2. 配置

```bash
cp config.example.yaml config.yaml
```

`security` 段整体为**必填** —— 任一密钥缺失都会导致启动失败。所有值均为十六进制编码的 Ed25519 密钥：

| 键 | 用途 |
| --- | --- |
| `registercodeEnckey` / `registercodeDeckey` | 注册码签名 / 校验 |
| `passwordEnckey` / `passwordDeckey` | 密码层密钥（预留） |
| `jwtAccessTokenEnckey` / `jwtAccessTokenDeckey` | access token 签名 / 校验 |
| `jwtRefreshTokenEnckey` / `jwtRefreshTokenDeckey` | refresh token 签名 / 校验 |

可用以下 Go 代码生成 Ed25519 密钥对：

```go
package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

func main() {
	pub, priv, _ := ed25519.GenerateKey(nil)
	fmt.Println("private (seed):", hex.EncodeToString(priv.Seed()))
	fmt.Println("private (full):", hex.EncodeToString(priv))
	fmt.Println("public:         ", hex.EncodeToString(pub))
}
```

> **测试 profile** — 若工作目录或可执行文件旁存在 `config.test.yaml`，其优先级高于 `config.yaml`，适合本地开发。

### 3. 运行

```bash
go mod tidy
go run ./cmd/main.go
```

默认监听 `127.0.0.1:4703`。SQLite 数据库会在首次启动时自动创建并迁移。

## 配置参考

| 配置段 | 字段 | 默认值 | 必填 | 说明 |
| --- | --- | --- | --- | --- |
| `main` | `isDebugMode` | `true` | 否 | 调试日志与行为 |
| `http` | `port` | `4703` | 否 | 监听端口 |
| `http` | `address` | `127.0.0.1` | 否 | 绑定地址 |
| `http.rateLimit` | `enabled` | `true` | 否 | 是否启用速率限制 |
| `http.rateLimit` | `rate` | `10` | 否 | 每秒补充的令牌数 |
| `http.rateLimit` | `burst` | `30` | 否 | 桶容量 |
| `http.rateLimit` | `expiresIn` | `3m` | 否 | 条目过期时间 |
| `database` | `isDebugMode` | `false` | 否 | GORM 详细日志 |
| `database` | `DSN` | `./data.db` | **是** | SQLite 文件路径 |
| `log` | `isColored` | `false` | 否 | 彩色日志输出 |
| `security` | `*Enckey` / `*Deckey` | — | **是** | 十六进制 Ed25519 密钥对（见上文） |
| `security` | `jwtAccessTokenExpireTime` | `15m` | **是** | access token 有效期 |
| `security` | `jwtRefreshTokenExpireTime` | `168h` | **是** | refresh token 有效期 |
| `security` | `jwtRevokedTokensRetainTime` | `72h` | **是** | 已吊销 token 保留时长 |

支持的时间单位：`h`、`m`、`s`、`ms`、`us`、`ns`。

## API

所有路由均以 `/api/v1` 为前缀。

### `POST /api/v1/register`

公开接口。注册账号，需要有效的注册码。

```json
{
  "registerway": "legacy",
  "username": "alice",
  "nickname": "Alice",
  "password": "secret",
  "registercode": "<payload_hex>.<sig_hex>"
}
```

`200 OK`：

```json
{ "message": "ok" }
```

### `POST /api/v1/login`

公开接口。凭账号密码换取 JWT 对。

```json
{
  "loginway": "legacy",
  "username": "alice",
  "password": "secret"
}
```

`200 OK`：

```json
{
  "access_token": "<jwt>",
  "refresh_token": "<jwt>"
}
```

### `GET /api/v1/user/{username}/me`

需认证 —— 请求头 `Authorization: Bearer <access_token>`。返回当前用户信息。路径中的 `{username}` 必须与认证身份（JWT claims）一致，否则 `403`。

`200 OK`：

```json
{
  "uid": 12345,
  "username": "alice",
  "nickname": "Alice",
  "email": "",
  "register_time": "2026-08-17T12:00:00Z",
  "role": "user",
  "github_id": null,
  "github_login": null
}
```

### `POST /api/v1/admin/sudo`

需认证且角色为 `admin`。占位的管理员操作。

## 项目结构

```
cmd/                     入口与服务组装
internal/
  abstract/interface/    各 provider 契约（crypto、db、repo、service）
  api/v1/                请求 / 响应 DTO
  config/                配置结构体与加载器（confx）
  consts/                角色、注册/登录方式、日志表达式
  crypto/                Ed25519、JWT 与注册码 provider
  dbware/                GORM 连接与仓储
  errs/                  统一错误体系
  handler/               HTTP handler（薄层）
  middleware/            injector、速率限制、JWT 认证、角色控制
  model/                 GORM 模型与共享 DTO
  router/                Echo 路由装配
  service/               业务逻辑
  utils/                 hash、hex、logger、rand、trace
etc/                     内部笔记与文档
```

分层调用链：`router` → `handler` → `service` → `dbware` → `model`；`middleware`、`crypto`、`errs`、`utils` 为横切包。

## 更新日志

- **v0.5.4** — 随机 `UserID` 身份；注册码原子认领；适应性收尾
- **v0.5.3** — 注册码逻辑修复
- **v0.5.2** — 配置依赖注入修复
- **v0.5.0** — 异常体系；trace 体系
- **v0.4.1** — 认证中间件与 `echo-jwt`
- **v0.4.0** — 迁移至 Echo v5；错误处理重构
- **v0.3.x** — 登录流程
- **v0.1.0** — 初始结构 