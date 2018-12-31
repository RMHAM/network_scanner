package main

import (
	"log"
	"net"
	"bufio"
	"time"
	"strings"
	//"os"
)

func ssh_fingerprint(ip string) string {
	//RFC4253 sec 4.2
	uri := net.JoinHostPort(ip, "22")
	conn, err := net.DialTimeout("tcp", uri, 3*time.Second)
	if err != nil {
		log.Printf("ssh_fingerprint: %s", err)
		return("")
	} else {
		defer conn.Close()
	}
	id_string, err := bufio.NewReader(conn).ReadString('\r')
	id_ary := strings.SplitN(id_string, " ", 2)
	if len(id_ary) < 2 {
		log.Printf("%s listens on port 22 but no useful fingerprint: %s", ip, id_string)
		return ""
	} else {
		return(id_ary[1])
	}
}

//func main() {
//	ip := os.Args[1]
//
//	log.Printf(ssh_fingerprint(ip))
//}