<script lang="ts">
    import {imageSrc} from "./stores";
    import {fade} from "svelte/transition";
</script>

<div>
    {#key $imageSrc}
        <img alt="main frame" transition:fade src={$imageSrc}>
    {/key}
</div>

<style>
    div {
        display: flex;
        position: fixed;
        left: 0px;
        top: 0px;
        width: 100%;
        height: 100%;
        justify-content: center;
        align-items: center;
    }

    img {
        width: auto;
        height: auto;
        max-width: 100%;
        max-height: 100%;
        position: fixed;
    }

    @media (orientation: landscape) {
        img {
            height: 100%;
        }
    }

    @media (orientation: portrait) {
        img {
            width: 100%;
        }
    }
</style>
