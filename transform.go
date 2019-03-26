package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

//ReadSource get source from source.yml
func ReadSource() Source {
	raw, err := ioutil.ReadFile("source.yml")
	if err != nil {
		fmt.Println(err)
	}
	var source Source
	err = yaml.Unmarshal(raw, source)
	Checkerr(err)
	return source
}

//GetSurgeConf download providers' surgeconf
func GetSurgeConf(provider string) string {
	client := &http.Client{}
	request, err := http.NewRequest("GET", provider, nil)
	request.Header.Add("User-Agent", "Surge/1166 CFNetwork/955.1.2 Darwin/18.0.0")
	Checkerr(err)
	respon, err := client.Do(request)
	defer respon.Body.Close()
	if err != nil {
		fmt.Println("获取托管失败", err)
		result.Network = append(result.Network, provider)
	}
	body, err := ioutil.ReadAll(respon.Body)
	Checkerr(err)
	return string(body)
}

// GetSurgeProxies from SurgeConf
func GetSurgeProxies(conf string) []string {
	re, err := regexp.Compile("\\[Proxy\\]([\\s\\S]*?)\\[Proxy Group\\]")
	if err == nil {
		submatch := re.FindSubmatch([]byte(conf))
		if len(submatch) == 2 {
			return strings.Split(string(submatch[1]), "\n")
		}
		return nil
	}
	return nil
}

//FormatProxy Get proxy from surge subscribe
func FormatProxy(surge string) Proxy {
	ssregexp, _ := regexp.Compile("(.*?)\\s*=\\s*custom,(.*?),(.*?),(.*?),(.*?),") //找到所有节点信息,滤出DIRECT和格式不规范的信息
	obfsHostRegexp, _ := regexp.Compile("obfs-host\\s*=\\s*(.*?)(?:,|$)")
	obfsTypeRegexp, _ := regexp.Compile("obfs\\s*=\\s*(.*?)(?:,|$)")
	var proxy Proxy
	surgeproxy := ssregexp.FindSubmatch([]byte(surge))
	if len(surgeproxy) == 6 {
		proxy.Server = strings.TrimSpace(string(surgeproxy[2]))
		proxy.Port, _ = strconv.Atoi(strings.TrimSpace(string(surgeproxy[3])))
		proxy.Password = strings.TrimSpace(string(surgeproxy[5]))
		proxy.Cipher = strings.TrimSpace(string(surgeproxy[4]))
		proxy.Name = strings.TrimSpace(string(surgeproxy[1]))
		proxy.Type = "ss"
		obfsType := obfsTypeRegexp.FindSubmatch([]byte(surge))
		if len(obfsType) == 2 {
			proxy.Obfs = strings.TrimSpace(string(obfsType[1]))
			obfsHost := obfsHostRegexp.FindSubmatch([]byte(surge))
			if len(obfsHost) == 2 {
				proxy.ObfsHost = strings.TrimSpace(string(obfsHost[1]))
			}
		}
	}
	return proxy
}
