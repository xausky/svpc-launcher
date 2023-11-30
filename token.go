package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
)

func ParserTokenUrl(tokenUrl string) (url.Values, error) {
	parse, err := url.Parse(tokenUrl)
	if err != nil {
		return nil, err
	}
	query, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return nil, err
	}
	if !query.Has("token") {
		return nil, errors.New("must has token params")
	}
	return query, nil
}

func InputTokenUrl() string {
	var tokenUrl string
	for {
		fmt.Println("登录信息无效，请重新输入授权地址，获取教程：https://blog.xausky.cn")
		reader := bufio.NewReaderSize(os.Stdin, 65536)
		tokenUrlBytes, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("授权地址有误请严格按照教程执行。", err)
			continue
		}
		tokenUrl = string(tokenUrlBytes)
		_, err = ParserTokenUrl(tokenUrl)
		if err != nil {
			fmt.Println("授权地址有误请严格按照教程执行。", err)
			continue
		} else {
			break
		}
	}
	return tokenUrl
}

func LoadToken(clean bool) url.Values {
	if clean {
		err := os.Remove("svpc-launcher-token.txt")
		if err != nil {
			panic(err)
		}
	}
	tokenUrlBytes, err := os.ReadFile("svpc-launcher-token.txt")
	var tokenUrl string
	if os.IsNotExist(err) {
		tokenUrl = InputTokenUrl()
		err = os.WriteFile("svpc-launcher-token.txt", []byte(tokenUrl), 0644)
		if err != nil {
			panic(err)
		}
	}
	if err != nil {
		panic(err)
	}
	tokenUrl = string(tokenUrlBytes)
	tokenUrlParams, err := ParserTokenUrl(tokenUrl)
	if err != nil {
		tokenUrl = InputTokenUrl()
		tokenUrlParams, _ = ParserTokenUrl(tokenUrl)
		err = os.WriteFile("svpc-launcher-token.txt", []byte(tokenUrl), 0644)
		if err != nil {
			panic(err)
		}
	}
	return tokenUrlParams
}
