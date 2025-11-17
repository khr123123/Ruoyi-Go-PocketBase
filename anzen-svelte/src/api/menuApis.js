import {pb} from "./userApis.js";

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
