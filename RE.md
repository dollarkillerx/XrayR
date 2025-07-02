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


-----BEGIN CERTIFICATE-----
MIIEmTCCA4GgAwIBAgIUGwT4wW+avstRsd7Sxv3fKqRnDkMwDQYJKoZIhvcNAQEL
BQAwgYsxCzAJBgNVBAYTAlVTMRkwFwYDVQQKExBDbG91ZEZsYXJlLCBJbmMuMTQw
MgYDVQQLEytDbG91ZEZsYXJlIE9yaWdpbiBTU0wgQ2VydGlmaWNhdGUgQXV0aG9y
aXR5MRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMRMwEQYDVQQIEwpDYWxpZm9ybmlh
MB4XDTI1MDcwMjA2MTkwMFoXDTQwMDYyODA2MTkwMFowYjEZMBcGA1UEChMQQ2xv
dWRGbGFyZSwgSW5jLjEdMBsGA1UECxMUQ2xvdWRGbGFyZSBPcmlnaW4gQ0ExJjAk
BgNVBAMTHUNsb3VkRmxhcmUgT3JpZ2luIENlcnRpZmljYXRlMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEAp4resC+xym4V2WoDsLw9Zq91JnKHmfZSc5aV
9eVMh2yVUQN1L1xzTvPhJmJcwhUUZzzG71vyIxXtj5b0Dyge3gCmiwyUNgHqzF4d
xMi1VqVGaib5q8lcxGoF94QVCHnE0dXftSHxDxBn+5Y1D3WAdII37UHYE2qMwKKO
Ff3X8u4Y+eHu0/mUX97lv7dzigii3lx0sUIKg6po6pfkSl1PSkF2gZsuZQsbSAr4
ZuOs+7LWNf1NYEuU+kaa2NHLKsZnwWMwmsePAihi29PHoVorQ9yCs11zGra9EYbM
V9szweYiY/SNyDx3QN0JMwpF2HY2D0vP2Wr0If5H3hrNaJHWgwIDAQABo4IBGzCC
ARcwDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMCBggrBgEFBQcD
ATAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBRRPULwaDP7db2mcjEdREHv+KGtHDAf
BgNVHSMEGDAWgBQk6FNXXXw0QIep65TbuuEWePwppDBABggrBgEFBQcBAQQ0MDIw
MAYIKwYBBQUHMAGGJGh0dHA6Ly9vY3NwLmNsb3VkZmxhcmUuY29tL29yaWdpbl9j
YTAcBgNVHREEFTATghFva3IuaGFja3NuZXdzLnRvcDA4BgNVHR8EMTAvMC2gK6Ap
hidodHRwOi8vY3JsLmNsb3VkZmxhcmUuY29tL29yaWdpbl9jYS5jcmwwDQYJKoZI
hvcNAQELBQADggEBAJ4vIFNda/i2VHKvPWS1OddxI6Do44v5G//rAAYaY8rMHvNE
ytZp8DlLwz74v/p2mBeJpEscXGo7g0tVNvCy3l5SSCoe+SsngweiA3dHZucjvf3J
OeLfMIxSs6RrHjnU+8X/WTiA0NJAdzL+KU5F9r7nnvK88xAPNlLvZ2BblOWy0bkz
K+m6nXUhEMyALMarb+ADjVOLcmTNMazYqrpK5ty9NPQiR8m9IGr1XYuggH6XwXBc
gGElix8wWI78hLrghQzk7eDQldVq2RY3o6/fxM1GoKEsecXOnOfPNPLRu0WWfwfr
9ruySQem3yl5g3bElQ9Hqc4tXulR28MGsT934ew=
-----END CERTIFICATE-----


-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCnit6wL7HKbhXZ
agOwvD1mr3UmcoeZ9lJzlpX15UyHbJVRA3UvXHNO8+EmYlzCFRRnPMbvW/IjFe2P
lvQPKB7eAKaLDJQ2AerMXh3EyLVWpUZqJvmryVzEagX3hBUIecTR1d+1IfEPEGf7
ljUPdYB0gjftQdgTaozAoo4V/dfy7hj54e7T+ZRf3uW/t3OKCKLeXHSxQgqDqmjq
l+RKXU9KQXaBmy5lCxtICvhm46z7stY1/U1gS5T6RprY0csqxmfBYzCax48CKGLb
08ehWitD3IKzXXMatr0RhsxX2zPB5iJj9I3IPHdA3QkzCkXYdjYPS8/ZavQh/kfe
Gs1okdaDAgMBAAECggEAH296bfNc9/7kwZjXHWuMV753DCI0GCqz+LUZ6Pu/eq9g
asVr0NmXVwqCRanx/oZnACHSr6mASp38J1Nq8ocUF5JmrTNMgb+lVvgXvMrt31Z+
Wj963AUZj07XVMQnECKMz9RBt6IpMSriX72ksgRZeELQjIcsrCPzSdeWMkEhKha9
WgWNBgeuInm4bnSVaKH8hVOTFAURgBkPTRc4E+N8kpQf7L+nFVDOt74iMY4/pCkI
T0FqDehDNkH1q3ZphrxN8XmQLk0/w/GEtazwHoPN7OQWfx3NoLDkxh6ExBIOOcp1
cNU85ePO+hNOLuciJZLDPSg4JcuvLvDDgsYx0fIWlQKBgQDmGtSCUySFcuERtbKi
OSxsr4ZqSNhtc+4Ea+1IEJ7/IewOPZYPMUIaOA3P+hoWjjHv2FfddTwL0ZZUtPkq
3H0V55VVVCMI51zEHliIe/uuBeftYnOsZdz0Gk48DaFHAKaHPZLR4uxaZpOkFM75
wRNPs0mzIvClmzVl9UwAJb7ZTwKBgQC6ZajGFd7posO+0iv1xpVB7cEOodP/csKJ
R1mhqt9FMe1Srdn3vZvY4YDZccOubGgq5/kUyXvKr6TLOgbtf9QdpuQr08PoMGZH
3n9Hhs2ki0l8Nu9xu+3ASLCTpn00mANB70qJcDu+gPJge0Z3LqdDQVG6Unzo2CAY
DwRF9LP6jQKBgG2I5hInu+Henq8MfD1m0/Pk1ipwBL6NigcGzUwZxWYT4segn3cZ
7qqGdmTDXHnSNIlExgbAkNXbbiFiiJ3TmaO8usSEqazKXclp3KETy+m6G/5PeFrf
nn+Aqi9CGbyv3ZXRRLuutu4NmEhqM2CGfUfaQ8IHZuCecpvXWddUIHZXAoGAT40H
hNLcuhXPOOd1/4TrOqH/3QDP+5u8zt32sPq+I0f1o22zMvpnQx6q4SKegxidNIKg
WXDMNmrUCDARuNbvbmBFzTapy7SsUkvxQlOT/H+9sxe3BXFRPNlJExrhIMsnyMIZ
q6KMvPdHeknifwcYN6nu+Xgu+ykSOXJfPoP7NBkCgYEAmWh+qyqcWp75BXwNYqJy
V8BtZr2FBYZELYq3vYnfz8JtVdWSomSHgdELWIUIuz08oOYuym3p3aQa1brHu4bY
986dvDsVkHzXa5iNuk9xd5qsuGa7j1UijHpwOrfhRQVFxruhXVxBIrZ8x6OKrSJg
PRaiCPRFFXRWtf2rP+xNywY=
-----END PRIVATE KEY-----

