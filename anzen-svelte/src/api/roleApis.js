import {pb} from "./userApis.js";

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