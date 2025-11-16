<!-- src/App.svelte -->
<script>
    import {wrap} from 'svelte-spa-router/wrap';
    import Router, {replace} from "svelte-spa-router";
    import LoginPage from './pages/LoginPage.svelte';
    import Toast from './components/Toast.svelte';
    import {toast} from './stores/toastStore';
    import Layout from './Layout.svelte';
    import {user} from "./stores/userStore.js";
    import {onMount} from 'svelte';
    import {getUserRouter} from "./api/menuApis.js";
    import {extractPermissions, myPermissions, myRouter} from "./stores/authStore.js";
    // 订阅 user（包含 token）
    let currentUser;
    user.subscribe(v => currentUser = v);

    // 路由守卫
    function authGuard(detail) {
        const hasToken = currentUser?.token && currentUser.token !== "";
        if (!hasToken && detail.location !== '/login') {
            replace('/login');
            return false;
        }
        if (hasToken && detail.location === '/login') {
            replace('/');
            return false;
        }
        return true;
    }

    // 主路由
    const routes = {
        '/login': wrap({component: LoginPage}),
        '/*': wrap({
            component: Layout,
            conditions: [authGuard]
        })
    }
    onMount(async () => {
        const hasToken = currentUser?.token && currentUser.token !== "";
        if (hasToken) {
            let res = await getUserRouter();
            myRouter.set(res.data || []);
            myPermissions.set(extractPermissions(res.data) || [])
        }
    });
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