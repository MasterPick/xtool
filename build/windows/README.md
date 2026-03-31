# Windows 构建资源目录

此目录存放 Windows 平台构建所需的资源文件。

## 文件说明

| 文件 | 说明 |
|------|------|
| `icon.ico` | 应用程序图标（16x16 ~ 256x256 多分辨率） |
| `installer.nsi` | NSIS 安装脚本（由 CI 自动生成） |
| `wails.exe.manifest` | Windows 应用程序清单（由 Wails 自动生成） |

## 图标要求

- 格式：ICO（包含 16/32/48/64/128/256px 尺寸）
- 位置：`build/windows/icon.ico`
- 如未提供，Wails 将使用默认图标

## NSIS 安装包

NSIS 安装脚本在 GitHub Actions 构建时自动生成，包含：
- 安装向导（中文界面）
- 创建桌面快捷方式
- 创建开始菜单项
- 写入注册表（支持控制面板卸载）
- 卸载程序

## 本地手动构建

```powershell
# 确保已安装 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 安装前端依赖
cd frontend && npm install && cd ..

# 构建 Windows 版本
wails build -platform windows/amd64

# 输出路径：build/bin/万能工具箱.exe
```
