package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	figure "github.com/common-nighthawk/go-figure"
	"github.com/onedss/httpFileServer/utils"
	"github.com/onedss/service"
)

type program struct {
	httpPort   int
	webroot    string
	httpServer *http.Server
}

func (p *program) Start(s service.Service) (err error) {
	log.Println("********** START **********")
	if utils.IsPortInUse(p.httpPort) {
		err = fmt.Errorf("HTTP port[%d] In Use", p.httpPort)
		return
	}
	p.StartHTTP()
	return
}

func (p *program) Stop(s service.Service) (err error) {
	defer log.Println("********** STOP **********")
	p.StopHTTP()
	return
}

func (p *program) StartHTTP() (err error) {
	p.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", p.httpPort),
		Handler:           http.FileServer(http.Dir(p.webroot)),
		ReadHeaderTimeout: 5 * time.Second,
	}
	link := fmt.Sprintf("http://%s:%d", utils.LocalIP(), p.httpPort)
	log.Println("http server start -->", link)
	go func() {
		if err := p.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("start http server error", err)
		}
		log.Println("http server end")
	}()
	return
}

func (p *program) StopHTTP() (err error) {
	if p.httpServer == nil {
		err = fmt.Errorf("HTTP Server Not Found")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = p.httpServer.Shutdown(ctx); err != nil {
		return
	}
	return
}

func main() {
	sec := utils.Conf().Section("service")
	svcConfig := &service.Config{
		Name:        sec.Key("name").MustString("HttpFileServer_Service"),
		DisplayName: sec.Key("display_name").MustString("HttpFileServer_Service"),
		Description: sec.Key("description").MustString("HttpFileServer_Service"),
	}

	rootPath := fmt.Sprintf("%s/www", utils.CWD())
	httpPort := utils.Conf().Section("http").Key("port").MustInt(8181)
	webroot := utils.Conf().Section("http").Key("webroot").MustString(rootPath)
	fmt.Println("Http File Server is Starting... ")
	fmt.Println("Listen Port:", httpPort, "")
	fmt.Println("ROOT PATH:", webroot, "")
	p := &program{
		httpPort: httpPort,
		webroot:  webroot,
	}
	var s, err = service.New(p, svcConfig)
	if err != nil {
		log.Println(err)
		utils.PauseExit()
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "install" || os.Args[1] == "stop" {
			figure.NewFigure("HttpFileServer", "", false).Print()
		}
		log.Println(svcConfig.Name, os.Args[1], "...")
		if err = service.Control(s, os.Args[1]); err != nil {
			log.Println(err)
			utils.PauseExit()
		}
		log.Println(svcConfig.Name, os.Args[1], "ok")
		return
	}
	figure.NewFigure("HttpFileServer", "", false).Print()
	if err = s.Run(); err != nil {
		log.Println(err)
		utils.PauseExit()
	}
}
