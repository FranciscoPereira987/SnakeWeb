<script setup>
import { Snack } from '@/commons/food';
import Snake, { Game, Position, Speed } from '@/commons/snake';
import { onMounted, ref } from 'vue';
import LooseScreen from './LoseScreen.vue';
import { snakeData } from '@/commons/store';
import { onDataUpdate, onServerMessage, SendMovement } from '@/commons/proto';



let startPositions = ref([])
for (let i = 0; i < 10; i++) {
        startPositions.value = startPositions.value.concat(new Position(i, 0))
}
const width = 50;
const height = 50;
const snake = ref(new Snake(snakeData.color, startPositions.value, 1, ''))
const activeSnack = ref(new Snack(width, height, 10))

const frameDuration = 1000 / 60

const instantiateSnake = () => {
    // TODO: Check how to do this so that it returns a random position where no snake is on
    let snakePositions = Array(10)
    for (let i = 0; i < 10; i++) {
        snakePositions[i] = new Position(startPositions.value[i].x, startPositions.value[i].y)
    }
    snake.value = new Snake(snakeData.color, snakePositions, 1, snakeData.name)
    snakeData.gameData = new Game(1, activeSnack.value, [snake.value])
}

const updateScreen = (e) => {
    onServerMessage(e.data, width, height)
}

const sendKeydownUpdate = (update) => {
    SendMovement(snakeData.connection, update, snakeData.name)
}

const animate = (canvas, last) => {
    
    if ((performance.now() - last) < frameDuration) {
        requestAnimationFrame(() => animate(canvas, last))
    }else {
        const ctx = canvas.getContext('2d')
        ctx.clearRect(0, 0, canvas.width, canvas.height)
        if (snake.value.stillAlive()) {
            snakeData.gameData.draw(ctx) 
        }else {
            canvas.width = 0
            canvas.height = 0
            return
        }
        requestAnimationFrame(() => animate(canvas, performance.now()))
    }
}

const manageInputs = (e) => {
    snake.value.keydown(e.key, sendKeydownUpdate) 
}

const load = () => {    
    //Basic Setup
    instantiateSnake()
    if (snakeData.connection != null) {
        snakeData.connection.onmessage = updateScreen
    }
    let canvas = document.getElementById('canvas1')
    if(canvas == null) {
        return
    }
    canvas.width = 500;
    canvas.height = 500;
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