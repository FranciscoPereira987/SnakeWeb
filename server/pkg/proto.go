package pkg

import (
	"context"
	"errors"
	"log"
)

const (
	/*
		Data sends all the message about the state of the game
			1. The whole board as a Frame
			2. The scores of each of the players
			3. The speeds of the players
	*/
	DATA = iota + 1
	/*
		Sends an update on the movement of a given player
	*/
	MOVEMENT_UPDATE
	/*
		Sends an update indicating the position, speed, color and name of a new player
	*/
	PLAYER_UPDATE
	/*
		Sends an update of the positioning of the food
	*/
	FOOD_UPDATE
	/*
		Ping to advance the game
	*/
	PING
)

/*
Wrapper around the message sent to and from
the websocket
*/
type MessageType struct {
	Oint    int
	Message interface{}
}

type FoodUpdateMessage struct {
	Food Position
}

type MovementUpdateMessage struct {
	Player       string
	NewDirection Direction
}

type PlayerUpdateMessage struct {
	Player    int
	Color     string
	Name      string
	Movement  MovementUpdateMessage
	Positions []Position
}

type DataUpdateMessage struct {
	Players []PlayerUpdateMessage
	Food    FoodUpdateMessage
}

type DataRequest struct {
	Op int
}

func getPlayerUpdateMessage(g *SnakeGame, name string) PlayerUpdateMessage {
	player := g.players[name]
	positions := make([]Position, len(player.Positions))
	for i, p := range player.Positions {
		positions[i] = *p
	}
	return PlayerUpdateMessage{
		Player: player.Id,
		Color:  player.Color,
		Name:   player.Name,
		Movement: MovementUpdateMessage{
			Player:       player.Name,
			NewDirection: player.HeadSpeed.Dir,
		},
		Positions: positions,
	}
}

func getFoodUpdateMessage(g *SnakeGame) FoodUpdateMessage {
	return FoodUpdateMessage{
		Food: *g.food,
	}

}

/*
Sends to a given player the information of all other players and
the location of the food on the game
*/
func SendDataTo(g *SnakeGame, c Connection) error {
	players := make([]PlayerUpdateMessage, 0)
	for _, p := range g.players {
		players = append(players, getPlayerUpdateMessage(g, p.Name))
	}
	data := DataUpdateMessage{
		players,
		getFoodUpdateMessage(g),
	}

	return c.Write(g.ctx, MessageType{
		DATA,
		data,
	})
}

/*
Broadcasts the position of the food to all players in a game
*/
func SendFoodUpdate(g *SnakeGame) (err error) {
	log.Println("Sending food update")
	message := getFoodUpdateMessage(g)
	for _, c := range g.connections {
		err = errors.Join(c.Write(g.ctx, MessageType{
			FOOD_UPDATE,
			message,
		}), err)
	}
	return
}

/*
Broadcasts an update of movement of a given player
*/
func SendMovementUpdate(g *SnakeGame, name string) error {

	message := MessageType{
		MOVEMENT_UPDATE,
		MovementUpdateMessage{
			g.players[name].Name,
			g.players[name].HeadSpeed.Dir,
		},
	}
	var result error
	for _, c := range g.connections {
		result = errors.Join(result, c.Write(g.ctx, message))
	}

	return result
}

/*
Broadcasts a player info to everyone
*/
func SendPlayerUpdate(g *SnakeGame, name string) (err error) {
	message := getPlayerUpdateMessage(g, name)
	for _, c := range g.connections {
		err = errors.Join(err, c.Write(g.ctx, MessageType{
			PLAYER_UPDATE,
			message,
		}))
	}
	return
}

/*
Broadcasts a Ping Message
*/
func SendPing(g *SnakeGame) (err error) {
	message := MessageType{PING, ""}
	for _, c := range g.connections {
		err = errors.Join(err, c.Write(g.ctx, message))
	}
	return
}

func ManageConnection(ctx context.Context, c Connection, channel chan MessageWrapper) {
	for {
		var wrapped MessageType
		if err := c.Read(ctx, &wrapped); err != nil {
			log.Printf("Error: %s\n", err)
			return
		}
		channel <- MessageWrapper{
			c,
			wrapped,
		}
	}
}

/*
Manages a message according to the protocol rules
*/
func ManageMessage(g *SnakeGame, m MessageWrapper) error {
	if m.Message.Oint == DATA {
		//Return the data the user need
		return nil
	}
	return errors.New("invalid operation from user recieved")
}
