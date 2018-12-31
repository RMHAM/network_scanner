//ben@kc2vjw.com

package main

import (
	"time"
	"fmt"
//	"log"
	"os"
	"strings"
	"net"
	"flag"
	"github.com/soniah/gosnmp"
)

var hosts map[string]int

func main() {
	entrypoint := flag.String("e", "entrypoint", "First host to query")
	community := flag.String("c", "public", "SNMP community string")
	devdb := flag.String("d", "http://10.30.20.7/cgi-bin/dev/dbapi", "devdb api endpoint")
	flag.Parse()
	if *entrypoint == "entrypoint" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	
	hosts = make(map[string]int)
	merge_hosts(*entrypoint, *community)
	for ip, _ := range(hosts) {
		hn := getname(ip)
		osstr := get_os(ip, *community)
		fmt.Printf("%s: %s: %s\n", ip, hn, osstr)
		touch_host(*devdb, ip, osstr)
	}
}

func getname(ip string) string {
		hn, err := net.LookupAddr(ip)
		if err != nil {
			hn = []string{"unknown"}
		}
		return(hn[0])
}
func merge_hosts(tgt string, community string) {
	newhosts := get_arp_table(tgt, community)
	for ip, _ := range newhosts {
		_, have := hosts[ip]
		hosts[ip] = hosts[ip]+1
		if !have {
			fmt.Printf("")
			fmt.Printf("discovered: %s (%s): %d\n", ip, getname(ip), hosts[ip])
			merge_hosts(ip, community)
		}
	}
}

func get_os(tgt string, community string) string {
	fingerp := snmp_fingerprint(tgt, community)
	if(len(fingerp) > 2) {
		return(fingerp)
	}
	fingerp = ssh_fingerprint(tgt)
	if(len(fingerp) > 2) {
		return(fingerp)
	}
	return("unknown")
}

func get_arp_table(tgt string, community string) map[string]int {
	hosts := make(map[string]int)

	arpoid := ".1.3.6.1.2.1.4.22.1.2"
	snmphost := &gosnmp.GoSNMP {
		Target:  tgt,
		Port: gosnmp.Default.Port,
		Community: community,
		Version: gosnmp.Version2c,
		Timeout: time.Duration(2*time.Second),
		//Logger: log.New(os.Stdout, "", 0),
		Logger: nil,
	}

	err := snmphost.Connect()
	if err == nil {
		defer snmphost.Conn.Close()
	} else {
		return(hosts)
	}

	result, _ := snmphost.WalkAll(arpoid)
	for _, r := range result {
		outoid := strings.Split(r.Name, ".")
		ip := strings.Join(outoid[len(outoid)-4:len(outoid)], ".")
		hosts[ip] = hosts[ip]+1
	}
	return(hosts)
}
