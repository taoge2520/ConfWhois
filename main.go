// ConfWhois project main.go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode"
)

var (
	SuffixMap map[string]chan string
	//	IpMap     map[string]chan Ipuse
	//Domain_conf []string
	Iplist  []string
	err     error
	servers map[string]string
	enter   chan string
	file    *os.File
	Lenip   int
)

const (
	STRICT = 30

	BATCH = 100000
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

		if IsChineseChar(value) { //中文域名转码
			value, err = ToASCII(value)
			if err != nil {
				fmt.Println(err)
			}
		}

		if !Check_domain(value) {
			fmt.Println("非法域名！")
			return
		}
		part := strings.Split(value, ".")
		suffix := part[len(part)-1]
		//fmt.Println(suffix, len(SuffixMap[suffix]))
		if _, ok := SuffixMap[suffix]; ok {

			SuffixMap[suffix] <- value
		} else {
			fmt.Println(value, "未启动该域名后缀对应的线程,默认处理方式")
			SuffixMap["curr"] <- value
		}

	}
}

func Producer() { //查询入口，可视情况调整更改
	maxid, err := Sql_getcount("domain_7")
	//maxid := 10000
	fmt.Println("max id of datas is:", maxid)
	if err != nil {
		fmt.Println(err)
		return
	}

	page := maxid / BATCH
	if maxid%BATCH > 0 {
		page++
	}

	//var jiaoyan int
	for i := 0; i < page; i++ {

		datas, err := Getdomain1(i*BATCH, (i+1)*BATCH, "domain_7")

		if len(datas) == 0 {
			continue
		}
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range datas {

			enter <- v

		}
	}

}

func Create_checker() {
	conf, err := Get_conf_suffix()
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range conf {
		ch := make(chan string, 1000)
		SuffixMap[v.suffix] = ch
		for i := 0; i < v.routine; i++ {
			ip := Ipuse{Iplist[i], 0}
			go Checker(v, ip)
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("checker ready!")
}
func Listener() { //监控并定时报告当前协程数量
	fmt.Println("linstener ready!")
	for {
		fmt.Println("now goroutine is:", runtime.NumGoroutine())
		time.Sleep(3 * time.Minute)
		conf, err := Get_conf_suffix()
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range conf {
			if _, ok := SuffixMap[v.suffix]; ok {
				continue
			} else {
				for i := 0; i < v.routine; i++ {
					ip := Ipuse{Iplist[i], 0}
					go Checker(v, ip)
				}

				time.Sleep(100 * time.Millisecond)
			}
		}
		fmt.Println("now goroutine is:", runtime.NumGoroutine())
	}
}
func Checker(info suffix_data, ip Ipuse) { //测试ing

	//	IpMap = make(map[string]chan Ipuse)

	//	for _, v := range Iplist {
	//		var ip_unit Ipuse
	//		ip_unit.Count = 0
	//		ip_unit.Ip = v
	//		ch1 := make(chan Ipuse, 30)
	//		IpMap[info.suffix] = ch1
	//		IpMap[info.suffix] <- ip_unit
	//	}
	//ip := <-IpMap[info.suffix]
	if len(Iplist) == 0 {
		fmt.Println("length of iplist is zero,this routine open fail!")
		return
	}

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

		if ip.Count >= info.limit {
			fmt.Println("change ip!!!!!!!!!!!!")
			ip = Get_ipuse(ip)

		}
		fmt.Println("use ip is :", ip.Ip, "check:", value)
		result, err := GetWhois(info.carry+value, ip.Ip)
		if err != nil {
			fmt.Println(value, err)
			time.Sleep(10 * time.Second) //等待结束在重试一次
			result, err = GetWhois(info.carry+value, ip.Ip)
			if err != nil {
				continue
			}
		}
		if strings.Contains(result, info.domain_name) { //其他非正常情况下，重试一次
			time.Sleep(time.Duration(info.wait) * time.Millisecond)
			result, err = GetWhois(info.carry+value, ip.Ip)
			if err != nil {
				continue
			}
		}

		file.WriteString(value + "---->" + result + "\n\n")
		time.Sleep(time.Duration(info.wait) * time.Millisecond)
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
	Lenip = len(h)
	if Lenip == 0 {
		os.Exit(0)
	}
	for _, v := range h {
		Iplist = append(Iplist, v)
	}
	fmt.Println("iplist is :", Iplist)
	SRC_DB, err = Open_db(myConfig.Mymap["conf=conn249"]) //249db
	if err != nil {
		fmt.Println(err)
	}
	fileName := "save.txt"
	file, err = os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}

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
func Get_ipuse(pre_ip Ipuse) (ipu Ipuse) {
	for {
		seed := Random_number(Lenip)
		ipu.Count = 0
		ipu.Ip = Iplist[seed]
		if ipu.Ip == pre_ip.Ip {
			continue
		}
		return
	}
	return
}

func Random_number(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max)

}
