package db

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	Name     string
	Password string
	Conn     *websocket.Conn
}

type Message struct {
	From    string
	Content string
	Date    string
}

type ResponseToClient struct {
	Success bool
	Status  string
	Title   string
	Body    any
}

var Mux sync.Mutex

func NewMessage(from, content, date string) *Message {
	return &Message{From: from, Content: content, Date: date}
}

func NewUser(name string, password string) *User {
	return &User{Name: name, Password: password}
}

type DataBase struct {
	Users    map[string]*User
	Messages []*Message
	Length   int
}

func (d DataBase) Store(key string, value *User) {
	Mux.Lock()
	d.Users[key] = value
	d.Length++
	Mux.Unlock()
	d.ShowAllUsers()
}

func (d DataBase) Get(key string) *User {
	return d.Users[key]
}

func (d DataBase) ShowAllUsers() {
	for _, user := range d.Users {
		fmt.Println("========")
		fmt.Println(user.Name)
		fmt.Println("--------")

	}
}

func (d DataBase) Delete(key string) {
	Mux.Lock()
	delete(d.Users, key)
	//d.Length--
	Mux.Unlock()
}

func CloseConnection(conn *websocket.Conn) {
	if conn != nil {
		err := conn.Close()
		if err != nil {
			fmt.Println("error when closing connection")
		}
	}
}

func (d DataBase) SetConnection(key string, conn *websocket.Conn) {
	if d.Users[key] != nil {
		Mux.Lock()
		d.Users[key].Conn = conn
		Mux.Unlock()
		return
	}

	err := conn.WriteMessage(websocket.TextMessage, []byte("connection without user"))
	if err != nil {
		fmt.Println("err ocuured when sendeng conn doasent have users")
	}
}

func (d DataBase) PublishNewMessage(msg *Message) {
	res := ResponseToClient{Success: true, Status: "200", Title: "new-message", Body: msg}
	msgAsJSON, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Faild to convert to json")
	}

	for _, user := range d.Users {

		if user.Conn == nil {
			continue
		}
		err = user.Conn.WriteJSON(msgAsJSON)

		if err != nil {
			fmt.Println("cannot send message")
		}
	}
}

func (d DataBase) UpdateUsersList() {
	msg := ResponseToClient{Success: true, Status: "200", Title: "update-users", Body: d.Users}

	res, err := json.Marshal(msg)

	if err != nil {
		return
	}

	for _, user := range d.Users {
		if user.Conn == nil {
			continue
		}
		err = user.Conn.WriteJSON(res)
		if err != nil {
			fmt.Println("error acuured when updating list")
		}
	}
}

func (d DataBase) StoreMessage(from, content string) {
	msg := NewMessage(from, content, time.TimeOnly)
	d.Messages = append(d.Messages, msg)
	d.PublishNewMessage(msg)
}
