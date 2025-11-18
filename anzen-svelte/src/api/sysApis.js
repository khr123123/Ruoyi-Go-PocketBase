// sysApis.js　　PocketBase API 封装
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
    return pb.authStore.record;
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

export async function saveUser(userData) {
    if (userData.id) {
        // 更新用户
        return await pb.collection('users').update(userData.id, userData);
    } else {
        // 创建用户
        return await pb.collection('users').create(userData);
    }
}

/**
 * 获取用户列表带排序
 * @param {number} page
 * @param {number} perPage
 * @param {string} sort
 * @param {string} filter
 */
export async function listUser(page = 1, perPage = 20, sort = '-created', filter = '') {
    return await pb.collection('users').getList(page, perPage, {
        sort,       // 排序字段，例如："-created", "email", "name"
        filter,     // PocketBase filter，例："email ~ 'gmail'"
        fields: 'id,email,emailVisibility,verified,name,avatar,created,updated,role', //你要的字段
        expand: 'role', //如果你想取角色表
    });
}

export async function listMenu(page = 1, perPage = 20, sort = '-created', filter = '') {
    return await pb.collection('sys_menu').getList(page, perPage, {
        sort, filter, fields: 'id,menuName,menuType,orderNum,permission,parentId,url,icon,created,updated',
    });
}

export async function getUserRouter() {
    const user = pb.authStore.record;
    return await pb.collection("users").getOne(user.id, {
        expand: "role,role.permission",
    });
}

/**
 *  获取菜单列表带排序
 * @param page
 * @param perPage
 * @param sort
 * @param filter
 */
export async function listRole(page = 1, perPage = 20, sort = '-created', filter = '') {
    return await pb.collection('sys_role').getList(page, perPage, {
        sort,       // 排序字段，例如："-created", "email", "name"
        filter,     // PocketBase filter，例："email ~ 'gmail'"
        fields: 'id,roleName,permission,created,updated',
    });
}


/**
 * 创建菜单
 * @param data
 */
export async function createMenu(data) {
    return await pb.collection('sys_menu').create(data);
}

/**
 * 更新菜单
 * @param id
 * @param data
 */
export async function updateMenu(id, data) {
    return await pb.collection('sys_menu').update(id, data);
}

/**
 * 删除菜单
 * @param id
 */
export async function deleteMenu(id) {
    return await pb.collection('sys_menu').delete(id);
}


/**
 * 新增或更新角色
 * @param {object} role
 */
export async function saveRole(role) {
    if (role.id) {
        return await pb.collection('sys_role').update(role.id, role);
    } else {
        return await pb.collection('sys_role').create(role);
    }
}

/**
 * 删除角色
 * @param {string} id
 */
export async function deleteRole(id) {
    return await pb.collection('sys_role').delete(id);
}

/**
 * 更新角色权限
 * @param {string} id
 * @param {Array<string>} permissions
 */
export async function saveRolePermissions(id, permissions) {
    return await pb.collection('sys_role').update(id, {permission: permissions});
}


pb.autoCancellation(false);
// 导出 PocketBase 实例（可直接用）
export {pb};