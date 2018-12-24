package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Conf file named update.json
type Conf struct {
	Providers []string `json:"providers"`
	Filter    []string `json:"filter"`
	FilterOut []string `json:"filterout"`
}

// Server type
type Server struct {
	Method     string `json:"method"`
	Password   string `json:"password"`
	Plugin     string `json:"plugin"`
	PluginArgs string `json:"plugin_args"`
	PluginOpts string `json:"plugin_opts"`
	Remarks    string `json:"remarks"`
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Timeout    int    `json:"timeout"`
}

// SSGui GUI json
type SSGui struct {
	AutoCheckUpdate        bool     `json:"autoCheckUpdate"`
	AvailabilityStatistics bool     `json:"availabilityStatistics"`
	CheckPreRelease        bool     `json:"checkPreRelease"`
	Configs                []Server `json:"configs"`
	Enabled                bool     `json:"enabled"`
	Global                 bool     `json:"global"`
	Hotkey                 struct {
		RegHotkeysAtStartup   bool   `json:"RegHotkeysAtStartup"`
		ServerMoveDown        string `json:"ServerMoveDown"`
		ServerMoveUp          string `json:"ServerMoveUp"`
		ShowLogs              string `json:"ShowLogs"`
		SwitchAllowLan        string `json:"SwitchAllowLan"`
		SwitchSystemProxy     string `json:"SwitchSystemProxy"`
		SwitchSystemProxyMode string `json:"SwitchSystemProxyMode"`
	} `json:"hotkey"`
	Index            int  `json:"index"`
	IsDefault        bool `json:"isDefault"`
	IsVerboseLogging bool `json:"isVerboseLogging"`
	LocalPort        int  `json:"localPort"`
	LogViewer        struct {
		BackgroundColor string `json:"BackgroundColor"`
		Font            string `json:"Font"`
		TextColor       string `json:"TextColor"`
		ToolbarShown    bool   `json:"toolbarShown"`
		TopMost         bool   `json:"topMost"`
		WrapText        bool   `json:"wrapText"`
	} `json:"logViewer"`
	PacURL       string `json:"pacUrl"`
	PortableMode bool   `json:"portableMode"`
	Proxy        struct {
		ProxyPort    int    `json:"proxyPort"`
		ProxyServer  string `json:"proxyServer"`
		ProxyTimeout int    `json:"proxyTimeout"`
		ProxyType    int    `json:"proxyType"`
		UseProxy     bool   `json:"useProxy"`
	} `json:"proxy"`
	SecureLocalPac bool        `json:"secureLocalPac"`
	ShareOverLan   bool        `json:"shareOverLan"`
	Strategy       interface{} `json:"strategy"`
	UseOnlinePac   bool        `json:"useOnlinePac"`
}

// ReadConf read update.json
func ReadConf() Conf {
	cb, err := ioutil.ReadFile("update.json")
	if err != nil {
		fmt.Println(err)
	}
	var conf Conf
	json.Unmarshal(cb, &conf)
	return conf
}

// ReadSSGui read gui json
func ReadSSGui() SSGui {
	cb, err := ioutil.ReadFile("gui-config.json")
	if err != nil {
		fmt.Println(err)
	}
	var gui SSGui
	json.Unmarshal(cb, &gui)
	return gui
}

// SurgeFromConf match surge urls
func SurgeFromConf(conf string) []string {
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

// Surge2SS convert surge style url to ss-gui format
func Surge2SS(surge string) Server {
	regex, _ := regexp.Compile("(.*?)\\s*=\\s*custom,(.*?),(.*?),(.*?),(.*?),")
	obfsRegex, _ := regexp.Compile("obfs-host\\s*=\\s*(.*?)(?:,|$)")
	obfsTypeRegex, _ := regexp.Compile("obfs\\s*=\\s*(.*?)(?:,|$)")
	var res Server
	params := regex.FindSubmatch([]byte(surge))
	if len(params) == 6 {
		res.Server = strings.TrimSpace(string(params[2]))
		res.ServerPort, _ = strconv.Atoi(strings.TrimSpace(string(params[3])))
		res.Password = strings.TrimSpace(string(params[5]))
		res.Method = strings.TrimSpace(string(params[4]))
		res.Remarks = strings.TrimSpace(string(params[1]))
		res.Timeout = 5
		obfsType := obfsTypeRegex.FindSubmatch([]byte(surge))
		if len(obfsType) == 2 {
			res.Plugin = "obfs-local"
			res.PluginOpts = "obfs=" + strings.TrimSpace(string(obfsType[1]))
			obfsParams := obfsRegex.FindSubmatch([]byte(surge))
			if len(obfsParams) == 2 {
				res.PluginOpts += ";obfs-host=" + strings.TrimSpace(string(obfsParams[1]))
			}
		}
	}
	return res
}

// Result of convertion
type Result struct {
	Success []string
	Fromat  []string
	Network []string
}

func main() {
	conf := ReadConf()
	gui := ReadSSGui()
	providers := conf.Providers
	fmt.Println(fmt.Sprintf("成功读取到%d个托管配置，开始下载...", len(providers)))
	filters := conf.Filter
	filterouts := conf.FilterOut
	fmt.Println(fmt.Sprintf("接受关键字：%s", strings.Join(filters, " | ")))
	fmt.Println(fmt.Sprintf("不接受关键字： %s", strings.Join(filterouts, " | ")))
	var result Result
	var remotes []string
	var wg sync.WaitGroup
	wg.Add(len(providers))
	for i := 0; i < len(providers); i++ {
		go func(url string) {
			defer wg.Done()
			client := &http.Client{}
			request, err := http.NewRequest("GET", url, nil)
			request.Header.Add("User-Agent", "Surge/1166 CFNetwork/955.1.2 Darwin/18.0.0")
			if err != nil {
				fmt.Println(err)
			}
			resp, err := client.Do(request)
			if err != nil {
				fmt.Println("获取托管失败", err)
				result.Network = append(result.Network, url)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			remotes = append(remotes, string(body))
		}(providers[i])
	}
	wg.Wait()
	var servers []Server
	for k := 0; k < len(remotes); k++ {
		urls := SurgeFromConf(remotes[k])
		if urls == nil {
			result.Fromat = append(result.Fromat, providers[k])
			continue
		}
		result.Success = append(result.Success, providers[k])
		for i := 0; i < len(urls); i++ {
			res := Surge2SS(urls[i])
			if res.Remarks != "" {
				if (len(filters) <= 0 || filters == nil) && (len(filterouts) <= 0 || filterouts == nil) {
					servers = append(servers, res)
					continue
				}
				for out := 0; out < len(filterouts); out++ {
					if on, _ := regexp.MatchString(filterouts[out], res.Remarks); on {
						goto FILTEROUTIT
					}
				}
				for j := 0; j < len(filters); j++ {
					if m, _ := regexp.MatchString(filters[j], res.Remarks); m {
						servers = append(servers, res)
						break
					}
				}
			FILTEROUTIT:
			}
		}
	}
	fmt.Println(fmt.Sprintf("\n----------------\n成功获取：\n - %s\n格式错误：\n - %s\n网络错误：\n - %s\n----------------\n", strings.Join(result.Success, "\n - "), strings.Join(result.Fromat, "\n - "), strings.Join(result.Network, "\n - ")))
	gui.Configs = servers
	outputJSON, _ := json.Marshal(gui)
	writeFileErr := ioutil.WriteFile("gui-config.json", outputJSON, 0644)
	if writeFileErr == nil {
		fmt.Println(fmt.Sprintf("服务器更新完毕，合计更新%d个节点", len(servers)))
		fmt.Println("请重启Shadowsocks客户端或进入节点列表点击确定")
	} else {
		fmt.Println("配置文件写入失败")
	}
}
