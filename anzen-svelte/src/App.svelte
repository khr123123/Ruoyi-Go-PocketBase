<script>
    import {wrap} from 'svelte-spa-router/wrap';
    import Router from "svelte-spa-router";
    import LoginPage from './pages/LoginPage.svelte';
    import Toast from './components/Toast.svelte';
    import {toast} from './stores/toastStore';
    import Layout from './Layout.svelte';
    // 主路由
    const routes = {
        '/login': wrap({component: LoginPage}),
        '/*': wrap({component: Layout}) // Layout 包裹所有其他路由
    }
</script>
<Router {routes}/>
<!-- 🔥 全局 Toast -->
{#if $toast}
    <Toast message={$toast.message} type={$toast.type} onClose={() => toast.set(null)}/>
{/if}
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