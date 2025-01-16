/*
    Game protocol definition
*/

import { Snack } from "./food"
import Snake, { Position, Speed } from "./snake"
import { snakeData } from "./store"

const DATA = 1
const MOVEMENT_UPDATE = 2
const PLAYER_UPDATE = 3
const FOOD_UPDATE = 4
const PING = 5

const FoodUpdateMessage = {
    Food: {
        X: 0,
        Y: 0,
    }
}

const MovementUpdateMessage = {
    Player: "",
    NewDirection: 0
}

const PlayerUpdateMessage = {
    Player: 0,
    Color: "",
    Name: "",
    Movement: MovementUpdateMessage,
    Positions: []
}

const DataUpdateMessage = {
    Players: [],
    Food: FoodUpdateMessage
}

function getPlayerPositions(positions) {
    return positions.map((p) => {
        return new Position(p.X, p.Y)
    })
}

function getPlayerSpeed(speed) {
    switch(speed.NewDirection) {
        case 0:
            return new Speed(0, -1)
        case 1:
            return new Speed(0, 1)
        case 2:
            return new Speed(-1, 0)
        case 3:
            return new Speed(1, 0)
    }
}

export function SendMovement(to, direction, player) {
    console.log("Here")
    to.send(JSON.stringify({
    Oint: MOVEMENT_UPDATE,    
    Message: {
        Player: player,
        NewDirection: direction,
        }
    }))
}

/*
    Manages data sent from the main server
*/
export function onServerMessage(data, width, length) {
    const parsed = JSON.parse(data)
    if (parsed.Oint == PING) {
        snakeData.gameData.Advance(width, length)
    }else if (parsed.Oint == DATA) {
        onDataUpdate(parsed.Message)
    }else if (parsed.Oint == MOVEMENT_UPDATE) {
        onMovementUpdate(parsed.Message)
    }else if (parsed.Oint == FOOD_UPDATE) {
        onFoodUpdate(parsed.Message)
    }
    else {
        console.log("Invalid message recieved: " + parsed.Oint)
    }
}

export function onFoodUpdate(data) {
    const food = new Snack(50, 50, 10)
    food.x = data.Food.X
    food.y = data.Food.Y
    snakeData.gameData.food = food
}

/*
    Manages changes made when a player changed its movement
*/
export function onMovementUpdate(data) {
    const player = data.Player 
    const newSpeed = getPlayerSpeed(data)
    snakeData.gameData.ChangePlayerDirection(player, newSpeed)
}

/*
    Manages the changes in data
*/
export function onDataUpdate(data) {
    
    const allPlayers = data.Players.map((p) => {
        return new Snake(p.Color,
            getPlayerPositions(p.Positions), 
            getPlayerSpeed(p.Movement),
            p.Name)
   }) 

   const food = new Snack(50, 50, 10)
   food.x = data.Food.Food.X
   food.y = data.Food.Food.Y
   
   snakeData.gameData.players = allPlayers
   snakeData.gameData.food = food
}
