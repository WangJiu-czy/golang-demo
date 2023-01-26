package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const (
	name     = "serve"
	version  = "0.0.2"
	revision = "HEAD"
)

func main() {
	addr := flag.String("a", ":5000", "address to serve(host:port)")
	prefix := flag.String("p", "/", "prefix path under")
	root := flag.String("r", ".", "root path to serve")
	certFile := flag.String("cf", "", "tls cert file")
	keyFile := flag.String("kf", "", "tls key file")
	dumpPost := flag.Bool("dumpPost", false, "dump post data")
	showVersion := flag.Bool("v", false, "show version")
	flag.Parse()
	if *showVersion {
		fmt.Printf("%s %s (rev: %s/%s)\n", name, version, revision, runtime.Version())
		return
	}
	var err error
	*root, err = filepath.Abs(*root)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("serving %s as %s on %s", *root, *prefix, *addr)
	//StripPrefix将给定的前缀删除,来获取FileServer提供的文件,http.Dir使用操作系统的文件系统实现
	http.Handle(*prefix, http.StripPrefix(*prefix, http.FileServer(http.Dir(*root))))
	//ServeHTTP 将请求分派给其模式与请求 URL 最匹配的处理程序
	mux := http.DefaultServeMux.ServeHTTP
	var logger http.HandlerFunc
	//下面这一段只是在记录日志,mux()也就是ServeHTTP函数,将请求分配到url对应的路由上面
	if *dumpPost {
		logger = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RemoteAddr + " " + r.Method + " " + r.URL.String())
			//请求体内容拷贝到os.Stderr
			io.Copy(os.Stderr, r.Body)
			os.Stderr.Write([]byte{'\n'})
			mux(w, r)
		})
	} else {
		logger = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RemoteAddr + " " + r.Method + " " + r.URL.String())
			mux(w, r)
		})
	}
	if *certFile != "" && *keyFile != "" {
		err = http.ListenAndServeTLS(*addr, *certFile, *keyFile, logger)
	} else {
		err = http.ListenAndServe(*addr, logger)
	}
	if err != nil {
		log.Fatalln(err)
	}

}
