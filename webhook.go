package iot

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func hasKey(dict map[string][]string, key string) bool {
	if _, ok := dict[key]; ok {
		return true
	}
	return false
}

func checkIdKey(id string, key string) (bool, error) {
	db, err := sql.Open("sqlite3", "account/passwd.db")
	if err != nil {
		return false, err
	}
	rows, err := db.Query(`SELECT key from "CERTS" where name=?`, id)
	defer rows.Close()
	var hash string
	rows.Next()
	if err := rows.Scan(&hash); err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(key))
	if err == nil {
		return true, nil
	}
	return false, nil
}

func Webhook(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form
	id := ""
	key := ""
	if hasKey(form, "id") {
		id = form["id"][0]
	} else {
		return
	}
	if hasKey(form, "key") {
		key = form["key"][0]
	} else {
		return
	}
	authSuccess, err := checkIdKey(id, key)
	if err != nil {
		log.Println(err)
	}
	if authSuccess {
		delete(form, "id")
		delete(form, "key")
		formString, err := json.Marshal(form)
		if err != nil {
			log.Println(err)
		}
		Publish(string(formString))
	} else {
		return
	}
}

func HttpsServer() {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(DOMAIN),
		Cache:      autocert.DirCache("/usr/local/etc/webhook-certs"),
	}

	http.HandleFunc("/", Webhook)

	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}
	go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
	log.Fatal(server.ListenAndServeTLS("", ""))
}
