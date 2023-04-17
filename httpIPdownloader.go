package httpIPdownloader

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

// DownloadFile 用于下载文件
func DownloadFile(url string, filename string, ip string, ua string) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // 跳过证书验证

	// 构建请求
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// 设置 User-Agent
	if ua != "" {
		request.Header.Set("User-Agent", ua)
	}

	// 自定义 DNS 解析
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, _ := net.SplitHostPort(addr)
			if ip != "" {
				host = ip
			}

			ipAddr, err := net.ResolveIPAddr("ip4", host)
			if err != nil {
				return nil, err
			}

			conn, err := net.DialTimeout("tcp", net.JoinHostPort(ipAddr.IP.String(), port), time.Second*30)
			if err != nil {
				return nil, err
			}
			conn.SetDeadline(time.Now().Add(30 * time.Second))
			return conn, nil
		},
	}

	client := http.Client{
		Transport: transport,
		Timeout:   time.Second * 30, // 请求超时时间 30 秒
	}

	// 发送请求
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
