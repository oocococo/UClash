## shadowsocks-update-go

使用GO语言实现的[shadowsocks-update](https://github.com/Fndroid/shadowsocks-update)

使用相同格式的配置文件

## 下载使用

1. [下载ssu.exe](https://github.com/Fndroid/shadowsocks-update-go/releases)，放置在Shadowsocks目录下
2. 创建配置文件``update.json``:
    ```json
    {
        "providers": [
            "https://xxx.xxx.com",
            "https://yyy.yyy.com"
        ],
        "filter": [
            "HK",
            "TW"
        ]
    }
    ```
    > 网址为Surge托管地址
    
    > filter为保留关键字
3. 运行
  - 双击运行：更新完成即自动关闭
  - 命令行运行：更新完成可显示结果
    
    ``./ssu.exe``