<script lang="ts">
    import {Route, Router} from "svelte-routing";
    import Home from "./routes/Home.svelte";
    import Frame from "./routes/Frame.svelte";

    export let url = "";
</script>

<Router url="{url}">
    <Route path="frame" component="{Frame}"/>
    <Route path="/">
        <Home/>
    </Route>
</Router>
