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


export function buildTree(list) {
    const map = new Map();
    const roots = [];

    // 预创建 children 并去重
    list.forEach(item => {
        if (!map.has(item.id)) {
            item.children = [];
            map.set(item.id, item);
        }
    });

    // 组装树结构
    map.forEach(item => {
        if (item.parentId && map.has(item.parentId)) {
            const parent = map.get(item.parentId);

            if (!parent.children.some(child => child.id === item.id)) {
                parent.children.push(item);
            }
        } else {
            roots.push(item);
        }
    });

    // 按 orderNum 排序整个树
    sortTree(roots);

    return roots;
}

// --- 按 orderNum 排序（递归）
function sortTree(nodes) {
    nodes.sort((a, b) => {
        const sa = a.orderNum ?? 0;
        const sb = b.orderNum ?? 0;
        return sa - sb;
    });

    nodes.forEach(node => {
        if (node.children?.length) {
            sortTree(node.children);
        }
    });
}

