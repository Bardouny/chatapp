package main

import (
	"example/go_funds/db"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader websocket.Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var mydb *db.DataBase = &db.DataBase{
	Users:    make(map[string]*db.User),
	Messages: make([]*db.Message, 700),
}

func main() {
	http.HandleFunc("/", SignUPHandler)
	http.HandleFunc("/chat", ChatInterFaceHandler)
	http.HandleFunc("/ws", websocketHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ERROR at running server .")
	}
	fmt.Println("Server is Running in PORT 8000")
}

// signup end point

func SignUPHandler(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("index.html"))
	tmp.Execute(w, nil)
}

func ChatInterFaceHandler(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("index2.html"))
	data := make(map[string]string)
	name := r.URL.Query().Get("name")
	password := r.URL.Query().Get("password")

	data["name"] = name
	data["password"] = password
	if mydb.Users[name] == nil {
		mydb.Store(name, &db.User{Name: name, Password: password})
	}
	tmp.Execute(w, data)
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("name")
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("faild to upgrade")
		return
	}
	if mydb.Get(userName) == nil {
		fmt.Println("User not found")
		return
	}
	if mydb.Get(userName).Conn != nil {
		fmt.Println("User already connected")
		mydb.Get(userName).Conn.Close()
		mydb.SetConnection(userName, nil)
	}

	mydb.SetConnection(userName, conn)

	go servUser(mydb.Get(userName))
	mydb.UpdateUsersList()
}

func servUser(user *db.User) {
	if user == nil {
		fmt.Println("nil conn has been detected")
		return
	}
	defer db.CloseConnection(user.Conn)

	for {
		_, message, err := user.Conn.ReadMessage()
		if err != nil {
			fmt.Println("Error acuured when reading message")
			break
		}
		mydb.StoreMessage(user.Name, string(message))
	}

	//delete(mydb.Users, user.Name)
	mydb.UpdateUsersList()
}
