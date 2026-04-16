# Volcano TTS Adaptor

本 fork 为 one-api 添加了 Volcano TTS（火山引擎语音合成）支持，用于 XiaChong 项目的云端 TTS 服务。

## 新增功能

- **Channel Type**: `VolcanoTTS` (ID: 57)
- **API Type**: `VolcanoTTS`
- **支持模型**: `volcano-tts`, `seed-tts-1.1`

## 配置方法

### 1. 添加渠道

在 one-api 管理面板中：

1. **渠道管理 → 添加渠道**
2. **类型**: 选择 `Volcano TTS`
3. **名称**: 自定义名称
4. **密钥**: 填入 JSON 格式凭证：
   ```json
   {"app_id": "your-app-id", "access_key": "your-access-key", "resource_id": "volc.service_type.10029"}
   ```
5. **模型**: `volcano-tts,seed-tts-1.1`
6. **分组**: 根据需要设置

### 2. 使用 API

与 OpenAI TTS API 兼容：

```bash
curl -X POST https://your-one-api-server/v1/audio/speech \
  -H "Authorization: Bearer sk-your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "volcano-tts",
    "input": "你好，世界！",
    "voice": "nova"
  }' \
  --output output.mp3
```

### 3. 支持的 Voice 映射

| OpenAI Voice | Volcano Speaker |
|--------------|-----------------|
| alloy | zh_female_shuangkuaisisi_moon_bigtts |
| echo | zh_male_chunhou_moon_bigtts |
| fable | zh_female_qiaopi_moon_bigtts |
| onyx | zh_male_wennuan_moon_bigtts |
| nova | zh_female_shuangkuaisisi_moon_bigtts |
| shimmer | zh_female_tianmei_moon_bigtts |

也可以直接使用 Volcano 原生 speaker ID。

## 文件结构

```
relay/adaptor/volcanoTTS/
├── adaptor.go       # Adaptor 接口实现
├── config.go        # 凭证配置解析
├── request.go       # 请求转换逻辑
├── response.go      # 响应结构定义
└── tts_handler.go   # 流式 TTS 处理
```

## 修改的文件

| 文件 | 修改 |
|------|------|
| `relay/channeltype/define.go` | 添加 `VolcanoTTS = 57` |
| `relay/apitype/define.go` | 添加 `VolcanoTTS` |
| `relay/adaptor.go` | 注册适配器 |
| `controller/user.go` | CreateUser 返回用户 ID，支持创建时设置 group/quota |
| `controller/token.go` | 新增 `AdminAddToken` 函数 |
| `router/api.go` | 新增 `/api/admin/token` 路由 |

## 外部集成 API

本 fork 新增了以下 API 供外部系统（如 XiaChong Workers）使用：

### POST /api/user (已修改)

创建用户，现在返回用户 ID：

```json
// Request
{
  "username": "xc_user123",
  "password": "random_password",
  "display_name": "用户昵称",
  "group": "premium",
  "quota": 50000
}

// Response
{
  "success": true,
  "data": {
    "id": 42,
    "username": "xc_user123"
  }
}
```

### POST /api/admin/token (新增)

为指定用户创建 API Token（需要管理员权限）：

```json
// Request
{
  "user_id": 42,
  "name": "xiachong-tts",
  "unlimited_quota": true,
  "expired_time": -1
}

// Response
{
  "success": true,
  "data": {
    "key": "sk-xxxxxxxxxxxxxxxx",
    "user_id": 42,
    "name": "xiachong-tts"
  }
}
```

## 获取 Volcano TTS 凭证

1. 访问 [火山引擎控制台](https://console.volcengine.com/)
2. 开通语音合成服务
3. 创建应用，获取 `app_id` 和 `access_key`
4. `resource_id` 默认为 `volc.service_type.10029`

## 相关项目

- [XiaChong](https://github.com/Neon-Wang/XiaChong) - 桌面数字宠物应用
- [one-api](https://github.com/songquanpeng/one-api) - 原始项目
