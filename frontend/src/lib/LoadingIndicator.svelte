<!--copied from https://dev.to/promatia/svelte-loading-state-indicator-promises-5c5n-->
<div class="LoadingIndicator">
    {text}
</div>
<script>
    import {onMount} from 'svelte'

    let text = ''

    async function loop() {
        if (text.length === 3) text = ''
        text += '.'
        await new Promise(resolve => setTimeout(resolve, 300))
        loop()
    }

    onMount(() => {
        loop()
    })

</script>
<style>
    .LoadingIndicator {
        display: flex;
        background: rgba(0, 0, 0, 0.2);
        justify-content: center;
        align-items: center;
        text-align: center;
        position: absolute;
        top: 0;
        bottom: 0;
        right: 0;
        left: 0;
        z-index: 20;
        font-size: 60px;
    }
</style>