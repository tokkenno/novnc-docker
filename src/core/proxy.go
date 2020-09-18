package core

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsToTCP(wsConn *websocket.Conn, tcpConn net.Conn) chan error {
	done := make(chan error, 2)
	go func() {
		defer wsConn.Close()
		defer tcpConn.Close()
		for {
			t, m, err := wsConn.ReadMessage()
			if err != nil {
				done <- err
				return
			}
			if t == websocket.BinaryMessage {
				_, err = tcpConn.Write(m)
				if err != nil {
					done <- err
					return
				}
			} else {
				log.Println("invalid message", t, m)
			}
		}
		done <- nil
	}()
	return done
}

func tcpToWs(tcpConn net.Conn, wsConn *websocket.Conn) chan error {
	done := make(chan error, 2)
	go func() {
		defer wsConn.Close()
		defer tcpConn.Close()
		data := make([]byte, 4096)
		for {
			l, err := tcpConn.Read(data)
			if err != nil {
				done <- err
				return
			}
			err = wsConn.WriteMessage(websocket.BinaryMessage, data[0:l])
			if err != nil {
				done <- err
				return
			}
		}
		done <- nil
	}()
	return done
}

func HandleProxyConnection(config ServerConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsupgrader.Upgrade(w, r, http.Header{"Sec-WebSocket-Protocol": {"binary"}})
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		log.Println(fmt.Sprintf("Websocket connected. Remote: %s", r.RemoteAddr))

		conn2, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
		if err != nil {
			log.Println("TCP connection error:", err)
			conn.WriteJSON(map[string]interface{}{"error": "connect failed"})
			return
		}
		defer conn2.Close()

		done1 := tcpToWs(conn2, conn)
		done2 := wsToTCP(conn, conn2)

		// wait
		log.Println("WebSocket to TCP connection closed:", <-done2)
		log.Println("TCP to WebSocket connection closed:", <-done1)
		log.Println("Disconnected")
	}
}
