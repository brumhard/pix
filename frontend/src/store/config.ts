import {writable, Writable} from "svelte/store";

interface Config {
    delay: number
}

export let config: Writable<Config> = writable({delay: 5})