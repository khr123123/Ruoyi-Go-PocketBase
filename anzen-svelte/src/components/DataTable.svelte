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

<div class="bg-white border border-gray-300 rounded overflow-hidden">
    <!-- Header -->
    <div class="flex items-center justify-between px-5 py-3 border-b border-gray-200 bg-gray-50">
        <div class="flex items-center gap-2 bg-white border border-gray-300 px-3 py-1.5 rounded w-full max-w-md">
            <Search size={14} class="text-gray-400"/>
            <input
                  class="bg-transparent outline-none text-sm flex-1 text-gray-900 placeholder-gray-400"
                  placeholder={searchPlaceholder}
                  bind:value={search}
            />
        </div>

        <div class="flex gap-2">
            {#if showAddButton}
                <button
                      use:permission={actions["add"]}
                      on:click={() =>dispatch('add')}
                      class="flex items-center gap-1.5 px-3 py-1.5 text-sm rounded bg-gray-800 text-white hover:bg-gray-700 transition-colors">
                    <span class="text-base leading-none font-light">+</span> {addButtonText}
                </button>
            {/if}
        </div>
    </div>

    <!-- Table Container with fixed columns -->
    <div class="table-container">
        <table class="data-table">
            <thead>
            <tr class="bg-gray-50 text-gray-700 text-xs font-medium border-b border-gray-200">
                <!-- Fixed: Checkbox Column -->
                <th class="sticky-checkbox">
                    <input
                          type="checkbox"
                          bind:checked={selectAll}
                          on:change={toggleSelectAll}
                          class="w-3.5 h-3.5 rounded border-gray-300 text-gray-800 focus:ring-gray-500"
                    />
                </th>

                {#each columns as column, idx}
                    {#if visibleColumns[column.key]}
                        <th
                              class="px-4 py-3 text-left {column.sortable !== false ? 'cursor-pointer hover:bg-gray-100' : ''} {idx === 0 ? 'sticky-first-col' : ''}"
                              on:click={() => column.sortable !== false && changeSort(column.key)}>
                            <div class="flex items-center gap-1.5">
                                {#if column.icon}
                                    <svelte:component this={column.icon} size={13} class="text-gray-500"/>
                                {/if}
                                <span>{column.label}</span>
                                {#if column.sortable !== false}
                                    <svg class="w-3 h-3 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
                                        <path d="M5 8l5-5 5 5H5z"/>
                                    </svg>
                                {/if}
                            </div>
                        </th>
                    {/if}
                {/each}

                <!-- Fixed: Actions Column -->
                <th class="sticky-actions">
                    <div class="flex items-center justify-end gap-2">
                        <span>Actions</span>
                        <button
                              class="flex items-center text-gray-500 hover:text-gray-700 transition-colors"
                              on:click={openColumnMenu}>
                            <ChevronDown size={14} class="transition-transform {showColumnMenu ? 'rotate-180' : ''}"/>
                        </button>
                    </div>
                </th>
            </tr>
            </thead>

            <tbody class="text-sm">
            {#each data as row, rowIdx}
                <tr class="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                    <!-- Fixed: Checkbox -->
                    <td class="sticky-checkbox">
                        <input
                              type="checkbox"
                              checked={selectedRows.has(row.id)}
                              on:change={() => toggleRow(row.id)}
                              class="w-3.5 h-3.5 rounded border-gray-300 text-gray-800 focus:ring-gray-500"
                        />
                    </td>

                    {#each columns as column, idx}
                        {#if visibleColumns[column.key]}
                            <td class="px-4 py-3 {column.class || ''} {idx === 0 ? 'sticky-first-col' : ''} text-gray-900">
                                {@html renderCell(row, column)}
                            </td>
                        {/if}
                    {/each}

                    <!-- Fixed: Actions -->
                    <td class="sticky-actions">
                        <div class="inline-flex gap-1.5 justify-end w-full">
                            {#if actions["edit"]}
                                <button
                                      use:permission={actions["edit"]}
                                      on:click={() =>dispatch('edit', row) }
                                      class="px-2.5 py-1 text-xs bg-blue-50 text-blue-700 rounded hover:bg-blue-100 transition-colors border border-blue-200">
                                    Edit
                                </button>
                            {/if}
                            {#if actions["delete"]}
                                <button
                                      on:click={() => dispatch('delete', row)}
                                      use:permission={actions["delete"]}
                                      class="px-2.5 py-1 text-xs bg-red-50 text-red-700 rounded hover:bg-red-100 transition-colors border border-red-200">
                                    Delete
                                </button>
                            {/if}
                            {#if actions["assignPerm"]}
                                <button
                                      on:click={() => dispatch('assignPerm', row)}
                                      use:permission={actions["assignPerm"]}
                                      class="px-2.5 py-1 text-xs bg-amber-50 text-amber-700 rounded hover:bg-amber-100 transition-colors border border-amber-200">
                                    Assign
                                </button>
                            {/if}
                        </div>
                    </td>
                </tr>
            {/each}
            </tbody>
        </table>
    </div>

    <!-- Footer -->
    <div class="px-5 py-3 border-t border-gray-200 bg-gray-50 flex items-center justify-between text-xs">
        <div class="text-gray-600">
            Total found: <strong class="text-gray-900 font-semibold">{total}</strong> record(s)
        </div>

        <!-- Pagination -->
        <div class="flex items-center gap-2">
            <button
                  class="px-3 py-1.5 border border-gray-300 rounded bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed text-gray-700 text-xs transition-colors"
                  disabled={page === 1}
                  on:click={prevPage}>
                ← Prev
            </button>
            <span class="text-xs text-gray-600 px-2">Page <strong class="text-gray-900">{page}</strong></span>
            <button
                  class="px-3 py-1.5 border border-gray-300 rounded bg-white hover:bg-gray-50 text-gray-700 text-xs transition-colors"
                  on:click={nextPage}>
                Next →
            </button>
        </div>

        <button
              class="p-1.5 flex items-center justify-center rounded bg-white border border-gray-300 hover:bg-gray-50 transition-colors">
            <Download size={14} class="text-gray-600"/>
        </button>
    </div>
</div>

<!-- Column Control Dropdown -->
{#if showColumnMenu}
    <div class="fixed inset-0 z-40" on:click={() => showColumnMenu = false}></div>
    <div
          class="fixed z-50 bg-white border border-gray-200 rounded shadow-xl min-w-[200px] py-1"
          style="left: {menuButtonRect?.left - 150}px; top: {menuButtonRect?.bottom + 5}px;"
          on:click={stopPropagation}>
        <div class="px-3 py-2 text-xs font-semibold text-gray-500 border-b border-gray-100">
            Toggle columns
        </div>
        <div class="max-h-60 overflow-y-auto">
            {#each columns as column}
                <label class="flex items-center gap-2 px-3 py-2 hover:bg-gray-50 cursor-pointer text-sm">
                    <input
                          type="checkbox"
                          bind:checked={visibleColumns[column.key]}
                          class="w-3.5 h-3.5 rounded border-gray-300 text-gray-800 focus:ring-gray-500"
                    >
                    <span class="text-gray-700">{column.label}</span>
                </label>
            {/each}
        </div>
    </div>
{/if}

<style>
    /* Table Container with Horizontal Scroll */
    .table-container {
        overflow-x: auto;
        overflow-y: auto;
        max-height: 600px;
        position: relative;
    }

    .data-table {
        width: 100%;
        border-collapse: separate;
        border-spacing: 0;
    }

    /* Fixed Checkbox Column (Left) */
    .sticky-checkbox {
        position: sticky;
        left: 0;
        z-index: 20;
        background-color: white;
        padding: 12px 16px;
        text-align: center;
        width: 48px;
        border-right: 1px solid #e5e7eb;
    }

    thead .sticky-checkbox {
        background-color: #f9fafb;
        z-index: 30;
    }

    /* Fixed First Data Column */
    /*.sticky-first-col {*/
    /*    position: sticky;*/
    /*    left: 48px;*/
    /*    z-index: 20;*/
    /*    background-color: white;*/
    /*    border-right: 1px solid #e5e7eb;*/
    /*}*/

    thead .sticky-first-col {
        background-color: #f9fafb;
        z-index: 30;
    }

    /* Fixed Actions Column (Right) */
    .sticky-actions {
        position: sticky;
        right: 0;
        z-index: 20;
        background-color: white;
        padding: 12px 16px;
        text-align: right;
        min-width: 100px;
        border-left: 1px solid #e5e7eb;
    }

    thead .sticky-actions {
        background-color: #f9fafb;
        z-index: 30;
    }

    /* Row hover effect */
    tbody tr:hover .sticky-checkbox,
    tbody tr:hover .sticky-first-col,
    tbody tr:hover .sticky-actions {
        background-color: #f9fafb;
    }

    /* Scrollbar styling */
    .table-container::-webkit-scrollbar {
        height: 8px;
        width: 8px;
    }

    .table-container::-webkit-scrollbar-track {
        background: #f1f1f1;
    }

    .table-container::-webkit-scrollbar-thumb {
        background: #d1d5db;
        border-radius: 4px;
    }

    .table-container::-webkit-scrollbar-thumb:hover {
        background: #9ca3af;
    }
</style>