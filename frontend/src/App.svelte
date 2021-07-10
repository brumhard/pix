<script lang="ts">
    import ImageFrame from './lib/ImageFrame.svelte'
    import {config} from "./store/config";

    $config.delay = 3;
</script>

<main>
    <ImageFrame delay={$config.delay}/>
</main>

<style>
</style>
