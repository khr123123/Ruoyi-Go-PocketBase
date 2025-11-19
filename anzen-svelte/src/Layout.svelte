<script>
    import Router, {link, location, push} from "svelte-spa-router";
    import {wrap} from "svelte-spa-router/wrap";
    import About from './pages/About.svelte';
    import * as Icons from './lib/icons/index.js';
    import {Settings, SettingsAlt} from './lib/icons/index.js';
    import {clearUser, user} from './stores/userStore.js';
    import {showToast} from './stores/toastStore.js';
    import {clearPermissions, myRouter} from './stores/authStore.js';
    import {logout} from "./api/sysApis.js";
    import {slide} from 'svelte/transition';
    import active from 'svelte-spa-router/active';

    const AVATAR_PREFIX = "http://127.0.0.1:8090/api/files/_pb_users_auth_/";

    let searchQuery = '';
    let filteredRouter = [];
    let expandedMenus = {};

    // 动态路由生成
    const modules = import.meta.glob('./pages/**/*.svelte');

    $: layoutRoutes = generateRoutes($myRouter);

    function generateRoutes(menuData) {
        const routes = {
            '/about': wrap({component: About}),
        };
        if (!menuData) return routes;

        menuData.forEach(menu => {
            menu.children?.forEach(child => {
                if (child.url && child.menuType === 'C') {
                    routes[child.url] = wrap({asyncComponent: modules[`./pages${child.url}.svelte`]});
                }
            });
        });

        return routes;
    }

    // 初始化展开状态 & 过滤菜单
    $: if ($myRouter?.length) {
        expandedMenus = $myRouter.reduce((acc, menu) => ({...acc, [menu.id]: true}), {});
        filteredRouter = filterRouter($myRouter, searchQuery);
    }

    // 搜索过滤
    $: filteredRouter = filterRouter($myRouter, searchQuery);

    function filterRouter(routerData, q) {
        if (!routerData) return [];
        const key = (q || '').trim().toLowerCase();
        if (!key) return routerData;

        return routerData
              .map(menu => {
                  const m = {...menu};
                  if (m.children?.length) {
                      const children = m.children.filter(ch =>
                            (ch.menuName || '').toLowerCase().includes(key) ||
                            (ch.url || '').toLowerCase().includes(key)
                      );
                      if (children.length) m.children = children;
                  }
                  if ((m.menuName || '').toLowerCase().includes(key) || (m.children?.length)) return m;
              })
              .filter(Boolean);
    }

    function toggleExpand(menuId) {
        expandedMenus = {...expandedMenus, [menuId]: !expandedMenus[menuId]};
    }

    function handleLogout() {
        logout();
        clearUser();
        clearPermissions();
        showToast("Logout successful!", "info");
        setTimeout(() => push('/login'), 800);
    }
</script>

<div class="flex h-screen bg-gray-50">
    <!-- Sidebar -->
    <aside class="flex flex-col w-48 bg-white border-r border-gray-200 shadow-sm">
        <!-- Logo & 搜索 -->
        <div class="p-4 border-b border-gray-200">
            <div class="flex items-center justify-center mb-3 text-gray-800">
                <img src="../public/logo.jpg" alt="RY" class="w-12 rounded-sm">
                <span class="ml-2 text-lg font-mono tracking-wide">RUOYI-PB</span>
            </div>

            <div class="relative">
                <input
                      type="text"
                      placeholder="Search menus..."
                      class="w-full px-3 py-2 text-sm border rounded-lg bg-gray-50 border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                      bind:value={searchQuery}
                />
                {#if searchQuery}
                    <button
                          class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                          on:click={() => searchQuery = ''}
                          aria-label="Clear search"
                    >✕
                    </button>
                {/if}
            </div>
        </div>

        <!-- 菜单列表 -->
        <nav class="flex-1 p-2 overflow-y-auto">
            {#if filteredRouter.length}
                {#each filteredRouter as item (item.id)}
                    {#if item.menuType === 'M'}
                        <div class="mb-2">
                            <!-- 父菜单 -->
                            <button
                                  class="flex items-center w-full gap-3 px-3 py-2 text-sm font-medium text-gray-700 rounded-lg hover:bg-gray-100 transition cursor-pointer"
                                  on:click={() => {
                          if(item.children.length===0 && item.url.startsWith("http")){
                              window.open(item.url)
                          }
                          toggleExpand(item.id)
                      }}
                                  aria-expanded={expandedMenus[item.id] ? "true" : "false"}
                            >
                <span class="flex items-center justify-center w-6 text-gray-500">
                  <svelte:component this={Icons[item.icon]} size={16} className="text-black"/>
                </span>
                                <span class="flex-1 text-left truncate">{item.menuName}</span>
                                {#if item.children.length !== 0}
                                    <svg class="w-4 h-4 text-gray-400 transform transition-transform"
                                         style="rotate: {expandedMenus[item.id] ? '90deg' : '0deg'}" viewBox="0 0 24 24"
                                         fill="none"
                                         stroke="currentColor" stroke-width="2">
                                        <path d="M9 6l6 6-6 6"/>
                                    </svg>
                                {/if}
                            </button>

                            <!-- 子菜单 -->
                            {#if expandedMenus[item.id] && item.children?.length}
                                <div class="ml-3 mt-1 space-y-1 border-l-2 border-gray-300 pl-2"
                                     transition:slide>
                                    {#each item.children as child (child.id)}
                                        {#if child.menuType === 'C'}
                                            <a
                                                  href={child.url}
                                                  use:link
                                                  use:active={{ path: child.url, className: 'bg-gray-100' }}
                                                  class="flex items-center gap-3 px-3 py-2 text-sm rounded-lg transition-all hover:bg-gray-100"
                                            >
                                                <svelte:component this={Icons[child.icon]} size={16}
                                                                  className="text-black"/>
                                                <span class="truncate">{child.menuName}</span>
                                            </a>
                                        {/if}
                                    {/each}
                                </div>
                            {/if}
                        </div>
                    {/if}
                {/each}
            {:else}
                <div class="py-6 text-sm text-center text-gray-400">
                    {searchQuery ? 'No matching menus' : 'No menus available'}
                </div>
            {/if}

            <!-- 系统菜单 -->
            <div class="mt-4 pt-3 border-t border-gray-200 space-y-1">
                <div class="px-3 py-2 text-xs font-semibold uppercase tracking-wider text-gray-500">System</div>
                <a
                      href="/about"
                      use:link
                      use:active={{ path: "/about", className: 'bg-gray-100' }}
                      class="flex items-center gap-3 px-3 py-2 text-sm rounded-lg transition-all hover:bg-gray-100"
                >
                    <Settings size={16}  className="text-black"/>
                    <span class="truncate">About</span>
                </a>
            </div>
        </nav>

        <!-- 底部用户信息 -->
        {#if $user?.id}
            <div class="p-3 border-t border-gray-200 bg-gray-50">
                <div class="flex items-center gap-3 px-2 py-2 rounded-lg hover:bg-gray-100 transition cursor-pointer">
                    <div class="w-9 h-9 overflow-hidden rounded-full ring-2 ring-gray-200">
                        <img src={`${AVATAR_PREFIX}${$user.id}/${$user.avatar}`} alt="Avatar"
                             class="w-full h-full object-cover"/>
                    </div>
                    <div class="flex-1 min-w-0">
                        <div class="text-sm font-medium text-gray-800 truncate">{$user.name}</div>
                        <div class="text-xs text-gray-500 truncate">{$user.email || 'User'}</div>
                    </div>
                </div>
            </div>
        {/if}
    </aside>

    <!-- 主内容 -->
    <main class="flex flex-col flex-1 overflow-hidden">
        <!-- 顶部导航 -->
        <header class="flex items-center justify-between px-6 py-4 bg-white border-b border-gray-200 shadow-sm">
            <div class="flex items-center gap-2 text-[15px] py-3 text-gray-600">
                <span>{$location.split("/")[1] || 'home'}</span>
                {#if $location.split("/")[2]}
                    <span class="text-gray-300">/</span>
                    <span class="font-semibold text-gray-900">{$location.split("/")[2]}</span>
                {/if}
            </div>
            <div class="flex items-center gap-3">
                <button
                      class="p-2 border border-gray-300 rounded-lg text-gray-600 hover:text-gray-900 hover:border-gray-400 hover:bg-gray-50 transition"
                      title="Settings">
                    <SettingsAlt size={20}/>
                </button>
                <button
                      class="px-4 py-2 font-medium text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-800 hover:text-white hover:border-gray-800 transition"
                      on:click={handleLogout}>
                    Logout
                </button>
            </div>
        </header>

        <!-- 路由内容 -->
        <div class="flex-1 p-6 overflow-y-auto bg-gray-50">
            <div class="max-w-7xl mx-auto">
                {#if layoutRoutes}
                    <Router routes={layoutRoutes}/>
                {/if}
            </div>
        </div>
    </main>
</div>
