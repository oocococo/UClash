package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	var config Config
	config.DNS.Enhancedmode = "fake-ip"
	config.Experimental.IgnoreResolveFail = true
	//var wg sync.WaitGroup
	var proxiesnumber int
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
		surgeproxies := GetSurgeProxies(GetSurgeConf(source.Providers[i]), source.Providers[i])
		for s := 0; s < len(surgeproxies); s++ {
			newproxy := FormatProxy(surgeproxies[s])
			if newproxy.Name != "" {
				newproxy.Udp = true
				config.Proxy = append(config.Proxy, newproxy)
				proxiesname = append(proxiesname, newproxy.Name)
				proxiesnumber += 1
			}
		}

	}
	fmt.Println("网络错误:", result.Network)
	fmt.Println("格式错误:", result.Fromat)
	fmt.Println("共读取", proxiesnumber, "个节点")
	//generate group
	for g := 0; g < len(source.Grouplist); g++ {
		var black int
		var white int
		var autogroup Group
		var afterdemand []string
		autogroup.Name = source.Grouplist[g].Name
		autogroup.Interval = source.Grouplist[g].Interval
		autogroup.Type = source.Grouplist[g].Type
		autogroup.URL = source.Grouplist[g].URL
		fmt.Print(autogroup.Name, "组,")
		needs := strings.Split(source.Grouplist[g].Demand, ",")
		fmt.Print("允许关键词", needs, ",")
		trash := strings.Split(source.Grouplist[g].Abandon, ",")
		fmt.Print("不允许关键词", trash, ",")
		for p := 0; p < len(proxiesname); p++ {
			if needs[0] == "" {
				afterdemand = proxiesname
				white = len(proxiesname)
				break
			} else {
				for n := 0; n < len(needs); n++ {
					if need, _ := regexp.MatchString(needs[n], proxiesname[p]); need {
						afterdemand = append(afterdemand, proxiesname[p])
						white += 1
						break
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
				flag := 0
				for a := 0; a < len(trash); a++ {
					//fmt.Println(trash[a], "!!!")
					if neednt, _ := regexp.MatchString(trash[a], afterdemand[p]); neednt {

						flag = 1

						break
					}
				}
				if flag == 0 {
					autogroup.Proxies = append(autogroup.Proxies, afterdemand[p])
					black += 1
				}
			}
		}
		fmt.Println("共更新了", black, "个节点")
		config.ProxyGroup = append(config.ProxyGroup, autogroup)
	}

	outputconfig, err := yaml.Marshal(config)
	clasherr := ioutil.WriteFile("config.yml", outputconfig, 0644)
	Checkerr(clasherr)
}
