package main
import(
	_"fmt"
	"log"
	"database/sql"
	"os/exec"
	_ "github.com/go-sql-driver/mysql"
)

}
func connecttoDB(user,password,hostname,port,database string ) string {
	return user +":"+password +"@" + "tcp(" + hostname + ":" +port +")/"+ database

}

func changeStatustoOK(service string, server string) string {
	return "UPDATE monitor SET statusCode = 0 where Hostname = '" + server + "' AND Service = '" + service +"'"
	
}
func changeStatustoFail(service string, server string) string {
	return "UPDATE monitor SET statusCode = 1 where Hostname = '" + server + "' AND Service = '" + service +"'"
	
}
func connect() *sql.DB {
	db, _ := sql.Open("mysql", connecttoDB("root","minhduc7b","192.168.141.204","3306","monitor_byGo"))
	return db
}
func checkServiceRunning(service string, server string, lastStatusService uint, db *sql.DB) {
	serviceName := "./exitCode.sh " + service + " " + server + " ;echo $?"
	StatusCode := exec.Command("sh", "-c", serviceName)
	statusCode, _ := StatusCode.Output()
	sttCode := string(statusCode)
	if sttCode == "0\n" && (lastStatusService != 0) {
		log.Printf("Service %s on %s: Dead -> Running", service,server)
		slect, _ :=db.Query(changeStatustoOK(service,server))
		defer slect.Close()
	} else if sttCode != "0\n" && (lastStatusService != 1) {
		log.Printf("Service %s on %s: Running ->  Dead", service,server)
		slect, _ :=db.Query(changeStatustoFail(service,server))
		defer slect.Close()
	}

}

func main(){

db := connect()
defer db.Close()
for {
slect, _ := db.Query("SELECT * FROM monitor")
defer slect.Close()
type Tag struct {
	HOSTNAME string
	SERVICE string
	STATUSCODE uint
}
for slect.Next(){
	var tag Tag
	_ = slect.Scan(&tag.HOSTNAME,&tag.SERVICE,&tag.STATUSCODE)
	checkServiceRunning(tag.SERVICE, tag.HOSTNAME, tag.STATUSCODE, db)
	
}
}

