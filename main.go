package main

import (
	"flag"
	"io/ioutil"
	"log"
)

var tokenPath = flag.String("token", "./token.txt", "token path")
var savePath = flag.String("path", "./albums", "albums save path")

func main() {
	log.Printf("读取配置 \n")
	token, err := readToken(*tokenPath)
	if err != nil {
		log.Fatalf("读取token失败 path: %s, err: %s", *tokenPath, err)
	}

	log.Printf("开始任务 \n")

	downloader := NewDownloader(*savePath, token)
	if err := downloader.Run(); err != nil {
		log.Fatalf("运行失败：%s \n", err)
	}

	log.Printf("任务完成 \n")
}

func readToken(path string) (string, error) {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
