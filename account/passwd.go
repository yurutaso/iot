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
		log.Fatal(fmt.Errorf("Usage: passwd username password"))
	}
	username := args[1]
	password := args[2]
	hashpassword, err := hashPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite3", "passwd.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "CERTS" ("name" VARCHAR(255), "key" VARCHAR(255))`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`REPLACE INTO "CERTS" ("name", "key") VALUES (?, ?)`, username, hashpassword)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
}
