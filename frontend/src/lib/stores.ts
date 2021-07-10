import {readable} from "svelte/store";

// TODO: check for what ; is used since it seems to work without it
export const imageSrc = readable("", set => {
    // this function is called once, when the first subscriber to the store arrives
    let socket = new WebSocket(import.meta.env.VITE_WS_URL)

    socket.onclose = () => console.log("closed websocket connection")
    socket.onopen = () => console.log("opened websocket connection")
    socket.onmessage = (e) => {
        let reader = new FileReader()
        reader.onloadend = () => set("data:image/png;base64," + reader.result)
        reader.readAsText(e.data)
    }

    // the function we return here will be called when the last subscriber
    // unsubscribes from the store (hence there's 0 subscribers left)
    return socket.close
})