## config

### basice config

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
      ApiHost: "https://your-v2board-domain.com"
      ApiKey: "your-api-key-here"
      NodeID: 1
      NodeType: V2ray
    ControllerConfig:
      ListenIP: 0.0.0.0
      SendIP: 0.0.0.0
      UpdatePeriodic: 60
```

| 协议 | 证书需求 | 说明 |
|------|----------|------|
| **Trojan** | 必须 | 基于TLS，无证书无法工作 |
| **VLESS + XTLS** | 必须 | XTLS需要TLS握手 |
| **V2ray/VMess** | 可选 | 有证书更安全，无证书也能工作 |
| **VLESS** | 可选 | 有证书更安全，无证书也能工作 |
| **Shadowsocks** | 不需要 | 使用自定义加密 |
| **Shadowsocks-Plugin** | 不需要 | 基于Shadowsocks |
| **REALITY** | 不需要真实证书 | 使用虚假证书 |

### 1. **特殊协议 - REALITY**

#### **REALITY协议**
- **证书要求**: **不需要真实证书**
- **原因**: REALITY使用虚假证书，但需要配置REALITY参数
- **配置示例**:
```yaml
NodeType: V2ray
ApiConfig:
  EnableVless: true
ControllerConfig:
  EnableREALITY: true
  REALITYConfigs:
    Dest: "www.microsoft.com:443"
    ServerNames: ["www.microsoft.com"]
    PrivateKey: "your-private-key"
    ShortIds: [""]
  CertConfig:
    CertMode: none  # REALITY不需要真实证书
```

## 证书配置模式说明

### 1. **none模式**
```yaml
CertConfig:
  CertMode: none  # 不使用证书
```
- 适用于: Shadowsocks、无TLS的V2ray/VMess/VLESS

### 2. **file模式**
```yaml
CertConfig:
  CertMode: file
  CertFile: /path/to/cert.crt
  KeyFile: /path/to/key.key
```
- 适用于: 已有证书文件的情况