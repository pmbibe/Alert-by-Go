package main

import (
	"fmt"
	"net/http"
	_"log"
	"os/exec"
	"os"
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
	// serviceName := "./exitCode.sh " + service
	StatusCode := exec.Cmd{
		Path: "./exitCode.sh",
		Args: []string{"./exitCode.sh" ,service},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	output, _ := StatusCode.Output()
	fmt.Println(output)

	}


func main(){
	checkServiceRunning("nginx")
	// lastStatus == 200
	// for {
	// checkServerJIRA("http://192.168.141.204/")
	// }

}