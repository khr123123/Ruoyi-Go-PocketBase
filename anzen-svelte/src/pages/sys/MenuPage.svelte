<script>
    import {onMount} from 'svelte';
    import {hasPermission} from '../../stores/authStore.js';
    import {showToast} from '../../stores/toastStore.js';
    import TreeTable from "../../components/TreeTable.svelte";
    import {listMenu} from "../../api/sysApis.js";
    import {user} from "../../stores/userStore.js";

    let menus = [];
    let rawMenus = [];
    let showDialog = false;
    let currentMenu = null;
    let parentMenuOptions = [];

    /**
     * ========================================
     * 列配置
     * ========================================
     * TreeTable 的 columns 参数配置示例
     */
    const columns = [
        {
            key: 'menuName',        // 数据字段名
            label: 'MenuPage Name',     // 表头显示
            sortable: true,         // 支持排序
            visible: true           // 默认显示
        },
        {
            key: 'icon',
            label: 'Icon',
            sortable: false,        // 不支持排序
            class: 'text-xs font-mono text-gray-600'
        },
        {
            key: 'menuType',
            label: 'Type',
            sortable: true,
            // 自定义渲染函数
            render: (v) => {
                const types = {M: 'Directory', C: 'Menu', F: 'Button'};
                const colors = {M: 'blue', C: 'green', F: 'yellow'};
                return `<span class="px-2 py-1 text-xs rounded bg-${colors[v] || 'gray'}-100 text-${colors[v] || 'gray'}-700">${types[v] || v || '-'}</span>`;
            }
        },
        {
            key: 'orderNum',
            label: 'Sort',
            sortable: true,
            render: (v) => `<span class="px-2 py-1 text-xs rounded bg-gray-100 text-gray-700">${v ?? '-'}</span>`
        },
        {
            key: 'permission',
            label: 'Permission',
            sortable: false,
            class: 'text-xs font-mono',
            render: (v) => v ? `<code class="text-blue-600">${v}</code>` : '<span class="text-gray-400">-</span>'
        },
        {
            key: 'url',
            label: 'Route',
            sortable: false,
            class: 'text-xs',
            render: (v) => v ? `<code class="text-green-600">${v}</code>` : '<span class="text-gray-400">-</span>'
        }
    ];

    /**
     * ========================================
     * 数据处理函数
     * ========================================
     */

    /**
     * buildTree - 构建树形结构
     * 将扁平数据转换为树形结构
     */
    function buildTree(flatList) {
        if (!flatList || flatList.length === 0) return [];

        const map = {};
        const roots = [];

        // 创建节点映射
        flatList.forEach(item => {
            map[item.id] = {...item, children: []};
        });

        // 构建父子关系
        flatList.forEach(item => {
            const node = map[item.id];
            if (!item.parentId || item.parentId === '0') {
                roots.push(node);
            } else if (map[item.parentId]) {
                map[item.parentId].children.push(node);
            } else {
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

    /**
     * buildParentOptions - 构建父菜单选项
     * 用于编辑对话框的父菜单选择器
     */
    function buildParentOptions() {
        parentMenuOptions = [{id: '0', menuName: 'Root MenuPage'}];
        rawMenus.forEach(menu => {
            if (menu.menuType === 'M' || menu.menuType === 'C') {
                parentMenuOptions.push(menu);
            }
        });
    }

    /**
     * ========================================
     * API 交互函数
     * ========================================
     */

    /**
     * loadMenus - 加载菜单列表
     */
    async function loadMenus(params = {}) {
        try {
            const {search = '', sort = '-created', page: p = 1} = params;
            const filter = search
                  ? `menuName ~ "${search}" || permission ~ "${search}" || url ~ "${search}"`
                  : '';

            const res = await listMenu(1, 500, sort, filter);
            rawMenus = res.items || [];
            menus = buildTree(rawMenus);
            buildParentOptions();
        } catch (error) {
            console.error('Load menu error:', error);
            showToast('Failed to load menus', 'error');
        }
    }

    /**
     * saveMenu - 保存菜单（新增或编辑）
     */
    async function saveMenu() {
        if (!currentMenu.menuName || !currentMenu.menuType) {
            showToast('Please fill in required fields', 'error');
            return;
        }
        try {
            const url = currentMenu.id
                  ? `http://127.0.0.1:8090/api/collections/sys_menu/records/${currentMenu.id}`
                  : 'http://127.0.0.1:8090/api/collections/sys_menu/records';
            const method = currentMenu.id ? 'PATCH' : 'POST';
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
                    'Authorization': `Bearer ${$user.token}`
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

    /**
     * ========================================
     * 事件处理函数
     * ========================================
     */
    function handleSearch(event) {
        loadMenus({search: event.detail.search, page: 1});
    }

    function handleSort(event) {
        loadMenus({sort: event.detail.sort, page: event.detail.page});
    }

    function handlePageChange(event) {
        loadMenus({page: event.detail.page});
    }

    /**
     * handleAdd - 处理添加按钮点击
     */
    function handleAdd() {
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

    /**
     * handleEdit - 处理编辑按钮点击
     */
    function handleEdit(event) {
        const menu = event.detail;

        if (!hasPermission('sys:menu:edit')) {
            showToast('No permission', 'error');
            return;
        }

        currentMenu = {...menu};
        delete currentMenu.children;
        showDialog = true;
    }

    /**
     * handleDelete - 处理删除按钮点击
     */
    async function handleDelete(event) {
        const menu = event.detail;

        if (!hasPermission('sys:menu:remove')) {
            showToast('No permission', 'error');
            return;
        }

        if (!confirm(`Are you sure to delete menu "${menu.menuName}"?`)) {
            return;
        }

        try {
            const response = await fetch(
                  `http://127.0.0.1:8090/api/collections/sys_menu/records/${menu.id}`,
                  {
                      method: 'DELETE',
                      headers: {
                          'Authorization': `Bearer ${localStorage.getItem('token')}`
                      }
                  }
            );

            if (response.ok) {
                showToast('Deleted successfully', 'success');
                loadMenus();
            } else {
                showToast('Delete failed', 'error');
            }
        } catch (error) {
            console.error('Delete error:', error);
            showToast('Operation failed', 'error');
        }
    }

    onMount(() => {
        loadMenus();
    });
</script>

<!--
    ========================================
    TreeTable 组件使用示例
    ========================================

    基础用法：
    <TreeTable
        data={treeData}           // 必需：树形数据
        columns={columnConfig}    // 必需：列配置
    />

    完整配置：
    <TreeTable
        data={treeData}
        columns={columnConfig}
        searchPlaceholder="Search..."
        addButtonText="Add MenuPage"
        showAddButton={true}
        showCheckbox={true}
        actions={['edit', 'delete']}
        onAdd={handleAdd}
        on:edit={handleEdit}
        on:delete={handleDelete}
        on:search={handleSearch}
        on:sort={handleSort}
    />
-->
<TreeTable
      data={menus}
      {columns}
      searchPlaceholder="Search menu name, permission, route..."
      addButtonText="Add Menu"
      showAddButton={true}
      showApiPreview={false}
      showCheckbox={true}
      actions={{"edit":"sys:menu:edit","add":"sys:menu:add",}}
      autoExpandOnSearch={true}
      onAdd={handleAdd}
      on:edit={handleEdit}
      on:delete={handleDelete}
      on:search={handleSearch}
      on:sort={handleSort}
/>

<!-- 编辑对话框 -->
{#if showDialog}
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 w-full max-w-md max-h-[90vh] overflow-y-auto">
            <h3 class="text-lg font-semibold mb-4">
                {currentMenu.id ? 'Edit MenuPage' : 'Add MenuPage'}
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
                            <option value={option.id} disabled={currentMenu.id === option.id}>
                                {option.menuName}
                            </option>
                        {/each}
                    </select>
                </div>

                <div>
                    <label class="block text-sm font-medium mb-1">Icon</label>
                    <input
                          type="text"
                          bind:value={currentMenu.icon}
                          placeholder="User, Settings, List..."
                          class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
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