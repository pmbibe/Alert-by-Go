package main
import(
	_"fmt"
	"log"
	"database/sql"
	"os/exec"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func startBOT() (tgbotapi.UpdatesChannel,*tgbotapi.BotAPI )  {
	bot, err := tgbotapi.NewBotAPI("1153037633:AAHp5oGyFvTncdN_9hkhoNyEQpuM4cwYnns")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	return updates, bot
}
func sendMsg(Msg string) {

	updates, bot := startBOT()
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Chat.ID == 817269876 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.Text = Msg
			bot.Send(msg)
		}

		}


	}

func checkServiceRunning(service string, server string, lastStatusService uint, db *sql.DB) {
	serviceName := "./exitCode.sh " + service + " " + server + " ;echo $?"
	StatusCode := exec.Command("sh", "-c", serviceName)
	statusCode, _ := StatusCode.Output()
	sttCode := string(statusCode)
	if sttCode == "0\n" && (lastStatusService != 0) {
		log.Printf("Service %s on %s is Running", service,server)
		sendMsg("Service " + service + " on " + service + " is Dead -> Running")
		slect, _ :=db.Query(changeStatustoOK(service,server))
		defer slect.Close()
	} else if sttCode != "0\n" && (lastStatusService != 1) {
		log.Printf("Service %s on %s is Dead", service,server)
		sendMsg("Service " + service + " on " + service + " is Running -> Dead")
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
