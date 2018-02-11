package iot

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

const (
	ACCOUNT_DB string = `/usr/local/etc/webhook-account/passwd.db`
)

type WebhookHandler struct {
	broker *Broker
	topic  string
}

func NewWebhookHandler(broker *Broker) *WebhookHandler {
	return &WebhookHandler{broker: broker}
}

func (handler *WebhookHandler) SetTopic(topic string) {
	handler.topic = topic
}

func hasKey(dict map[string][]string, key string) bool {
	if _, ok := dict[key]; ok {
		return true
	}
	return false
}

func checkIdKey(id string, key string) (bool, error) {
	db, err := sql.Open("sqlite3", ACCOUNT_DB)
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

func (handler *WebhookHandler) PublishPost(w http.ResponseWriter, r *http.Request) {
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
		handler.broker.Publish(handler.topic, string(formString))
	} else {
		return
	}
}
