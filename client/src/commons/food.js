import { SQUARE_TRANSFORMATION } from "./snake"

export class Snack {
    constructor(width, height, value) {
        this.x = Math.floor(Math.random() * width)
        this.y = Math.floor(Math.random() * height)
        this.value = value
    }

    draw(context) {
        context.fillStyle = "yellow"
        context.fillRect(this.x * SQUARE_TRANSFORMATION, this.y *SQUARE_TRANSFORMATION, 10, 10)
    }

    eaten(snakeHead) {
        if(snakeHead == null) return
        let diffX = (snakeHead.x-this.x) == 0
        let diffY = (snakeHead.y - this.y) == 0
        return (diffX && diffY)
    }
}