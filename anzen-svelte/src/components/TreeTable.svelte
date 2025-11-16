<script>
    import { createEventDispatcher, onMount } from 'svelte';
    import {ChevronDown, Code, Search} from "../lib/icons/index.js";

    // Props
    export let data = [];
    export let total = 0;
    export let page = 1;
    export let showAddButton = true;
    export let addButtonText = "New record";
    export let showApiPreview = true;
    export let columns = [];
    export let searchPlaceholder = 'Search...';
    export let showCheckbox = true;
    export let actions = ['edit', 'delete', 'more'];
    export let autoExpandOnSearch = true;
    const dispatch = createEventDispatcher();

    let search = "";
    let expandedIds = new Set();
    let selectedIds = new Set();
    let allExpanded = false;
    let flatData = [];
    let nodeMap = new Map();
    let parentMap = new Map();
    // 内部状态
    let sort = "-created";
    let showColumnMenu = false;
    let menuButtonRect = null;
    let selectedRows = new Set();
    let selectAll = false;

    // --------- sanitize & helpers ----------
    function ensureChildrenArray(node) {
        // if children is a JSON string, parse it; if undefined, set to []
        if (typeof node.children === 'string') {
            try {
                node.children = JSON.parse(node.children || '[]');
            } catch (e) {
                node.children = [];
            }
        }
        if (!Array.isArray(node.children)) {
            node.children = [];
        }
        // recursively sanitize children
        node.children.forEach(c => ensureChildrenArray(c));
    }

    function sanitizeData(nodes) {
        if (!Array.isArray(nodes)) return [];
        // work on shallow copies to avoid unexpected external mutation
        const cloned = nodes.map(n => ({ ...n }));
        cloned.forEach(n => ensureChildrenArray(n));
        return cloned;
    }

    // --------- build maps ----------
    function buildMaps(nodes, parent = null) {
        if (!Array.isArray(nodes)) return;
        nodes.forEach(node => {
            nodeMap.set(node.id, node);
            parentMap.set(node.id, parent);
            if (node.children && node.children.length) buildMaps(node.children, node.id);
        });
    }

    // --------- flat render ----------
    function renderTreeRows(nodes, level = 0) {
        const rows = [];
        if (!Array.isArray(nodes)) return rows;
        nodes.forEach(node => {
            // push shallow copy with level marker (does not mutate original node)
            rows.push({ ...node, _level: level });
            if (node.children && node.children.length > 0 && expandedIds.has(node.id)) {
                rows.push(...renderTreeRows(node.children, level + 1));
            }
        });
        return rows;
    }

    // --------- toggles ----------
    function toggleExpand(id) {
        const newSet = new Set(expandedIds);
        if (newSet.has(id)) newSet.delete(id);
        else newSet.add(id);
        expandedIds = newSet;
        allExpanded = checkAllExpanded();
    }

    function checkAllExpanded() {
        const nodesWithChildren = Array.from(nodeMap.values()).filter(n => Array.isArray(n.children) && n.children.length);
        if (nodesWithChildren.length === 0) return false;
        return nodesWithChildren.every(n => expandedIds.has(n.id));
    }

    function toggleExpandAll() {
        allExpanded = !allExpanded;
        if (allExpanded) {
            const newSet = new Set();
            function dfs(nodes) {
                if (!Array.isArray(nodes)) return;
                nodes.forEach(n => {
                    if (n.children && n.children.length) {
                        newSet.add(n.id);
                        dfs(n.children);
                    }
                });
            }
            dfs(data);
            expandedIds = newSet;
        } else {
            expandedIds = new Set();
        }
    }

    // --------- search ----------
    function handleSearch() {
        dispatch('search', { search: search });
        if (!search) return;
        const matched = [];
        const q = search.toLowerCase();
        nodeMap.forEach((node, id) => {
            const text = ((node.menuName || node.name || '') + ' ' + (node.permission || '') + ' ' + (node.url || '')).toLowerCase();
            if (text.includes(q)) matched.push(id);
        });

        if (autoExpandOnSearch) {
            const newSet = new Set(expandedIds);
            matched.forEach(id => {
                let p = parentMap.get(id);
                while (p) {
                    newSet.add(p);
                    p = parentMap.get(p);
                }
            });
            expandedIds = newSet;
        }
    }

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
    // --------- checkbox selection (propagate down/up) ----------
    function toggleSelect(id) {
        const checked = !selectedIds.has(id);
        const newSel = new Set(selectedIds);
        if (checked) newSel.add(id);
        else newSel.delete(id);

        // propagate down
        const node = nodeMap.get(id);
        if (node) {
            function setDown(n, val) {
                if (val) newSel.add(n.id);
                else newSel.delete(n.id);
                if (n.children && n.children.length) n.children.forEach(c => setDown(c, val));
            }
            setDown(node, checked);
        }

        // propagate up
        function updateUp(currId) {
            const pId = parentMap.get(currId);
            if (!pId) return;
            const parent = nodeMap.get(pId);
            if (!parent) return;
            const allChildrenSelected = parent.children.every(c => newSel.has(c.id));
            if (allChildrenSelected) newSel.add(pId);
            else newSel.delete(pId);
            updateUp(pId);
        }
        updateUp(id);

        selectedIds = newSel;
        dispatch('selectionChange', { selected: Array.from(selectedIds) });
    }

    function isIndeterminate(id) {
        const node = nodeMap.get(id);
        if (!node || !node.children || node.children.length === 0) return false;
        const some = node.children.some(c => selectedIds.has(c.id) || isIndeterminate(c.id));
        const all = node.children.every(c => selectedIds.has(c.id) && !isIndeterminate(c.id));
        return some && !all;
    }

    // --------- actions ----------
    function onAdd() { dispatch('add'); }
    function onEdit(row) { dispatch('edit', row); }
    function onDelete(row) { dispatch('delete', row); }
    function onMore(row) { dispatch('more', row); }

    // --------- cell value ----------
    function getCellValue(row, column) {
        const value = row[column.key];
        if (column.render && typeof column.render === 'function') {
            return column.render(value, row);
        }
        return value ?? '-';
    }

    // --------- reactive: when data changes, sanitize & rebuild maps & flatData ----------
    $: {
        // sanitize data to ensure children are arrays
        data = sanitizeData(data);
        // rebuild maps
        nodeMap = new Map();
        parentMap = new Map();
        buildMaps(data);
        // prune selected / expanded to existing nodes
        const ids = new Set(Array.from(nodeMap.keys()));
        selectedIds = new Set(Array.from(selectedIds).filter(id => ids.has(id)));
        expandedIds = new Set(Array.from(expandedIds).filter(id => ids.has(id)));
        allExpanded = checkAllExpanded();
        // rebuild flatData
        flatData = renderTreeRows(data);
    }

    onMount(() => {
        // nothing else
    });
</script>


<div class="bg-white rounded-xl shadow-md overflow-hidden">
    <!-- Header -->
    <div class="flex items-center justify-between px-5 py-4 border-b border-gray-200">
        <div class="flex items-center gap-2 bg-gray-100 border border-gray-200 px-3 py-1.5 rounded-md w-full max-w-md">
            <Search size={16}/>
            <input
                    class="bg-transparent outline-none text-sm flex-1"
                    placeholder={searchPlaceholder}
            />
        </div>

        <div class="flex gap-2">
            {#if showApiPreview}
                <button class="flex items-center gap-1 px-3 py-1.5 text-sm rounded-md border border-gray-300 hover:bg-gray-100 transition">
                    <Code size={16}/> API Preview
                </button>
            {/if}

            {#if showAddButton}
                <button
                        on:click={() => onAdd && onAdd()}
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
                {#if showCheckbox}
                    <th class="px-4 py-3 text-left text-sm font-medium text-gray-700 w-12">
                        <input type="checkbox" checked={selectAll} on:change={toggleSelectAll} />
                    </th>
                {/if}

                {#each columns as column}
                    {#if visibleColumns[column.key]}
                        <th class="px-4 py-3 text-left text-sm font-medium text-gray-700 {column.sortable !== false ? 'cursor-pointer' : ''}"
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

                <th class="px-4 py-3 text-center text-sm font-medium text-gray-700">
                    Actions
                </th>

                <th class="px-3 w-10">
                    <button class="flex items-center gap-1 text-gray-600 hover:text-black transition-colors"
                            on:click={openColumnMenu}>
                        <ChevronDown size={16} class="transition-transform {showColumnMenu ? 'rotate-180' : ''}"/>
                    </button>
                </th>
            </tr>
            </thead>


            <tbody>
            {#each flatData as row}
                <tr class="border-b hover:bg-gray-50 transition">
                    {#if showCheckbox}
                        <td class="px-4 py-3 text-sm">
                            <input type="checkbox"
                                   checked={selectedIds.has(row.id)}
                                   aria-checked={isIndeterminate(row.id) ? 'mixed' : selectedIds.has(row.id)}
                                   on:change={() => toggleSelect(row.id)} />
                        </td>
                    {/if}

                    {#each columns as column, idx}
                        <td class={`px-4 py-3 text-sm ${column.class || ''}`}>
                            {#if idx === 0}
                                <div class="flex items-center gap-2" style="padding-left: {row._level*20}px">
                                    {#if row.children?.length}
                                        <button on:click={() => toggleExpand(row.id)}
                                                class="w-5 h-5 flex items-center justify-center hover:bg-gray-200 rounded transition">
                                            {expandedIds.has(row.id) ? "▼" : "▶"}
                                        </button>
                                    {:else}
                                        <span class="w-5"></span>
                                    {/if}
                                    <span>{@html getCellValue(row, column)}</span>
                                </div>
                            {:else}
                                {@html getCellValue(row, column)}
                            {/if}
                        </td>
                    {/each}

                    <td class="px-4 py-3 text-center">
                        <div class="inline-flex gap-2">
                            {#each actions as act}
                                {#if act === 'edit'}
                                    <button on:click={() => onEdit(row)}
                                            class="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition">
                                        Edit
                                    </button>
                                {:else if act === 'delete'}
                                    <button on:click={() => onDelete(row)}
                                            class="px-3 py-1 text-sm bg-red-100 text-red-700 rounded hover:bg-red-200 transition">
                                        Delete
                                    </button>
                                {:else if act === 'more'}
                                    <button on:click={() => onMore(row)}
                                            class="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition">
                                        ···
                                    </button>
                                {/if}
                            {/each}
                        </div>
                    </td>
                </tr>
            {/each}
            </tbody>

        </table>
    </div>

    {#if flatData.length === 0}
        <div class="p-8 text-center text-gray-500">No data available</div>
    {/if}
</div>

<style>
    /* 可以在这里添加 indeterminate checkbox JS（如果需要） */
</style>
