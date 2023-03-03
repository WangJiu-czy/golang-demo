package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Address() string
	IsAlive() bool
	Server(rw http.ResponseWriter, r *http.Request)
}
type simpleServer struct {
	addr string
	//反向代理对象
	proxy *httputil.ReverseProxy
}

func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)
	handleErr(err)
	return &simpleServer{addr: addr, proxy: httputil.NewSingleHostReverseProxy(serverUrl)}
}
func handleErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

type LocaBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

//初始化负责均衡器
func NewLoadBalancer(port string, servers []Server) *LocaBalancer {

	return &LocaBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

//获取下一个可用的服务器
func (lb *LocaBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}

func (lb *LocaBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {
	server := lb.getNextAvailableServer()
	fmt.Printf("forwarding request to addresss:%q\n", server.Address())
	server.Server(rw, r)
}

func (s simpleServer) Address() string {
	return s.addr
}

func (s simpleServer) IsAlive() bool {
	return true
}

func (s simpleServer) Server(rw http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(rw, r)
}

func main() {
	servers := []Server{
		newSimpleServer("http://class.seig.edu.cn:7001/sise/"),
		newSimpleServer("http://www.bing.com"),
		newSimpleServer("http://www.baidu.com"),
	}
	lb := NewLoadBalancer("8000", servers)
	handleRedirect := func(rw http.ResponseWriter, r *http.Request) {
		lb.serveProxy(rw, r)
	}
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("serving requests at 'localhost:%s'\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
