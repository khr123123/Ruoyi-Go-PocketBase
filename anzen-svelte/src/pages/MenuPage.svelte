<!-- src/pages/MenuPage.svelte -->
<script>
    import TreeTable from '../components/TreeTable.svelte';
    import {onMount} from 'svelte';
    import {hasPermission} from '../stores/authStore';
    import {showToast} from '../stores/toastStore';
    import {listMenu} from "../api/menuApis.js";
    import {permission} from '../utils/permissionDirective.js';

    let menus = [];
    let rawMenus = []; // 保存原始扁平数据
    let showDialog = false;
    let currentMenu = null;
    let parentMenuOptions = []; // 父菜单选项

    // Column definitions
    const columns = [
        {key: 'menuName', label: 'MenuName'},
        {
            key: 'icon',
            label: 'Icon',
            class: 'text-xs font-mono'
        },
        {
            key: 'menuType',
            label: 'Type',
            render: (v) => {
                const types = {M: 'Directory', C: 'Menu', F: 'Button'};
                const colors = {M: 'blue', C: 'green', F: 'yellow'};
                return `<span class="px-2 py-1 text-xs rounded bg-${colors[v] || 'gray'}-100 text-${colors[v] || 'gray'}-700">${types[v] || v || '-'}</span>`;
            }
        },
        {
            key: 'orderNum',
            label: 'Sort',
            render: (v) => `<span class="px-2 py-1 text-xs rounded bg-gray-100 text-gray-700">${v}</span>`
        },
        {key: 'permission', label: 'Permission', class: 'text-xs font-mono'},
        {key: 'url', label: 'Route', class: 'text-xs'}
    ];

    // 构建树形结构
    function buildTree(flatList) {
        if (!flatList || flatList.length === 0) return [];

        const map = {};
        const roots = [];

        // 深拷贝并创建映射
        flatList.forEach(item => {
            map[item.id] = {...item, children: []};
        });

        // 构建树
        flatList.forEach(item => {
            const node = map[item.id];
            if (!item.parentId || item.parentId === '0') {
                roots.push(node);
            } else if (map[item.parentId]) {
                map[item.parentId].children.push(node);
            } else {
                // 父节点不存在，当作根节点
                roots.push(node);
            }
        });

        // 递归排序
        function sortNodes(nodes) {
            nodes.sort((a, b) => (a.orderNum || 0) - (b.orderNum || 0));
            nodes.forEach(node => {
                if (node.children && node.children.length > 0) {
                    sortNodes(node.children);
                }
            });
        }

        sortNodes(roots);
        return roots;
    }

    // Load menu list
    async function loadMenus(params = {}) {
        try {
            const {search = ''} = params;
            const filter = search ? `menuName ~ "${search}" || permission ~ "${search}" || url ~ "${search}"` : '';
            // 获取所有菜单（不分页）
            const res = await listMenu(1, 500, 'orderNum', filter);
            rawMenus = res.items || [];
            menus = buildTree(rawMenus);
        } catch (error) {
            console.error('Load menu error:', error);
            showToast('Failed to load menus', 'error');
        }
    }

    function handleSearch(event) {
        loadMenus({search: event.detail.search});
    }

    // Add menu
    function handleAdd() {
        if (!hasPermission('sys:menu:add')) {
            showToast('No permission', 'error');
            return;
        }
        currentMenu = {
            menuName: '',
            menuType: 'M',
            parentId: '0',
            orderNum: 0,
            url: '',
            icon: '',
            permission: '',
            status: 'show'
        };
        showDialog = true;
    }

    // Edit menu
    function handleEdit(menu) {
        if (!hasPermission('sys:menu:edit')) {
            showToast('No permission', 'error');
            return;
        }
        currentMenu = {...menu};
        delete currentMenu.children; // 移除 children 字段
        showDialog = true;
    }

    // Save menu
    async function saveMenu() {
        // 验证必填项
        if (!currentMenu.menuName || !currentMenu.menuType) {
            showToast('Please fill in required fields', 'error');
            return;
        }

        try {
            const url = currentMenu.id
                ? `http://127.0.0.1:8090/api/collections/sys_menu/records/${currentMenu.id}`
                : 'http://127.0.0.1:8090/api/collections/sys_menu/records';
            const method = currentMenu.id ? 'PATCH' : 'POST';

            // 准备提交数据
            const submitData = {
                menuName: currentMenu.menuName,
                menuType: currentMenu.menuType,
                parentId: currentMenu.parentId || '0',
                orderNum: Number(currentMenu.orderNum) || 0,
                url: currentMenu.url || '',
                icon: currentMenu.icon || '',
                permission: currentMenu.permission || '',
                status: currentMenu.status || 'show'
            };

            const response = await fetch(url, {
                method,
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify(submitData)
            });

            if (response.ok) {
                showToast('Saved successfully', 'success');
                showDialog = false;
                loadMenus();
            } else {
                const error = await response.json();
                showToast(`Save failed: ${error.message}`, 'error');
            }
        } catch (error) {
            console.error('Save error:', error);
            showToast('Operation failed', 'error');
        }
    }

    onMount(() => {
        loadMenus();
    });
</script>

    <TreeTable
            data={menus}
            columns={columns}
            searchPlaceholder="Search menu name, permission..."
            addButtonText="Add Menu"
            on:search={handleSearch}
    />

<!-- Menu Edit Dialog -->
{#if showDialog}
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 w-full max-w-md max-h-[90vh] overflow-y-auto">
            <h3 class="text-lg font-semibold mb-4">
                {currentMenu.id ? 'Edit Menu' : 'Add Menu'}
            </h3>

            <div class="space-y-4">
                <div>
                    <label class="block text-sm font-medium mb-1">Menu Name *</label>
                    <input
                            type="text"
                            bind:value={currentMenu.menuName}
                            placeholder="Enter menu name"
                            class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                </div>

                <div>
                    <label class="block text-sm font-medium mb-1">Menu Type *</label>
                    <select
                            bind:value={currentMenu.menuType}
                            class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                        <option value="M">Directory (目录)</option>
                        <option value="C">Menu (菜单)</option>
                        <option value="F">Button (按钮)</option>
                    </select>
                </div>

                <div>
                    <label class="block text-sm font-medium mb-1">Parent Menu</label>
                    <select
                            bind:value={currentMenu.parentId}
                            class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                        {#each parentMenuOptions as option}
                            <option value={option.id}
                                    disabled={currentMenu.id === option.id}>
                                {option.menuName}
                            </option>
                        {/each}
                    </select>
                    <p class="text-xs text-gray-500 mt-1">选择父级菜单</p>
                </div>

                <div>
                    <label class="block text-sm font-medium mb-1">Icon</label>
                    <input
                            type="text"
                            bind:value={currentMenu.icon}
                            placeholder="User, Settings, List..."
                            class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    <p class="text-xs text-gray-500 mt-1">图标名称</p>
                </div>

                {#if currentMenu.menuType !== 'F'}
                    <div>
                        <label class="block text-sm font-medium mb-1">Route Path</label>
                        <input
                                type="text"
                                bind:value={currentMenu.url}
                                placeholder="/sys/menu/index"
                                class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                    </div>
                {/if}

                <div>
                    <label class="block text-sm font-medium mb-1">Permission String</label>
                    <input
                            type="text"
                            bind:value={currentMenu.permission}
                            placeholder="sys:menu:query"
                            class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    <p class="text-xs text-gray-500 mt-1">权限标识</p>
                </div>

                <div>
                    <label class="block text-sm font-medium mb-1">Sort Order</label>
                    <input
                            type="number"
                            bind:value={currentMenu.orderNum}
                            class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                </div>
            </div>

            <div class="flex justify-end gap-2 mt-6">
                <button
                        on:click={() => showDialog = false}
                        class="px-4 py-2 border rounded hover:bg-gray-100 transition"
                >
                    Cancel
                </button>
                <button
                        on:click={saveMenu}
                        class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
                >
                    Save
                </button>
            </div>
        </div>
    </div>
{/if}

<!-- 测试权限指令 -->
<div class="fixed bottom-4 right-4 flex gap-2">
    <button use:permission={'sys:user:query'}
            class="px-4 py-2 bg-green-600 text-white rounded shadow">
        有权限可见
    </button>
    <button use:permission={'sys:user:edit'}
            class="px-4 py-2 bg-blue-600 text-white rounded shadow">
        编辑权限
    </button>
    <button use:permission={{any: ['sys:user:edit', 'sys:user:remove']}}
            class="px-4 py-2 bg-yellow-600 text-white rounded shadow">
        编辑或删除
    </button>
</div>