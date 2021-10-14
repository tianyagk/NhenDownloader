package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"nhen_tool/nhen"
)

func main() {
	fmt.Println("Launch Nhentai-Downloader, entry 'help' for more information.")
	var config map[string]string = make(map[string]string)

	// init params
	config["proxies"] = "http://localhost:7890"
	config["language"] = "chinese"
	config["maxTryTimes"] = "5"

	for {
		fmt.Print(">> ")

		reader := bufio.NewReader(os.Stdin)
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		Command(cmdString, &config)
	}
}

func Command(cmdString string, config *map[string]string) {
	// drop suffix of command string
	cmdString = strings.TrimSuffix(cmdString, "\n")
	arrCommandStr := strings.Fields(cmdString)
	switch arrCommandStr[0] {
	case "help":
		fmt.Println("help			| Query help")
		fmt.Println("exit			| Exit NhenDownloader")
		fmt.Println("lang	-language	| Chose from ['chinese', 'japanese', 'english',etc.]")
		fmt.Println("proxy	-host:port	| Like 'http://localhost:7890'")
		fmt.Println("maxTryTimes	-tryNum	| Default max try num: 5")
		fmt.Println("recent			| Show recent popular manga with id")
		fmt.Println("download	-id	| 376478 -> 'https://nhentai.net/g/376478/'")
	case "lang":
		lang := arrCommandStr[1][1:]

		conf_ := *config
		conf_["language"] = lang
		*config = conf_

		fmt.Printf("Setting language: %s\n", lang)
	case "proxy":
		proxy := arrCommandStr[1][1:]

		conf_ := *config
		conf_["proxies"] = proxy
		*config = conf_

		fmt.Printf("Setting proxy: %s\n", proxy)
	case "exit":
		fmt.Println("Exit Downloader")
		os.Exit(0)
	case "recent":
		nhen.Recent(config)
	case "download":
		id := arrCommandStr[1][1:]
		nhen.DownloadByID(config, id)
	}
}
