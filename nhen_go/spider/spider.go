package spider

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func initClient(conf map[string]string) (*http.Client, error) {
	//adding the proxy settings to the Transport object
	proxyStr := conf["proxies"]
	transport := &http.Transport{}
	if proxyStr != "" {
		//creating the proxyURL
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			return nil, err
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}
	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}
	return client, nil
}

func parseHTML(url string, client *http.Client) (soup.Root, error) {
	// Get Parsed HTML String from url
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36 Edg/94.0.992.37")
	resp, _ := client.Do(req)
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return soup.HTMLParse(""), err
	}
	doc := soup.HTMLParse(string(body))
	return doc, nil
}

func saveImage(url string, root string, id int, client *http.Client) error {
	// Get Parsed HTML String from url
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36 Edg/94.0.992.37")
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	if len(body)<= 600 {
		return errors.New("download fail")
	} else {
		err = ioutil.WriteFile(root+fmt.Sprintf("%04d", id)+".jpg", body, 0755)
		if err != nil {
			fmt.Printf("ioutil.WriteFile failure, err=[%v]\n", err)
		}
		return err
	}
}

func trySaveImage(tryNum int, url string, mangaName string, id int, client *http.Client, ch chan int) {
	path := "./galleries/"+mangaName+"/"
	ch <- 0
	for idx := 0; idx < tryNum; idx++ {
		err := saveImage(url, path, id, client)
		if err == nil {
			_ = <-ch
			if len(ch)==0 {
				integrityCheck(mangaName)
			}
			return
		} else {
			fmt.Println(err)
		}
	}
	_ = <-ch
	fmt.Println("Connection timed out, exceeding the maximum try times " + fmt.Sprintf("%d", tryNum))
}

func integrityCheck(mangaName string) {
	fmt.Println("download finish",mangaName,".\n", ">> ")
}


func Recent(conf map[string]string) error {
	jsonFile, err := os.Open("id2tag.json")
	byteFile, _ := ioutil.ReadAll(jsonFile)
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)

	if err != nil {
		return err
	}
	var id2tag map[int]string
	err = json.Unmarshal(byteFile, &id2tag)

	client, err := initClient(conf)
	if err != nil {
		return err
	}

	recentUrl := "https://nhentai.net/language/" + conf["language"]
	doc, err := parseHTML(recentUrl, client)
	if err != nil {
		return err
	}

	captions := doc.Find("div", "id", "content").FindAll("div", "class", "caption")
	gallery := doc.Find("div", "id", "content").FindAll("div", "class", "gallery")
	for ids, caption := range captions {
		fmt.Println(ids, "|", caption.Text(), "|")
		fmt.Print("	")
		for _, str := range strings.Fields(gallery[ids].Attrs()["data-tags"]) {
			tagId,_:= strconv.Atoi(str)
			if len(id2tag[tagId])>0 {
				fmt.Print(id2tag[tagId], ", ")
			}
		}
		fmt.Print("\n")
	}
	return err
}

func DownloadByID(conf map[string]string, id string) error {
	client, err := initClient(conf)
	if err != nil {
		return err
	}

	downloadUrl := "https://nhentai.net/g/" + id
	doc, err := parseHTML(downloadUrl, client)
	if err != nil {
		return err
	}

	coverUrl := doc.Find("div", "id", "cover").Find("noscript").Text()
	galleries := strings.Split(coverUrl, "/")[4]
	// fmt.Println(galleries)

	nameBefore := doc.Find("h1", "class", "title").Find("span", "class", "before").Text()
	namePretty := doc.Find("h1", "class", "title").Find("span", "class", "pretty").Text()
	nameAfter := doc.Find("h1", "class", "title").Find("span", "class", "after").Text()

	mangaName := nameBefore + namePretty + nameAfter
	mangaName = strings.ReplaceAll(mangaName, " ", "_")
	mangaName = strings.ReplaceAll(mangaName, "?", "")
	mangaName = strings.ReplaceAll(mangaName, "*", "")
	mangaName = strings.ReplaceAll(mangaName, ":", "")
	mangaName = strings.ReplaceAll(mangaName, "<", "")
	mangaName = strings.ReplaceAll(mangaName, ">", "")
	mangaName = strings.ReplaceAll(mangaName, "|", "")
	mangaName = strings.ReplaceAll(mangaName, "/", "")

	page := doc.Find("section", "id", "tags").FindAll("div")[7].Find("a", "class", "tag").Find("span", "class", "name").Text()
	pageNum, _ := strconv.Atoi(page)

	// make dir for download
	err = os.Mkdir("./galleries/"+mangaName, os.ModePerm)
	if err != nil {
		return err
	}

	// 使用 chan 控制并发图片下载
	fmt.Println("- ", mangaName, " in download queue.")
	maxOccurs, err := strconv.Atoi(conf["maxOccurs"])
	ch := make(chan int, maxOccurs)
	for idx := 1; idx <= pageNum; idx++ {
		imgUrl := "https://i.nhentai.net/galleries/" + galleries + "/" + fmt.Sprintf("%d", idx) + ".jpg"
		go trySaveImage(5, imgUrl, mangaName, idx, client, ch)
	}
	return err
}
