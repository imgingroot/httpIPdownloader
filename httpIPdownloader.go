package httpIPdownloader

import (
	"context"
	"crypto/tls"
	"fmt"
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
		ResponseHeaderTimeout: time.Second * 30, // 等待服务器响应头的时间
		ExpectContinueTimeout: time.Second * 30, // 发送带有Expect: 100-continue的请求时等待服务器响应的时间

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
			return conn, nil
		},
	}

	client := http.Client{
		Transport: transport,
		Timeout:   0,
	}

	// 发送请求
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response error with status code = %d", resp.StatusCode)
	}

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
