# Linux 构建资源目录

此目录存放 Linux 平台构建所需的资源文件。

## 文件说明

| 文件 | 说明 |
|------|------|
| `icon.png` | 应用程序图标（512x512 PNG） |

## .deb 包结构

GitHub Actions 构建时自动生成 .deb 包，目录结构如下：

```
xtool_VERSION_amd64/
├── DEBIAN/
│   ├── control       # 包元信息
│   ├── postinst      # 安装后脚本
│   └── prerm         # 卸载前脚本
└── usr/
    ├── bin/
    │   └── xtool    # 主程序
    └── share/
        ├── applications/
        │   └── xtool.desktop   # 应用菜单图标
        └── doc/
            └── xtool/
                ├── README.md
                └── LICENSE
```

## 安装方式

```bash
# DEB 包安装（Ubuntu/Debian）
sudo dpkg -i xtool_*_amd64.deb

# 或使用 apt（自动处理依赖）
sudo apt install ./xtool_*_amd64.deb

# tar.gz 通用安装
tar -xzf xtool_linux_amd64.tar.gz
./xtool
```

## 系统依赖

运行需要：
- `libgtk-3-0`
- `libwebkit2gtk-4.0-37`（Ubuntu 22.04 可能需要 webkit2gtk-4.1）

```bash
# Ubuntu/Debian 安装依赖
sudo apt-get install -y libgtk-3-0 libwebkit2gtk-4.0-37
```
