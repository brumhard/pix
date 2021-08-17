import { Readable, readable } from "svelte/store";

export const getImageSrc = (delay: number): Readable<string> => {
    return readable("", set => {
        // this function is called once, when the first subscriber to the store arrives
        let socket = new WebSocket(`ws://${window.location.host}${import.meta.env.VITE_WS_PATH}?delay=${delay.toString()}`)

        socket.onmessage = (e) => {
            let reader = new FileReader()
            reader.onloadend = () => set("data:image/png;base64," + reader.result)
            reader.readAsText(e.data)
        }

        // the function we return here will be called when the last subscriber
        // unsubscribes from the store (hence there's 0 subscribers left)
        return () => socket.close()
    })
}