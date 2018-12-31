package main

import (
	"time"
	"github.com/soniah/gosnmp"
	//"os"
	//"log"
)

func get_snmp_os(tgt string, community string) string {

	osoid := ".1.3.6.1.2.1.47.1.1.1.1.2"
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
		return("")
	}

	result, _ := snmphost.WalkAll(osoid)
	for _, r := range result {
		return(string(r.Value.([]byte)))
	}
	return("")
}

func get_snmp_osv1(tgt string, community string) string {

	osoid := ".1.3.6.1.2.1.1"
	snmphost := &gosnmp.GoSNMP {
		Target:  tgt,
		Port: gosnmp.Default.Port,
		Community: community,
		Version: gosnmp.Version1,
		Timeout: time.Duration(2*time.Second),
		//Logger: log.New(os.Stdout, "", 0),
		Logger: nil,
	}

	err := snmphost.Connect()
	if err == nil {
		defer snmphost.Conn.Close()
	} else {
		return("")
	}

	result, _ := snmphost.WalkAll(osoid)
	for _, r := range result {
		return(string(r.Value.([]byte)))
	}
	return("")
}

func snmp_fingerprint(ip string, community string) string {
	result := ""
	if(community != ""){
		result = get_snmp_os(ip, community)
	}
	if(len(result) > 1) {
		return(result)
	}

	result = get_snmp_os(ip, "public")
	if(len(result) > 1) {
		return(result)
	}

	if(community != "") {
		result = get_snmp_osv1(ip, community)
	}
	if(len(result) > 1) {
		return(result)
	}
	result = get_snmp_osv1(ip, "public")
	if(len(result) > 1) {
		return(result)
	}
	return("")
}

//func main() {
//	var community string
//	if(len(os.Args) > 2) {
//		community = os.Args[2]	
//	} else {
//		community = ""
//	}
//	log.Printf(snmp_fingerprint(os.Args[1], community))
//}