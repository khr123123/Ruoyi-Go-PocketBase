<!-- src/App.svelte -->
<script>
    import {wrap} from 'svelte-spa-router/wrap';
    import Router, {replace} from "svelte-spa-router";
    import LoginPage from './pages/LoginPage.svelte';
    import Layout from './Layout.svelte';
    import Toast from './components/Toast.svelte';
    import {toast} from './stores/toastStore';
    import {user} from "./stores/userStore.js";
    import {myPermissions, myRouter} from "./stores/authStore.js";
    import {getUserRouter} from "./api/sysApis.js";
    import {buildTree, extractMenus, extractPermissions} from "./utils/menuUtils.js";
    import ConfirmModal from "./components/ConfirmModal.svelte";

    // 路由守卫
    function authGuard(detail) {
        const token = $user?.token;
        const isLogin = detail.location === '/login';
        if (!token && !isLogin) {
            replace('/login');
            return false;
        }
        if (token && isLogin) {
            replace('/');
            return false;
        }
        return true;
    }

    // 路由配置
    const routes = {
        '/login': wrap({component: LoginPage}),
        '/*': wrap({
            component: Layout,
            conditions: [authGuard]
        })
    }

    // token 每次变化时自动加载权限与菜单
    $: if ($user?.token) {
        loadRouterAndPermission();
    }

    async function loadRouterAndPermission() {
        const data = await getUserRouter();
        const roles = data?.expand?.role || [];
        myPermissions.set(extractPermissions(roles));
        myRouter.set(buildTree(extractMenus(roles)));
    }
</script>

<Router {routes}/>

{#if $toast}
    <Toast message={$toast.message} type={$toast.type} onClose={() => toast.set(null)}/>
{/if}
<ConfirmModal/>