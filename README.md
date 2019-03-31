###### 本程序可以从surge订阅链接中获取节点,并按照source.yml文件的配置将节点和节点组写入config.yml文件中,同时保留原有的general,dns,rule部分

*source.yml配置说明*
```yml
providers: 
  - link1
  - link2
grouplist:
  - name: lb
    type: load-balance
    demand: "专线,香港,台湾"
    abandon: "游戏"
    url: http://www.gstatic.com/generate_204
    interval: 300
  - name: Netflix
    type: url-test
    demand: "netflix"
    url: http://www.gstatic.com/generate_204
    interval: 300
  - Proxy
    type: select
```

##### before
```yml
port: 7890
socks-port: 7891
allow-lan: false
mode: Rule
log-level: info
external-controller: 127.0.0.1:9090
 dns: your dns
Proxy:
whatever
Proxy Group:
whatever
Rule:
- your rule
```
##### after
```yml
port: 7890
socks-port: 7891
allow-lan: false
mode: Rule
log-level: info
external-controller: 127.0.0.1:9090
 dns: your dns
Proxy:
all proxies from surge subscribe
Proxy Group:
filtered proxies and groups as you wish

Rule:
- your rule
```
