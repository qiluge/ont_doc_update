package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestLinkReplace(t *testing.T) {
	fileContent, _ := ioutil.ReadFile("testnet.md")
	//fileContent:=[]byte("English | [中文](testnet_CN.md)")
	linkReg := regexp.MustCompile(`\[.*?\]\([^#].*.*?\)`)
	result := linkReg.FindAll(fileContent, -1)
	fmt.Printf("%q\n", result)
}

func TestEncoder(t *testing.T) {
	name := "1234百万红包"
	words := ([]rune)(name)
	fmt.Printf("Before :%s\n", string(words))
	afterWords := strings.Replace(string(words), string([]rune("百万")), string([]rune("两个")), -1)
	fmt.Printf("After :%s\n", string(afterWords))
}

func TestDownload(t *testing.T) {
	url := "https://raw.githubusercontent.com/ontio/ontology-DID/master/docs/cn/ONTID_protocol_spec_cn.md"
	result, _ := download(url)
	ioutil.WriteFile("test.md", result, 0644)
}

func TestCmd(t *testing.T) {
	cmd := exec.Command("mv", "README.md", "temp.md")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	fmt.Println(err)
	fmt.Print(out.String())
}
