# Changelog

所有版本的重要变更记录。

格式遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)

## [1.1.0] - 2026-03-31

### 新增
- 🚀 **多平台安装包支持**
  - Windows：`.exe` 单文件 + NSIS 安装向导（含桌面快捷方式、控制面板卸载）
  - Linux：`.deb` 安装包（含应用菜单图标）+ `tar.gz` 通用压缩包
  - macOS：Intel `.dmg` + Apple Silicon (M1/M2/M3) `.dmg` 双架构
- ⚡ **GitHub Actions CI/CD 升级**
  - 三平台并行构建，互不依赖
  - 自动生成 SHA256 校验和文件
  - Release 页面展示详细安装说明
- 📦 **Makefile 本地构建支持**
  - `make dev` 启动热重载开发模式
  - `make build-linux` 生成 .deb + tar.gz
  - `make build-macos` 生成双架构 .dmg
- 📋 **平台构建文档**：各平台详细构建说明

### 优化
- GitHub Release 页面美化，含功能展开列表
- 构建产物 SHA256 校验和汇总到单文件

## [1.0.0] - 2026-03-31

### 新增
- 🛠️ **开发工具模块**：JSON 格式化/压缩/校验/转义，XML 格式化，Base64 编解码，URL 编解码，哈希计算（MD5/SHA1/SHA256），文本对比/替换/统计，UUID 批量生成，时间戳转换，正则测试，代码片段管理
- 💻 **系统工具模块**：进程管理（查看/终止），端口占用查看/释放，系统信息（CPU/内存/磁盘），批量文件重命名
- 📱 **日常工具模块**：标准+科学计算器，长度/重量/温度/速度换算，备忘录
- 🌐 **网络工具模块**：Ping 测试，内网 IP 扫描，HTTP 接口测试，Hosts 文件编辑
- 🎨 **主题系统**：深色/浅色/自动/蓝/绿/紫/橙/云母效果 8 种主题
- ⚙️ **设置中心**：主题切换、字体大小、窗口行为配置
- 🪟 **Windows 11 云母效果**：原生半透明背景

### 技术栈
- Go 1.21 + Wails v2.9
- Vue 3.4 + TypeScript + TailwindCSS
- SQLite（本地数据）
- gopsutil（系统信息）
