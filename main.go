// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"io"

	"github.com/gookit/goutil/netutil"
	"net/http"
	"os"
	"time"
)

var (
	log = logrus.New()
)

func init(){


	log.Formatter = new(logrus.TextFormatter)                     //default
	//log.Formatter.(*logrus.TextFormatter).DisableColors = true    // remove colors
	//log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove Timestamp from test output
	log.Formatter= &nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
		TimestampFormat: time.RFC3339,
	}
	log.Level = logrus.TraceLevel
	log.Out = os.Stdout

	writer3, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, writer3))

	//log.WithFields(logrus.Fields{
	//	"animal": "walrus",
	//	"number": 0,
	//}).Info("Chat Server Inited.")
	log.Info("Chat Server Inited.")
}


var addr = flag.String("addr", netutil.InternalIP()+":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	log.Info("Listen Serving On: ",*addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
