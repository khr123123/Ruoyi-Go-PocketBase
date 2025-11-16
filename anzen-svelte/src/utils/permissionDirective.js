// src/utils/permissionDirective.js
import {hasAll, hasAny, hasPermission} from '../stores/authStore';

/**
 * 权限指令 - 根据权限控制元素显示
 * 用法：
 * <button use:permission={'sys:user:add'}>添加</button>
 * <button use:permission={{perm: 'sys:user:edit'}}>编辑</button>
 * <button use:permission={{any: ['sys:user:edit', 'sys:user:remove']}}>编辑或删除</button>
 * <button use:permission={{all: ['sys:user:edit', 'sys:user:export']}}>编辑且导出</button>
 */
export function permission(node, params) {
    let hasAuth = false;
    // 解析参数
    if (typeof params === 'string') {
        // 单个权限字符串
        hasAuth = hasPermission(params);
    } else if (typeof params === 'object' && params !== null) {
        if (params.perm) {
            // {perm: 'sys:user:add'}
            hasAuth = hasPermission(params.perm);
        } else if (params.any) {
            // {any: ['sys:user:edit', 'sys:user:remove']}
            hasAuth = hasAny(params.any);
        } else if (params.all) {
            // {all: ['sys:user:edit', 'sys:user:export']}
            hasAuth = hasAll(params.all);
        }
    }

    if (!hasAuth) {
        node.style.display = 'none';
    }
    return {
        update(newParams) {
            // 更新权限参数时重新检查
            let newHasAuth = false;

            if (typeof newParams === 'string') {
                newHasAuth = hasPermission(newParams);
            } else if (typeof newParams === 'object' && newParams !== null) {
                if (newParams.perm) {
                    newHasAuth = hasPermission(newParams.perm);
                } else if (newParams.any) {
                    newHasAuth = hasAny(newParams.any);
                } else if (newParams.all) {
                    newHasAuth = hasAll(newParams.all);
                }
            }

            node.style.display = newHasAuth ? '' : 'none';
        },
        destroy() {
            // 清理工作（如需要）
        }
    };
}

/**
 * 更严格的权限指令 - 从 DOM 中完全移除元素
 */
export function permissionStrict(node, params) {
    let hasAuth = false;
    let placeholder = null;

    // 解析参数（同上）
    if (typeof params === 'string') {
        hasAuth = hasPermission(params);
    } else if (typeof params === 'object' && params !== null) {
        if (params.perm) {
            hasAuth = hasPermission(params.perm);
        } else if (params.any) {
            hasAuth = hasAny(params.any);
        } else if (params.all) {
            hasAuth = hasAll(params.all);
        }
    }

    // 如果没有权限，从 DOM 移除
    if (!hasAuth) {
        placeholder = document.createComment('permission-denied');
        node.parentNode?.replaceChild(placeholder, node);
    }

    return {
        destroy() {
            if (placeholder && placeholder.parentNode) {
                placeholder.parentNode.removeChild(placeholder);
            }
        }
    };
}
