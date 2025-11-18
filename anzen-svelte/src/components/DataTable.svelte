<!-- src/components/DataTable.svelte -->
<script>
    import {ChevronDown, Download, Search} from '../lib/icons';
    import {createEventDispatcher} from "svelte";
    import {permission} from "../utils/permissionDirective.js";

    const dispatch = createEventDispatcher();
    // Props
    export let data = [];
    export let total = 0;
    export let columns = [];
    export let page = 1;
    export let searchPlaceholder = "Search...";
    export let showAddButton = true;
    export let addButtonText = "New record";
    export let actions = {}
    // 内部状态
    let search = "";
    let sort = "-created";
    let showColumnMenu = false;
    let menuButtonRect = null;
    let selectedRows = new Set();
    let selectAll = false;

    // 列显示控制
    let visibleColumns = {};
    $: {
        visibleColumns = {};
        columns.forEach(col => {
            visibleColumns[col.key] = col.visible !== false;
        });
    }

    // 搜索变化时触发
    $: if (search !== undefined) {
        dispatch('search', {search, page: 1});
    }

    // 排序切换
    function changeSort(field) {
        if (sort === field) sort = "-" + field;
        else if (sort === "-" + field) sort = field;
        else sort = field;
        dispatch('sort', {sort, page});
    }

    // 分页
    function prevPage() {
        if (page > 1) {
            page -= 1;
            dispatch('pageChange', {page});
        }
    }

    function nextPage() {
        page += 1;
        dispatch('pageChange', {page});
    }

    // 列显示控制
    function toggleColumn(column) {
        visibleColumns[column] = !visibleColumns[column];
    }

    function openColumnMenu(event) {
        showColumnMenu = true;
        menuButtonRect = event.currentTarget.getBoundingClientRect();
    }

    function stopPropagation(event) {
        event.stopPropagation();
    }

    // 行选择
    function toggleSelectAll() {
        if (selectAll) {
            selectedRows = new Set();
        } else {
            selectedRows = new Set(data.map(item => item.id));
        }
        selectAll = !selectAll;
    }

    function toggleRow(id) {
        if (selectedRows.has(id)) {
            selectedRows.delete(id);
        } else {
            selectedRows.add(id);
        }
        selectedRows = selectedRows;
        selectAll = selectedRows.size === data.length;
    }

    // 渲染单元格内容
    function renderCell(row, column) {
        if (column.render) {
            return column.render(row[column.key], row);
        }
        return row[column.key];
    }
</script>

<div class="bg-white rounded-xl shadow-md overflow-hidden">
    <!-- Header -->
    <div class="flex items-center justify-between px-5 py-4 border-b border-gray-200">
        <div class="flex items-center gap-2 bg-gray-100 border border-gray-200 px-3 py-1.5 rounded-md w-full max-w-md">
            <Search size={16}/>
            <input
                  class="bg-transparent outline-none text-sm flex-1"
                  placeholder={searchPlaceholder}
                  bind:value={search}
            />
        </div>

        <div class="flex gap-2">
            {#if showAddButton}
                <button
                      use:permission={actions["add"]}
                      on:click={() =>dispatch('add')}
                      class="flex items-center gap-2 px-3 py-1.5 text-sm rounded-md bg-black text-white hover:bg-gray-800 transition">
                    <span class="text-lg leading-none">+</span> {addButtonText}
                </button>
            {/if}
        </div>
    </div>

    <!-- Table -->
    <div class="overflow-x-auto max-h-[500px] overflow-y-auto">
        <table class="w-full text-sm">
            <thead>
            <tr class="bg-gray-100 text-gray-600 border-b border-gray-200">
                <th class="p-3 w-10">
                    <input type="checkbox" bind:checked={selectAll} on:change={toggleSelectAll}/>
                </th>

                {#each columns as column}
                    {#if visibleColumns[column.key]}
                        <th class="p-3 {column.sortable !== false ? 'cursor-pointer' : ''}"
                            on:click={() => column.sortable !== false && changeSort(column.key)}>
                            <div class="flex items-center gap-1">
                                {#if column.icon}
                                    <svelte:component this={column.icon} size={14}/>
                                {/if}
                                {column.label}
                            </div>
                        </th>
                    {/if}
                {/each}
                <th class="px-4 py-3 text-center text-sm font-medium text-gray-700">Actions</th>
                <th class="p-3 w-10">
                    <button
                          class="flex items-center gap-1 text-gray-600 hover:text-black transition-colors"
                          on:click={openColumnMenu}>
                        <ChevronDown size={16} class="transition-transform {showColumnMenu ? 'rotate-180' : ''}"/>
                    </button>
                </th>
            </tr>
            </thead>

            <tbody>
            {#each data as row}
                <tr class="border-b hover:bg-gray-50">
                    <td class="p-3">
                        <input type="checkbox" checked={selectedRows.has(row.id)}
                               on:change={() => toggleRow(row.id)}/>
                    </td>

                    {#each columns as column}
                        {#if visibleColumns[column.key]}
                            <td class="p-3 {column.class || ''}">
                                {@html renderCell(row, column)}
                            </td>
                        {/if}
                    {/each}

                    <td class="px-4 py-3 text-center">
                        <div class="inline-flex gap-2">
                            {#if actions["edit"]}
                                <button
                                      use:permission={actions["edit"]}
                                      on:click={() =>dispatch('edit', row) }
                                      class="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition">
                                    Edit
                                </button>
                            {/if}
                            {#if actions["delete"]}
                                <button on:click={() => dispatch('delete', row)}
                                        use:permission={actions["delete"]}
                                        class="px-3 py-1 text-sm bg-red-100 text-red-700 rounded hover:bg-red-200 transition">
                                    Delete
                                </button>
                            {/if}
                            {#if actions["assignPerm"]}
                                <button
                                      on:click={() => dispatch('assignPerm', row)}
                                      use:permission={actions["assignPerm"]}
                                      class="px-3 py-1 text-sm font-medium rounded shadow-sm bg-amber-100 text-amber-700 hover:bg-amber-200 hover:text-amber-800"
                                >
                                    Assign
                                </button>
                            {/if}

                        </div>
                    </td>
                    <td class="px-3"></td>
                </tr>
            {/each}
            </tbody>
        </table>
    </div>

    <!-- Footer -->
    <div class="px-5 py-3 border-t bg-gray-50 flex justify-between text-sm">
        <div>Total found: <strong>{total}</strong></div>
        <div class="flex gap-2">
            <button class="w-8 h-8 flex items-center justify-center rounded bg-white border hover:bg-gray-100">
                <Download size={16}/>
            </button>
        </div>
    </div>

    <!-- Pagination -->
    <div class="flex items-center justify-center gap-3 py-4">
        <button
              class="px-3 py-1.5 border rounded bg-gray-100 hover:bg-gray-200 disabled:opacity-50"
              disabled={page === 1}
              on:click={prevPage}>
            Prev
        </button>
        <span class="text-sm">Page {page}</span>
        <button
              class="px-3 py-1.5 border rounded bg-gray-100 hover:bg-gray-200"
              on:click={nextPage}>
            Next
        </button>
    </div>
</div>

<!-- 列控制下拉菜单 -->
{#if showColumnMenu}
    <div class="fixed inset-0 z-40" on:click={() => showColumnMenu = false}></div>
    <div
          class="fixed z-50 bg-white border border-gray-200 rounded-lg shadow-xl min-w-[180px] py-2"
          style="left: {menuButtonRect?.left - 150}px; top: {menuButtonRect?.bottom + 5}px;"
          on:click={stopPropagation}>
        <div class="px-3 py-2 text-xs font-semibold text-gray-500 border-b border-gray-100">
            Toggle Columns
        </div>
        <div class="max-h-60 overflow-y-auto">
            {#each columns as column}
                <label class="flex items-center gap-2 px-3 py-2 hover:bg-gray-50 cursor-pointer">
                    <input type="checkbox" bind:checked={visibleColumns[column.key]} class="rounded border-gray-300">
                    <span class="text-sm">{column.label}</span>
                </label>
            {/each}
        </div>
    </div>
{/if}