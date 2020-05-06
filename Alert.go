package main

import (
	"database/sql"
	_ "fmt"
	"log"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var updates tgbotapi.UpdatesChannel
var bot *tgbotapi.BotAPI

func startBOT() (tgbotapi.UpdatesChannel, *tgbotapi.BotAPI) {
	bot, err := tgbotapi.NewBotAPI("1153037633:AAHp5oGyFvTncdN_9hkhoNyEQpuM4cwYnns")
	if err != nil {
		log.Panic(err)
	}
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
		log.Printf("Service %s on %s is Dead", service, server)

		slect, _ := db.Query(changeStatustoFail(service, server))
		defer slect.Close()

		sendMsg("Service " + service + " on " + server + " is Running -> Dead")
	}

}

func main() {
	updates, bot = startBOT()
	db := connect()
	defer db.Close()
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
