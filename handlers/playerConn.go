package handlers

import (
	"errors"
	"log"

	"github.com/JohnnyS318/go-poker/models"
	"github.com/gorilla/websocket"
)

type PlayerConn struct {
	conn    *websocket.Conn
	Out     chan []byte
	In      chan *models.Event
	Events  map[string]EventHandler
	OnClose func(error)
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WS Message Error: %v", err)
			}
			p.OnClose(err)
			break
		}

		event, err := models.NewEventFromRaw(message)

		if err != nil {
			log.Printf("Error parsing message %v", err)
			continue
		}

		select {
		case p.In <- event:
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
				p.conn.Close()
				p.OnClose(errors.New("Error during message chanelling"))
				return
			}

			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				p.conn.WriteMessage(websocket.CloseMessage, make([]byte, 0))
				p.conn.Close()
				p.OnClose(err)
				return
			}

			w.Write(message)

			w.Close()
		}
	}
}
