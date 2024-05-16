# iKuai Bypass

将在线规则添加到爱快的`自定义运营商`和`域名分流`，并自动更新。目的请参见 [爱快路由器+Padavan实现无感科学的外出与回家](https://www.bilibili.com/read/cv4140884) 和 [爱快 & OpenWrt 分流拓扑，旁路由模式从此扔进垃圾堆](https://www.right.com.cn/forum/thread-8252571-1-1.html)

## 如何使用

1. 编写配置文件，命名为`config.yml`
2. 进入爱快自带 docker 中，点击`镜像管理`->`添加`，选择`镜像库下载`，搜索`ztc1997/ikuai-bypass`，下载`TAG`为`latest`的镜像
3. 点击`容器列表`->`添加`，`选择镜像文件`选择`ztc1997/ikuai-bypass:latest`，打开`高级选项`，添加一个`挂载目录`，`源路径`填写放置配置文件的路径，`目标路径`填写`/etc/ikuai-bypass`，内存64M即可，其它根据需要自行填写
4. 保存并启用

### 配置文件模板

```yaml
## 爱快管理页面的 URL，结尾不要加 "/"，
## 如不填写，则使用第一个接口的网关地址作为爱快地址，
## 如果在爱快自带 docker 运行，网关就是爱快地址，可以不写
# ikuai-url: http://192.168.1.1 
username: admin # 爱快用户名
password: pass  # 爱快密码
cron: 0 4 * * * # 执行更新的周期 crontab
custom-isp:     # 自定义运营商
  - name: 国内  # 自定义运营商名称
    ## 自定义运营商 cidr 列表网址，每行一个，超过5000行会自动分为多个，ipv6 地址会被删除
    url: https://cdn.jsdelivr.net/gh/Loyalsoldier/geoip@release/text/cn.txt
  - name: Telegram
    url: https://cdn.jsdelivr.net/gh/Loyalsoldier/geoip@release/text/telegram.txt
ip-group:     # IP分组
  - name: 国内  # IP分组名称，分为多个时在名称后面拼接“序号”，如“国内_1”
    ## IP分组 cidr 列表网址，每行一个，超过1000行会自动分为多个，ipv6 地址会被删除
    url: https://cdn.jsdelivr.net/gh/Loyalsoldier/geoip@release/text/cn.txt
  - name: Telegram
    url: https://cdn.jsdelivr.net/gh/Loyalsoldier/geoip@release/text/telegram.txt
stream-domain:      # 域名分流
  - interface: wan2 # 分流线路
    src-addr: 192.168.1.100-192.168.1.250   # 分流的源地址
    ## 域名列表网址，每行一个，超过1000行会自动分为多个
    url: https://cdn.jsdelivr.net/gh/Loyalsoldier/v2ray-rules-dat@release/gfw.txt
  - interface: wan2
    src-addr: 192.168.1.100-192.168.1.250
    url: https://cdn.jsdelivr.net/gh/Loyalsoldier/v2ray-rules-dat@release/greatfire.txt
stream-ipport:      # 端口分流，与IP分组组合使用
  - type: 0 # 分流方式：0-外网线路，1-下一跳网关
    interface: wan2 # 分流线路
    src-addr: 192.168.1.100-192.168.1.250   # 分流的源地址
    ## 域名列表网址，每行一个，超过1000行会自动分为多个
    ip-group: 国内 # 前面填写的IP分组名称，多个用逗号隔开
  - type: 1 # 分流方式：0-外网线路，1-下一跳网关
    nexthop: 192.168.1.2 # 下一跳网关
    src-addr: 192.168.1.100-192.168.1.250
    ip-group: Telegram
```
