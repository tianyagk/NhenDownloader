package nhen

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
)

func init_client(conf map[string]string) *http.Client {
	//adding the proxy settings to the Transport object
	proxyStr := conf["proxies"]
	transport := &http.Transport{}
	if proxyStr != "" {

		//creating the proxyURL
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			fmt.Println(err)
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}
	return client
}

func getParseHTML(url string, client *http.Client) soup.Root {
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
	doc := soup.HTMLParse(string(body))
	return doc
}

func saveImage(url string, rootpath string, id int, client *http.Client) error {
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

	idstr := fmt.Sprintf("%04d", id)
	err = ioutil.WriteFile(rootpath+idstr+".jpg", body, 0755)
	if err != nil {
		fmt.Printf("ioutil.WriteFile failure, err=[%v]\n", err)
	}
	return err
}

func trySaveImage(tryNum int, url string, rootpath string, id int, client *http.Client) {
	for idx := 0; idx < tryNum; idx++ {
		err := saveImage(url, rootpath, id, client)
		if err == nil {
			return
		}
	}
	fmt.Println("Connection timed out, exceeding the maximum number of attempts " + fmt.Sprintf("%d", tryNum))
}

func Recent(config *map[string]string) {
	conf := *config
	client := init_client(conf)

	url := "https://nhentai.net/language/" + conf["language"]

	doc := getParseHTML(url, client)
	captions := doc.Find("div", "id", "content").FindAll("div", "class", "caption")
	for ids, caption := range captions {
		fmt.Println(ids, "|", caption.Text(), "|")
	}
}

func DownloadByID(config *map[string]string, id string) {
	conf := *config
	client := init_client(conf)

	url := "https://nhentai.net/g/" + id
	doc := getParseHTML(url, client)

	cover_url := doc.Find("div", "id", "cover").Find("noscript").Text()
	galleries := strings.Split(cover_url, "/")[4]
	// fmt.Println(galleries)

	name_before := doc.Find("h1", "class", "title").Find("span", "class", "before").Text()
	name_pretty := doc.Find("h1", "class", "title").Find("span", "class", "pretty").Text()
	name_after := doc.Find("h1", "class", "title").Find("span", "class", "after").Text()

	manga_name := name_before + name_pretty + name_after

	page := doc.Find("section", "id", "tags").FindAll("div")[7].Find("a", "class", "tag").Find("span", "class", "name").Text()
	pageNum, _ := strconv.Atoi(page)

	// make dir for download
	err := os.Mkdir("./galleries/"+strings.ReplaceAll(manga_name, " ", "_"), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("- ", strings.ReplaceAll(manga_name, " ", "_"), " in download queue.")

	for idx := 1; idx <= pageNum; idx++ {
		img_url := "https://i.nhentai.net/galleries/" + galleries + "/" + fmt.Sprintf("%d", idx) + ".jpg"

		// overwrite
		go trySaveImage(5, img_url, "./galleries/"+strings.ReplaceAll(manga_name, " ", "_")+"/", idx, client)
	}
}
