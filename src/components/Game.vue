<script setup>
import { Snack } from '@/commons/food';
import Snake, { Position, Speed } from '@/commons/snake';



let startPositions = []
for (let i = 0; i < 100; i++) {
    startPositions = startPositions.concat(new Position(i, 0))
}
const width = 500;
const height = 500;
const startSnaky = new Snake('green', startPositions, 10)
let activeSnack = new Snack(width, height, 10)

const frameDuration = 1000 / 60

const animate = (canvas, last) => {
    
    if ((performance.now() - last) < frameDuration) {
        requestAnimationFrame(() => animate(canvas, last))
    }else {
        const ctx = canvas.getContext('2d')
        ctx.clearRect(0, 0, canvas.width, canvas.height)
        startSnaky.draw(ctx)
        startSnaky.update(canvas.width, canvas.height)
        if (startSnaky.eaten(activeSnack)) {
            activeSnack = new Snack(width, height, 10)
        }
        activeSnack.draw(ctx)
        requestAnimationFrame(() => animate(canvas, performance.now()))
    }
}

const manageInputs = (e) => {
    startSnaky.keydown(e.key) 
}

const load = () => {
    //Basic Setup
    const canvas = document.getElementById('canvas1')
    canvas.width = 500;
    canvas.height = 500
    requestAnimationFrame(() => animate(canvas, performance.now()))
}


window.addEventListener('load', load)
window.addEventListener('keydown', manageInputs)
</script>


<template>
    <canvas id="canvas1">

    </canvas>
</template>