package pkg

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Connection interface {
	Read(ctx context.Context, t interface{}) error
	Write(ctx context.Context, t interface{}) error
	Close(s websocket.StatusCode, reason string) error
	Run(ctx context.Context)
}

/*
Wrapper around a websocket connection
*/
type WebsocketConn struct {
	conn    *websocket.Conn
	Channel chan MessageWrapper
}

type SnakeCreateInfo struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

const (
	EMPTY = iota
	FOOD
)

type NewGameInfo struct {
	Name string `json:"name"`
}

/*
Frame indicating the state of the game
each square is represented by a Byte
if the square is empty => 0
if food is on the square => 1
if player n is on the square => 1 + n
*/
type Frame struct {
	Data   string `json:"data"`
	Scores []int  `json:"scores"`
}

/*
Update that a new player has got online
*/
type NewPlayerInfo struct {
	SnakeCreateInfo
	Id int
}

/*
Wrapper to send the connction alongside the Message
*/
type MessageWrapper struct {
	Conn    Connection
	Message MessageType
}

/*
Returns the index for the frame of the given position
*/
func translateInto(position *Position, zone *Zone) int {
	return position.X + zone.x*position.Y
}

func BuildFrame(game *SnakeGame) *Frame {
	data := make([]byte, game.zone.x*game.zone.y)
	scores := make([]int, len(game.players))
	data[translateInto(game.food, game.zone)] = FOOD
	for _, p := range game.players {
		scores[p.Id] = p.Score
		for _, position := range p.Positions {
			data[translateInto(position, game.zone)] = byte(p.Id) + 1
		}
	}
	return &Frame{
		hex.EncodeToString(data),
		scores,
	}
}

func NewWebsocketConn(conn *websocket.Conn) *WebsocketConn {
	return &WebsocketConn{
		conn,
		nil,
	}
}

func routine(w *WebsocketConn, ctx context.Context) {
	go func() {
		defer close(w.Channel)
		for {
			var message MessageType
			if err := wsjson.Read(ctx, w.conn, &message); err != nil {
				log.Printf("Could not recover message from socket: %s\n", err)
				return
			}
			w.Channel <- MessageWrapper{
				w,
				message,
			}
		}
	}()
}

/*
If Run is called, Read should no longer be called. The subroutine will stop when the context is Done
and the socket channel will be closed
*/
func (w *WebsocketConn) Run(ctx context.Context) {
	routine(w, ctx)
}

func (w *WebsocketConn) Read(ctx context.Context, t interface{}) error {
	return wsjson.Read(ctx, w.conn, t)
}

func (w *WebsocketConn) Write(ctx context.Context, t interface{}) error {
	return wsjson.Write(ctx, w.conn, t)
}

func (w *WebsocketConn) Close(s websocket.StatusCode, reason string) error {
	defer func() error {
		if r := recover(); r != nil {
			return fmt.Errorf("recovered from panic because of: %v", r)
		}
		return nil
	}()
	return w.conn.Close(s, reason)
}
