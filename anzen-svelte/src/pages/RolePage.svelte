<!-- src/pages/RolePage.svelte -->
<script>
    import DataTable from '../components/DataTable.svelte';
    import {onMount} from 'svelte';
    import {hasPermission} from '../stores/authStore';
    import {showToast} from '../stores/toastStore';
    import {listRole} from "../api/roleApis.js";

    let roles = [];
    let total = 0;
    let page = 1;
    let perPage = 20;
    let showDialog = false;
    let showPermissionDialog = false;
    let currentRole = null;
    let allMenus = [];
    let selectedPermissions = [];

    // Column definitions with operations
    const columns = [
        {key: 'id', label: 'ID', sortable: true, class: 'font-mono text-xs'},
        {key: 'roleName', label: 'Role Name', sortable: true},
        {key: 'roleKey', label: 'Role Key', class: 'font-mono text-sm'},
        {
            key: 'status',
            label: 'Status',
            render: (v) => {
                const isNormal = v === 'normal';
                return `<span class="px-2 py-1 text-xs rounded ${isNormal ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}">
                    ${isNormal ? 'Enabled' : 'Disabled'}
                </span>`;
            }
        },
        {
            key: 'permission',
            label: 'Permissions',
            render: (v) => {
                if (!v || v.length === 0) return '<span class="text-gray-400 text-xs">No permissions</span>';
                return `<span class="px-2 py-1 text-xs rounded bg-blue-100 text-blue-700">${v.length} items</span>`;
            }
        },
        {
            key: 'created',
            label: 'Created',
            class: 'text-xs text-gray-600',
            render: (v) => new Date(v).toLocaleString('zh-CN')
        }
    ];

    // Custom action buttons
    const customActions = [
        {
            label: 'Permissions',
            icon: '🔐',
            class: 'text-purple-600 hover:text-purple-800',
            show: () => hasPermission('system:role:edit'),
            onClick: handleConfigPermission
        }
    ];

    // Load roles
    async function loadRoles(params = {}) {
        try {
            const {search = '', sort = '-created', page: p = 1} = params;
            const filter = search ? `roleName ~ "${search}" || roleKey ~ "${search}"` : '';
            const res = await listRole(p, perPage, sort, filter);

            roles = res.items;
            total = res.totalItems;
            page = p;
        } catch (error) {
            console.error('Load roles error:', error);
            showToast('Failed to load roles', 'error');
        }
    }

    // Load all menus
    async function loadAllMenus() {
        try {
            const response = await fetch('http://127.0.0.1:8090/api/getAllMenuTree', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });
            const result = await response.json();
            if (result.code === 200) {
                allMenus = flattenMenuTree(result.data);
            }
        } catch (error) {
            console.error('Load menus error:', error);
        }
    }

    function flattenMenuTree(tree, result = []) {
        tree.forEach(node => {
            result.push(node);
            if (node.children && node.children.length > 0) {
                flattenMenuTree(node.children, result);
            }
        });
        return result;
    }

    function handleSearch(event) {
        loadRoles({search: event.detail.search, page: 1});
    }

    function handleSort(event) {
        loadRoles({sort: event.detail.sort, page: event.detail.page});
    }

    function handlePageChange(event) {
        loadRoles({page: event.detail.page});
    }

    function handleAdd() {
        currentRole = {
            roleName: '',
            roleKey: '',
            status: 'normal',
            permission: [],
            remark: ''
        };
        showDialog = true;
    }

    function handleEdit(role) {
        currentRole = {...role};
        showDialog = true;
    }

    function handleConfigPermission(role) {
        currentRole = {...role};
        selectedPermissions = role.permission || [];
        showPermissionDialog = true;
    }

    function togglePermission(menuId) {
        const index = selectedPermissions.indexOf(menuId);
        if (index > -1) {
            selectedPermissions = selectedPermissions.filter(id => id !== menuId);
        } else {
            selectedPermissions = [...selectedPermissions, menuId];
        }
    }

    async function saveRole() {
        if (!currentRole.roleName || !currentRole.roleKey) {
            showToast('Role name and role key are required', 'error');
            return;
        }

        try {
            const url = currentRole.id
                ? `http://127.0.0.1:8090/api/collections/sys_role/records/${currentRole.id}`
                : 'http://127.0.0.1:8090/api/collections/sys_role/records';
            const method = currentRole.id ? 'PATCH' : 'POST';

            const response = await fetch(url, {
                method,
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify(currentRole)
            });

            if (response.ok) {
                showToast('Saved successfully', 'success');
                showDialog = false;
                loadRoles();
            } else {
                const error = await response.json();
                showToast(`Save failed: ${error.message}`, 'error');
            }
        } catch (error) {
            console.error('Save error:', error);
            showToast('Operation failed', 'error');
        }
    }

    async function savePermissions() {
        try {
            currentRole.permission = selectedPermissions;

            const response = await fetch(
                `http://127.0.0.1:8090/api/collections/sys_role/records/${currentRole.id}`,
                {
                    method: 'PATCH',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    },
                    body: JSON.stringify({permission: selectedPermissions})
                }
            );

            if (response.ok) {
                showToast('Permissions saved successfully', 'success');
                showPermissionDialog = false;
                loadRoles();
            } else {
                showToast('Failed to save permissions', 'error');
            }
        } catch (error) {
            console.error('Save permissions error:', error);
            showToast('Operation failed', 'error');
        }
    }

    async function handleDelete(role) {
        if (!confirm(`Are you sure to delete role "${role.roleName}"?`)) {
            return;
        }

        try {
            const response = await fetch(
                `http://127.0.0.1:8090/api/collections/sys_role/records/${role.id}`,
                {
                    method: 'DELETE',
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    }
                }
            );

            if (response.ok) {
                showToast('Deleted successfully', 'success');
                loadRoles();
            } else {
                showToast('Delete failed', 'error');
            }
        } catch (error) {
            console.error('Delete error:', error);
            showToast('Operation failed', 'error');
        }
    }

    onMount(() => {
        loadRoles();
        loadAllMenus();
    });
</script>

<DataTable
        data={roles}
        {total}
        {columns}
        {page}
        searchPlaceholder="Search role name, role key..."
        addButtonText="Add Role"
        showOperations={true}
        {customActions}
        canEdit={hasPermission('system:role:edit')}
        canDelete={hasPermission('system:role:delete')}
        onAdd={handleAdd}
        onEdit={handleEdit}
        onDelete={handleDelete}
        on:search={handleSearch}
        on:sort={handleSort}
        on:pageChange={handlePageChange}
/>

<!-- Role Edit Dialog -->
{#if showDialog}
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 w-full max-w-md">
            <h3 class="text-lg font-semibold mb-4">
                {currentRole.id ? 'Edit Role' : 'Add Role'}
            </h3>

            <div class="space-y-4">
                <div>
                    <label class="block text-sm font-medium mb-1">Role Name *</label>
                    <input
                            type="text"
                            bind:value={currentRole.roleName}
                            placeholder="Enter role name"
                            class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                </div>

                <div>
                    <label class="block text-sm font-medium mb-1">Role Key *</label>
                    <input
                            type="text"
                            bind:value={currentRole.roleKey}
                            placeholder="e.g., admin, user"
                            class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                </div>

                <div>
                    <label class="block text-sm font-medium mb-1">Status</label>
                    <select
                            bind:value={currentRole.status}
                            class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                        <option value="normal">Enabled</option>
                        <option value="disabled">Disabled</option>
                    </select>
                </div>

                <div>
                    <label class="block text-sm font-medium mb-1">Remark</label>
                    <textarea
                            bind:value={currentRole.remark}
                            placeholder="Role description..."
                            rows="3"
                            class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    ></textarea>
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
                        on:click={saveRole}
                        class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
                >
                    Save
                </button>
            </div>
        </div>
    </div>
{/if}

<!-- Permission Dialog -->
{#if showPermissionDialog}
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
            <h3 class="text-lg font-semibold mb-4">
                Configure Permissions - {currentRole.roleName}
            </h3>

            <div class="space-y-2 max-h-96 overflow-y-auto border rounded p-4">
                {#each allMenus as menu}
                    <label class="flex items-center gap-2 p-2 hover:bg-gray-50 rounded cursor-pointer">
                        <input
                                type="checkbox"
                                checked={selectedPermissions.includes(menu.id)}
                                on:change={() => togglePermission(menu.id)}
                                class="rounded border-gray-300"
                        />
                        <span class="flex-1">
                            {#if menu.icon}
                                <i class="{menu.icon} mr-2 text-gray-600"></i>
                            {/if}
                            {menu.menuName}
                        </span>
                        <span class="text-xs text-gray-500 font-mono">
                            {menu.permission || '-'}
                        </span>
                        <span class="text-xs px-2 py-1 rounded bg-gray-100">
                            {menu.menuType === 'M' ? 'Dir' : menu.menuType === 'C' ? 'Menu' : 'Button'}
                        </span>
                    </label>
                {/each}
            </div>

            <div class="mt-4 p-3 bg-blue-50 rounded">
                <p class="text-sm text-blue-800">
                    Selected: <strong>{selectedPermissions.length}</strong> permissions
                </p>
            </div>

            <div class="flex justify-end gap-2 mt-6">
                <button
                        on:click={() => showPermissionDialog = false}
                        class="px-4 py-2 border rounded hover:bg-gray-100 transition"
                >
                    Cancel
                </button>
                <button
                        on:click={savePermissions}
                        class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
                >
                    Save Permissions
                </button>
            </div>
        </div>
    </div>
{/if}