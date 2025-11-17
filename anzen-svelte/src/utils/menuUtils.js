// 提取权限字符串
export function extractPermissions(roles) {
    const set = new Set();

    roles?.forEach(role => {
        role.expand?.permission?.forEach(p => {
            if (p.permission) set.add(p.permission);
        });
    });

    return [...set];
}


// 提取菜单数据（只保留 M、C 类型）
export function extractMenus(roles) {
    return roles.flatMap(role =>
        role.expand.permission
            .filter(p => ["M", "C"].includes(p.menuType))
            .map(p => ({...p}))
    );
}


// 把菜单列表转成树结构（带去重）
export function buildTree(list) {
    const map = new Map();
    const roots = [];

    // 预创建 children 并去重
    list.forEach(item => {
        if (!map.has(item.id)) {  // 如果 map 中没有这个 id，就添加
            item.children = [];
            map.set(item.id, item);
        }
    });

    map.forEach(item => {
        if (item.parentId && map.has(item.parentId)) {
            const parent = map.get(item.parentId);
            // 避免重复加入
            if (!parent.children.some(child => child.id === item.id)) {
                parent.children.push(item);
            }
        } else {
            roots.push(item);
        }
    });

    return roots;
}
