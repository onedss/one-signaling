package main

import (
	"context"
	"fmt"
	"github.com/onedss/one-signaling/figure"
	"github.com/onedss/one-signaling/keyboard"
	"github.com/onedss/one-signaling/models"
	"github.com/onedss/one-signaling/routers"
	"github.com/onedss/one-signaling/service"
	"github.com/onedss/one-signaling/utils"
	"log"
	"net/http"
	"os"
	"time"
)

type program struct {
	httpPort   int
	httpServer *http.Server
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

func (p *program) StartHTTP() (err error) {
	p.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", p.httpPort),
		Handler:           routers.Router,
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
func (p *program) Start(s service.Service) (err error) {
	log.Println("********** START **********")
	if utils.IsPortInUse(p.httpPort) {
		err = fmt.Errorf("HTTP port[%d] In Use", p.httpPort)
		return
	}
	err = models.Init()
	if err != nil {
		return
	}
	err = routers.Init()
	if err != nil {
		return
	}
	p.StartHTTP()
	if !utils.Debug {
		log.Println("log files -->", utils.LogDir())
		log.SetOutput(utils.GetLogWriter())
	}
	go func() {
		for range routers.API.RestartChan {
			p.StopHTTP()
			utils.ReloadConf()
			p.StartHTTP()
		}
	}()
	return
}

func (p *program) Stop(s service.Service) (err error) {
	defer log.Println("********** STOP **********")
	defer utils.CloseLogWriter()
	p.StopHTTP()
	models.Close()
	return
}
func main() {
	log.SetPrefix("[OneSignaling] ")
	log.SetFlags(log.LstdFlags)
	if utils.Debug {
		log.SetFlags(log.Lshortfile | log.LstdFlags)
	}
	sec := utils.Conf().Section("service")
	svcConfig := &service.Config{
		Name:        sec.Key("name").MustString("OneSignaling_Service"),
		DisplayName: sec.Key("display_name").MustString("OneSignaling_Service"),
		Description: sec.Key("description").MustString("OneSignaling_Service"),
	}

	httpPort := utils.Conf().Section("http").Key("port").MustInt(50008)
	p := &program{
		httpPort: httpPort,
	}
	var server, err = service.New(p, svcConfig)
	if err != nil {
		log.Println(err)
		utils.PauseExit()
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "install" || os.Args[1] == "stop" {
			figure.NewFigure("OneSignaling", "", false).Print()
		}
		log.Println(svcConfig.Name, os.Args[1], "...")
		if err = service.Control(server, os.Args[1]); err != nil {
			log.Println(err)
			PauseExit()
		}
		log.Println(svcConfig.Name, os.Args[1], "ok")
		return
	}
	figure.NewFigure("OneSignaling", "", false).Print()
	if err = server.Run(); err != nil {
		log.Println(err)
		PauseExit()
	}
}

func PauseExit() {
	log.Println("Press any to exit")
	keyboard.GetSingleKey()
	os.Exit(0)
}
