package main

import (
	"fmt"
	"net/http"
	_"log"
	"os/exec"
	_"os"
	// "strings"
	// "strconv"
)
var lastStatus uint


func checkServerJIRA(url string) {

	resp, err := http.Get(url)
	defer func (){
		if r:=recover(); r!=nil  && (lastStatus == 200){
			fmt.Println("OK -> CRITICAL")
			 lastStatus = 404
		}
	}()
	
	
	if err != nil {
		panic("CHECK YOUR SERVER NOW")
	}
	defer resp.Body.Close()
	if (resp.StatusCode == 200) && (lastStatus != 200 ){
		fmt.Println("CRITICAL -> OK")
		lastStatus = 200
	}
}

func checkServiceRunning(service string) {
	serviceName := "./exitCode.sh " + service + " ;echo $?"
	StatusCode := exec.Command("sh", "-c", serviceName)
	statusCode, _ := StatusCode.Output()
	sttCode := string(statusCode)
	// int_Code, _ := strconv.ParseUint(sttCode, 10, 32)
	if sttCode == "0\n" {
		fmt.Print("OK")
	} else {
		fmt.Print("NOT OK")
	}



	}


func main(){
	// checkServiceRunning("nginx")
	checkServiceRunning("httpd")
	// lastStatus == 200
	// for {
	// checkServerJIRA("http://192.168.141.204/")
	// }

}
