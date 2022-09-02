package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func ListenForMail(){
	go func(){
		for{
			msg:= <- app.MailChan
			SendMsg(msg)
		}
	}()
}

func SendMsg(m models.MailData){
	server:= mail.NewSMTPClient()
	server.KeepAlive= false
	server.Host="localhost"
	server.Port=1025
	server.ConnectTimeout= 10*time.Second
	server.SendTimeout= 10*time.Second

	client,err:= server.Connect()
	if err!=nil{
		infoLog.Println(err)
	}

	email:=mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == ""{
		email.SetBody(mail.TextHTML,m.Content)
	}else{
		data,err:= ioutil.ReadFile(fmt.Sprintf("./email-template/%s", m.Template))
		if err!= nil{
			app.ErrorLog.Println(err)
		}
		mailTemplate:= string(data)
		msgToSend:=strings.Replace(mailTemplate,"[%body%]",m.Content,1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	err= email.Send(client)
	if err!= nil{
		log.Println(err)
	}else{
		log.Println("Email sent")
	}


}