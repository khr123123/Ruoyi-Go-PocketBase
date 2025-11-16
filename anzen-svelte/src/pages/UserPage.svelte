<!-- src/pages/UserPage.svelte -->
<script>
    import DataTable from '../components/DataTable.svelte';
    import {listUser} from '../api/userApis.js';
    import {ArrowDown, Check, Circle, Image, Menu} from '../lib/icons';

    const AVATAR_PREFIX = "http://127.0.0.1:8090/api/files/_pb_users_auth_/";

    let users = [];
    let total = 0;
    let page = 1;
    let perPage = 10;

    // 定义列配置
    const columns = [
        {key: 'id', label: 'id', icon: ArrowDown, sortable: true},
        {key: 'email', label: 'email', icon: Menu, sortable: true},
        {
            key: 'emailVisibility',
            label: 'emailVisibility',
            icon: Circle,
            render: (value) => `<span class="px-2 py-0.5 rounded-full text-xs ${value ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'}">${value ? 'true' : 'false'}</span>`
        },
        {
            key: 'verified',
            label: 'verified',
            icon: Check,
            render: (value) => `<span class="px-2 py-0.5 rounded-full text-xs ${value ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'}">${value ? 'true' : 'false'}</span>`
        },
        {key: 'name', label: 'name', sortable: true},
        {
            key: 'avatar',
            label: 'avatar',
            icon: Image,
            sortable: false,
            render: (value, row) => `<img class="w-8 h-8 rounded-full object-cover" src="${AVATAR_PREFIX}${row.id}/${value}" alt=""/>`
        },
        {key: 'created', label: 'created', class: 'text-xs text-gray-600'},
        {key: 'updated', label: 'updated', class: 'text-xs text-gray-600'}
    ];

    async function loadUsers(params = {}) {
        const {search = '', sort = '-created', page: p = 1} = params;
        const filter = search ? `email ~ "${search}" || name ~ "${search}"` : '';
        const res = await listUser(p, perPage, sort, filter);
        users = res.items;
        total = res.totalItems;
        page = p;
    }

    function handleSearch(event) {
        loadUsers({search: event.detail.search, page: 1});
    }

    function handleSort(event) {
        loadUsers({sort: event.detail.sort, page: event.detail.page});
    }

    function handlePageChange(event) {
        loadUsers({page: event.detail.page});
    }

    loadUsers();
</script>

<DataTable data={users}
           {total}
           {columns}
           {page}
           searchPlaceholder="Search email or name..."
           on:search={handleSearch}
           on:sort={handleSort}
           on:pageChange={handlePageChange}
/>