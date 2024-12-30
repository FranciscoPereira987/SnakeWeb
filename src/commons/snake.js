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

    collided(otherPosition) {
        let collX = Math.abs(otherPosition.x-this.x) == 0
        let collY = Math.abs(otherPosition.y-this.y)  == 0
        return collX && collY
    }
}

const BODY_PART_SEPARATION = 10

/*
    The snake should not grow in size.
    The size is increased when the snake accelerates
    You can accelerate the snake and thus you increase the size
*/
export default class Snake {
    constructor(color, positions, speed) {
        this.width = BODY_PART_SEPARATION
        this.height = BODY_PART_SEPARATION
        this.size = positions.length
        this.color = color
        this.positions = positions
        this.speedModulus = speed   
        this.speed = new Speed(-1, 0)
        this.dead = false
        this.score = 0
        this.name = "snake"
    }

    //Updates the position of the snack and checks if it kills itself
    update(maxWidth, maxHeight) {
        if (this.dead && this.size > 0) {
            this.positions.pop()
            this.size--
        }
        if (this.dead) return
        let newHeadPosition = new Position(this.positions[0].x, this.positions[0].y)
        newHeadPosition.move(new Speed(this.speedModulus*this.speed.x, this.speedModulus*this.speed.y), maxWidth, maxHeight)
        for (let i = this.size-1; i >= (this.speedModulus)/BODY_PART_SEPARATION; i--) {
            this.positions[i].update(this.positions[i-(this.speedModulus)/BODY_PART_SEPARATION])
            this.dead = this.dead | this.positions[i].collided(newHeadPosition) 
        }
        for (let i = (this.speedModulus)/BODY_PART_SEPARATION-1; i >= 0; i--) {
            this.positions[i].update(this.positions[0])
            let speedUpdate = new Speed((this.speedModulus-i) * this.speed.x, (this.speedModulus-i) * this.speed.y)
            this.positions[i].move(speedUpdate, maxWidth, maxHeight)
            if (i != 0) {
                this.dead = this.dead | this.positions[i].collided(newHeadPosition)
            }
        }
    }

    draw(context) {
       this.positions.forEach(position => {
        context.fillStyle = this.dead ? "red" : this.color
        position.draw(context, this.width, this.height, this.color)
       })
       context.font = "18px arial"
       context.fillStyle = "black"
       context.fillText(this.name, this.positions[0].x, this.positions[0].y)
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