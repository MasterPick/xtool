# macOS 构建资源目录

此目录存放 macOS 平台构建所需的资源文件。

## 文件说明

| 文件 | 说明 |
|------|------|
| `icon.icns` | 应用程序图标（ICNS 格式，含多分辨率） |

## DMG 安装镜像

GitHub Actions 构建时自动创建 .dmg 文件：

| 文件名 | 适用架构 |
|--------|----------|
| `universal-toolbox_macos_amd64.dmg` | Intel Mac |
| `universal-toolbox_macos_arm64.dmg` | Apple Silicon (M1/M2/M3) |

## 安装方式

1. 双击 `.dmg` 文件挂载镜像
2. 将 `万能工具箱.app` 拖入 `Applications` 文件夹
3. 在启动台中找到应用并运行

## 首次运行说明

macOS Gatekeeper 可能阻止运行未签名应用，解决方式：

```bash
# 方式1：右键点击 .app，选择"打开"
# 方式2：终端命令移除隔离标志
sudo xattr -rd com.apple.quarantine /Applications/万能工具箱.app
```

## 本地构建

```bash
# 安装 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 构建 Intel 版
wails build -platform darwin/amd64

# 构建 Apple Silicon 版
wails build -platform darwin/arm64

# 创建 DMG
hdiutil create -volname "万能工具箱" \
  -srcfolder "build/bin/万能工具箱.app" \
  -ov -format UDZO \
  "build/bin/universal-toolbox_macos.dmg"
```
