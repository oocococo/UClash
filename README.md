providers:
  - https://良辰云.app/link/uRsVMAOGNUUg9FrJ?is_ss=1
  - https://www.cordcloud.cc/link/HVdeIfeRyzerwVij?is_ss=1
grouplist:
  - name: Proxy
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
    demand: "香港,台湾"
    abandon: "美国"
  - name: Netflix
    type: select
    demand: "台湾"
  - name: Spotify
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
    demand: "美国"
