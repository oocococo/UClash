package main

import (
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	var config Config
	rawconfig, err := ioutil.ReadFile("config.yml")
	Checkerr(err)
	err = yaml.Unmarshal(rawconfig, config)
	Checkerr(err)
	var proxiesname []string
	//read surge sub and transform into proxy
	source := ReadSource()
	for i := 0; i < len(source.Providers); i++ {
		surgeproxies := GetSurgeProxies(GetSurgeConf(source.Providers[i]))
		for i := 0; i < len(surgeproxies); i++ {
			go func(rawproxy string) {
				newproxy := FormatProxy(rawproxy)
				config.Proxy = append(config.Proxy, newproxy)
				proxiesname = append(proxiesname, newproxy.Name)
			}(surgeproxies[i])
		}
	}
	//generate group

	for g := 0; g < len(source.Grouplist); g++ {
		go func(wantedlist GroupList) {
			var autogroup Group
			var afterdemand []string
			needs := strings.Split(wantedlist.Demand, ",")
			trash := strings.Split(wantedlist.Abandon, ",")
			autogroup.Name = wantedlist.Name
			autogroup.Interval = wantedlist.Interval
			autogroup.Type = wantedlist.Type
			autogroup.URL = wantedlist.URL
			for p := 0; p < len(proxiesname); p++ {
				if needs == nil {
					afterdemand = proxiesname
					break
				} else {
					for n := 0; n < len(needs); n++ {
						if need, _ := regexp.MatchString(needs[n], proxiesname[p]); need {
							afterdemand = append(afterdemand, proxiesname[p])
						}
					}
				}
			}

			for p := 0; p < len(afterdemand); p++ {
				if trash == nil {
					autogroup.Proxies = afterdemand
					break
				} else {
					for a := 0; a < len(trash); a++ {
						if neednt, _ := regexp.MatchString(trash[a], afterdemand[p]); !neednt {
							autogroup.Proxies = append(autogroup.Proxies, afterdemand[p])
						}
					}
				}
			}
		}(source.Grouplist[g])
	}
	outputconfig, err := yaml.Marshal(config)
	clasherr := ioutil.WriteFile("config.yml", outputconfig, 0644)
	Checkerr(clasherr)
}
