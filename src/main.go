package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"src/core"
	"strings"
)

func main() {
	if config, error := core.ReadConfigFile(); error == nil {
		workDir, _ := os.Getwd()

		clientPath := path.Join(workDir, "client")
		if config.ClientPath != "" {
			clientPath = config.ClientPath
		}
		clientFs := http.FileServer(http.Dir(clientPath))
		http.Handle("/", clientFs)
		log.Println(fmt.Sprintf("Serving client from %s", clientPath))

		noVncPath := path.Join(workDir, "novnc")
		if config.ClientPath != "" {
			noVncPath = config.NoVNCPath
		}
		noVncFs := http.FileServer(http.Dir(noVncPath))
		http.Handle("/vnc/", http.StripPrefix("/vnc/", noVncFs))
		log.Println(fmt.Sprintf("Serving noVNC from %s", noVncPath))

		for i := range config.Servers {
			config.Servers[i].Proxy = "/ws/" + strings.Replace(strings.ToLower(config.Servers[i].Name), " ", "_", -1)
			http.HandleFunc(config.Servers[i].Proxy, core.HandleProxyConnection(config.Servers[i]))
			log.Println(fmt.Sprintf("Websocket for <%s> installed on %s", config.Servers[i].Name, config.Servers[i].Proxy))
		}

		http.HandleFunc("/servers", core.HandleServerList(config.Servers))

		if config.Port < 1 {
			config.Port = 8084
		}
		listenAddr := ":" + fmt.Sprint(config.Port)
		log.Println(fmt.Sprintf("HTTP server listening on %s", listenAddr))
		http.ListenAndServe(listenAddr, nil)
	} else {
		log.Fatalln("Configuration load error")
	}
}
