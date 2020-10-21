package handlers

import (
	"log"

	"github.com/JohnnyS318/go-poker/models"
	"github.com/gorilla/websocket"
)

type PlayerConn struct {
	conn   *websocket.Conn
	Out    chan []byte
	In     chan *models.Event
	Events map[string]EventHandler
}

type EventHandler func(*models.Event)

func NewPlayerConn(conn *websocket.Conn) *PlayerConn {
	return &PlayerConn{
		conn:   conn,
		Out:    make(chan []byte),
		In:     make(chan *models.Event),
		Events: make(map[string]EventHandler),
	}
}

func (p *PlayerConn) reader() {
	defer func() {
		p.conn.Close()
	}()
	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {

			log.Printf("Err reader", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WS Message Error: %v", err)
			}
			event := models.NewEvent("LEAVE_LOBBY", "")
			if action, ok := p.Events[event.Name]; ok {
				action(event)
			}
			break
		}

		event, err := models.NewEventFromRaw(message)

		if err != nil {
			log.Printf("Error parsing message %v", err)
			continue
		}

		log.Printf("Got Event %v", event)

		select {
		case p.In <- event:
			log.Printf("Chanelling Event %v", event)
		default:
		}

		if action, ok := p.Events[event.Name]; ok {
			action(event)
		}

	}
}

func (p *PlayerConn) writer() {
	for {
		select {
		case message, ok := <-p.Out:
			if !ok {
				log.Printf("Closing conn due to sending error")
				p.conn.WriteMessage(websocket.CloseMessage, make([]byte, 0))
				return
			}

			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			w.Close()
		}
	}
}

func (p *PlayerConn) On(eventName string, action EventHandler) *PlayerConn {
	p.Events[eventName] = action
	return p
}
