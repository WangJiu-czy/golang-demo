package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func formHandler(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		fmt.Fprintf(writer, "ParseForm() err:%v", err)
		return
	}
	fmt.Fprintln(writer, "POST request successful")
	name := request.FormValue("name")
	addr := request.FormValue("address")
	fmt.Fprintf(writer, "Name=%s\n", name)
	fmt.Fprintf(writer, "Address=%s", addr)
}

func helloHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/hello" {
		http.Error(writer, "404 not fount", http.StatusNotFound)
		return
	}
	if request.Method != "GET" {
		http.Error(writer, "method is not supported", http.StatusNotFound)
		return
	}
	fmt.Fprintf(writer, "hello!")
}
