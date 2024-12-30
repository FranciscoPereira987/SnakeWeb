<script setup>
import { Snack } from '@/commons/food';
import Snake, { Position, Speed } from '@/commons/snake';
import { onMounted, ref } from 'vue';
import LooseScreen from './LoseScreen.vue';
import { snakeData } from '@/commons/store';


let startPositions = ref([])
for (let i = 0; i < 10; i++) {
        startPositions.value = startPositions.value.concat(new Position(i, 0))
}
const width = 500;
const height = 500;
const snake = ref(new Snake(snakeData.color, startPositions.value, 10, ''))
const activeSnack = ref(new Snack(width, height, 10))

const frameDuration = 1000 / 60

const instantiateSnake = () => {
    // TODO: Check how to do this so that it returns a random position where no snake is on
    let snakePositions = Array(10)
    for (let i = 0; i < 10; i++) {
        snakePositions[i] = new Position(startPositions.value[i].x, startPositions.value[i].y)
    }
    snake.value = new Snake(snakeData.color, snakePositions, 10, snakeData.name)
}

const animate = (canvas, last) => {
    
    if ((performance.now() - last) < frameDuration) {
        requestAnimationFrame(() => animate(canvas, last))
    }else {
        const ctx = canvas.getContext('2d')
        ctx.clearRect(0, 0, canvas.width, canvas.height)
        if (snake.value.stillAlive()) {
            snake.value.draw(ctx)
            snake.value.update(canvas.width, canvas.height)
            if (snake.value.eaten(activeSnack.value)) {
                activeSnack.value = new Snack(width, height, 10)
            }
        }else {
            canvas.width = 0
            canvas.height = 0
            return
        }
        activeSnack.value.draw(ctx)
        requestAnimationFrame(() => animate(canvas, performance.now()))
    }
}

const manageInputs = (e) => {
    snake.value.keydown(e.key) 
}

const load = () => {    
    //Basic Setup
    instantiateSnake()
    let canvas = document.getElementById('canvas1')
    if(canvas == null) {
        return
    }
    canvas.width = 500;
    canvas.height = 500
    requestAnimationFrame(() => animate(canvas, performance.now()))
}




onMounted(load)
window.addEventListener('keydown', manageInputs)
</script>


<template>
    <canvas id="canvas1">
    </canvas>
    <LooseScreen v-if="!snake.stillAlive()" :reinitiate="load" :snake-name="snake.name" :final-score="snake.score" />
</template>