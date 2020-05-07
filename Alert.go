package main



import (

	"database/sql"

	_ "fmt"

	_ "fmt"

	

	"net/http"

	"os/exec"

	_ "time"



	_ "github.com/go-sql-driver/mysql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

)



var updates tgbotapi.UpdatesChannel

var bot *tgbotapi.BotAPI



func startBOT() (tgbotapi.UpdatesChannel, *tgbotapi.BotAPI) {

	bot, _ := tgbotapi.NewBotAPI("1153037633:AAHp5oGyFvTncdN_9hkhoNyEQpuM4cwYnns")



	u := tgbotapi.NewUpdate(0)

	updates, _ := bot.GetUpdatesChan(u)

	return updates, bot

}

func sendMsg(Msg string) {



	msg := tgbotapi.NewMessage(817269876, "")

	msg.Text = Msg

	bot.Send(msg)



}



func checkServiceRunning(service string, server string, lastStatusService uint, db *sql.DB) {

	serviceName := "./exitCode.sh " + service + " " + server + " ;echo $?"

	StatusCode := exec.Command("sh", "-c", serviceName)

	statusCode, _ := StatusCode.Output()

	sttCode := string(statusCode)

	if sttCode == "0\n" && (lastStatusService != 0) {



		slect, _ := db.Query(changeStatustoOK(service, server))

		defer slect.Close()



		sendMsg("Service " + service + " on " + server + " is Dead -> Running")

	} else if sttCode != "0\n" && (lastStatusService != 1) {

		



		slect, _ := db.Query(changeStatustoFail(service, server))

		defer slect.Close()



		sendMsg("Service " + service + " on " + server + " is Running -> Dead")

	}



}



var last_status uint



func check_server_JIRA() {



	resp, _ := http.Get("http://192.168.141.209/")

	defer func() {

		if r := recover(); r != nil && (last_status == 200) {

			
			sendMsg("CHECK YOUR JIRA SERVER NOW")

			last_status = 404

		}

	}()


	defer resp.Body.Close()

	if (resp.StatusCode == 200) && (last_status != 200) {

		sendMsg("JIRA SERVER IS RUNNING")

		last_status = 200

	}

}



func main() {

	last_status = 200

	updates, bot = startBOT()

	db := connect()

	defer db.Close()



	go func() {

		for {



			slect, _ := db.Query("SELECT * FROM monitor")



			defer slect.Close()

			type Tag struct {

				HOSTNAME   string

				SERVICE    string

				STATUSCODE uint

			}

			for slect.Next() {

				var tag Tag

				_ = slect.Scan(&tag.HOSTNAME, &tag.SERVICE, &tag.STATUSCODE)

				checkServiceRunning(tag.SERVICE, tag.HOSTNAME, tag.STATUSCODE, db)



			}

		}

	}()



	for {

		check_server_JIRA()

	}



}



func connecttoDB(user, password, hostname, port, database string) string {

	return user + ":" + password + "@" + "tcp(" + hostname + ":" + port + ")/" + database



}



func changeStatustoOK(service string, server string) string {

	return "UPDATE monitor SET statusCode = 0 where Hostname = '" + server + "' AND Service = '" + service + "'"



}

func changeStatustoFail(service string, server string) string {

	return "UPDATE monitor SET statusCode = 1 where Hostname = '" + server + "' AND Service = '" + service + "'"



}

func connect() *sql.DB {

	db, _ := sql.Open("mysql", connecttoDB("root", "minhduc7b", "192.168.141.204", "3306", "monitor_byGo"))

	return db

}



