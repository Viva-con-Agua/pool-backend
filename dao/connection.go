package dao

import "github.com/Viva-con-Agua/vcago"

var Database = vcago.NewMongoDB("pool-user")
var MailSend = vcago.NewMailSend()
