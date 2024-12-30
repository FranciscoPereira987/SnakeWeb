<script setup>
import { snakeData } from '@/commons/store';
import router from '@/router';

const startGame = () => {
    
    router.push("/game")
}

const colors = [
    {color: "#4CAF50", name: "Green"},
    {color: "#9C27B0", name: "Violet"},
    {color: "#E91E63", name: "Pink"},
    {color: "#FFFFFF", name: "white"}
]
</script>


<template>
    <div class="main-welcome-div">
        <v-card class="welcome-card" title="Snake Web">
            <template v-slot:default>
                <v-text-field label="Player Name" v-model:model-value="snakeData.name"></v-text-field>
                <v-chip-group v-model:model-value="snakeData.color">
                   <v-chip v-for="color in colors" :color="color.color" :text="color.name" :key="color.color" :value="color.color"/> 
                </v-chip-group>
            </template>
            <template v-slot:actions>
                <v-btn text="Start Game" @click="startGame" block/>
            </template>
        </v-card>
    </div>
</template>


<style>
.main-welcome-div {
    display: flex;
    flex: 1;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%
}

</style>