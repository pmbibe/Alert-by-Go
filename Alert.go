package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

var lastStatusJIRA uint

func checkServerJIRA(url string) {

	resp, err := http.Get(url)
	defer func() {
		if r := recover(); r != nil && (lastStatusJIRA == 200) {
			fmt.Println("OK -> CRITICAL")
			lastStatusJIRA = 404
		}
	}()

	if err != nil {
		panic("CHECK YOUR SERVER NOW")
	}
	defer resp.Body.Close()
	if (resp.StatusCode == 200) && (lastStatusJIRA != 200) {
		fmt.Println("CRITICAL -> OK")
		lastStatusJIRA = 200
	}
}

var lastStatusService uint

func checkServiceRunning(service string, server string) {
	serviceName := "./exitCode.sh " + service + " " + server + " ;echo $?"
	StatusCode := exec.Command("sh", "-c", serviceName)
	statusCode, _ := StatusCode.Output()
	sttCode := string(statusCode)
	if sttCode == "0\n" && (lastStatusService != 0) {
		log.Print("Service is running")
		lastStatusService = 0
	} else if sttCode != "0\n" && (lastStatusService != 1) {
		log.Print("Service is Dead")
		lastStatusService = 1
	}

}

// func openFiletoReadService() {

// }
func openFiletoReadService() {
	type Server struct {
		NAMESERVER string
		SERVICE    []string
	}
	data, err := ioutil.ReadFile("test.json")
	if err != nil {
		fmt.Print(err)
	}

	var obj []Server
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}
	for i := 0; i < len(obj); i++ {
		for j := 0; j < len(obj[i].SERVICE); j++ {
			checkServiceRunning(obj[i].SERVICE[j], obj[i].NAMESERVER)
		}

	}

}

func main() {
	for {
		openFiletoReadService()
	}

}
