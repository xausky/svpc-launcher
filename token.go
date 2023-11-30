package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func ParserTokenUrl(tokenUrl string) (url.Values, error) {
	if !strings.HasPrefix(tokenUrl, "https") {
		return nil, errors.New("地址应该以 https 开头")
	}
	parse, err := url.Parse(tokenUrl)
	if err != nil {
		return nil, err
	}
	query, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return nil, err
	}
	if !query.Has("token") {
		return nil, errors.New("地址必须有 token 参数，请完整复制")
	}
	return query, nil
}

func InputTokenUrl() string {
	var tokenUrl string
	reader := bufio.NewReaderSize(os.Stdin, 65536)
	for {
		fmt.Println("登录信息无效，请重新输入授权地址，获取教程：https://blog.xausky.cn")
		tokenUrlBytes, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("授权地址有误请严格按照教程执行。", err)
			continue
		}
		tokenUrl = strings.TrimSpace(string(tokenUrlBytes))
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
	file := os.ExpandEnv("$APPDATA\\svpc-launcher-token.txt")
	if clean {
		err := os.Remove(file)
		if err != nil {
			panic(err)
		}
	}
	tokenUrlBytes, err := os.ReadFile(file)
	var tokenUrl string
	if os.IsNotExist(err) {
		tokenUrl = InputTokenUrl()
		err = os.WriteFile(file, []byte(tokenUrl), 0644)
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
		err = os.WriteFile(file, []byte(tokenUrl), 0644)
		if err != nil {
			panic(err)
		}
	}
	tokenUrlParams.Del("gv")
	tokenUrlParams.Del("cv")
	return tokenUrlParams
}
