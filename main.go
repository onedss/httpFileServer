package main

import (
	"fmt"
	"github.com/onedss/httpFileServer/utils"
	"net/http"
)

func main() {
	rootPath := fmt.Sprintf("%s/www", utils.CWD())
	httpPort := utils.Conf().Section("http").Key("port").MustInt(8181)
	webroot := utils.Conf().Section("http").Key("webroot").MustString(rootPath)
	address := fmt.Sprintf(":%d", httpPort)
	fmt.Println("Starting Http File Server... ")
	fmt.Println("Listen Port:", httpPort, "")
	fmt.Println("ROOT PATH:", webroot, "")
	http.ListenAndServe(address, http.FileServer(http.Dir(webroot)))
}
