# http_custom_ip_dl

# httpdownloader

一个可以自定义 IP 和 User-Agent 的 Go 语言 HTTP 下载器。

已知问题：
仅仅处理了http 200的情况，不支持301、302

## 安装

使用以下命令下载安装：

go get github.com/imgingroot/httpIPdownloader


## 使用

```go
package main

import (
	"fmt"
	"github.com/imgingroot/httpIPdownloader"
)

func main() {
	url := "https://codeload.github.com/imgingroot/jsonl_check_fillter/zip/refs/heads/main"
	filename := "file.zip"
	ip := "140.82.112.9" // 如果不需要自定义 IP 就设置为空字符串
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36" // 如果不需要自定义 User-Agent 就设置为空字符串

	err := httpIPdownloader.DownloadFile(url, filename, ip, ua)
	if err != nil {
		fmt.Println("下载文件出错：", err)
	} else {
		fmt.Printf("已将文件下载至 %s\n", filename)
	}
}
```