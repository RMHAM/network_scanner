package main

import (
	//"flag"
	"net/http"
	"net/url"
	"log"
	"strings"
	"bufio"
	"bytes"
	"io/ioutil"
)

var subnetcache map[string]bool

func init() {
	subnetcache = make(map[string]bool)
}

func slash24(subnet string) string {
	//1.2.3.4 -> 1.2.3
	octets := strings.Split(subnet, ".")
	if len(octets) > 3 {
		return(strings.Join(octets[0:3], "."))
	} else {
		return(subnet)
	}
}

func check_subnet(api string, address string) bool {
	//cache the easy stuff
	if subnetcache[slash24(address)] == true {
		return(true)
	}

	v := url.Values{}
	v.Add("MODE", "scan")
	v.Add("NET", slash24(address))
	log.Printf("Using endpoint: " + api)
	requrl := api + "?" + v.Encode()

	log.Printf("Getting: %s\n", requrl)
	res, err := http.Get(requrl)
	if err != nil {
		log.Fatal(err)
	} else {
		defer res.Body.Close()
	}
	txt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("%s", txt)

	scanner := bufio.NewScanner(bytes.NewReader(txt))
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == "ALL" {
			subnetcache[slash24(address)] = true
			return(true)
		}
		if strings.TrimSpace(scanner.Text()) == strings.TrimSpace(address) {
			return(true)
		}
		log.Println(scanner.Text())
	}
	return(false)
}

func touch_host(api string, address string, desc string) {
	//don't admit that we found disallowed hosts
	if !check_subnet(api, address) {
		return
	}
	v := url.Values{}
	v.Add("MODE", "seen")
	v.Add("IP", address)
	if(desc != "") {
		v.Add("INFO", desc)
	} else {
		v.Add("INFO", "TBD")
	}
	requrl := api + "?" + v.Encode()
	log.Printf("Getting: %s\n", requrl)
	res, err := http.Get(requrl)
	if err != nil {
		log.Fatal(err)
	} else {
		defer res.Body.Close()
	}

	txt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("%s: %s\n", address, txt)
	}
}

/*func main() {
	var api string
	apif := flag.String("a", "api", "api endpoint")
	if *apif != "api" {
		api = *apif
	} else {
		api = "http://10.30.20.7/cgi-bin/dev/dbapi"
	}
	flag.Parse()
	//if(false) {
	_ = check_subnet(api, "192.168.108.1")
	touch_host(api, "192.168.108.1", "API TESTING")
	//}
	//log.Printf("%s\n", slash24("192.168.1.1"))

}*/