<script>
    import {onMount} from 'svelte';
    import {showToast} from '../../stores/toastStore.js';
    import TreeTable from "../../components/TreeTable.svelte";
    import MenuForm from "../../components/MenuForm.svelte";
    import {createMenu, deleteMenu, listMenu, updateMenu} from "../../api/sysApis.js";
    import {defaultMenuData, menuColumns} from "../../config/menuConfig.js";
    import {confirmDialog} from "../../stores/confirmStore.js";

    let menus = [];
    let rawMenus = [];
    let showDialog = false;
    let currentMenu = null;
    let parentMenuOptions = [];

    /** 扁平数据转树形结构 */
    function buildTree(flatList) {
        if (!flatList?.length) return [];

        const map = {};
        const roots = [];

        flatList.forEach(item => map[item.id] = {...item, children: []});

        flatList.forEach(item => {
            const node = map[item.id];
            if (!item.parentId || item.parentId === '0') roots.push(node);
            else if (map[item.parentId]) map[item.parentId].children.push(node);
            else roots.push(node);
        });

        const sortNodes = nodes => {
            nodes.sort((a, b) => (a.orderNum || 0) - (b.orderNum || 0));
            nodes.forEach(n => n.children?.length && sortNodes(n.children));
        };

        sortNodes(roots);
        return roots;
    }

    /** 构建父菜单选项（用于选择父菜单） */
    function buildParentOptions(list) {
        const options = [{id: '0', menuName: 'Root Menu'}];
        list.forEach(menu => {
            if (menu.menuType === 'M' || menu.menuType === 'C') options.push(menu);
        });
        return options;
    }

    //加载菜单列表
    async function loadMenus({search = '', sort = '-created'} = {}) {
        try {
            const filter = search
                  ? `menuName ~ "${search}" || permission ~ "${search}" || url ~ "${search}"` : '';
            const res = await listMenu(1, 500, sort, filter);

            rawMenus = res.items || [];
            menus = buildTree(rawMenus);
            parentMenuOptions = buildParentOptions(rawMenus);
        } catch (err) {
            console.error('Load menus failed:', err);
            showToast('Failed to load menus', 'error');
        }
    }

    // 保存菜单（新增或更新
    async function saveMenu(event) {
        const menuData = event.detail;
        const submitData = {
            menuName: menuData.menuName,
            menuType: menuData.menuType,
            parentId: menuData.parentId || '0',
            orderNum: Number(menuData.orderNum) || 0,
            url: menuData.url || '',
            icon: menuData.icon || '',
            permission: menuData.permission || '',
            status: menuData.status || 'show'
        };
        try {
            if (menuData.id) await updateMenu(menuData.id, submitData);
            else await createMenu(submitData);
            showToast('Saved successfully', 'success');
            showDialog = false;
            loadMenus();
        } catch (err) {
            console.error('Save menu failed:', err);
            if (err.status === 400) {
                showToast('Save failed: No permission', 'error');
            } else {
                showToast(`${err.message}`, 'error');
            }
        }
    }

    // 添加菜单
    function handleAdd() {
        currentMenu = {...defaultMenuData};
        showDialog = true;
    }

    // 编辑菜单
    function handleEdit(event) {
        const menu = event.detail;
        currentMenu = {...menu};
        delete currentMenu.children;
        showDialog = true;
    }
    // 删除菜单
    async function handleDelete(event) {
        const menu = event.detail;
        const ok = await confirmDialog(`Are you sure to delete menu "${menu.menuName}"?`);
        if (!ok) return;
        try {
            await deleteMenu(menu.id);
            showToast('Deleted successfully', 'success');
            loadMenus();
        } catch (err) {
            console.error('Delete menu failed:', err);
            showToast('Operation failed', 'error');
        }
    }

    // 搜索和排序
    const handleSearch = (event) => loadMenus({search: event.detail.search});
    const handleSort = (event) => loadMenus({sort: event.detail.sort});
    // 挂载搜索
    onMount(() => loadMenus());
</script>

<TreeTable
      data={menus}
      columns={menuColumns}
      searchPlaceholder="Search menu name, permission, route..."
      addButtonText="Add Menu"
      showAddButton={true}
      showCheckbox={false}
      actions={{ edit: "sys:menu:edit", delete: "sys:menu:delete",add:"sys:menu:add" }}
      autoExpandOnSearch={true}
      on:add={handleAdd}
      on:edit={handleEdit}
      on:delete={handleDelete}
      on:search={handleSearch}
      on:sort={handleSort}
/>

{#if showDialog}
    <MenuForm
          menu={currentMenu}
          parentOptions={parentMenuOptions}
          isEdit={!!currentMenu.id}
          on:save={saveMenu}
          on:cancel={() => showDialog = false}
          on:error={(e) => showToast(e.detail, 'error')}
    />
{/if}
