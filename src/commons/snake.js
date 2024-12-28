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
        if (this.x >= maxWidth) {
            this.x -= maxWidth
        }else if (this.x < 0) {
            this.x += maxWidth
        }
        if (this.y >= maxHeight) {
            this.y -= maxHeight
        }else if (this.y < 0) {
            this.y += maxHeight
        }
    }

    update(newPosition) {
        //Updates the position and changes speed
        this.x = newPosition.x
        this.y = newPosition.y
    }

    move(speed, maxWidth, maxHeight) {
        this.applyLimits(maxWidth, maxHeight)
        this.x += speed.x
        this.y += speed.y
    }

    draw(context, width, height) {
        context.fillRect(this.x, this.y, width, height)
    }
}

/*
    The snake should not grow in size.
    The size is increased when the snake accelerates
    You can accelerate the snake and thus you increase the size
*/
export default class Snake {
    constructor(color, positions, speed) {
        this.width = 10
        this.height = 10
        this.size = positions.length
        this.color = color
        this.positions = positions
        this.speedModulus = speed   
        this.speed = new Speed(-1, 0)
    }

    update(maxWidth, maxHeight) {
        for (let i = this.size-1; i >= this.speedModulus; i--) {
            this.positions[i].update(this.positions[i-this.speedModulus])
        }
        for (let i = this.speedModulus-1; i >= 0; i--) {
            this.positions[i].update(this.positions[0])
            let speedUpdate = new Speed((this.speedModulus-i) * this.speed.x, (this.speedModulus-i) * this.speed.y)
            this.positions[i].move(speedUpdate, maxWidth, maxHeight)
        }
    }

    draw(context) {
       this.positions.forEach(position => {
        context.fillStyle = this.color
        position.draw(context, this.width, this.height, this.color)
       })
    }

    grow(value) {
        for (let i = 0; i < value; i++){
            let offsetX = this.positions[this.size-2].x - this.positions[this.size-1]
            let offsetY = this.positions[this.size-2].y - this.positions[this.size-1]
            let newPosition = new Position(this.positions[this.size-1]-offsetX, this.positions[this.size-1]-offsetY)
            this.positions = this.positions.concat(newPosition)  
            this.size++
        }
    }

    eaten(snack) {
        let eaten = snack.eaten(this.positions[0])
        if (eaten) {
            this.grow(snack.value)
        }
        return eaten
    }

    keydown(key) {
        switch (key) {
            case "ArrowUp":
            case "w":
            case "W":
                this.speed = new Speed(0, -1)
                break;
            case "ArrowDown":
            case "s":
            case "S":
                this.speed = new Speed(0, 1)
                break;
            case "ArrowLeft":
            case "a":
            case "A":
                this.speed = new Speed(-1, 0)
                break;
            case "ArrowRight":
            case "d":
            case "D":
                this.speed = new Speed(1, 0)
                break
        }
    } 
}