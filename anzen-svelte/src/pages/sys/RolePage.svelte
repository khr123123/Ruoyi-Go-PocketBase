<!-- src/pages/RolePage.svelte -->
<script>
    import {onMount} from 'svelte';
    import {showToast} from '../../stores/toastStore.js';
    import {deleteRole, listMenu, listRole, saveRole, saveRolePermissions} from "../../api/sysApis.js";
    import DataTable from "../../components/DataTable.svelte";
    import {ArrowDown, Image, Menu} from "../../lib/icons/index.js";
    import Drawer from "../../components/Drawer.svelte";
    import {confirmDialog} from "../../stores/confirmStore.js";

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
        {key: 'id', label: 'ID', icon: ArrowDown, sortable: true, class: 'font-mono'},
        {key: 'roleName', label: 'Role Name', icon: Menu, sortable: true, class: "text-xs"},
        {
            key: 'permission',
            label: 'Permissions', icon: Image,
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
        },
        {
            key: 'updated',
            label: 'Updated',
            class: 'text-xs text-gray-600',
            render: (v) => new Date(v).toLocaleString('zh-CN')
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

    const handleSearch = (event) => loadRoles({search: event.detail.search, page: 1});
    const handleSort = (event) => loadRoles({sort: event.detail.sort, page: event.detail.page});
    const handlePageChange = (event) => loadRoles({page: event.detail.page});

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

    function handleEdit(event) {
        const role = event.detail
        currentRole = {...role};
        showDialog = true;
    }

    function handelAssignPerm(event) {
        const role = event.detail
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

    async function saveRoleHandler() {
        try {
            await saveRole(currentRole);
            showToast('Saved successfully', 'success');
            showDialog = false;
            loadRoles({page});
        } catch (error) {
            console.error('Save error:', error);
            showToast('Operation failed', 'error');
        }
    }

    async function savePermissionsHandler() {
        try {
            await saveRolePermissions(currentRole.id, selectedPermissions);
            showToast('Permissions saved successfully', 'success');
            showPermissionDialog = false;
            loadRoles({page});
        } catch (error) {
            console.error('Save permissions error:', error);
            showToast('Operation failed', 'error');
        }
    }

    async function handleDelete(event) {
        const role = event.detail
        const ok = await confirmDialog(`Are you sure to delete role "${role.roleName}"?`);
        if (!ok) return;
        try {
            await deleteRole(role.id);
            showToast('Deleted successfully', 'success');
            loadRoles({page});
        } catch (error) {
            console.error('Delete error:', error);
            showToast('Operation failed', 'error');
        }
    }


    onMount(async () => {
        loadRoles();
        const res = await listMenu(1, 500);
        allMenus = res.items || [];
    });
</script>

<DataTable
      data={roles}
      {total}
      {columns}
      {page}
      searchPlaceholder="Search role name, role key..."
      addButtonText="Add Role"
      actions={{ edit: "sys:role:edit", delete: "sys:role:delete",add:"sys:role:add",assignPerm:"sys:role:perm" }}
      on:add={handleAdd}
      on:assignPerm={handelAssignPerm}
      on:edit={handleEdit}
      on:delete={handleDelete}
      on:search={handleSearch}
      on:sort={handleSort}
      on:pageChange={handlePageChange}
/>
<Drawer show={showDialog} title={currentRole?.id ? 'Edit Role' : 'Add Role'} position="right"
        on:close={() => showDialog=false}>
    <div class="space-y-4">
        <div>
            <label class="block text-sm font-medium mb-1">Role Name *</label>
            <input type="text" bind:value={currentRole.roleName} placeholder="Enter role name"
                   class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"/>
        </div>
        <div class="flex justify-end gap-2 mt-4">
            <button on:click={() => showDialog=false} class="px-4 py-2 border rounded hover:bg-gray-100 transition">
                Cancel
            </button>
            <button on:click={saveRoleHandler}
                    class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition">Save
            </button>
        </div>
    </div>
</Drawer>

<Drawer show={showPermissionDialog} title={`Configure Permissions - ${currentRole?.roleName}`} position="right"
        width="600px" on:close={() => showPermissionDialog=false}>
    <div class="space-y-2 max-h-[62vh] overflow-y-auto border rounded p-4">
        {#each allMenus as menu}
            <label class="flex items-center gap-2 p-2 hover:bg-gray-50 rounded cursor-pointer">
                <input type="checkbox" checked={selectedPermissions.includes(menu.id)}
                       on:change={() => togglePermission(menu.id)}
                       class="rounded border-gray-300"/>
                <span class="flex-1">
                    {#if menu.icon}
                        <i class="{menu.icon} mr-2 text-gray-600"></i>
                    {/if}
                    {menu.menuName}
                </span>
                <span class="text-xs text-gray-500 font-mono">{menu.permission || '-'}</span>
                <span class="text-xs px-2 py-1 rounded bg-gray-100">
                    {menu.menuType === 'M' ? 'Dir' : menu.menuType === 'C' ? 'MenuPage' : 'Button'}
                </span>
            </label>
        {/each}
    </div>

    <div class="mt-4 p-3 bg-blue-50 rounded">
        <p class="text-sm text-blue-800">
            Selected: <strong>{selectedPermissions.length}</strong> permissions
        </p>
    </div>

    <div class="flex justify-end gap-2 mt-4">
        <button on:click={() => showPermissionDialog=false}
                class="px-4 py-2 border rounded hover:bg-gray-100 transition">Cancel
        </button>
        <button on:click={savePermissionsHandler}
                class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition">Save Permissions
        </button>
    </div>
</Drawer>
