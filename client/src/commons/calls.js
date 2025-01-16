import { snakeData } from "./store"

export const LOCAL_SERVER = "http://127.0.0.1:8080/game/"

const GAME_LIST = "get"
const GAME_CREATE = "start"
const GAME_JOIN = "join"


export async function GetGames(server) {
    const url = server + GAME_LIST
    const response = await fetch(url)
    return await response.json()
}

export async function CreateGame(server, gameName) {
    const url = server + GAME_CREATE
    await fetch(url, {
        method: "POST",
        body: JSON.stringify({
            name: gameName
        })
    })
}

export async function JoinGame(server, gameName) {
    const url = server + GAME_JOIN + "/" + gameName
    const socket = new WebSocket(url)
    socket.onopen = (_) => {
        socket.send(JSON.stringify({
            name: snakeData.name,
            color: snakeData.color,
        }))
    }
    snakeData.connection = socket
}