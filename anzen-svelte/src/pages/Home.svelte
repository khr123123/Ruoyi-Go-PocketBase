<script>
    import {ArrowDown, Bars4, Check, ChevronDown, Circle, Code, Download, Image, Menu, Search} from '../lib/icons';
    import {listUser} from '../api/userApis.js';

    let users = [];
    let total = 0;
    const AVATAR_PREFIX = "http://127.0.0.1:8090/api/files/_pb_users_auth_/";
    // 控制参数
    let page = 1;
    let perPage = 10;
    let search = "";
    let sort = "-created";
    // Filter state
    let filterField = null;
    let filterType = "contains";

    // 列显示控制
    let showColumnMenu = false;
    let visibleColumns = {
        id: true,
        email: true,
        emailVisibility: true,
        verified: true,
        name: true,
        avatar: true,
        created: true,
        updated: true
    };

    // 下拉菜单位置
    let menuButtonRect = null;

    // 构建过滤器
    function buildFilter() {
        if (!search) return "";
        if (!filterField) {
            return `email ~ "${search}" || name ~ "${search}"`;
        }
        if (filterType === "contains") {
            return `${filterField} ~ "${search}"`;
        }
        if (filterType === "starts") {
            return `${filterField} ~ "^${search}"`;
        }
        if (filterType === "ends") {
            return `${filterField} ~ "${search}$"`;
        }
        return "";
    }

    // 重新加载
    async function loadUsers() {
        const res = await listUser(page, perPage, sort, buildFilter());
        users = res.items;
        total = res.totalItems;
    }

    // 排序切换
    function changeSort(field) {
        if (sort === field) sort = "-" + field;
        else if (sort === "-" + field) sort = field;
        else sort = field;
        loadUsers();
    }

    // 切换列显示
    function toggleColumn(column) {
        visibleColumns[column] = !visibleColumns[column];
    }

    // 打开下拉菜单
    function openColumnMenu(event) {
        showColumnMenu = true;
        menuButtonRect = event.currentTarget.getBoundingClientRect();
    }

    // 阻止事件冒泡
    function stopPropagation(event) {
        event.stopPropagation();
    }

    // watch
    $: if (search !== undefined) loadUsers();
    loadUsers();
</script>

<!-- 容器 -->
<div class="bg-white rounded-xl shadow-md overflow-hidden">

    <!-- Header -->
    <div class="flex items-center justify-between px-5 py-4 border-b border-gray-200">
        <div class="flex items-center gap-2 bg-gray-100 border border-gray-200 px-3 py-1.5 rounded-md w-full max-w-md">
            <Search size={16}/>
            <input
                    class="bg-transparent outline-none text-sm flex-1"
                    placeholder="Search email or name..."
                    bind:value={search}
            />
        </div>

        <div class="flex gap-2">
            <button class="flex items-center gap-1 px-3 py-1.5 text-sm rounded-md border border-gray-300 hover:bg-gray-100 transition">
                <Code size={16}/> API Preview
            </button>
            <button class="flex items-center gap-2 px-3 py-1.5 text-sm rounded-md bg-black text-white hover:bg-gray-800 transition">
                <span class="text-lg leading-none">+</span> New record
            </button>
        </div>
    </div>

    <!-- Table 区域（带滚动条） -->
    <div class="overflow-x-auto max-h-[500px] overflow-y-auto">
        <table class="w-full text-sm">
            <thead>
            <tr class="bg-gray-100 text-gray-600 border-b border-gray-200">
                <th class="p-3 w-10"><input type="checkbox"/></th>

                {#if visibleColumns.id}
                    <th class="p-3 cursor-pointer" on:click={() => changeSort('id')}>
                        <div class="flex items-center gap-1">
                            <ArrowDown size={14}/>
                            id
                        </div>
                    </th>
                {/if}

                {#if visibleColumns.email}
                    <th class="p-3 cursor-pointer" on:click={() => changeSort('email')}>
                        <div class="flex items-center gap-1">
                            <Menu size={14}/>
                            email
                        </div>
                    </th>
                {/if}

                {#if visibleColumns.emailVisibility}
                    <th class="p-3 cursor-pointer" on:click={() => changeSort('emailVisibility')}>
                        <div class="flex items-center gap-1">
                            <Circle size={14}/>
                            emailVisibility
                        </div>
                    </th>
                {/if}

                {#if visibleColumns.verified}
                    <th class="p-3 cursor-pointer" on:click={() => changeSort('verified')}>
                        <div class="flex items-center gap-1">
                            <Check size={14}/>
                            verified
                        </div>
                    </th>
                {/if}

                {#if visibleColumns.name}
                    <th class="p-3 cursor-pointer" on:click={() => changeSort('name')}>
                        <div class="flex items-center gap-1">
                            <span class="font-bold">T</span> name
                        </div>
                    </th>
                {/if}

                {#if visibleColumns.avatar}
                    <th class="p-3">
                        <div class="flex items-center gap-1">
                            <Image size={14}/>
                            avatar
                        </div>
                    </th>
                {/if}

                {#if visibleColumns.created}
                    <th class="p-3 cursor-pointer" on:click={() => changeSort('created')}>
                        <div class="flex items-center gap-1">
                            <ArrowDown size={14}/>
                            created
                        </div>
                    </th>
                {/if}

                {#if visibleColumns.updated}
                    <th class="p-3 cursor-pointer" on:click={() => changeSort('updated')}>
                        <div class="flex items-center gap-1">
                            <ArrowDown size={14}/>
                            updated
                        </div>
                    </th>
                {/if}

                <!-- 列控制下拉菜单按钮 -->
                <th class="p-3 w-10">
                    <button
                            class="flex items-center gap-1 text-gray-600 hover:text-black transition-colors"
                            on:click={openColumnMenu}
                    >
                        <ChevronDown size={16} class="transition-transform {showColumnMenu ? 'rotate-180' : ''}"/>
                    </button>
                </th>
            </tr>
            </thead>

            <tbody>
            {#each users as user}
                <tr class="border-b hover:bg-gray-50">
                    <td class="p-3"><input type="checkbox"/></td>
                    {#if visibleColumns.id}
                        <td class="p-3 text-xs font-mono flex items-center gap-1 text-gray-600">
                            {user.id}
                        </td>
                    {/if}

                    {#if visibleColumns.email}
                        <td class="p-3">{user.email}</td>
                    {/if}

                    {#if visibleColumns.emailVisibility}
                        <td class="p-3">
                        <span class={`px-2 py-0.5 rounded-full text-xs ${
                            user.emailVisibility ? "bg-green-100 text-green-700" : "bg-gray-100 text-gray-600"
                        }`}>
                            {user.emailVisibility ? "true" : "false"}
                        </span>
                        </td>
                    {/if}

                    {#if visibleColumns.verified}
                        <td class="p-3">
                        <span class={`px-2 py-0.5 rounded-full text-xs ${
                            user.verified ? "bg-green-100 text-green-700" : "bg-gray-100 text-gray-600"
                        }`}>
                            {user.verified ? "true" : "false"}
                        </span>
                        </td>
                    {/if}

                    {#if visibleColumns.name}
                        <td class="p-3">{user.name}</td>
                    {/if}

                    {#if visibleColumns.avatar}
                        <td class="p-3">
                            <img
                                    class="w-8 h-8 rounded-full object-cover"
                                    src={AVATAR_PREFIX + user.id + "/" + user.avatar}
                            />
                        </td>
                    {/if}

                    {#if visibleColumns.created}
                        <td class="p-3 text-xs text-gray-600">{user.created}</td>
                    {/if}

                    {#if visibleColumns.updated}
                        <td class="p-3 text-xs text-gray-600">{user.updated}</td>
                    {/if}

                    <td class="p-3">
                        <button class="text-gray-600 hover:text-black">
                            <Bars4/>
                        </button>
                    </td>
                </tr>
            {/each}
            </tbody>
        </table>
    </div>

    <!-- Footer -->
    <div class="px-5 py-3 border-t bg-gray-50 flex justify-between text-sm">
        <div> Total found: <strong>{total}</strong></div>

        <div class="flex gap-2">
            <button class="w-8 h-8 flex items-center justify-center rounded bg-white border hover:bg-gray-100">
                <Download size={16}/>
            </button>
            <button class="w-8 h-8 flex items-center justify-center rounded bg-white border hover:bg-gray-100">⚙️
            </button>
            <button class="w-8 h-8 flex items-center justify-center rounded bg-white border hover:bg-gray-100">❓
            </button>
            <button class="w-8 h-8 flex items-center justify-center rounded bg-white border hover:bg-gray-100">ℹ️
            </button>
        </div>
    </div>

    <!-- Pagination -->
    <div class="flex items-center justify-center gap-3 py-4">
        <button
                class="px-3 py-1.5 border rounded bg-gray-100 hover:bg-gray-200"
                on:click={() => { page = Math.max(1, page - 1); loadUsers(); }}>
            Prev
        </button>

        <span class="text-sm">Page {page}</span>

        <button
                class="px-3 py-1.5 border rounded bg-gray-100 hover:bg-gray-200"
                on:click={() => { page += 1; loadUsers(); }}>
            Next
        </button>
    </div>

</div>

<!-- 在文档最外层渲染下拉菜单 -->
{#if showColumnMenu}
    <!-- 背景遮罩层 -->
    <div class="fixed inset-0 z-40" on:click={() => showColumnMenu = false}></div>

    <!-- 下拉菜单 -->
    <div
            class="fixed z-50 bg-white border border-gray-200 rounded-lg shadow-xl min-w-[180px] py-2"
            style="left: {menuButtonRect?.left - 150}px; top: {menuButtonRect?.bottom + 5}px;"
            on:click={stopPropagation}
    >
        <div class="px-3 py-2 text-xs font-semibold text-gray-500 border-b border-gray-100">
            Toggle Columns
        </div>
        <div class="max-h-60 overflow-y-auto">
            <label class="flex items-center gap-2 px-3 py-2 hover:bg-gray-50 cursor-pointer">
                <input type="checkbox" bind:checked={visibleColumns.id} class="rounded border-gray-300">
                <span class="text-sm">id</span>
            </label>
            <label class="flex items-center gap-2 px-3 py-2 hover:bg-gray-50 cursor-pointer">
                <input type="checkbox" bind:checked={visibleColumns.email} class="rounded border-gray-300">
                <span class="text-sm">email</span>
            </label>
            <label class="flex items-center gap-2 px-3 py-2 hover:bg-gray-50 cursor-pointer">
                <input type="checkbox" bind:checked={visibleColumns.emailVisibility} class="rounded border-gray-300">
                <span class="text-sm">emailVisibility</span>
            </label>
            <label class="flex items-center gap-2 px-3 py-2 hover:bg-gray-50 cursor-pointer">
                <input type="checkbox" bind:checked={visibleColumns.verified} class="rounded border-gray-300">
                <span class="text-sm">verified</span>
            </label>
            <label class="flex items-center gap-2 px-3 py-2 hover:bg-gray-50 cursor-pointer">
                <input type="checkbox" bind:checked={visibleColumns.name} class="rounded border-gray-300">
                <span class="text-sm">name</span>
            </label>
            <label class="flex items-center gap-2 px-3 py-2 hover:bg-gray-50 cursor-pointer">
                <input type="checkbox" bind:checked={visibleColumns.avatar} class="rounded border-gray-300">
                <span class="text-sm">avatar</span>
            </label>
            <label class="flex items-center gap-2 px-3 py-2 hover:bg-gray-50 cursor-pointer">
                <input type="checkbox" bind:checked={visibleColumns.created} class="rounded border-gray-300">
                <span class="text-sm">created</span>
            </label>
            <label class="flex items-center gap-2 px-3 py-2 hover:bg-gray-50 cursor-pointer">
                <input type="checkbox" bind:checked={visibleColumns.updated} class="rounded border-gray-300">
                <span class="text-sm">updated</span>
            </label>
        </div>
    </div>
{/if}