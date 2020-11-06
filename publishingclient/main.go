package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/phamvinhdat/httpclient"
	"github.com/phamvinhdat/tribe/pkg/config"
	"github.com/phamvinhdat/tribe/pkg/try"
	"github.com/phamvinhdat/tribe/publishingclient/msgservice"
	"github.com/phamvinhdat/tribe/publishingclient/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// load config
	config.LoadConfig("./config")

	// http client
	httpClient := httpclient.NewClient()

	// tryer
	timeInterval := viper.GetDuration("messageserver.interval")
	if timeInterval <= 0 {
		panic("time interval is invalid")
	}
	tryer := try.New(
		try.WithTimeout(time.Hour*24*365), // a year
		try.WithInterval(timeInterval),
	)

	// service
	msgServerURL := viper.GetString("messageserver.url")
	if len(msgServerURL) == 0 {
		panic("message server url invalid")
	}
	msgService := msgservice.New(msgServerURL, httpClient)
	service.New(tryer, msgService).Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	close(quit)
	logrus.Info("shutting down")
}
