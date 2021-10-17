package main

import (
	"NhenDownloader/spider"
	"fmt"
	"os"
	"strings"

	"github.com/tianyagk/CliToolkit"
)

var config map[string]string = make(map[string]string)

func main() {
	// Init and Setup Command Client with Function Mapper
	CommandClient := CliToolkit.Command{
		Use:    "NHentai Downloader",
		Intro:  "NHentai Downloader",
		Short:  "Nhentai manga downloader, entry help for more information",
		Long:   "long:",
		Prompt: ">> ",
	}

	// init params values
	config["proxies"] = "http://localhost:7890"
	config["language"] = "chinese"
	config["maxTryTimes"] = "5"
	config["maxOccurs"] = "30"

	err := os.Mkdir("./galleries/", os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	FuncMap := make(map[string]CliToolkit.Event)
	FuncMap["recent"] = CliToolkit.Event{DoFunc: doRecent, Description: "Show recent popular manga", Flag: "-r", ErrorHandler: CliToolkit.DefaultErrorHandler}
	FuncMap["download"] = CliToolkit.Event{DoFunc: doDownloadByID, Description: "Download manga by id", Flag: "-d", ErrorHandler: CliToolkit.DefaultErrorHandler}

	FuncMap["proxy"] = CliToolkit.Event{DoFunc: setProxy, Description: "Setting proxy address", Flag: "-p", ErrorHandler: CliToolkit.DefaultErrorHandler}
	FuncMap["lang"] = CliToolkit.Event{DoFunc: setLang, Description: "Setting default language", Flag: "-l", ErrorHandler: CliToolkit.DefaultErrorHandler}

	CommandClient.FuncMap = FuncMap
	CommandClient.Run()
}

// Define your command func here
func setProxy(str string, _ CliToolkit.Command) error {
	// https://ip.ihuan.me  小幻HTTP在线代理
	config["proxies"] = strings.TrimSuffix(str, "\r")
	return nil
}

func setLang(str string, _ CliToolkit.Command) error {
	config["language"] = str
	return nil
}

func doRecent(_ string, _ CliToolkit.Command) error {
	err := spider.Recent(config)
	return err
}

func doDownloadByID(id string, _ CliToolkit.Command) error {
	err := spider.DownloadByID(config, id)
	if err != nil {
		return err
	}
	return nil
}
