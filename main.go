package main

import (
	"encoding/json"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/ontio/ontology/common/log"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type item struct {
	OriginalLink string
	NewPostion   string
}

type DocMap []item

const DOC_REP_PATH = "../../ontio/documentation/"

var linkMap = make(map[string]string, 0)

func main() {
	log.InitLog(log.InfoLog, log.Stdout, log.PATH)

	// read doc map
	docMapFileContent, err := ioutil.ReadFile("doc-map.json")
	if err != nil {
		log.Errorf("read doc map err: %s!", err)
		return
	}
	var docMap DocMap
	err = json.Unmarshal(docMapFileContent, &docMap)
	if err != nil {
		log.Errorf("unmarshal doc map err: %s!", err)
		return
	}
	// read link map
	linkMapFileContent, err := ioutil.ReadFile("link-map.json")
	if err != nil {
		log.Errorf("read link map err: %s!", err)
		return
	}
	err = json.Unmarshal(linkMapFileContent, &linkMap)
	if err != nil {
		log.Errorf("unmarshal link map err: %s!", err)
		return
	}
	var wg = new(sync.WaitGroup)
	for _, v := range docMap {
		url := v.OriginalLink
		if !strings.Contains(url, "http") {
			continue
		}
		filePath := DOC_REP_PATH + v.NewPostion
		wg.Add(1)
		go handleFile(url, filePath, wg)
	}
	wg.Wait()
	jsonBytes, _ := json.Marshal(&linkMap)
	jsonString := string(jsonBytes)
	jsonString = strings.Replace(jsonString, ":\"", ":\n\"", -1)
	jsonString = strings.Replace(jsonString, ",", ",\n", -1)
	ioutil.WriteFile("link-map.json", []byte(jsonString), 0644)
}

func handleFile(url, filePath string, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Infof("start download: url is %s!", url)
	originContent, err := download(url)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("download success: url is %s, destination is %s.", url, filePath)
	newContent := handleContent(string(originContent), url, filePath)
	log.Info("handle content success:", filePath)
	err = writeToNewPosition([]byte(newContent), filePath)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("write to new position success:", filePath)
}

func download(url string) ([]byte, error) {
	client := http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		err = fmt.Errorf("download file connection err: %s, url is %s!", err, url)
		return []byte{}, err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("download file err: url is %s, response is %s!", url, resp.Status)
		return []byte{}, err
	}
	utf8Decoder := mahonia.NewDecoder("utf-8")
	decoderReader := utf8Decoder.NewReader(resp.Body)
	originContent, err := ioutil.ReadAll(decoderReader)
	if err != nil {
		err = fmt.Errorf("read respong body err: %s", err)
		return []byte{}, err
	}
	return originContent, nil
}

func handleContent(originContent, downloadUrl, filePath string) string {

	// get info from file path
	filePath = strings.Replace(filePath, DOC_REP_PATH, "", 1)
	splitPath := strings.Split(filePath, "/")
	firstDir := splitPath[2]  // doc_zh, doc_en
	secondDir := splitPath[3] // Ontology, DID, Dapp, Introduction, SDKs
	fileName := splitPath[4]

	// replace relative link
	linkMapKeyPrefix := strings.Replace(downloadUrl, "https://raw.githubusercontent.com/ontio/", "", 1)
	linkMapKeyPrefix = strings.Replace(linkMapKeyPrefix, "master/", "", 1)
	newContent := replaceLink(originContent, linkMapKeyPrefix)

	// handle title and version info
	newContent = handleTitleAndVersion(newContent)

	headerUrl := strings.Replace(downloadUrl, "raw.githubusercontent", "github", 1)
	headerUrl = strings.Replace(headerUrl, "/master/", "/blob/master/", 1)
	header := constructHeader(headerUrl, firstDir, secondDir, fileName)
	newContent = header + newContent
	return newContent
}

func replaceLink(originContent, prefix string) string {
	// remove godoc, go report card, travis, gitter link
	goDocLinkReg := regexp.MustCompile(`\[!.*\]\(.*\)`)
	goDocLinkRegReuslt := goDocLinkReg.FindAllString(originContent, -1)
	for _, v := range goDocLinkRegReuslt {
		originContent = strings.Replace(originContent, v, "", -1)
	}
	// extract link(relative path)
	linkReg := regexp.MustCompile(`\[.*?\]\([^#].*?[^html]\)`)
	result := linkReg.FindAllString(originContent, -1)

	for _, extractLink := range result {
		// extractLink is [xxxx](aaa.md)
		leftIndex := strings.Index(extractLink, "(")
		rightIndex := strings.Index(extractLink, ")")
		originLink := extractLink[leftIndex+1 : rightIndex]
		linkMapKey := prefix + extractLink
		if strings.HasPrefix(originLink, "http") || strings.Contains(originLink, "html") {
			resp, err := http.Get(originLink)
			if err != nil {
				log.Errorf("check link connection err: %s, url is %s!", err, originLink)
				continue
			}
			if resp.StatusCode == 404 {
				log.Errorf("check link 404 err, file is %s, link is %s!", prefix, originLink)
			}
			if resp.Body != nil {
				resp.Body.Close()
			}
			continue
		}
		if newLink, ok := linkMap[linkMapKey]; ok {
			originContent = strings.Replace(originContent, originLink, newLink, 1)
		} else {
			linkMap[linkMapKey] = ""
		}
	}
	return originContent
}

func handleTitleAndVersion(originContent string) string {

	versionReg := regexp.MustCompile(`(?m:^.*[v|V]ersion\s\d\.\d\.\d.*$)`)
	version := versionReg.FindAllString(originContent, 1)
	var versionInfo string
	if len(version) == 0 {
		versionInfo = "<p align=\"center\" class=\"version\">Version 1.0.0 </p>"
	}

	titleReg := regexp.MustCompile(`(?m:^#\s.*$)`)
	title := titleReg.FindAllString(originContent, 1)
	if len(title) == 1 {
		titleInfo := "<h1 align=\"center\">" + title[0][2:] + "</h1>"
		replaceString := titleInfo + "\n\n" + versionInfo + "\n\n"
		originContent = strings.Replace(originContent, title[0], replaceString, 1)
	}
	return originContent
}

func constructHeader(url, firstDir, secondDir, fileName string) string {
	var sidebar string
	if strings.Contains(firstDir, "zh") {
		sidebar = secondDir + "_zh"
	} else {
		sidebar = secondDir + "_en"
	}
	permalink := strings.Replace(fileName, ".md", ".html", 1)
	folder := firstDir + "/" + secondDir
	// add header
	header := "---\n"
	header += "title:\n"
	header += "keywords: sample homepage\n"
	header += "sidebar: " + sidebar + "\n"
	header += "permalink: " + permalink + "\n"
	header += "folder: " + folder + "\n"
	header += "giturl: " + url + "\n"
	header += "---\n\n"
	return header
}

func writeToNewPosition(fileContent []byte, filePath string) error {

	err := ioutil.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		return fmt.Errorf("writeToNewPosition, write to file %s err: %s, content length is %d",
			filePath, err, len(fileContent))
	}
	return nil
}
