// -------------------------
//
// Copyright 2015, undiabler
//
// git: github.com/undiabler/golang-whois
//
// http://undiabler.com
//
// Released under the Apache License, Version 2.0
//
//--------------------------

package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"time"

	//"github.com/astaxie/beego"
)

//Simple connection to whois servers with default timeout 5 sec
func GetWhois(domain string, ip string) (string, error) {
	return GetWhoisTimeout(domain, ip, time.Second*5)

}

//Connection to whois servers with various time.Duration
func GetWhoisTimeout(domain string, ip string, timeout time.Duration) (result string, err error) {

	var (
		parts []string

		buffer     []byte
		connection net.Conn
	)
	//conf_ip := beego.AppConfig.String("ip")
	parts = strings.Split(domain, ".")
	if len(parts) < 2 {
		err = fmt.Errorf("Domain(%s) name is wrong!", domain)
		return
	}
	//last part of domain is zome
	/*zone := parts[len(parts)-1]

	server, ok := servers[zone]

	if !ok {
		err = fmt.Errorf("No such server for zone %s. Domain %s.", zone, domain)
		return
	}
	ip := <-ipch
	fmt.Println("!!!!!!!!!!!!!!!!!!", ip)
	connection, err = Dial("tcp", ip, server)
	//connection, err = net.DialTimeout("tcp", net.JoinHostPort(server, "43"), timeout)
	*/
	zone := parts[len(parts)-1]

	server, ok := servers[zone]

	if !ok {
		err = fmt.Errorf("No such server for zone %s. Domain %s.", zone, domain)
		return
	}

	//fmt.Println("!!!!!!!!!!!!!!!!!!", ip)
	sip := net.ParseIP(ip)
	//fmt.Println("use ip :", ip)
	ns, err := net.LookupHost(server)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(ns) == 0 {
		return
	}
	dip := net.ParseIP(ns[0])
	srcAddr := &net.TCPAddr{IP: sip, Port: 0}
	dstAddr := &net.TCPAddr{IP: dip, Port: 43}

	connection, err = net.DialTCP("tcp", srcAddr, dstAddr)
	if err != nil {
		//return net.Conn error
		return
	}

	defer connection.Close()

	connection.Write([]byte(domain + "\r\n"))

	buffer, err = ioutil.ReadAll(connection)

	if err != nil {
		return
	}

	result = string(buffer[:])
	result = strings.ToLower(result)
	//	ipch <- ip
	return
}
func GetWhois2(domain string, remote string, ip string, timeout time.Duration) (result string, err error) { //注册商server查询可以使用

	var (
		parts []string

		buffer     []byte
		connection net.Conn
	)
	//conf_ip := beego.AppConfig.String("ip")
	parts = strings.Split(domain, ".")
	if len(parts) < 2 {
		err = fmt.Errorf("Domain(%s) name is wrong!", domain)
		return
	}
	//	ip := <-ipch

	sip := net.ParseIP(ip)
	fmt.Println("use ip :", ip)
	ns, err := net.LookupHost(remote)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(ns) == 0 {
		fmt.Println("no found local ip!")
		return
	}
	dip := net.ParseIP(ns[0])
	srcAddr := &net.TCPAddr{IP: sip, Port: 0}
	dstAddr := &net.TCPAddr{IP: dip, Port: 43}
	connection, err = net.DialTCP("tcp", srcAddr, dstAddr)
	//connection, err = Dial("tcp", ip, remote)
	//connection, err = net.DialTimeout("tcp", net.JoinHostPort(remote, "43"), timeout)

	if err != nil {
		//return net.Conn error
		return
	}

	defer connection.Close()

	connection.Write([]byte(domain + "\r\n"))

	buffer, err = ioutil.ReadAll(connection)

	if err != nil {
		return
	}

	result = string(buffer[:])
	//ipch <- ip
	return
}
func Dial(network string, local string, remote string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   500 * time.Millisecond, //超时设置
		KeepAlive: 1 * time.Second,
	}
	local = local + ":0" //端口0,系统会自动分配本机端口
	switch network {
	case "udp":
		addr, err := net.ResolveUDPAddr(network, local)
		if err != nil {
			return nil, err
		}
		dialer.LocalAddr = addr
	case "tcp":
		addr, err := net.ResolveTCPAddr(network, local)
		if err != nil {
			return nil, err
		}
		dialer.LocalAddr = addr
	}
	return dialer.Dial(network, remote+":43")
}
