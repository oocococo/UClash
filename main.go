package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

func main() {
	var config Config
	var wg sync.WaitGroup

	rawconfig, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Println("can't read the config.yml", err)
	}
	err = yaml.Unmarshal(rawconfig, &config)
	if err != nil {
		fmt.Println("配置格式错误", err)
	}
	config.ProxyGroup = nil
	config.Proxy = nil
	var proxiesname []string
	//read surge sub and transform into proxy
	source := ReadSource()
	for i := 0; i < len(source.Providers); i++ {
		surgeproxies := GetSurgeProxies(GetSurgeConf(source.Providers[i]))
		wg.Add(len(surgeproxies))
		for i := 0; i < len(surgeproxies); i++ {
			go func(rawproxy string) {
				defer wg.Done()
				if rawconfig != nil {
					newproxy := FormatProxy(rawproxy)
					config.Proxy = append(config.Proxy, newproxy)
					proxiesname = append(proxiesname, newproxy.Name)

				}
				//fmt.Println("成功识别第", i, "个节点")
			}(surgeproxies[i])
		}
		wg.Wait()
		fmt.Println("成功读取所有节点")
	}
	//generate group
	for g := 0; g < len(source.Grouplist); g++ {
		var black int
		var white int
		var autogroup Group
		var afterdemand []string
		needs := strings.Split(source.Grouplist[g].Demand, ",")
		fmt.Println("允许关键词", needs)
		trash := strings.Split(source.Grouplist[g].Abandon, ",")
		fmt.Println("不允许关键词", trash)
		autogroup.Name = source.Grouplist[g].Name
		autogroup.Interval = source.Grouplist[g].Interval
		autogroup.Type = source.Grouplist[g].Type
		autogroup.URL = source.Grouplist[g].URL
		for p := 0; p < len(proxiesname); p++ {
			if needs == nil {
				afterdemand = proxiesname
				break
			} else {
				for n := 0; n < len(needs); n++ {
					if need, _ := regexp.MatchString(needs[n], proxiesname[p]); need {
						afterdemand = append(afterdemand, proxiesname[p])
						white += 1
					}
				}
			}
		}
		//fmt.Print(afterdemand)
		for p := 0; p < len(afterdemand); p++ {
			if trash[0] == "" {
				autogroup.Proxies = afterdemand
				black = white
				break
			} else {
				for a := 0; a < len(trash); a++ {
					//fmt.Println(trash[a], "!!!")
					if neednt, _ := regexp.MatchString(trash[a], afterdemand[p]); !neednt {
						autogroup.Proxies = append(autogroup.Proxies, afterdemand[p])
						black += 1
					}
				}
			}
		}
		fmt.Println(autogroup.Name, "组更新了", black, "个节点")
		config.ProxyGroup = append(config.ProxyGroup, autogroup)
	}

	outputconfig, err := yaml.Marshal(config)
	clasherr := ioutil.WriteFile("config.yml", outputconfig, 0644)
	Checkerr(clasherr)
}
