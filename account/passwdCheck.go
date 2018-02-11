package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

func hashPassword(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	return string(bytes), err
}

func CheckPasswordHash(p, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
	return err == nil
}

func main() {
	args := os.Args
	argc := len(args)
	if argc != 3 {
		log.Fatal(fmt.Errorf("Usage: passwdCheck username password"))
	}
	username := args[1]
	password := args[2]
	db, err := sql.Open("sqlite3", "passwd.db")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(`SELECT key from "CERTS" where name=?`, username)
	if err != nil {
		log.Fatal(err)
	}
	var hash string
	for rows.Next() {
		if err != nil {
			log.Fatal(err)
		}
		if err := rows.Scan(&hash); err != nil {
			log.Fatal(err)
			fmt.Println("NO")
		} else {
			if CheckPasswordHash(password, hash) {
				fmt.Println("OK")
			}
		}
	}
}
