# XrayR 环境变量配置说明

## 概述

现在 XrayR 支持通过环境变量来配置参数。当配置文件不存在时，程序会自动从环境变量读取配置并生成配置文件。

## 支持的环境变量

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `XRAYR_API_HOST` | `https://your-v2board-domain.com` | 面板API地址 |
| `XRAYR_API_KEY` | `your-api-key-here` | 面板API密钥 |
| `XRAYR_NODE_ID` | `1` | 节点ID |
| `XRAYR_PANEL_TYPE` | `V2board` | 面板类型 |
| `XRAYR_NODE_TYPE` | `V2ray` | 节点类型 |
| `XRAYR_LISTEN_IP` | `0.0.0.0` | 监听IP地址 |
| `XRAYR_SEND_IP` | `0.0.0.0` | 发送IP地址 |
| `XRAYR_LOG_LEVEL` | `info` | 日志级别 |

## 使用方法

### 1. 使用环境变量启动

```bash
# 设置环境变量
export XRAYR_API_HOST="https://your-panel.com"
export XRAYR_API_KEY="your-secret-key"
export XRAYR_NODE_ID="123"
export XRAYR_PANEL_TYPE="V2board"
export XRAYR_NODE_TYPE="V2ray"

# 启动程序
./XrayR
```

### 2. 使用 .env 文件

创建 `.env` 文件：
```bash
XRAYR_API_HOST=https://your-panel.com
XRAYR_API_KEY=your-secret-key
XRAYR_NODE_ID=123
XRAYR_PANEL_TYPE=V2board
XRAYR_NODE_TYPE=V2ray
XRAYR_LISTEN_IP=0.0.0.0
XRAYR_SEND_IP=0.0.0.0
XRAYR_LOG_LEVEL=info
```

然后使用 `source` 命令加载：
```bash
source .env
./XrayR
```

### 3. 使用 Docker 环境变量

```bash
docker run -d \
  -e XRAYR_API_HOST="https://your-panel.com" \
  -e XRAYR_API_KEY="your-secret-key" \
  -e XRAYR_NODE_ID="123" \
  -e XRAYR_PANEL_TYPE="V2board" \
  -e XRAYR_NODE_TYPE="V2ray" \
  -p 443:443 \
  xrayr/xrayr
```

## 生成的配置文件

程序会根据环境变量生成以下格式的配置文件：

```yaml
Log:
  Level: info

ConnectionConfig:
  Handshake: 4
  ConnIdle: 300
  UplinkOnly: 2
  DownlinkOnly: 4
  BufferSize: 64

Nodes:
  - PanelType: "V2board"
    ApiConfig:
      ApiHost: "https://your-panel.com"
      ApiKey: "your-secret-key"
      NodeID: 123
      NodeType: V2ray
    ControllerConfig:
      ListenIP: 0.0.0.0
      SendIP: 0.0.0.0
      UpdatePeriodic: 60
```

## 配置优先级

1. **命令行参数** (`-c` 指定配置文件) - 最高优先级
2. **环境变量** - 中等优先级
3. **默认值** - 最低优先级

## 面板类型支持

- `V2board`
- `SSpanel`
- `NewV2board`
- `PMpanel`
- `Proxypanel`
- `V2RaySocks`
- `GoV2Panel`
- `BunPanel`

## 节点类型支持

- `V2ray`
- `Vmess`
- `Vless`
- `Shadowsocks`
- `Trojan`
- `Shadowsocks-Plugin`

## 示例配置

### V2board 面板
```bash
export XRAYR_API_HOST="https://your-v2board.com"
export XRAYR_API_KEY="your-v2board-api-key"
export XRAYR_NODE_ID="1"
export XRAYR_PANEL_TYPE="V2board"
export XRAYR_NODE_TYPE="V2ray"
```

### SSpanel 面板
```bash
export XRAYR_API_HOST="https://your-sspanel.com"
export XRAYR_API_KEY="your-sspanel-api-key"
export XRAYR_NODE_ID="2"
export XRAYR_PANEL_TYPE="SSpanel"
export XRAYR_NODE_TYPE="V2ray"
```

## 注意事项

1. 环境变量只在配置文件不存在时生效
2. 如果配置文件已存在，环境变量不会覆盖配置文件
3. 所有环境变量都是可选的，未设置时使用默认值
4. 生成的配置文件会保存在 `config/config.yml` 