// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/strutil"
	"time"

	//"golang.org/x/exp/slices"
)

type TextStyle struct{
	Color string	`json:"color"`
	BgColor string	`json:"bgcolor"`
	Font  string	`json:"font"`

}
type IconStyle struct {
	Name string `json:"name"`
	Style string `json:"style"`
}

var(
	random_text_styles = []TextStyle{
		{Color: "red",Font: "normal bold 16px Helvetica"},
		{Color: "black",Font: "normal bold 16px Helvetica"},
		{Color: "blue",Font: "normal bold 16px Helvetica"},
		{Color: "green",Font: "normal bold 16px Helvetica"},
		{Color: "Salmon",Font: "normal bold 16px Helvetica"},
		{Color: "Maroon",Font: "normal bold 16px Helvetica"},
		{Color: "Navy",Font: "normal bold 16px Helvetica"},
		{Color: "Purple",Font: "normal bold 16px Helvetica"},
	}

	rand_text_icons = []IconStyle{
		// https://www.w3schools.com/icons/fontawesome5_icons_sports.asp
		{Name: "fas fa-skating",Style: "font-size:24px;color:greenyellow"},
		{Name: "fas fa-futbol",Style: "font-size:24px;color:greenyellow"},
		{Name: "fas fa-snowboarding",Style: "font-size:24px;color:greenyellow"},
		{Name: "fas fa-football-ball",Style: "font-size:24px;color:greenyellow"},
		{Name: "fas fa-basketball-ball",Style: "font-size:24px;color:greenyellow"},

		{Name: "fas fa-frog",Style: "font-size:24px;color:cornflowerblue"},
		{Name: "fas fa-umbrella",Style: "font-size:24px;color:cornflowerblue"},
		{Name: "fas fa-fish",Style: "font-size:24px;color:coral"},
	}

	msglimit_max uint32= 100
)

type Message struct{
	Ok        string
	Timestamp time.Time `json:"-"`
	Time 	string  `json:"time"`
	Nick    string    `json:"nick"`
	Session string	`json:"session"`
	Passwd  string	`json:"passwd"`
	Text 	string	`json:"text"`
	Style		TextStyle	`json:"style"`
	Icon 	IconStyle `json:"icon"`
	Self	bool	`json:"self"`
}

func (m Message)to_json() ([]byte,error)  {
	m.Time = fmt.Sprintf("%02d:%02d",m.Timestamp.Hour(),m.Timestamp.Minute())
	return json.Marshal(m)
}

type Session struct{
	name  string
	first	time.Time
	kick 	time.Time	// 最近一次消息时间
	messages 	[]*Message 	//
	msglimit	uint32		// 最大保留消息条数
	clients  map[*Client] bool
}


// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan struct{ a []byte; b *Client}

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	sessions 	map[string] *Session
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan struct{ a []byte;b *Client }),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		sessions:   make(map[string] *Session),
	}
}

func ( h *Hub) onMessage(client *Client, bytes []byte){
	var msg Message
	if err := json.Unmarshal(bytes,&msg); err!=nil{
		log.Error("message umarshall failed!" ,client.conn.RemoteAddr().String(), string(bytes))
		return
	}
	session_id := msg.Session + msg.Passwd
	if len(strutil.Trim(msg.Session)) == 0{
		return
	}
	ss := h.sessions[session_id]
	if len(msg.Text) ==0 {

		if ss !=nil && client.session != session_id{ // pull history message
			for _,msg := range ss.messages{
				bytes ,_ = msg.to_json()
				//fmt.Println(msg,bytes)
				client.send <- bytes  // 将历史消息重新发给client
			}
			// enter new session
			if _, ok:=ss.clients[client];!ok{
				ss.clients[client] = true
			}
			// leave other sessions
			for _,s := range h.sessions{
				if s == ss{
					continue
				}
				if _,ok:= s.clients[ client ] ; ok {
					delete(s.clients,client)
				}
			}
		}
		client.session = session_id
		return
	}

	msg.Timestamp = time.Now()

	if ss == nil{
		ss = &Session{first: time.Now() , name: session_id,kick: time.Now(),
				messages: []*Message{},
				clients: map[*Client] bool{},
				msglimit: msglimit_max,
			}
	}
	h.sessions[session_id] = ss
	if _, ok:=ss.clients[client];!ok{
		ss.clients[client] = true
	}


	ss.messages = append(ss.messages, &msg)
	if ss.msglimit < uint32(len(ss.messages)) {
		num := uint32(len(ss.messages)) - ss.msglimit
		ss.messages = append(ss.messages[:0], ss.messages[0:num]...)
	}

	msg.Style = arrutil.GetRandomOne(random_text_styles).(TextStyle)
	msg.Icon = arrutil.GetRandomOne(rand_text_icons).(IconStyle)
	var e error

	if bytes,e = msg.to_json(); e!=nil{
	//if bytes,e = json.Marshal(&msg); e!=nil{
		log.Error("Message Marshall failed: ", e.Error())
		return
	}

	sender := client
	for client = range ss.clients{
		if sender == client{
			msgcopy := msg;
			msgcopy.Self = true;
			bytes,_ = msgcopy.to_json()
		}else{
			bytes,_ = msg.to_json()
		}

		select {
			case client.send <- bytes:
		default:
			close(client.send)
			delete(h.clients, client)
			delete(ss.clients,client)
		}
	}

}

// Welcome
func (h *Hub) onClientIncoming(client *Client){
	var msg Message
	msg.Nick = "robot"
	msg.Style = arrutil.GetRandomOne(random_text_styles).(TextStyle)
	msg.Timestamp = time.Now()
	msg.Text = "Welcome To Chat Land..."
	if data,e:= json.Marshal(&msg); e == nil {
		client.send <- data
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.onClientIncoming(client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				for _,ss:= range h.sessions{
					if _, ok:=ss.clients[client];ok{
						delete(ss.clients,client)
					}

				}
			}
		case message := <-h.broadcast:
			h.onMessage(message.b , message.a)
		}
	}
}
