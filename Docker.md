# Docker 

### build

```
docker build -t dollarkiller/xrayr:latest .
```

### deployment

docker-compose

```yaml
services:
  xrayr:
    image: dollarkiller/xrayr:latest
    container_name: xrayr
    volumes:
      - ./xrayr/config.yml:/etc/XrayR/config.yml:ro
      - ./ssl:/etc/xrayr-ssl:ro
    network_mode: host
    restart: unless-stopped
```

config.yml

```yaml
Log:
  Level: info
  AccessPath: ""
  ErrorPath: ""

DnsConfigPath: ""

ConnectionConfig:
  Handshake: 4
  ConnIdle: 300
  UplinkOnly: 2
  DownlinkOnly: 5
  BufferSize: 64

Nodes:
  - PanelType: "V2board"
    ApiConfig:
      ApiHost: "https://api.hisuyou.org"
      ApiKey: "a60a866d-97b4-411c-ae97-a208e1ed4f09"
      NodeID: 313
      NodeType: Trojan
      Timeout: 30
      EnableVless: false
      EnableXTLS: false
      SpeedLimit: 0
      DeviceLimit: 0
    ControllerConfig:
      ListenIP: 0.0.0.0
      ListenPort: 443
      EnableDNS: false
      TransportConfig:
        Enable: true
        Type: tcp
        Path: /api/upload
        Host: news.ycombinator.com
      CertConfig:
        Enable: true
        CertMode: file       # 使用本地文件模式 :contentReference[oaicite:2]{index=2}
        CertDomain: news.ycombinator.com
        CertFile: /etc/xrayr-ssl/cert.crt
        KeyFile: /etc/xrayr-ssl/private.key
```