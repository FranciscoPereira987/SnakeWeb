package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"websocket/pkg"

	"github.com/coder/websocket"
	"github.com/rs/cors"
)

var opts = websocket.AcceptOptions{
	OriginPatterns: []string{"localhost:*", "127.0.0.1:*"},
}

func HandleGameCreation(server *pkg.SnakeServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var info pkg.NewGameInfo
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&info); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Invalid body: %s", err)))
			return
		}
		if err := server.Send(pkg.GameMessage{
			Op: pkg.CREATE_GAME,
			Id: info.Name,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error while creating new game: %s", err)))
		}
	}
}

func HandleGetActiveGames(server *pkg.SnakeServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		games := server.GetActiveGames()
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(games); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error getting games: %s", err)))
		}
	}
}

func HandleJoinGame(server *pkg.SnakeServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		c, err := websocket.Accept(w, r, &opts)
		if err != nil {
			log.Printf("Error joining game: %s", err)
		}
		log.Println("Joining someone to a game")
		if err := server.Send(pkg.GameMessage{
			Op:   pkg.JOIN_GAME,
			Id:   id,
			Conn: pkg.NewWebsocketConn(c),
		}); err != nil {
			log.Printf("Error joining game: %s\n", err)
		}

	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mux := http.NewServeMux()
	server := pkg.NewGame(ctx)

	mux.HandleFunc("/game/start", HandleGameCreation(server))
	mux.HandleFunc("/game/get", HandleGetActiveGames(server))
	mux.HandleFunc("/game/join/{id}", HandleJoinGame(server))

	server.Start()

	http.ListenAndServe(":8080", cors.Default().Handler(mux))
}
