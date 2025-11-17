import {hasAll, hasAny, hasPermission} from "../stores/authStore.js";

export function permission(node, params) {
    let placeholder = null;

    function check() {
        if (typeof params === 'string') return hasPermission(params);
        if (typeof params === 'object' && params !== null) {
            if (params.perm) return hasPermission(params.perm);
            if (params.any) return hasAny(params.any);
            if (params.all) return hasAll(params.all);
        }
        return false;
    }

    // 只挂载时检查一次，不操作根节点
    if (!check() && node.parentNode) {
        placeholder = document.createComment('permission-denied,FUCK YOU');
        node.parentNode.replaceChild(placeholder, node);
    }
    return {
        destroy() {
            placeholder = null;
        }
    };
}
