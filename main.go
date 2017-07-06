// ConfWhois project main.go
package main

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode"
)

var (
	SuffixMap map[string]chan string
	IpMap     map[string]chan Ipuse
	//Domain_conf []string
	Iplist  []string
	err     error
	servers map[string]string
	enter   chan string
)

const (
	STRICT = 30
)

type Ipuse struct {
	Ip    string
	Count int
}

func main() {
	err = Loadconf()
	if err != nil {
		fmt.Println(err)
		return
	}
	Create_checker() //创建工作者
	time.Sleep(1 * time.Second)
	fmt.Println("init success!")
	go Producer() //生产者
	go Listener() //监控者

	Distribution() //分拣者

}
func Distribution() {
	for {
		value := <-enter

		value1 := value
		if IsChineseChar(value) { //中文域名转码
			value1, err = ToASCII(value)
			if err != nil {
				fmt.Println(err)
			}
		}

		if !Check_domain(value1) {
			fmt.Println("非法域名！")
			return
		}
		part := strings.Split(value1, ".")
		suffix := part[len(part)-1]
		//fmt.Println(suffix, len(SuffixMap[suffix]))
		if _, ok := SuffixMap[suffix]; ok {

			SuffixMap[suffix] <- value1
		} else {
			fmt.Println("未启动该域名后缀对应的线程,默认处理方式")
			SuffixMap["curr"] <- value1
		}

	}
}

func Producer() { //查询入口，可视情况调整更改
	enter <- "apple.xxx" //test
	//	domains := []string{"baidu.biz", "app.biz", "baidu.cc", "hl.cc",
	//		"baidu.cn", "qq.cn", "baidu.com", "dns.com", "jmu.edu", "stanford.edu",
	//		"sdf.info", "ten.info", "baidu.ltd", "golang.org", "baidu.pub", "baidu.top"}
	//	for _, v := range domains {
	//		enter <- v
	//	}
}

func Create_checker() {
	conf, err := Get_conf_suffix()
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range conf {
		go Checker(v)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("checker ready!")
}
func Listener() { //监控并定时报告当前协程数量
	fmt.Println("linstener ready!")
	for {
		time.Sleep(3 * time.Minute)
		conf, err := Get_conf_suffix()
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range conf {
			if _, ok := SuffixMap[v.suffix]; ok {
				continue
			} else {
				go Checker(v)

				time.Sleep(100 * time.Millisecond)
			}
		}
		fmt.Println("now goroutine is:", runtime.NumGoroutine())
	}
}
func Checker(info suffix_data) {

	ch := make(chan string, 30)
	SuffixMap[info.suffix] = ch

	IpMap = make(map[string]chan Ipuse)

	for _, v := range Iplist {
		var ip_unit Ipuse
		ip_unit.Count = 0
		ip_unit.Ip = v
		ch1 := make(chan Ipuse, 30)
		IpMap[info.suffix] = ch1
		IpMap[info.suffix] <- ip_unit
	}
	ip := <-IpMap[info.suffix]
	//fmt.Println(info.suffix, "goroutine is ready get ip is:", ip.Ip)
	var analize []string
	analize = append(analize, info.domain_name)
	analize = append(analize, "registrar iana id:") //按照顺序,和解析结构体一致原则
	analize = append(analize, info.domain_status)
	analize = append(analize, info.name_server)
	analize = append(analize, info.update)
	analize = append(analize, info.create)
	analize = append(analize, info.expiration)

	for {
		value := <-SuffixMap[info.suffix]
		//fmt.Println(info.suffix, "get!", value)
		if ip.Count >= info.limit {
			ip.Count = 0
			IpMap[info.suffix] <- ip
			ip = <-IpMap[info.suffix]

		}
		result, err := GetWhois(info.carry+value, ip.Ip)
		if err != nil {
			fmt.Println(err)
		}
		ip.Count++

		//此处开始做msg处理部分

		data, err := get_data(value, result, analize, info.deal, info.name_server)
		if err != nil {
			fmt.Println(err)
		}
		var re Deal_data
		re.Domain = value
		re.Msg = data
		Exch <- re
	}
}
func Loadconf() (err error) { //程序启动时加载并初始化相关参数

	enter = make(chan string, 100)
	SuffixMap = make(map[string]chan string)
	//初始化配置文件,加载配置的ip文件
	myConfig := new(Config)
	myConfig.InitConfig("./configip.txt")
	h := strings.Split(myConfig.Mymap["conf=ip"], ",")
	for _, v := range h {
		Iplist = append(Iplist, v)
	}

	//	Domain_conf, err = Get_conf_analysis()
	//	if err != nil {
	//		return
	//	}

	servers = make(map[string]string)
	data, err := Get_conf_server()
	if err != nil {
		return
	}
	for _, v := range data {
		servers[v.suffix] = v.server
	}
	Local_db, err = Open_db(myConfig.Mymap["conf=conn"])
	if err != nil {
		fmt.Println(err)
	}
	return

}

func Check_domain(str string) bool {

	if m, _ := regexp.MatchString(`[a-zA-Z0-9][-a-zA-Z0-9]{0,62}\.([a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\.?`, str); !m {
		return false
	}
	return true

}
func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}
