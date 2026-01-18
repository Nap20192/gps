package ws

import (
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	conn      *websocket.Conn
	manager   *Manager
	inbound   chan []byte
	outbound  chan []byte
	closeOnce sync.Once
	id        uuid.UUID
}

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

func NewClient(id uuid.UUID, conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		id:       id,
		conn:     conn,
		manager:  manager,
		inbound:  make(chan []byte),
		outbound: make(chan []byte),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.close()
	}()

	c.conn.SetReadLimit(512)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return
	}

	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			slog.Info("Client disconnected", "client_id", c.id)
			return
		}
		c.inbound <- message
	}
}

func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.close()
	}()

	for {
		select {
		case event, ok := <-c.outbound:

			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte("conn shutting down"))
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, event); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(pingInterval))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) close() {
	c.closeOnce.Do(func() {
		close(c.outbound)
		close(c.inbound)
		_ = c.conn.Close()
		c.manager.removeClient(c.id)
	})
}
