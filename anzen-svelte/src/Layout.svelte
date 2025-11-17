<script>
    import Router, {location, push} from "svelte-spa-router";
    import {wrap} from 'svelte-spa-router/wrap';
    import About from './pages/About.svelte';
    import {ChevronDown, Logo, Settings, SettingsAlt} from './lib/icons';
    import {clearUser, user} from './stores/userStore.js';
    import {logout} from "./api/sysApis.js";
    import {showToast} from "./stores/toastStore.js";
    import Setting from './pages/Setting.svelte';
    import {clearPermissions, myRouter} from "./stores/authStore.js";

    const AVATAR_PREFIX = "http://127.0.0.1:8090/api/files/_pb_users_auth_/";

    let searchQuery = '';
    let filteredRouter = [];

    // 展开状态管理 - 默认全部展开
    let expandedMenus = {};

    // 生成动态路由
    $: layoutRoutes = generateRoutes($myRouter);
    const modules = import.meta.glob('./pages/**/*.svelte');

    function generateRoutes(menuData) {
        const routes = {};
        // 添加动态菜单路由
        if (menuData && menuData.length > 0) {
            menuData.forEach(menu => {
                if (menu.children && menu.children.length > 0) {
                    menu.children.forEach(child => {
                        if (child.url && child.menuType === 'C') {
                            const path = `./pages${child.url}.svelte`;
                            routes[child.url] = wrap({
                                asyncComponent: modules[path]
                            });
                        }
                    });
                }
            });
        }
        // 添加固定路由
        routes['/about'] = wrap({component: About});
        routes['/setting'] = wrap({component: Setting});
        return routes;
    }

    // 初始化展开状态 - 默认全部展开
    $: if ($myRouter && $myRouter.length > 0 && Object.keys(expandedMenus).length === 0) {
        const allExpanded = {};
        $myRouter.forEach(item => {
            if (item.menuType === 'M' && item.children && item.children.length > 0) {
                allExpanded[String(item.id)] = true;
            }
        });
        expandedMenus = allExpanded;
    }

    // 搜索过滤
    $: {
        if (searchQuery.trim() === '') {
            filteredRouter = $myRouter;
        } else {
            const query = searchQuery.toLowerCase();
            filteredRouter = $myRouter.map(menu => {
                if (menu.menuType === 'M' && menu.children) {
                    const filteredChildren = menu.children.filter(child =>
                        child.menuName.toLowerCase().includes(query) ||
                        child.url?.toLowerCase().includes(query)
                    );

                    if (filteredChildren.length > 0 || menu.menuName.toLowerCase().includes(query)) {
                        return {
                            ...menu,
                            children: filteredChildren
                        };
                    }
                }
                return null;
            }).filter(Boolean);

            // 搜索时自动展开所有匹配的菜单
            if (filteredRouter.length > 0) {
                const expanded = {};
                filteredRouter.forEach(item => {
                    if (item.menuType === 'M') {
                        expanded[String(item.id)] = true;
                    }
                });
                expandedMenus = expanded;
            }
        }
    }

    // 切换菜单展开状态
    function toggleMenu(menuId) {
        const key = String(menuId);
        expandedMenus = {
            ...expandedMenus,
            [key]: !expandedMenus[key]
        };
    }

    // 检查菜单是否展开
    function isMenuExpanded(menuId) {
        return expandedMenus[String(menuId)] === true;
    }

    // 处理登出
    const handleLogout = () => {
        logout();
        clearUser();
        clearPermissions();
        showToast("Logout successful!", "info");
        setTimeout(() => push('/login'), 800);
    };

    // 检查路由是否激活
    const isActive = (path) => {
        console.log("fuck")
        console.log(path)
        console.log($location)
        return $location === path;
    }

    // 获取页面标题
    function getPageTitle(location) {
        const titleMap = {
            '/sys/user/index': 'User Management',
            '/sys/role/index': 'Role Management',
            '/sys/menu/index': 'Menu Management',
            '/sys/log/index': 'Log Management',
            '/sys/monitor/index': 'Server Monitor',
            '/sys/db/index': 'Database Management',
            '/sys/redis/index': 'Redis Monitor',
            '/setting': 'Settings',
            '/about': 'About'
        };
        return titleMap[location] || 'Dashboard';
    }
</script>

<div class="flex h-screen bg-gray-50">
  <!-- Sidebar -->
  <aside class="flex flex-col bg-white border-r border-gray-200 w-42 shadow-sm">
    <!-- Logo 区域 -->
    <div class="p-5 border-b border-gray-200">
      <div class="flex items-center justify-center mb-4 text-gray-800">
        <Logo size={32}/>
        <span class="text-lg font-semibold tracking-wide ml-2">RUOYI-PB</span>
      </div>

      <div class="relative">
        <input
                type="text"
                placeholder="Search menus..."
                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent focus:bg-white transition-all"
                bind:value={searchQuery}
        />
        {#if searchQuery}
          <button
                  class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                  on:click={() => searchQuery = ''}
          >
            ✕
          </button>
        {/if}
      </div>
    </div>

    <!-- 导航菜单 -->
    <nav class="flex-1 overflow-y-auto p-3 scrollbar-thin scrollbar-thumb-gray-300 scrollbar-track-transparent">
      {#if filteredRouter && filteredRouter.length > 0}
        {#each filteredRouter as item (item.id)}
          {#if item.menuType === 'M'}
            <div class="mb-1">
              <!-- 父级菜单 -->
              <button
                      class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-gray-700 text-sm font-medium hover:bg-gray-100 transition-all duration-200 group"
                      type="button"
                      on:click={() => toggleMenu(item.id)}
              >
                                <span class="flex items-center justify-center w-5 text-gray-500 group-hover:text-gray-700">
                                    {#if item.icon}
                                        <!-- 这里可以根据 icon 名称动态渲染图标 -->
                                        <span class="text-lg">{item.icon === 'Settings' ? '⚙️' : '📁'}</span>
                                    {:else}
                                        📁
                                    {/if}
                                </span>
                <span class="flex-1 text-left">{item.menuName}</span>
                {#if item.children && item.children.length > 0}
                                    <span
                                            class="transform transition-transform duration-200 text-gray-400 {isMenuExpanded(item.id) ? 'rotate-180' : ''}"
                                    >
                                        <ChevronDown size={16}/>
                                    </span>
                {/if}
              </button>

              <!-- 子菜单 -->
              {#if item.children && item.children.length > 0 && isMenuExpanded(item.id)}
                <div class="ml-2 mt-1 space-y-0.5 border-l-2 border-gray-200 pl-2">
                  {#each item.children as child (child.id)}
                    {#if child.menuType === 'C'}
                      <a
                              href="#{child.url}"
                              class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-gray-600 hover:text-gray-900 hover:bg-gray-50 transition-all duration-200 {isActive(child.url) ? 'bg-blue-50 text-blue-600 font-medium border-l-2 border-blue-500 -ml-[2px]' : ''}"
                      >
                                                <span class="flex items-center justify-center w-5 text-xs">
                                                    {child.icon === 'User' ? '👤' :
                                                        child.icon === 'CustomerService' ? '👔' :
                                                            child.icon === 'List' ? '📋' :
                                                                child.icon === 'CalendarClock' ? '📅' :
                                                                    child.icon === 'Eye' ? '👁️' : '•'}
                                                </span>
                        <span>{child.menuName}</span>
                      </a>
                    {/if}
                  {/each}
                </div>
              {/if}
            </div>
          {/if}
        {/each}
      {:else}
        <div class="text-center text-gray-400 text-sm py-8">
          {searchQuery ? 'No matching menus' : 'No menus available'}
        </div>
      {/if}

      <!-- 固定系统菜单 -->
      <div class="mt-6 pt-4 border-t border-gray-200">
        <div class="text-xs text-gray-500 uppercase tracking-wider font-semibold px-3 py-2">
          System
        </div>
        <a
                href="#/setting"
                class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-gray-600 hover:text-gray-900 hover:bg-gray-50 transition-all duration-200 {isActive('/setting') ? 'bg-blue-50 text-blue-600 font-medium' : ''}"
        >
                    <span class="flex items-center justify-center w-5">
                        <Settings size={18}/>
                    </span>
          <span>Settings</span>
        </a>
        <a
                href="#/about"
                class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-gray-600 hover:text-gray-900 hover:bg-gray-50 transition-all duration-200 {isActive('/about') ? 'bg-blue-50 text-blue-600 font-medium' : ''}"
        >
                    <span class="flex items-center justify-center w-5">
                        <Settings size={18}/>
                    </span>
          <span>About</span>
        </a>
      </div>
    </nav>

    <!-- 底部用户信息 -->
    <div class="p-4 border-t border-gray-200 bg-gray-50">
      {#if $user?.id}
        <div class="flex items-center gap-3 px-2 py-2 rounded-lg hover:bg-gray-100 transition-colors cursor-pointer">
          <div class="w-9 h-9 rounded-full overflow-hidden ring-2 ring-gray-200">
            <img
                    src="{AVATAR_PREFIX}{$user.id}/{$user.avatar}"
                    alt="Avatar"
                    class="w-full h-full object-cover"
            />
          </div>
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium text-gray-800 truncate">
              {$user.name}
            </div>
            <div class="text-xs text-gray-500 truncate">
              {$user.email || 'User'}
            </div>
          </div>
        </div>
      {/if}
    </div>
  </aside>

  <!-- 主内容区 -->
  <main class="flex-1 flex flex-col overflow-hidden">
    <!-- 顶部导航栏 -->
    <header class="bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between shadow-sm">
      <div class="flex items-center gap-2 text-sm">
        <span class="text-gray-500">Collections</span>
        <span class="text-gray-300">/</span>
        <span class="text-gray-800 font-medium">
                    {getPageTitle($location)}
                </span>
      </div>

      <div class="flex items-center gap-3">
        <button
                class="p-2 border border-gray-300 rounded-lg text-gray-600 hover:text-gray-900 hover:border-gray-400 hover:bg-gray-50 transition-all"
                title="Settings"
        >
          <SettingsAlt size={20}/>
        </button>

        <button
                class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 font-medium hover:bg-gray-800 hover:text-white hover:border-gray-800 transition-all duration-200"
                on:click={handleLogout}
        >
          Logout
        </button>
      </div>
    </header>

    <!-- 内容区域 -->
    <div class="flex-1 overflow-y-auto p-6 bg-gray-50">
      <div class="max-w-7xl mx-auto">
        <Router routes={layoutRoutes}/>
      </div>
    </div>
  </main>
</div>

<style>
    /* 自定义滚动条样式 */
    .scrollbar-thin::-webkit-scrollbar {
        width: 6px;
    }

    .scrollbar-thin::-webkit-scrollbar-track {
        background: transparent;
    }

    .scrollbar-thin::-webkit-scrollbar-thumb {
        background: #d1d5db;
        border-radius: 3px;
    }

    .scrollbar-thin::-webkit-scrollbar-thumb:hover {
        background: #9ca3af;
    }
</style>