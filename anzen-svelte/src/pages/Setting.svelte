<script>

    import {permission} from "../utils/permissionDirective.js"

    let menus = [
        {
            id: "1",
            name: "系统管理",
            children: [
                {
                    id: "11",
                    name: "用户管理",
                    children: [
                        {id: "111", name: "用户列表"},
                        {id: "112", name: "权限分配"}
                    ]
                },
                {
                    id: "12",
                    name: "菜单管理",
                    children: [
                        {id: "121", name: "新增菜单"},
                        {id: "122", name: "编辑菜单"}
                    ]
                }
            ]
        }
    ];

    let columns = [
        {label: "名称", key: "name"},
        {label: "路径", key: "path"}
    ];
</script>
<div>
    <button use:permission={''}>kong</button>
    <button use:permission={'sys:user:edit'}>sys:user:edit</button>
    <button use:permission={'sys:menu:add'}>sys:menu:add</button>
    <button use:permission={'sys:user:delete'}>sys:user:delete</button>
    <button use:permission={'sys:user:add'}>sys:user:add</button>
    <button use:permission={'sys:user:add1'}>sys:user:add1</button>
</div>
