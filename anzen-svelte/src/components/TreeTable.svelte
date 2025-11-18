<script>
    import {createEventDispatcher} from 'svelte';
    import {ArrowDown, ArrowUp, ChevronDown, Code, Search} from "../lib/icons/index.js";
    import {permission} from "../utils/permissionDirective.js"

    export let data = [];
    export let columns = [];
    export let searchPlaceholder = 'Search...';
    export let addButtonText = "New record";
    export let showAddButton = true;
    export let showApiPreview = false;
    export let showCheckbox = true;
    export let autoExpandOnSearch = true;
    export let actions = {};
    export let onAdd = null;
    export let onEdit = null;
    export let onDelete = null;

    const dispatch = createEventDispatcher();

    let search = "";
    let expandedIds = new Set();
    let selectedIds = new Set();
    let allExpanded = false;
    let flatData = [];
    let nodeMap = new Map();
    let parentMap = new Map();
    let currentSort = "";
    let showColumnMenu = false;
    let visibleColumns = {};

    function ensureChildrenArray(node) {
        if (typeof node.children === 'string') {
            try {
                node.children = JSON.parse(node.children || '[]');
            } catch {
                node.children = [];
            }
        }
        if (!Array.isArray(node.children)) node.children = [];
        node.children.forEach(c => ensureChildrenArray(c));
    }

    function sanitizeData(nodes) {
        if (!Array.isArray(nodes)) return [];
        const cloned = nodes.map(n => ({...n}));
        cloned.forEach(n => ensureChildrenArray(n));
        return cloned;
    }

    function buildMaps(nodes, parent = null) {
        if (!Array.isArray(nodes)) return;
        nodes.forEach(node => {
            nodeMap.set(node.id, node);
            parentMap.set(node.id, parent);
            if (node.children && node.children.length) buildMaps(node.children, node.id);
        });
    }

    function renderTreeRows(nodes, level = 0) {
        const rows = [];
        if (!Array.isArray(nodes)) return rows;
        nodes.forEach(node => {
            rows.push({...node, _level: level});
            if (node.children?.length && expandedIds.has(node.id)) rows.push(...renderTreeRows(node.children, level + 1));
        });
        return rows;
    }

    function toggleExpand(id) {
        const newSet = new Set(expandedIds);
        newSet.has(id) ? newSet.delete(id) : newSet.add(id);
        expandedIds = newSet;
        allExpanded = checkAllExpanded();
    }

    function checkAllExpanded() {
        const nodesWithChildren = Array.from(nodeMap.values()).filter(n => n.children?.length);
        if (nodesWithChildren.length === 0) return false;
        return nodesWithChildren.every(n => expandedIds.has(n.id));
    }

    function toggleExpandAll() {
        allExpanded = !allExpanded;
        if (allExpanded) {
            const newSet = new Set();

            function dfs(nodes) {
                nodes?.forEach(n => {
                    if (n.children?.length) {
                        newSet.add(n.id);
                        dfs(n.children);
                    }
                });
            }

            dfs(data);
            expandedIds = newSet;
        } else expandedIds = new Set();
    }

    function handleSearch() {
        dispatch('search', {search});
        if (!search) return;

        const matched = [];
        const q = search.toLowerCase();
        nodeMap.forEach((node, id) => {
            const searchText = [node.menuName || node.name || '', node.permission || '', node.url || ''].join(' ').toLowerCase();
            if (searchText.includes(q)) matched.push(id);
        });

        if (autoExpandOnSearch) {
            const newSet = new Set(expandedIds);
            matched.forEach(id => {
                let parentId = parentMap.get(id);
                while (parentId) {
                    newSet.add(parentId);
                    parentId = parentMap.get(parentId);
                }
            });
            expandedIds = newSet;
        }
    }

    function handleSort(field) {
        if (currentSort === field) currentSort = "-" + field;
        else if (currentSort === "-" + field) currentSort = field;
        else currentSort = field;
        dispatch('sort', {sort: currentSort});
    }

    function getSortIcon(key) {
        if (currentSort === key) return ArrowUp;
        if (currentSort === `-${key}`) return ArrowDown;
        return null;
    }

    function toggleColumn(key) {
        visibleColumns[key] = !visibleColumns[key];
        visibleColumns = {...visibleColumns};
    }

    function toggleSelect(id) {
        const checked = !selectedIds.has(id);
        const newSel = new Set(selectedIds);

        if (checked) newSel.add(id); else newSel.delete(id);

        const node = nodeMap.get(id);
        if (node) {
            function setChildrenState(n, isChecked) {
                if (isChecked) newSel.add(n.id); else newSel.delete(n.id);
                n.children?.forEach(child => setChildrenState(child, isChecked));
            }

            setChildrenState(node, checked);
        }

        function updateParentState(nodeId) {
            const pId = parentMap.get(nodeId);
            if (!pId) return;
            const parent = nodeMap.get(pId);
            if (!parent?.children) return;
            const allChildrenSelected = parent.children.every(c => newSel.has(c.id));
            if (allChildrenSelected) newSel.add(pId); else newSel.delete(pId);
            updateParentState(pId);
        }

        updateParentState(id);

        selectedIds = newSel;
        dispatch('selectionChange', {selected: Array.from(selectedIds)});
    }

    function isIndeterminate(id) {
        const node = nodeMap.get(id);
        if (!node?.children?.length) return false;
        const someSelected = node.children.some(c => selectedIds.has(c.id) || isIndeterminate(c.id));
        const allSelected = node.children.every(c => selectedIds.has(c.id) && !isIndeterminate(c.id));
        return someSelected && !allSelected;
    }

    function handleAdd() {
        if (onAdd) onAdd();
        dispatch('add');
    }

    function handleEdit(row) {
        if (onEdit) onEdit(row);
        dispatch('edit', row);
    }

    function handleDelete(row) {
        if (onDelete) onDelete(row);
        dispatch('delete', row);
    }

    function getCellValue(row, column) {
        const value = row[column.key];
        return column.render?.(value, row) ?? value ?? '-';
    }

    $: {
        visibleColumns = {};
        columns.forEach(col => {
            visibleColumns[col.key] = col.visible !== false;
        });
    }

    $: {
        data = sanitizeData(data);
        nodeMap = new Map();
        parentMap = new Map();
        buildMaps(data);
        const validIds = new Set(nodeMap.keys());
        selectedIds = new Set([...selectedIds].filter(id => validIds.has(id)));
        expandedIds = new Set([...expandedIds].filter(id => validIds.has(id)));
        allExpanded = checkAllExpanded();
        flatData = renderTreeRows(data);
    }
</script>

<div class="bg-white rounded-xl shadow-md overflow-hidden">
    <div class="flex items-center justify-between px-5 py-4 border-b border-gray-200">
        <div class="flex items-center gap-2 bg-gray-100 border border-gray-200 px-3 py-1.5 rounded-md w-full max-w-md">
            <Search size={16}/>
            <input bind:value={search} on:input={handleSearch} class="bg-transparent outline-none text-sm flex-1"
                   placeholder={searchPlaceholder}/>
        </div>

        <div class="flex gap-2">
            <button on:click={toggleExpandAll}
                    class="flex items-center gap-1 px-3 py-1.5 text-sm rounded-md border border-gray-300 hover:bg-gray-100 transition"
                    title={allExpanded ? 'Collapse All' : 'Expand All'}>
                {allExpanded ? '📂' : '📁'} {allExpanded ? 'Collapse' : 'Expand'}
            </button>
            {#if showApiPreview}
                <button
                      class="flex items-center gap-1 px-3 py-1.5 text-sm rounded-md border border-gray-300 hover:bg-gray-100 transition">
                    <Code size={16}/> API Preview
                </button>
            {/if}
            {#if showAddButton}
                <button on:click={handleAdd}
                        class="flex items-center gap-2 px-3 py-1.5 text-sm rounded-md bg-blue-600 text-white hover:bg-blue-700 transition">
                    <span class="text-lg leading-none">+</span> {addButtonText}
                </button>
            {/if}
        </div>
    </div>

    <div class="overflow-x-auto">
        <table class="w-full text-sm">
            <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
                {#if showCheckbox}
                    <th class="px-4 py-3 w-12"></th>
                {/if}
                {#each columns as column}
                    {#if visibleColumns[column.key]}
                        <th class="px-4 py-3 text-left text-sm font-medium text-gray-700 {column.sortable !== false ? 'cursor-pointer hover:bg-gray-100' : ''}"
                            on:click={() => column.sortable !== false && handleSort(column.key)}>
                            <div class="flex items-center gap-2">
                                {#if column.icon}
                                    <svelte:component this={column.icon} size={14}/>
                                {/if}
                                <span>{column.label}</span>
                                {#if column.sortable !== false}
                                    {#if getSortIcon(column.key)}
                                        <svelte:component this={getSortIcon(column.key)} size={12}
                                                          class="text-blue-600"/>
                                    {:else}
                                        <span class="text-gray-400 text-xs">⇅</span>
                                    {/if}
                                {/if}
                            </div>
                        </th>
                    {/if}
                {/each}
                <th class="px-4 py-3 text-center text-sm font-medium text-gray-700">Actions</th>
                <th class="px-3 w-10">
                    <button class="flex items-center gap-1 text-gray-600 hover:text-black transition-colors"
                            on:click={() => showColumnMenu = !showColumnMenu}>
                        <ChevronDown size={16} class="transition-transform {showColumnMenu ? 'rotate-180' : ''}"/>
                    </button>
                </th>
            </tr>
            </thead>

            <tbody>
            {#each flatData as row}
                <tr class="border-b hover:bg-gray-50 transition">
                    {#if showCheckbox}
                        <td class="px-4 py-3">
                            <input type="checkbox" checked={selectedIds.has(row.id)}
                                   indeterminate={isIndeterminate(row.id)} on:change={() => toggleSelect(row.id)}
                                   class="rounded"/>
                        </td>
                    {/if}

                    {#each columns as column, idx}
                        {#if visibleColumns[column.key]}
                            <td class="px-4 py-3 text-sm {column.class || ''}">
                                {#if idx === 0}
                                    <div class="flex items-center gap-2" style="padding-left: {row._level * 20}px">
                                        {#if row.children?.length}
                                            <button on:click={() => toggleExpand(row.id)}
                                                    class="w-5 h-5 flex items-center justify-center hover:bg-gray-200 rounded transition text-gray-600">
                                                {expandedIds.has(row.id) ? '▼' : '▶'}
                                            </button>
                                        {:else}<span class="w-5"></span>{/if}
                                        <span>{@html getCellValue(row, column)}</span>
                                    </div>
                                {:else}
                                    {@html getCellValue(row, column)}
                                {/if}
                            </td>
                        {/if}
                    {/each}

                    <td class="px-4 py-3 text-center">
                        <div class="inline-flex gap-2">
                            {#each Object.entries(actions) as [key, value]}
                                {#if key === 'edit'}
                                    <button
                                          use:permission={value}
                                          on:click={() => handleEdit(row)}
                                          class="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition">
                                        Edit
                                    </button>
                                {:else if key === 'delete'}
                                    <button on:click={() => handleDelete(row)}
                                            use:permission={value}

                                            class="px-3 py-1 text-sm bg-red-100 text-red-700 rounded hover:bg-red-200 transition">
                                        Delete
                                    </button>
                                {/if}
                            {/each}
                        </div>
                    </td>
                    <td class="px-3"></td>
                </tr>
            {/each}
            </tbody>
        </table>
    </div>

    {#if flatData.length === 0}
        <div class="p-8 text-center text-gray-500">
            <p class="text-lg">📭</p>
            <p class="mt-2">No data available</p>
        </div>
    {/if}
</div>

{#if showColumnMenu}
    <div class="fixed inset-0 z-40" on:click={() => showColumnMenu = false}></div>
    <div class="fixed z-50 bg-white rounded-lg shadow-lg border p-2 mt-1 right-4 top-32">
        <div class="text-xs font-medium text-gray-600 px-2 py-1 mb-1">Show Columns</div>
        {#each columns as column}
            <label class="flex items-center gap-2 px-2 py-1.5 hover:bg-gray-100 rounded cursor-pointer">
                <input type="checkbox" checked={visibleColumns[column.key]} on:change={() => toggleColumn(column.key)}
                       class="rounded"/>
                <span class="text-sm">{column.label}</span>
            </label>
        {/each}
    </div>
{/if}

<style>
    input[type="checkbox"]:indeterminate {
        background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 20 20'%3e%3cpath fill='none' stroke='%23fff' stroke-linecap='round' stroke-linejoin='round' stroke-width='3' d='M6 10h8'/%3e%3c/svg%3e");
        background-color: #3b82f6;
        border-color: #3b82f6;
    }
</style>
