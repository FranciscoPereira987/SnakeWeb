<script setup>
const props = defineProps(["final-score", "snake-name"])

console.log(props)

</script>


<template>
    <v-card class="losing-card" title="You Loose">
        <template v-slot="text">
            <p>Try Again as: {{ props.snakeName }}</p>
            <p>Final Score: {{ props.finalScore }}</p>
        </template>
    </v-card>
</template>


<style>

.losing-card {
   position: absolute;
   max-width: 50vb;
   height: 500;
   top: 50%;
   left:50%;
   transform: translate(-50%, -50%); 
}

</style>