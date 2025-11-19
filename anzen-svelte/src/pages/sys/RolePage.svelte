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
<Drawer show={showDialog} title={currentRole?.id ? 'Edit role record' : 'New role record'} position="right"
        on:close={() => showDialog=false}>
    <div class="space-y-5">
        <div>
            <label class="block text-xs font-medium text-gray-600 mb-1.5">
                roleName <span class="text-red-500">*</span>
            </label>
            <input
                  type="text"
                  bind:value={currentRole.roleName}
                  placeholder="e.g., Administrator, Editor"
                  class="w-full bg-white border border-gray-300 rounded px-3 py-2 text-sm text-gray-900 placeholder-gray-400
                       focus:border-gray-500 focus:outline-none focus:ring-1 focus:ring-gray-500"
            />
        </div>

        <div class="flex justify-end gap-2 pt-4 border-t border-gray-200">
            <button
                  on:click={() => showDialog=false}
                  class="px-4 py-2 border border-gray-300 rounded text-sm hover:bg-gray-50 transition-colors">
                Cancel
            </button>
            <button
                  on:click={saveRoleHandler}
                  class="px-4 py-2 bg-gray-800 text-white rounded text-sm hover:bg-gray-700 transition-colors">
                Save changes
            </button>
        </div>
    </div>
</Drawer>

<Drawer show={showPermissionDialog} title={`Permissions - ${currentRole?.roleName}`} position="right"
        width="600px" on:close={() => showPermissionDialog=false}>
    <div class="space-y-4">
        <div class="max-h-[calc(100vh-230px)] overflow-y-auto border border-gray-300 rounded bg-white">
            {#each allMenus as menu}
                <label
                      class="flex items-center gap-3 px-4 py-2.5 hover:bg-gray-50 cursor-pointer border-b border-gray-100 last:border-b-0">
                    <input
                          type="checkbox"
                          checked={selectedPermissions.includes(menu.id)}
                          on:change={() => togglePermission(menu.id)}
                          class="w-4 h-4 rounded border-gray-300 text-gray-800 focus:ring-gray-500"
                    />
                    <span class="flex-1 text-sm text-gray-900">
                        {#if menu.icon}
                            <i class="{menu.icon} mr-2 text-gray-500 text-xs"></i>
                        {/if}
                        {menu.menuName}
                    </span>
                    {#if menu.permission}
                        <span class="text-xs text-gray-500 font-mono bg-gray-100 px-2 py-0.5 rounded">{menu.permission}</span>
                    {/if}
                    <span class="text-xs px-2 py-0.5 rounded bg-gray-100 text-gray-600">
                        {menu.menuType === 'M' ? 'Dir' : menu.menuType === 'C' ? 'Page' : 'Btn'}
                    </span>
                </label>
            {/each}
        </div>

        <div class="px-4 py-3 bg-blue-50 rounded border border-blue-100">
            <p class="text-sm text-blue-900">
                Selected: <strong class="font-semibold">{selectedPermissions.length}</strong> permission(s)
            </p>
        </div>

        <div class="flex justify-end gap-2 pt-4 border-t border-gray-200">
            <button
                  on:click={() => showPermissionDialog=false}
                  class="px-4 py-2 border border-gray-300 rounded text-sm hover:bg-gray-50 transition-colors">
                Cancel
            </button>
            <button
                  on:click={savePermissionsHandler}
                  class="px-4 py-2 bg-gray-800 text-white rounded text-sm hover:bg-gray-700 transition-colors">
                Save changes
            </button>
        </div>
    </div>
</Drawer>