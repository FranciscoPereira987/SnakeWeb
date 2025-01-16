package pkg

import (
	"context"
	"errors"
	"log"
)

/*
Snake server manages all games
it initializes the corresponding goroutines
and manages them
*/
type SnakeServer struct {
	//Mapping of Game.ID => Game
	games map[string]Game
	//Main games channel
	gameChan chan GameMessage
	//Server Ctx
	ctx context.Context
}

func (s *SnakeServer) Send(m GameMessage) error {
	select {
	case <-s.ctx.Done():
		return errors.New("Server is not running")
	default:
		s.gameChan <- m
	}
	return nil
}

/*
Game OPs IDs for the management engine
*/
const (
	CREATE_GAME = iota
	GET_GAME
	JOIN_GAME
	STOP_GAME
)

/*
Message passing struct to manage operations
with the game mapping of the SnakeServer
*/
type GameMessage struct {
	Op int
	Id string
	// TODO: Change so as to be able to test this
	//OPTIONAL: Works with CREATE_GAME, JOIN_GAME
	Conn Connection
	// TODO: Change so as to be able to test this
	//OPTIONAL: Works with GET_GAME
	Game chan<- Game
}

func NewGame(ctx context.Context) *SnakeServer {
	return &SnakeServer{
		make(map[string]Game),
		make(chan GameMessage, 10),
		ctx,
	}
}

func (s *SnakeServer) StopGame(game string) {
	if g, ok := s.games[game]; ok {
		if err := g.EndGame(); err != nil {
			log.Printf("Error closing game: %s", err)
		}
	}
}

func (s *SnakeServer) CreateNewGame(id string) {
	if _, ok := s.games[id]; ok {
		return
	}
	game := NewSnakeGame(id, s.ctx)
	s.games[id] = game
	go game.Run(make(chan GameMessage))
}

func (s *SnakeServer) joinGame(m GameMessage) {
	if g, ok := s.games[m.Id]; ok {
		if err := g.Manage(m); err != nil {
			log.Printf("Error joining gamne: %s", err)
		}
	}
}

/*
Executes the operation encoded in the given message
*/
func (s *SnakeServer) executeMessage(m GameMessage) {
	switch m.Op {
	case STOP_GAME:
		s.StopGame(m.Id)
	case CREATE_GAME:
		log.Println("Creating new game")
		s.CreateNewGame(m.Id)
	case GET_GAME:
		panic("No reason to call this")
	case JOIN_GAME:
		log.Println("Joining someone to a game")
		s.joinGame(m)
	}
}

/*
Shuts the server down
*/
func (s *SnakeServer) shutdown() {
	for _, g := range s.games {
		g.Manage(GameMessage{
			Op: STOP_GAME,
		})
	}
}

/*
Manages the main goroutine of the server
*/
func mainRoutine(server *SnakeServer) {
	running := true
	for running {
		select {
		case <-server.ctx.Done():
			//Ctx is done and I should stop the server
			// TODO: Gracefully stop the server
			server.shutdown()
			running = false
		case g, ok := <-server.gameChan:
			if !ok {
				//Server is stoped
				// TODO: Gracefully stop the server
				server.shutdown()
			} else {
				// Switch over g for the different cases
				server.executeMessage(g)
			}
			running = ok
		}
	}
}

func (s *SnakeServer) Start() {
	go mainRoutine(s)
}

// Stops the server
// Can only be called once
func (s *SnakeServer) Stop() {
	close(s.gameChan)
}

func (s *SnakeServer) GetActiveGames() []string {
	games := make([]string, 0)
	for id, _ := range s.games {
		games = append(games, id)
	}
	return games
}
