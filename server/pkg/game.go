package pkg

import (
	"context"
	"errors"
	"log"
	"time"
)

/*
Interface for game's operations
*/
type Game interface {
	///Starts the game up
	Run(at chan GameMessage)
	//Adds a new player to the game
	AddPlayer(connection Connection) error
	//Terminates the game
	EndGame() error
	//Manage message passing
	Manage(message GameMessage) error
}

const DEFAULT_ZONE_WIDTH = 50
const DEFAULT_ZONE_HEIGHT = 50
const DEFAULT_PING = time.Second / 20

/*
Snake Game struct
*/
type SnakeGame struct {
	id string

	players map[string]*Snake

	food *Position

	zone *Zone

	connections []Connection

	running bool

	ctx context.Context

	cancel func()

	InfoChann chan GameMessage

	playersChan chan MessageWrapper
}

func NewSnakeGame(id string, parentContext context.Context) *SnakeGame {
	zone := NewZone(DEFAULT_ZONE_WIDTH, DEFAULT_ZONE_HEIGHT)
	food := RandomPosition()
	zone.IntoLimits(food)
	ctx, cancel := context.WithCancel(parentContext)
	return &SnakeGame{
		id,
		make(map[string]*Snake),
		food,
		zone,
		make([]Connection, 0),
		false,
		ctx,
		cancel,
		nil,
		make(chan MessageWrapper),
	}
}

func (g *SnakeGame) manageMessage(m GameMessage) {
	switch m.Op {
	case JOIN_GAME:
		log.Println("From Game: Joining...")
		g.AddPlayer(m.Conn)
	case STOP_GAME:
		if err := g.EndGame(); err != nil {
			log.Printf("Error stopping game: %s\n", err)
		}
	default:
		log.Printf("Invalid operation sent to game\n")
	}
}

func (g *SnakeGame) Advance() {
	eaten := false
	for _, s := range g.players {
		if !s.Alive && len(s.Positions) == 0 {
			delete(g.players, s.Name)
			continue
		}
		eaten = eaten || s.Move(g.zone, g.food)
		if eaten {
			g.food = RandomPosition()
			g.zone.IntoLimits(g.food)
		}
	}
	for _, c := range g.connections {
		SendDataTo(g, c)
	}

}

func (g *SnakeGame) manageRequest(c Connection, m MessageType) {
	if m.Oint != MOVEMENT_UPDATE {
		// TODO: Add other commands
		return
	}
	data := m.Message.(map[string]interface{})
	direction := int(data["NewDirection"].(float64))
	player := data["Player"].(string)
	if _, ok := g.players[player]; ok {
		g.players[player].ChangeSpeed(Speed{direction, 1})
	}
	if err := SendMovementUpdate(g, player); err != nil {
		log.Printf("Error broadcasting movement update: %s\n", err)
	}
}

/*
Runs the game, should be run inside of a goroutine
*/
func (g *SnakeGame) Run(at chan GameMessage) {
	g.running = true
	g.InfoChann = at
	for g.running {
		select {
		case <-time.After(DEFAULT_PING):
			//Advance the game
			g.Advance()
		case m, ok := <-at:
			if !ok {
				log.Println("Stopping game")
				break
			}
			log.Println("Gotten message")
			g.manageMessage(m)
		case <-g.ctx.Done():
			if err := g.EndGame(); err != nil {
				log.Printf("Error ending game at done context: %s\n", err)
			}
		case m, ok := <-g.playersChan:
			if !ok {
				break
			}
			log.Println("Recieved information from players Chan")
			g.manageRequest(m.Conn, m.Message)

		}

	}
	g.cancel()
}

func (g *SnakeGame) EndGame() error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovering from panic at end game %s; Reason => %s\n", g.id, r)
		}
	}()
	g.running = false
	close(g.InfoChann)
	return nil
}

/*
Starts the game by asking for the snake
name and color.
*/
func (g *SnakeGame) handshakeConnection(connection Connection) error {
	var info SnakeCreateInfo
	if err := connection.Read(g.ctx, &info); err != nil {
		log.Printf("Error handshaking client: %s\n", err)
		return err
	}
	if _, ok := g.players[info.Name]; ok {
		return errors.New("name already taken")
	}
	g.players[info.Name] = SpawnSnake(g.zone, info.Name, info.Color, len(g.players))
	log.Printf("Spawned snake at: %v\n", g.players[info.Name])
	return nil
}

func (g *SnakeGame) AddPlayer(connection Connection) error {
	if !g.running {
		return errors.New("Game is not running")
	}
	if err := g.handshakeConnection(connection); err != nil {
		log.Printf("Error handshaking client: %s\n", err)
		return err
	}
	g.connections = append(g.connections, connection)
	go ManageConnection(g.ctx, connection, g.playersChan)
	return nil
}

func (g *SnakeGame) Manage(m GameMessage) error {
	if !g.running {
		return errors.New("Game is not running")
	}
	g.InfoChann <- m
	return nil
}
