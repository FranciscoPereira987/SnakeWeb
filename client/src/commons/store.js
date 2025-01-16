import { reactive } from "vue";

export const snakeData = reactive({
    name: "pere",
    game: "peregame",
    color: "green",
    connection: null,
    //Snake Game that gets updated based on messages
    //exchanged with the server
    gameData: null
})