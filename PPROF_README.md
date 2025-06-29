# XrayR Pprof 性能分析

XrayR 现在支持通过 pprof 进行性能分析，帮助您监控和分析应用程序的性能。

## 功能特性

- 默认启用 pprof 服务器
- 默认监听地址：`127.0.0.1:5050`
- 支持 CPU、内存、goroutine 等性能分析
- 可通过配置文件自定义设置

## 配置说明

在 `config.yml` 文件中添加以下配置：

```yaml
Pprof:
  enabled: true  # 启用 pprof 服务器
  address: "127.0.0.1:5050"  # pprof 服务器地址和端口
```

### 配置参数

- `enabled`: 是否启用 pprof 服务器（默认：true）
- `address`: pprof 服务器监听地址（默认：127.0.0.1:5050）

## 使用方法

### 1. 启动 XrayR

启动 XrayR 后，pprof 服务器会自动在配置的地址上启动。

### 2. 访问 pprof 界面

在浏览器中访问：`http://127.0.0.1:5050/debug/pprof/`

### 3. 常用的 pprof 端点

- **CPU 分析**: `http://127.0.0.1:5050/debug/pprof/profile?seconds=30`
- **内存分析**: `http://127.0.0.1:5050/debug/pprof/heap`
- **Goroutine 分析**: `http://127.0.0.1:5050/debug/pprof/goroutine`
- **阻塞分析**: `http://127.0.0.1:5050/debug/pprof/block`
- **互斥锁分析**: `http://127.0.0.1:5050/debug/pprof/mutex`

### 4. 使用 go tool pprof 命令行工具

```bash
# CPU 分析（采样 30 秒）
go tool pprof "http://154.21.80.108:5050/debug/pprof/profile?seconds=30"

# 内存（堆）分析
go tool pprof "http://127.0.0.1:5050/debug/pprof/heap"

# Goroutine 分析
go tool pprof "http://127.0.0.1:5050/debug/pprof/goroutine"
```

### 5. 生成火焰图

```bash
# ✅ 安装 Graphviz（如尚未安装）
# macOS:
brew install graphviz

# 打开 CPU 火焰图界面（采样 30 秒）
go tool pprof -http=:8080 "http://154.21.80.108:5050/debug/pprof/profile?seconds=30"

# 打开内存火焰图界面
go tool pprof -http=:8080 "http://127.0.0.1:5050/debug/pprof/heap"

```

## 安全注意事项

⚠️ **重要**: pprof 服务器默认只监听本地地址（127.0.0.1），这是出于安全考虑。如果您需要从外部访问，请确保：

1. 修改 `address` 为 `0.0.0.0:5050`
2. 配置防火墙规则限制访问
3. 在生产环境中考虑使用反向代理和认证

## 故障排除

### 无法访问 pprof 服务器

1. 检查配置文件中的 `Pprof` 设置
2. 确认 `enabled` 为 `true`
3. 检查防火墙设置
4. 查看 XrayR 日志确认 pprof 服务器已启动

### 性能分析工具无法连接

1. 确认端口 5050 没有被其他程序占用
2. 检查网络连接
3. 尝试使用 `curl` 测试连接：`curl http://127.0.0.1:5050/debug/pprof/`

## 示例配置

```yaml
Log:
  Level: warning
  AccessPath: 
  ErrorPath: 
ConnectionConfig:
  Handshake: 4
  ConnIdle: 30
  UplinkOnly: 2
  DownlinkOnly: 4
  BufferSize: 64
Pprof:  # 默认设置可不填写
  enabled: true
  address: "127.0.0.1:5050"
Nodes:
  - PanelType: "V2board"
    # ... 其他节点配置
```

## 更多信息

- [Go pprof 官方文档](https://golang.org/pkg/net/http/pprof/)
- [性能分析最佳实践](https://golang.org/doc/diagnostics.html)
- [火焰图分析](https://github.com/brendangregg/FlameGraph) 