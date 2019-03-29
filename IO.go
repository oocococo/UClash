package main

//Source surge subscribe and Group information
type Source struct {
	Providers []string `yaml:"providers"`
	Grouplist []GroupList
}

//GroupList
type GroupList struct {
	Name     string `yaml:"name"`
	Type     string
	URL      string
	Interval int
	Abandon  string
	Demand   string
}

//Config Clash config format
type Config struct {
	Port      int  `yaml:"port"`
	SocksProt int  `yaml:"socks-port"`
	RedirPort int  `yaml:"redir-port"`
	AllowLan  bool `yaml:"allow-lan"`
	Mode      string
	LogLevel  string `yaml:"log-level"`
	External  string `yaml:"external-controller"`
	Secret    string `yaml:"secret"`
	DNS       struct {
		Enable       bool `yaml:"enable"`
		Ipv6         bool
		Listen       string
		Enhancedmode string `yaml:"enhanced-mode"`
		Nameserver   []string
		Fallback     []string
	}
	Proxy      []Proxy  `yaml:"Proxy"`
	ProxyGroup []Group  `yaml:"Proxy Group"`
	Rule       []string `yaml:"Rule"`
}

//Group Clash Proxy Group
type Group struct {
	Name     string `yaml:"name"`
	Type     string
	Proxies  []string
	URL      string
	Interval int
}

//Proxy clash proxy type
type Proxy struct {
	Type     string `yaml:"type"`
	Server   string `yaml:"server"`
	Port     int
	Password string
	Cipher   string
	Name     string
	Obfs     string
	ObfsHost string `yaml:"obfs-host"`
}

// Reslist of convertion
type Reslist struct {
	Sucess  []string
	Fromat  []string
	Network []string
}

var result Reslist
