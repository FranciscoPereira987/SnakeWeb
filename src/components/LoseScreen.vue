<script setup>
import router from '@/router';

const props = defineProps(["final-score", "snake-name", "reinitiate"])

const backToMain = () => {
    router.push('/')
}

</script>


<template>
    <v-card class="losing-card" title="You Lose">
        <template v-slot:text>
            <p>Final Score: {{ props.finalScore }}</p>
        </template>
        <template v-slot:actions>
            <v-btn variant="tonal" @click="backToMain">Go Back</v-btn>
            <v-btn variant="tonal" @click="props.reinitiate">Try Again as: {{ props.snakeName }}</v-btn>
        </template>
    </v-card>
</template>


<style>

.losing-card {
   position: absolute;
   display: flex;
   justify-content: space-around;
   max-width: 50vb;
   height: 500;
   top: 50%;
   left:50%;
   transform: translate(-50%, -50%); 
}

</style>