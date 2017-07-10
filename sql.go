package main

import (
	"database/sql"

	"time"

	_ "github.com/go-sql-driver/mysql"
)

type server_data struct {
	suffix string
	server string
}
type suffix_data struct {
	suffix        string
	limit         int
	wait          int
	carry         string
	deal          string
	routine       int
	domain_name   string
	create        string
	expiration    string
	update        string
	domain_status string
	name_server   string
}

var (
	SRC_DB *sql.DB
)

func Get_conf_analysis() (conflist []string, err error) {
	db, err := sql.Open("mysql", "root:root@/whois")
	if err != nil {
		return
	}
	defer db.Close()
	rows, err := db.Query("select content from conf_analysis")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		conflist = append(conflist, temp)
	}
	return
}
func Get_conf_server() (datas []server_data, err error) {
	db, err := sql.Open("mysql", "root:root@/whois")
	if err != nil {
		return
	}
	defer db.Close()
	rows, err := db.Query("select suffix,nserver from conf_server")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var temp server_data
		err = rows.Scan(&temp.suffix, &temp.server)
		datas = append(datas, temp)

	}
	return
}

func Get_conf_suffix() (datas []suffix_data, err error) {
	db, err := sql.Open("mysql", "root:root@/whois")
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query("select suffix,limited,wait,carry,deal,routine,create_at,expiration_at,update_at,domain_status,name_server,domain_name from conf_suffix")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var temp suffix_data
		err = rows.Scan(&temp.suffix, &temp.limit, &temp.wait, &temp.carry, &temp.deal, &temp.routine, &temp.create, &temp.expiration, &temp.update, &temp.domain_status, &temp.name_server, &temp.domain_name)
		if err != nil {
			return
		}
		datas = append(datas, temp)

	}
	return
}

type Domains_registrar struct { //数据库存储格式(domains_registrar)
	Domain          string
	Status          int
	RegistrarID     string
	DomainStatus    string
	CreationDate    string
	UpdatedDate     string
	ExpirationDate  string
	NameServers     string
	RegistrantEmail string
	Registrant      string
	Updated         int64
}
type Domain struct { //数据库存储格式（domain）
	Domain    string
	Msg1      string
	Msg2      string
	Create_at time.Time
	update_at time.Time
}

var (
	Local_db *sql.DB
)

//func init() {//合并到control.go中

//	Local_db, err = Open_db("root:mydns123@tcp(10.10.100.21:3306)/whois")
//	if err != nil {
//		log.Println(err)
//	}
//}
func Open_db(conn string) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", conn)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	return
}

//func Get_maxid(tb string) (id int, err error) {
//	err = Db.QueryRow("select max(id) from " + tb).Scan(&id)
//	if err != nil {
//		return
//	}
//	return
//}

//func Getdomain(begin int, end int, tb string) (re []string, err error) {

//	rows, err := Db.Query(`select tdomain from `+tb+` where health_status=0 and status in(1,2,3) and id in(select id  from `+tb+` where  id>? and id<=?) `, begin, end)

//	if err != nil {
//		return
//	}
//	defer rows.Close()

//	for rows.Next() {
//		var temp string
//		rows.Scan(&temp)
//		re = append(re, temp)
//	}
//	return

//}
func downtosql(domain string, data Domains_registrar) (err error) {
	t := time.Now()
	stamp := t.Unix()

	stmt, err := Local_db.Prepare(`insert domains_registrar1 set Domain=?,RegistrarID=?,DomainStatus=?,
	CreationDate=?,ExpirationDate=?,UpdatedDate=?,NameServers=?,RegistrantEmail=?,Registrant=?,Updated=?`)

	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		domain,
		data.RegistrarID,
		data.DomainStatus,
		data.CreationDate,
		data.ExpirationDate,
		data.UpdatedDate,
		data.NameServers,
		data.RegistrantEmail,
		data.Registrant,
		stamp)
	if err != nil {
		return
	}

	return
}
func Getdomain1(begin int, end int, tb string) (re []string, err error) {

	rows, err := SRC_DB.Query(`select tdomain  from `+tb+` where status in(1,2,3) and id>? and id<=? and health_status =0`, begin, end)
	if err != nil {
		return
	}
	defer rows.Close()
	var t string
	for rows.Next() {
		err = rows.Scan(&t)
		if err != nil {
			return
		}
		re = append(re, t)
	}
	return

}
func Sql_getcount(tb string) (id int, err error) {

	err = SRC_DB.QueryRow("SELECT max(id) from " + tb).Scan(&id)

	if err != nil {
		return
	}

	return
}
