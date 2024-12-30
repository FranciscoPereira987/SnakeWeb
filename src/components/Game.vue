<script setup>
import { Snack } from '@/commons/food';
import Snake, { Position, Speed } from '@/commons/snake';
import { ref } from 'vue';
import LooseScreen from './LooseScreen.vue';



let startPositions = []
for (let i = 0; i < 10; i++) {
    startPositions = startPositions.concat(new Position(i, 0))
}
const width = 500;
const height = 500;
const snake = ref(new Snake('green', startPositions, 10))
const activeSnack = ref(new Snack(width, height, 10))

const frameDuration = 1000 / 60

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
    const canvas = document.getElementById('canvas1')
    if(canvas == null) {
        return
    }
    canvas.width = 500;
    canvas.height = 500
    requestAnimationFrame(() => animate(canvas, performance.now()))
}


window.addEventListener('load', load)
window.addEventListener('keydown', manageInputs)
</script>


<template>
    <p>{{ snake.score }}</p>
    <canvas v-if="snake.stillAlive() && false" id="canvas1">
    </canvas>
    <LooseScreen v-else :snake-name="snake.name" :final-score="snake.score" />
</template>