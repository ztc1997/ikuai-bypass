package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/ztc1997/ikuai-bypass/api"
	"github.com/ztc1997/ikuai-bypass/router"
	"gopkg.in/yaml.v3"
)

var confPath = flag.String("c", "./config.yml", "配置文件路径")

var conf struct {
	IkuaiURL  string `yaml:"ikuai-url"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Cron      string `yaml:"cron"`
	CustomIsp []struct {
		Name string `yaml:"name"`
		URL  string `yaml:"url"`
	} `yaml:"custom-isp"`
	StreamDomain []struct {
		Interface string `yaml:"interface"`
		SrcAddr   string `yaml:"src-addr"`
		URL       string `yaml:"url"`
	} `yaml:"stream-domain"`
}

func main() {
	flag.Parse()

	err := readConf(*confPath)
	if err != nil {
		log.Println("读取配置文件失败：", err)
		return
	}

	update()

	if conf.Cron == "" {
		return
	}

	c := cron.New()
	_, err = c.AddFunc(conf.Cron, update)
	if err != nil {
		log.Println("启动计划任务失败：", err)
		return
	} else {
		log.Println("已启动计划任务")
	}
	c.Start()

	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
	}
}

func readConf(filename string) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return fmt.Errorf("in file %q: %v", filename, err)
	}
	return nil
}

func update() {
	err := readConf(*confPath)
	if err != nil {
		log.Println("更新配置文件失败：", err)
		return
	}

	baseurl := conf.IkuaiURL
	if baseurl == "" {
		gateway, err := router.GetGateway()
		if err != nil {
			log.Println("获取默认网关失败：", err)
			return
		}
		baseurl = "http://" + gateway
		log.Println("使用默认网关地址：", baseurl)
	}

	iKuai := api.NewIKuai(baseurl)

	err = iKuai.Login(conf.Username, conf.Password)
	if err != nil {
		log.Println("登陆失败：", err)
		return
	} else {
		log.Println("登录成功")
	}

	err = iKuai.DelIKuaiBypassCustomIsp()
	if err != nil {
		log.Println("移除旧的自定义运营商失败：", err)
	} else {
		log.Println("移除旧的自定义运营商成功")
	}
	for _, customIsp := range conf.CustomIsp {
		err = updateCustomIsp(iKuai, customIsp.Name, customIsp.URL)
		if err != nil {
			log.Printf("添加自定义运营商'%s@%s'失败：%s\n", customIsp.Name, customIsp.URL, err)
		} else {
			log.Printf("添加自定义运营商'%s@%s'成功\n", customIsp.Name, customIsp.URL)
		}
	}

	err = iKuai.DelIKuaiBypassStreamDomain()
	if err != nil {
		log.Println("移除旧的域名分流失败：", err)
	} else {
		log.Println("移除旧的域名分流成功")
	}
	for _, streamDomain := range conf.StreamDomain {
		err = updateStreamDomain(iKuai, streamDomain.Interface, streamDomain.SrcAddr, streamDomain.URL)
		if err != nil {
			log.Printf("添加域名分流 '%s@%s' 失败：%s\n", streamDomain.Interface, streamDomain.URL, err)
		} else {
			log.Printf("添加域名分流 '%s@%s' 成功\n", streamDomain.Interface, streamDomain.URL)
		}
	}
}

func updateCustomIsp(iKuai *api.IKuai, name, url string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ips := strings.Split(string(body), "\n")
	ips = removeIpv6(ips)
	ipGroups := group(ips, 5000)
	for _, ig := range ipGroups {
		ipGroup := strings.Join(ig, ",")
		iKuai.AddCustomIsp(name, ipGroup)
	}
	return
}

func updateStreamDomain(iKuai *api.IKuai, iface, srcAddr, url string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	domains := strings.Split(string(body), "\n")
	domainGroup := group(domains, 1000)
	for _, d := range domainGroup {
		domain := strings.Join(d, ",")
		iKuai.AddStreamDomain(iface, srcAddr, domain)
	}
	return
}

func removeIpv6(ips []string) []string {
	i := 0
	for _, ip := range ips {
		if !strings.Contains(ip, ":") {
			ips[i] = ip
			i++
		}
	}
	return ips[:i]
}

func group(arr []string, subGroupLength int64) [][]string {
	max := int64(len(arr))
	var segmens = make([][]string, 0)
	quantity := max / subGroupLength
	remainder := max % subGroupLength
	i := int64(0)
	for i = int64(0); i < quantity; i++ {
		segmens = append(segmens, arr[i*subGroupLength:(i+1)*subGroupLength])
	}
	if quantity == 0 || remainder != 0 {
		segmens = append(segmens, arr[i*subGroupLength:i*subGroupLength+remainder])
	}
	return segmens
}
