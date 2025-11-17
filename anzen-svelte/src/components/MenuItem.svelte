<script>
    import {push} from 'svelte-spa-router';
    import {createEventDispatcher} from 'svelte';

    export let menu;
    export let sidebarOpen = true;

    const dispatch = createEventDispatcher();

    let expanded = false;

    function toggle() {
        expanded = !expanded;
    }

    const isExternal = (url) => url?.startsWith('http') || url?.startsWith('https');
    const isActive = (url) => url && window.location.pathname === url;

    function handleClick() {
        if (isExternal(menu.url)) {
            window.open(menu.url, '_blank');
        } else {
            push(menu.url);
        }
    }
</script>

<div class="mb-1">
    {#if menu.menuType === 'M'}
        <!-- 父菜单 -->
        <button
                class="w-full flex items-center gap-3 px-3 py-2.5 rounded-md text-gray-700 text-sm hover:bg-blue-50 transition"
                on:click={toggle}
                type="button">
            {#if sidebarOpen}
                <span class="flex-1 text-left">{menu.menuName}</span>
            {/if}
            <svg
                    class="w-4 h-4 text-gray-400 transition-transform duration-200"
                    class:rotate-180={expanded}
                    fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
            </svg>
        </button>

        <!-- 子菜单列表 -->
        {#if expanded && menu.children?.length > 0}
            <div class="ml-2 mt-1 submenu-container">
                {#each menu.children as child (child.id)}
                    {#if child.menuType === 'C'}
                        <!-- 可点击的叶子菜单 -->
                        <button
                                class="w-full flex items-center gap-2 px-2 py-1.5 rounded text-gray-600 text-xs hover:bg-gray-100 transition"
                                class:bg-blue-100={isActive(child.url)}
                                class:text-blue-700={isActive(child.url)}
                                class:font-medium={isActive(child.url)}
                                on:click={() => {
                                if (isExternal(child.url)) {
                                    window.open(child.url, '_blank');
                                } else {
                                    push(child.url);
                                }
                            }}>
                            <span class="w-2 h-2 bg-gray-400 rounded-full"></span>
                            <span>{child.menuName}</span>
                        </button>
                    {:else if child.menuType === 'M'}
                        <!-- 嵌套父菜单（可扩展二级以上） -->
                        <div>
                            <button
                                    class="w-full flex items-center gap-2 px-2 py-1.5 rounded text-gray-600 text-xs hover:bg-gray-100 transition"
                                    on:click={() => {
                                    // 如果你也想让父菜单有跳转行为，可以加判断
                                    if (child.url) push(child.url);
                                }}>
                                <span class="w-2 h-2 bg-gray-400 rounded-full opacity-0"></span>
                                <span>{child.menuName}</span>
                            </button>
                            <!-- 可继续递归渲染更深层菜单 -->
                        </div>
                    {/if}
                {/each}
            </div>
        {/if}

    {:else if menu.menuType === 'C'}
        <!-- 叶子菜单项 -->
        <button
                class="w-full flex items-center gap-3 px-3 py-2.5 rounded-md text-gray-600 text-sm hover:bg-gray-100 transition"
                class:bg-blue-100={isActive(menu.url)}
                class:text-blue-700={isActive(menu.url)}
                class:font-medium={isActive(menu.url)}
                on:click={handleClick}>
            <span class="w-4 h-4 flex items-center justify-center text-gray-400">•</span>
            <span>{menu.menuName}</span>
        </button>
    {/if}
</div>

<style>
    .submenu-container {
        margin-top: 0.25rem;
        border-left: 2px solid #e5e7eb;
        padding-left: 0.75rem;
        animation: slideDown 0.2s ease-out;
    }

    @keyframes slideDown {
        from {
            opacity: 0;
            transform: translateY(-10px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
</style>