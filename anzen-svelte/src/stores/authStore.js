// src/stores/authStore.js
import { writable, get } from 'svelte/store';

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
export const myPermissions = writable(defaultPerms);

myPermissions.subscribe((value) => {
    localStorage.setItem(PERM_KEY, JSON.stringify(value));
});

/* -----------------------------
    从菜单树提取权限字符数组
----------------------------- */
export function extractPermissions(menuTree) {
    const perms = new Set();

    function traverse(nodes) {
        if (!Array.isArray(nodes)) return;

        for (const node of nodes) {
            if (node.permission && node.permission.trim() !== '') {
                perms.add(node.permission.trim());
            }
            if (Array.isArray(node.children) && node.children.length > 0) {
                traverse(node.children);
            }
        }
    }

    traverse(menuTree);

    return Array.from(perms);
}

/* -----------------------------
    权限判断函数
----------------------------- */
export function hasPermission(permission) {
    if (!permission || permission.trim() === '') return true;
    const perms = get(myPermissions);
    return perms.includes(permission.trim());
}

export function hasAny(permArray = []) {
    if (!permArray || permArray.length === 0) return true;
    const perms = get(myPermissions);
    return permArray.some(p => p && perms.includes(p.trim()));
}

export function hasAll(permArray = []) {
    if (!permArray || permArray.length === 0) return true;
    const perms = get(myPermissions);
    return permArray.every(p => p && perms.includes(p.trim()));
}

/* -----------------------------
    设置菜单和权限
----------------------------- */
export function setMenuAndPermissions(menuTree) {
    myRouter.set(menuTree);
    myPermissions.set(extractPermissions(menuTree));
}

/* -----------------------------
    清理本地数据（登出用）
----------------------------- */
export function clearPermissions() {
    localStorage.removeItem(ROUTER_KEY);
    localStorage.removeItem(PERM_KEY);
    myRouter.set([]);
    myPermissions.set([]);
}
