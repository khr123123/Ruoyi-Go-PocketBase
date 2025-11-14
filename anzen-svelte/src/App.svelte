<script>
    import {wrap} from 'svelte-spa-router/wrap';
    import Router from "svelte-spa-router";
    import LoginPage from './pages/LoginPage.svelte';
    import Layout from './Layout.svelte';
    // 主路由
    const routes = {
        '/login': wrap({component: LoginPage}),
        '/*': wrap({component: Layout}) // Layout 包裹所有其他路由
    }
</script>
<Router {routes}/>

<style global>
    :global(*) {
        box-sizing: border-box;
    }
    :global(body) {
        margin: 0;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
        background-color: #f7f7f7;
    }
</style>