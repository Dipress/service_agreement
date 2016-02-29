package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

type Service struct {
	title    string
	comment  string
	login    string
	password string
	ip       string
	note     string
}

var service = Service{}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	//Take the settings from the configuration file config.json
	cfg := GetConfiguration()

	//open Database
	db, err := sql.Open(cfg.DBDriver, cfg.DBUsername+":"+cfg.DBPassword+"@"+cfg.DBProtocol+"("+cfg.DBHost+":"+cfg.DBPort+")/"+cfg.DBName)
	check(err)
	defer db.Close()

	//create query string
	query := `SELECT contract.title, contract.comment, inet_serv_14.login, inet_serv_14.password, inet_serv_14.addressFrom, inet_serv_14.comment
	 					FROM contract LEFT JOIN inet_serv_14 ON contract.id = inet_serv_14.contractId
	 					WHERE contract.id = ?`

	rows, err := db.Query(query, 204)

	check(err)
	defer rows.Close()

	//Create a file that will store the data from the database
	file, err := os.Create(cfg.FileName)
	check(err)
	defer file.Close()

	io.WriteString(file, "^  IP адрес                              ^  Договор        ^  Владелец       ^ Доступ        ^ Примечания ^  "+"\n")

	for rows.Next() {
		err := rows.Scan(&service.title, &service.comment, &service.login, &service.password, &service.ip, &service.note)
		check(err)

		//Fill in the data file
		io.WriteString(file, "|[[http://"+humanityIPAddress(service.ip)+"|"+humanityIPAddress(service.ip)+"]]"+
			" | "+service.title+
			" | "+service.comment+
			" | "+service.login+":"+service.password+
			" | "+service.note+"|"+"\n")
	}

	//authenticate with the remote server
	sshConfig := &ssh.ClientConfig{
		User: cfg.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.SSHPassword),
		},
	}

	//create session
	client, err := ssh.Dial(cfg.SSHProtocol, cfg.SSHHost+":"+cfg.SSHPort, sshConfig)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	defer session.Close()

	//copy current file to remote server
	scp.CopyPath(cfg.FileName, cfg.RemotePath, session)

	//remove file
	os.Remove(cfg.FileName)
}

func humanityIPAddress(template string) string {
	ip := []byte(template)
	var strip string
	for index, value := range ip {
		str := fmt.Sprintf("%v", value)
		if index != len(ip)-1 {
			strip += str + "."
		} else {
			strip += str
		}
	}
	return strip
}
