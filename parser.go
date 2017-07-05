//https://help.aliyun.com/knowledge_detail/35772.html
package main

import (
	"fmt"
	"reflect"
	"strings"
)

/*type WhoisInfo struct {
	DomainName     string   `whois:"domain name"`
	RegistrarID    string   `whois:"registrar iana id"`
	Status         []string `whois:"domain status,state,status"`
	NameServer     []string `whois:"name server,nserver"`
	UpdatedDate    string   `whois:"updated date,update date"`
	CreationDate   string   `whois:"registration time,creation date"`
	ExpirationDate string   `whois:"expiration date,expiration time,expiry date"`

	WhoisServer     string `whois:"whois server"`
	RegistrantEmail string `whois:"registrant email,registrant contact email"`
	Registrant      string `whois:"registrant name,registrant"`
	//Registrar      RegistrarInfo //注册商
	//Registrant Info `whois:"registrant "` //注册者
	//Administrative Info   `whois:"admin "`      //管理联系人
	//Technical      Info   `whois:"tech "`       //技术联系人
	//Billing        Info   `whois:"billing "`    //付费联系人
	//DNSSEC         string `whois:"dnssec"`
}*/
type WhoisInfo struct {
	DomainName     string   `whois:"domain name"`
	RegistrarID    string   `whois:"registrar iana id"`
	Status         []string `whois:"domain status,state,status"`
	NameServer     []string `whois:"name server,nserver,name servers"`
	UpdatedDate    string   `whois:"updated date,update date,domain record last updated"`
	CreationDate   string   `whois:"registration time,creation date,domain record activated,domain registration date"`
	ExpirationDate string   `whois:"expiration date,expiration time,expiry date,registry expiry date,domain expires"`
	//WhoisServer     string   `whois:"whois server"`
	//RegistrantEmail string   `whois:"registrant email,registrant contact email"`
	//Registrant      string   `whois:"registrant name,registrant"`
}

func PrintTags() {
	w := new(WhoisInfo)
	v := reflect.ValueOf(w).Elem()
	t := v.Type()
	for index := 0; index < v.NumField(); index++ {
		vField := v.Field(index)
		tags := strings.Split(t.Field(index).Tag.Get("whois"), ",")
		// name := t.Field(index).Name
		kind := vField.Kind()
		if kind == reflect.Struct {
			Struct := vField
			tStruct := Struct.Type()
			for i := 0; i < Struct.NumField(); i++ {
				// field := Struct.Field(i)
				_tags := strings.Split(tStruct.Field(i).Tag.Get("whois"), ",")
				for _, tag := range tags {
					for _, _tag := range _tags {
						if tag != "" {
							fmt.Println(tag + " " + _tag)
						} else {
							fmt.Println(_tag)
						}

					}
				}
			}
		} else if kind == reflect.Slice {
			for _, tag := range tags {
				fmt.Println(tag)
			}
		} else {
			for _, tag := range tags {
				fmt.Println(tag)
			}
		}
	}
}

/*func Parse(result string) *WhoisInfo {
	w := new(WhoisInfo)
	v := reflect.ValueOf(w).Elem()
	t := v.Type()
	for index := 0; index < v.NumField(); index++ {
		vField := v.Field(index)
		tags := strings.Split(t.Field(index).Tag.Get("whois"), ",")

		kind := vField.Kind()
		if kind == reflect.Struct {
			Struct := vField
			tStruct := Struct.Type()

			for i := 0; i < Struct.NumField(); i++ {
				_tags := strings.Split(tStruct.Field(i).Tag.Get("whois"), ",")
			NextField:
				for _, tag := range tags {
					for _, _tag := range _tags {
						if tag != "" {
							// fmt.Println(tag + " " + _tag)
							if value, ok := getValue(result, tag+" "+_tag); ok {
								Struct.Field(i).SetString(value)
								break NextField
							}
						} else {
							// fmt.Println(_tag)
							if value, ok := getValue(result, _tag); ok {
								Struct.Field(i).SetString(value)
								break NextField
							}
						}
					}
				}
			}
		} else if kind == reflect.Slice {
			for _, tag := range tags {
				if value, ok := getValueSlice(result, tag); ok {
					vField.Set(reflect.ValueOf(value))
					break
				}
			}
		} else {
			for _, tag := range tags {
				if value, ok := getValue(result, tag); ok {
					vField.SetString(value)
					break
				}
			}
		}
	}
	return w
}*/
func Parse(result string, dconf []string) *WhoisInfo {

	w := new(WhoisInfo)
	v := reflect.ValueOf(w).Elem()
	//t := v.Type()

	for index := 0; index < v.NumField(); index++ {

		vField := v.Field(index)
		temp := dconf[index] //strings.Split(t.Field(index).Tag.Get("whois"), ",")
		tags := strings.Split(temp, ",")
		kind := vField.Kind()

		if kind == reflect.Slice {

			for _, v := range tags {
				if value, ok := getValueSlice(result, v); ok {
					vField.Set(reflect.ValueOf(value))
					break
				}
			}
		} else {
			for _, v := range tags {
				if value, ok := getValue(result, v); ok {
					vField.SetString(value)
					break
				}
			}
		}

	}
	return w
}

func parse_name_servers(result string, tag string) (ns []string) {
	start := strings.Index(result, tag)
	if start == -1 {
		return
	}
	start += len(tag) //+ ":"
	result = result[start:]
	end := strings.Index(result, "\n\n")
	if end == -1 {
		return
	}
	temp := result[:end]

	temp = strings.Trim(temp, "\n")
	h := strings.Split(temp, "\n")
	for _, v := range h {
		v = strings.Trim(v, " ")
		n := strings.Split(v, " ")
		if len(n) > 0 {
			ns = append(ns, n[0])
		} else {
			continue
		}
	}
	return
}
func getValue(result, tag string) (string, bool) {
	key := strings.TrimSpace(tag) //+ ":"
	start := strings.Index(result, key)
	if start < 0 {
		return "", false
	}
	start += len(key)
	end := strings.Index(result[start:], "\n")
	value := strings.TrimSpace(result[start : start+end])
	return value, true
}

func getValueSlice(result, tag string) (slice []string, ok bool) {
	key := tag //+ ":"
	for {
		start := strings.Index(result, key)
		if start < 0 {
			break
		}
		ok = true
		start += len(key)
		end := strings.Index(result[start:], "\n")
		value := strings.TrimSpace(result[start : start+end])
		slice = append(slice, value)
		result = result[start+end:]
	}
	return
}
