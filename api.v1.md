# Daizy Night API v1

## 概览

- Base URL: `http://127.0.0.1:4703`，由 `config.yaml` 的 `http.address` / `http.port` 决定
- 所有路由前缀 `/api/v1`，请求与响应体均为 JSON（`Content-Type: application/json`）
- 认证方式：`Authorization: Bearer <access_token>`，JWT 签名算法 EdDSA（ed25519）
- Token 有效期：access token 15m，refresh token 168h（`config.yaml` 的 `security` 段可调）

## 通用行为

### 错误响应

所有错误响应统一为 `{"message": string}`：

| 场景 | 状态码 | message |
|---|---|---|
| 业务错误（校验失败、凭据错误、注册码无效等） | 400/401/404 | 可读描述，如 `failure in user login; details::...` |
| 参数解析失败、验签失败、DB 错误等非业务错误 | 视错误而定 | 固定回退为 `internal server error` |
| 未捕获的未知错误 | 500 | `internal server error` |

### 限流

按客户端 IP 全局限流，超限返回 429 `{"message": "too many requests"}`。
默认：速率 2 req/s，突发 7，窗口 3m（`config.yaml` 的 `http.rateLimit`）。

## 端点

### POST /api/v1/register

注册新用户。公开端点，无需认证。

请求体：

| 字段 | 类型 | 约束 |
|---|---|---|
| `registerway` | string | `legacy`（唯一支持，`oauth-github` 未实现） |
| `username` | string | 1-15 字符，仅允许 `a-zA-Z0-9` |
| `nickname` | string | 1-15 字符 |
| `password` | string | 6-128 字符 |
| `registercode` | string | 注册码原文，格式 `<hex载荷>.<hex签名>` |

```json
{
  "registerway": "legacy",
  "username": "alice",
  "nickname": "Alice",
  "password": "secret123",
  "registercode": "ab12....cd.ef34....89"
}
```

响应：

- `200` `{"message": "ok"}`
- `400` 参数校验失败，或注册码无效/过期（`failure in ...`）
- `500` 用户名或注册码重复（当前实现下表现为 internal server error）

### POST /api/v1/login

登录并签发 token 对。公开端点。

请求体：

| 字段 | 类型 | 说明 |
|---|---|---|
| `loginway` | string | `legacy`（唯一支持，`oauth-github` 返回 400） |
| `username` | string | 用户名 |
| `password` | string | 密码 |
| `entrycode` | string | 预留字段，legacy 方式不使用 |

```json
{
  "loginway": "legacy",
  "username": "alice",
  "password": "secret123"
}
```

响应：

- `200`

```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ..."
}
```

- `400` 用户名或密码错误（`failure in user login; details::...`），或参数校验失败

### POST /api/v1/refresh-access-token

使用 refresh token 换发全新 token 对。公开端点，无需认证头。

请求体：

| 字段 | 类型 | 说明 |
|---|---|---|
| `refresh_token` | string | 登录或上次刷新获得的 refresh token |

```json
{
  "refresh_token": "eyJ..."
}
```

响应：

- `200`，返回全新 token 对

```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ..."
}
```

- `401` refresh token 无法在库中找到或已被吊销（含已使用 token 的重放）
- `400` 请求体不是合法 JSON
- `500` `refresh_token` 字段缺失、token 过期或验签失败（非业务错误，统一回退 `internal server error`）

轮换语义：

- 每次成功刷新都会吊销被使用的这一枚 refresh token，并签发全新 token 对；客户端必须整体覆盖保存，旧 refresh token 立即失效
- 吊销精确到单枚 token：多设备各自持有独立的 refresh token，一台设备刷新不影响其他设备的登录态

### GET /api/v1/user/me

获取当前认证用户的信息。需要认证。

请求头：`Authorization: Bearer <access_token>`

响应：

- `200`

```json
{
  "uid": 1527277,
  "username": "alice",
  "nickname": "Alice",
  "email": "",
  "register_time": "2026-08-31T12:00:00.000000Z",
  "role": "user",
  "github_id": null,
  "github_login": null
}
```

- `401` access token 缺失、无效或过期
- `404` token 有效但用户记录不存在

### POST /api/v1/admin/sudo

管理端点，占位实现。需要认证，且要求 `admin` 角色。

响应：

- `200` 空响应体
- `401` 认证失败
- `403` 已认证但角色非 `admin`

## Token 语义补充

- access token 载荷：`uid`、`username`、`role`、`iat`、`exp`；角色取值 `user` / `admin`
- refresh token 载荷另含随机 `jti`；服务端只存哈希，不存原文
- 已吊销的 token 记录保留 72h，超期由服务端自动清理
- 客户端保存的 refresh token 一旦被使用即作废，重放会得到 401，此时应引导用户重新登录
