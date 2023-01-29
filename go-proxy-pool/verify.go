package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var verifyIs = false

var ProxyPool []ProxyIp
var lock sync.Mutex
var mux2 sync.Mutex

var count int

func countAdd(i int) {
	mux2.Lock()
	count += i
	mux2.Unlock()
}
func countDel() {
	mux2.Lock()
	fmt.Println("\r 代理验证中: %d ", count)
	count--
	mux2.Unlock()
}
func Verify(pi *ProxyIp, wg *sync.WaitGroup, ch chan int, first bool) {
	defer func() {
		wg.Done()
		countDel()
		<-ch
	}()

	//	pr := pi.Ip + ":" + pi.Port
	//	startT := time.Now()

}
func VerifyHttp(pr string) bool {
	proxyUrl, proxyErr := url.Parse("http://" + pr)
	if proxyErr != nil {
		return false
	}
	tr := http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	tr.Proxy = http.ProxyURL(proxyUrl)
	client := http.Client{Timeout: 10 * time.Second, Transport: &tr}
	request, err := http.NewRequest("GET", "http://baid.com", nil)
	res, err := client.Do(request)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	dataBytes, _ := io.ReadAll(res.Body)
	result := string(dataBytes)
	if strings.Contains(result, "0;url=http://www.baidu.com") {
		return true
	}
	return false

}
func VerifyHttps(pr string) bool {
	destConn, err := net.DialTimeout("tcp", pr, 10*time.Second)
	if err != nil {
		return false
	}
	defer destConn.Close()
	req := []byte{67, 79, 78, 78, 69, 67, 84, 32, 119, 119, 119, 46, 98, 97, 105, 100, 117, 46, 99, 111, 109, 58, 52, 52, 51, 32, 72, 84, 84, 80, 47, 49, 46, 49, 13, 10, 72, 111, 115, 116, 58, 32, 119, 119, 119, 46, 98, 97, 105, 100, 117, 46, 99, 111, 109, 58, 52, 52, 51, 13, 10, 85, 115, 101, 114, 45, 65, 103, 101, 110, 116, 58, 32, 71, 111, 45, 104, 116, 116, 112, 45, 99, 108, 105, 101, 110, 116, 47, 49, 46, 49, 13, 10, 13, 10}
	//随便准备一个byte数组,将数据写入连接
	destConn.Write(req)
	bytes := make([]byte, 1024)
	//设置读取数据的超时时间,最后时间是当前时间加上10秒
	destConn.SetReadDeadline(time.Now().Add(10 * time.Second))
	read, err := destConn.Read(bytes)
	if strings.Contains(string(bytes[:read]), "200 Connection established") {
		return true
	}
	return false

}
func VerifySocket5(pr string) bool {
	destConn, err := net.DialTimeout("tcp", pr, 10*time.Second)
	if err != nil {
		return false
	}
	defer destConn.Close()
	req := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	destConn.Write(req)
	bytes := make([]byte, 1024)
	destConn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, err = destConn.Read(bytes)
	if err != nil {
		return false
	}
	if bytes[0] == 5 && bytes[1] == 255 {
		return true
	}
	return false
}

//TODO:还没好
