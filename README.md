# Snake Online

This is a very basic online (web-based) multiplayer snake game.

The server is written in Go and the front-end is written in Javascript with vue

This project was built as a way of learning about websockets, so the game is not the main focus here.

## Running the project

To run the front-end, execute the following command from the root folder of this project:

```bash
cd client && npm install && npm run dev
```

From the root folder of this project, execute the following command from the root folder of the project as well:

```bash
cd server && go run main.go
```
