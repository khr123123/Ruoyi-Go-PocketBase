import {pb} from "./userApis.js";

export async function listMenu(page = 1, perPage = 20, sort = '-created', filter = '') {
    return await pb.collection('sys_menu').getList(page, perPage, {
        sort,
        filter,
        fields: 'id,menuName,menuType,orderNum,permission,parentId,url,icon,created,updated',
    });
}


export async function getUserRouter() {
    return await pb.send('/api/getUserRouter')
}
