package socket

import "github.com/gorilla/websocket"

type Client struct {
	socket  *websocket.Conn
	send    chan []byte
	receive chan []byte
}

func (c *Client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.receive <- msg
	}
}
func (c *Client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		// TODO: this could be changed if anything else than images should be transmitted
		err := c.socket.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			return
		}
	}
}
