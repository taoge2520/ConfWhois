package main

import (
	"log"
	"strings"
)

type Deal_data struct {
	Domain string
	Msg    Domains_registrar
}

var Exch chan Deal_data

func init() {
	Exch = make(chan Deal_data, 200)
	go Exchange()

}

func Exchange() {
	for {
		value := <-Exch
		//		log.Println("debug", value.Domain, value.Msg)
		//		if value.Msg == "" {
		//			log.Println(value.Domain, ": no msg return")
		//			continue
		//		}
		//		data, err := get_data(value.Domain, value.Msg)
		//		if err != nil {
		//			log.Println(err)
		//		}
		//		//log.Println(data)
		err = downtosql(value.Domain, value.Msg) //存储
		if err != nil {
			log.Println(err)
		}
	}
}
func get_data(domain string, result string, dconf []string, tag string, nstag string) (data Domains_registrar, err error) {
	if result == "" {
		return
	}
	temp := Parse(result, dconf)
	data.Domain = temp.DomainName
	data.Status = 0
	data.RegistrarID = temp.RegistrarID
	status := strings.Join(temp.Status, "\n")
	data.DomainStatus = status
	data.CreationDate = temp.CreationDate
	data.ExpirationDate = temp.ExpirationDate
	data.UpdatedDate = temp.UpdatedDate
	if tag == "1" {

		temp.NameServer = parse_name_servers(result, nstag)

	}
	temp.NameServer = RemoveDuplicatesAndEmpty(temp.NameServer)
	ns := strings.Join(temp.NameServer, ",")
	data.NameServers = ns
	//	data.RegistrantEmail = temp.RegistrantEmail
	//	data.Registrant = temp.Registrant
	if data.DomainStatus == "" {
		data.DomainStatus = "registered .."
	}
	return

}
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
