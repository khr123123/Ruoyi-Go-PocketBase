// src/stores/authStore.js
import {get, writable} from 'svelte/store';

/* -----------------------------
    LocalStorage Keys
----------------------------- */
const ROUTER_KEY = "ruoyi-pb-router";
const PERM_KEY = "ruoyi-pb-perms";

/* -----------------------------
    初始化：菜单路由
----------------------------- */
const storedRouter = localStorage.getItem(ROUTER_KEY);
const defaultRouter = storedRouter ? JSON.parse(storedRouter) : [];
export const myRouter = writable(defaultRouter);

myRouter.subscribe((value) => {
    localStorage.setItem(ROUTER_KEY, JSON.stringify(value));
});

/* -----------------------------
    初始化：权限字符串数组
----------------------------- */
const storedPerms = localStorage.getItem(PERM_KEY);
const defaultPerms = storedPerms ? JSON.parse(storedPerms) : [];
export const permissions = writable(defaultPerms);

permissions.subscribe((value) => {
    localStorage.setItem(PERM_KEY, JSON.stringify(value));
});

/* -----------------------------
    从菜单树提取权限字符数组
----------------------------- */
export function extractPermissions(menuTree) {
    let perms = [];

    function traverse(nodes) {
        if (!nodes || !Array.isArray(nodes)) return;

        nodes.forEach(node => {
            // 提取当前节点的权限
            if (node.permission && node.permission.trim() !== '') {
                perms.push(node.permission.trim());
            }

            // 递归提取子节点权限
            if (node.children && Array.isArray(node.children) && node.children.length > 0) {
                traverse(node.children);
            }
        });
    }

    traverse(menuTree);

    // 去重并返回
    return [...new Set(perms)];
}

/* -----------------------------
    获取所有权限（用于调试）
----------------------------- */
export function getAllPermissions() {
    return get(permissions);
}

/* -----------------------------
    设置：菜单 & 权限
----------------------------- */
export function setMenuAndPermissions(menuTree) {
    myRouter.set(menuTree);
    const perms = extractPermissions(menuTree);
    permissions.set(perms);
    console.log('🔐 Extracted permissions:', perms);
}

/* -----------------------------
    权限判断：hasPermission
----------------------------- */
export function hasPermission(permission) {
    if (!permission || permission.trim() === '') return true; // 空权限默认允许

    const perms = get(permissions);
    return perms.includes(permission.trim());
}

/* -----------------------------
    多权限：任意一个满足
----------------------------- */
export function hasAny(permArray = []) {
    if (!permArray || permArray.length === 0) return true;

    const perms = get(permissions);
    return permArray.some(p => p && perms.includes(p.trim()));
}

/* -----------------------------
    多权限：必须全部满足
----------------------------- */
export function hasAll(permArray = []) {
    if (!permArray || permArray.length === 0) return true;

    const perms = get(permissions);
    return permArray.every(p => p && perms.includes(p.trim()));
}

/* -----------------------------
    清理本地数据（登出时用）
----------------------------- */
export function clearPermissions() {
    localStorage.removeItem(ROUTER_KEY);
    localStorage.removeItem(PERM_KEY);

    myRouter.set([]);
    permissions.set([]);
}