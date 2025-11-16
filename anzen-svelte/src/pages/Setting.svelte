<script>

    import TreeTable from "../components/TreeTable.svelte";

    let menus = [
        {
            id: "1",
            name: "系统管理",
            children: [
                {
                    id: "11",
                    name: "用户管理",
                    children: [
                        { id: "111", name: "用户列表" },
                        { id: "112", name: "权限分配" }
                    ]
                },
                {
                    id: "12",
                    name: "菜单管理",
                    children: [
                        { id: "121", name: "新增菜单" },
                        { id: "122", name: "编辑菜单" }
                    ]
                }
            ]
        }
    ];

    let columns = [
        { label: "名称", key: "name" },
        { label: "路径", key: "path" }
    ];
</script>

<TreeTable
        data={menus}
        columns={columns}
        onEdit={row => console.log("edit", row)}
        onDelete={row => console.log("delete", row)}
        onCheck={row => console.log("check", row)}
/>
