// userApis.js
// PocketBase API 封装（登录、注册、获取用户信息等）

import PocketBase from "pocketbase";

// 初始化 PocketBase 客户端
const pb = new PocketBase("http://127.0.0.1:8090");

// --- 用户相关 API ---

/**
 * 登录
 * @param {string} email
 * @param {string} password
 */
export async function login(email, password) {
    return await pb.collection("users").authWithPassword(email, password);
}

/**
 * 注册
 * @param {Object} user
 * user = { email, password, passwordConfirm, username }
 */
export async function register(user) {
    return await pb.collection("users").create(user);
}

/**
 * 获取当前登录用户信息
 */
export function getCurrentUser() {
    return pb.authStore.model;
}

/**
 * 退出登录
 */
export function logout() {
    pb.authStore.clear();
}

/**
 * 根据 ID 获取用户信息
 * @param {string} id
 */
export async function getUserById(id) {
    return await pb.collection("users").getOne(id);
}

/**
 * 更新用户信息
 * @param {string} id
 * @param {Object} data
 */
export async function updateUser(id, data) {
    return await pb.collection("users").update(id, data);
}

/**
 * 删除用户（需要管理员权限）
 * @param {string} id
 */
export async function deleteUser(id) {
    return await pb.collection("users").delete(id);
}

// 导出 PocketBase 实例（可直接用）
export {pb};