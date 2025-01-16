export class Speed {
    constructor(x, y) {
        this.x = x
        this.y = y
    }

    update(position) {
        position.x += this.x
        position.y += this.y
    }

}

export class Position {
    constructor(x, y) {
        this.x = x 
        this.y = y 
    }

    applyLimits(maxWidth, maxHeight) {
        while (this.x < 0 ) {
            this.x += maxWidth
        }
        while (this.y < 0) {
            this.y += maxHeight
        }
        this.x = this.x % maxWidth
        this.y = this.y % maxHeight
    }

    update(newPosition) {
        //Updates the position and changes speed
        this.x = newPosition.x
        this.y = newPosition.y
    }

    move(speed, maxWidth, maxHeight) {
        this.x += speed.x
        this.y += speed.y
        this.applyLimits(maxWidth, maxHeight)
    }

    draw(context, width, height) {
        context.fillRect(this.x * SQUARE_TRANSFORMATION, this.y * SQUARE_TRANSFORMATION, width, height)
    }

    collided(otherPosition) {
        let collX = Math.abs(otherPosition.x-this.x) == 0
        let collY = Math.abs(otherPosition.y-this.y)  == 0
        return collX && collY
    }
}

const BODY_PART_SEPARATION = 10
export const SQUARE_TRANSFORMATION = 10


export class Game {
    constructor(snakeId, food, players) {
        this.player = snakeId
        this.food = food
        this.players = players
        this.width = 50
        this.height = 50
    }

    ChangePlayerDirection(player, newDirection) {
        const playerSnake = this.players.filter((p) => p.name == player).pop()
        playerSnake.speed = newDirection
    }

    ChangeFoodPosition(newFood) {
        this.food = newFood
    }

    /*
        Advances the game, moving all the snakes
        in this case just one
    */
    Advance(maxWidth, maxLength) {
        this.players.forEach((snake) => {
            snake.update(maxWidth, maxLength)
            snake.eaten(this.food)
        })
    }   
    
    /*
        Draws the Game in the canvas
    */
    draw(context){
        this.players.forEach((snake) => snake.draw(context))
        this.food.draw(context)
   }
}

/*
    The snake should not grow in size.
    The size is increased when the snake accelerates
    You can accelerate the snake and thus you increase the size
*/
export default class Snake {
    constructor(color, positions, speed, name) {
        this.width = BODY_PART_SEPARATION
        this.height = BODY_PART_SEPARATION
        this.size = positions.length
        this.color = color
        this.positions = positions
        this.speedModulus = speed   
        this.speed = new Speed(-1, 0)
        this.dead = false
        this.score = 0
        this.name = name 
    }

    //Updates the position of the snack and checks if it kills itself
    update(maxWidth, maxHeight) {
        if (this.dead && this.size > 0) {
            this.positions.pop()
            this.size--
        }
        if (this.dead) return
        let newHeadPosition = new Position(this.positions[0].x, this.positions[0].y)
        newHeadPosition.move(this.speed, maxWidth, maxHeight)
        for (let i = this.size-1; i > 0; i--) {
            this.positions[i].update(this.positions[i - 1])
            this.dead = this.dead | this.positions[i].collided(newHeadPosition) 
        }
        this.positions[0] = newHeadPosition
    }   

    draw(context) {
       this.positions.forEach(position => {
        context.fillStyle = this.dead ? "red" : this.color
        position.draw(context, this.width, this.height, this.color)
       })
       context.font = "18px arial"
       context.fillStyle = "black"
       context.fillText(this.name, this.positions[0].x * SQUARE_TRANSFORMATION, this.positions[0].y * SQUARE_TRANSFORMATION)
    }

    stillAlive() {
        return !this.dead || this.size > 0
    }

    grow(value) {
        this.score += value
        for (let i = 0; i < (value)/BODY_PART_SEPARATION; i++){
            let offsetX = this.positions[this.size-2].x - this.positions[this.size-1]
            let offsetY = this.positions[this.size-2].y - this.positions[this.size-1]
            let newPosition = new Position(this.positions[this.size-1]-offsetX, this.positions[this.size-1]-offsetY)
            this.positions = this.positions.concat(newPosition)  
            this.size++
        }
    }

    //Returns true if a snack is eaten
    eaten(snack) {
        let eaten = snack.eaten(this.positions[0])
        if (eaten) {
            this.grow(snack.value)
        }
        return eaten
    }

    keydown(key, callback) {
        switch (key) {
            case "ArrowUp":
            case "w":
            case "W":
                callback(0)
                break;
            case "ArrowDown":
            case "s":
            case "S":
                callback(1)
                break;
            case "ArrowLeft":
            case "a":
            case "A":
                callback(2)
                break;
            case "ArrowRight":
            case "d":
            case "D":
                callback(3)
                break
        }
    } 
}