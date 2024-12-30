<script setup>
import { snakeData } from '@/commons/store';
import router from '@/router';

const startGame = () => {
    router.push("/game")
}
</script>


<template>
    <v-btn text="Start Game" @click="startGame" />
</template>