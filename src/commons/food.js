export class Snack {
    constructor(width, height, value) {
        this.x = Math.floor(Math.random() * width)
        this.y = Math.floor(Math.random() * height)
        this.value = value
    }

    draw(context) {
        context.fillStyle = "yellow"
        context.fillRect(this.x, this.y, 10, 10)
    }

    eaten(snakeHead) {
        if(snakeHead == null) return
        let diffX = Math.abs(snakeHead.x-this.x) < 10
        let diffY = Math.abs(snakeHead.y - this.y) < 10
        return (diffX && diffY)
    }
}