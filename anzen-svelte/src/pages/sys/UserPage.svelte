<!-- src/pages/UserPage.svelte -->
<script>
    import {onMount} from 'svelte';
    import {deleteUser, listRole, listUser, saveUser} from '../../api/sysApis.js';
    import {ArrowDown, Check, Circle, Image, Menu, Upload, X} from '../../lib/icons/index.js';
    import DataTable from "../../components/DataTable.svelte";
    import Drawer from "../../components/Drawer.svelte";
    import {showToast} from '../../stores/toastStore.js';
    import {confirmDialog} from "../../stores/confirmStore.js";
    import {permission} from "../../utils/permissionDirective.js";

    const AVATAR_PREFIX = "http://127.0.0.1:8090/api/files/_pb_users_auth_/";

    let users = [];
    let total = 0;
    let page = 1;
    let perPage = 10;

    let showDialog = false;
    let currentUser = null;

    let allRoles = [];
    let roleSelectValue = '';

    // 头像相关
    let avatarFile = null;
    let avatarPreview = null;
    let avatarInput; // 文件输入框引用

    $: selectedRoles = currentUser?.role?.map(id => allRoles.find(r => r.id === id)).filter(Boolean) || [];

    // 当前用户头像 URL
    $: currentAvatarUrl = currentUser?.id && currentUser?.avatar
          ? `${AVATAR_PREFIX}${currentUser.id}/${currentUser.avatar}`
          : null;

    const columns = [
        {key: 'id', label: 'ID', icon: ArrowDown, sortable: true, class: "font-mono"},
        {key: 'email', label: 'Email', icon: Menu, sortable: true},
        {
            key: 'emailVisibility',
            label: 'Email Visible',
            icon: Circle,
            render: (v) => `<span class="px-2 py-0.5 rounded-full text-xs ${v ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'}">${v ? 'Yes' : 'No'}</span>`
        },
        {
            key: 'verified',
            label: 'Verified',
            icon: Check,
            render: (v) => `<span class="px-2 py-0.5 rounded-full text-xs ${v ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'}">${v ? 'Yes' : 'No'}</span>`
        },
        {key: 'name', label: 'Name', sortable: true},
        {
            key: 'avatar',
            label: 'Avatar',
            icon: Image,
            render: (value, row) => {
                if (!value) {
                    return `<div class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center text-white text-sm font-semibold">${row.name?.charAt(0)?.toUpperCase() || row.email?.charAt(0)?.toUpperCase() || 'U'}</div>`;
                }
                return `<img class="w-8 h-8 rounded-full object-cover" src="${AVATAR_PREFIX}${row.id}/${value}" alt="${row.name || 'User'}" onerror="this.style.display='none';this.nextElementSibling.style.display='flex'"/><div style="display:none" class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center text-white text-sm font-semibold">${row.name?.charAt(0)?.toUpperCase() || row.email?.charAt(0)?.toUpperCase() || 'U'}</div>`;
            }
        },
        {
            key: 'created',
            label: 'Created',
            class: 'text-xs text-gray-600',
            render: (v) => new Date(v).toLocaleString('zh-CN')
        },
        {
            key: 'updated',
            label: 'Updated',
            class: 'text-xs text-gray-600',
            render: (v) => new Date(v).toLocaleString('zh-CN')
        }
    ];

    async function loadUsers(params = {}) {
        const {search = '', sort = '-created', page: p = 1} = params;
        const filter = search ? `email ~ "${search}" || name ~ "${search}"` : '';
        try {
            const res = await listUser(p, perPage, sort, filter);
            users = res.items || [];
            total = res.totalItems || 0;
            page = p;
        } catch (error) {
            console.error('Load users error:', error);
            showToast('加载用户列表失败: ' + (error?.message || '未知错误'), 'error');
        }
    }

    function handleSearch(event) {
        loadUsers({search: event.detail.search, page: 1});
    }

    function handleSort(event) {
        loadUsers({sort: event.detail.sort, page: event.detail.page});
    }

    function handlePageChange(event) {
        loadUsers({page: event.detail.page});
    }

    function handleAdd() {
        currentUser = {
            email: '',
            name: '',
            password: '',
            passwordConfirm: '',
            emailVisibility: true,
            verified: false,
            avatar: '',
            role: []
        };
        roleSelectValue = '';
        resetAvatarUpload();
        showDialog = true;
    }

    function handleEdit(event) {
        const user = event.detail;
        currentUser = {
            ...user,
            password: '',
            passwordConfirm: '',
            role: user.role || []
        };
        roleSelectValue = '';
        resetAvatarUpload();
        showDialog = true;
    }

    // 重置头像上传状态
    function resetAvatarUpload() {
        avatarFile = null;
        avatarPreview = null;
        if (avatarInput) {
            avatarInput.value = '';
        }
    }

    // 处理头像文件选择
    function handleAvatarSelect(event) {
        const file = event.target.files?.[0];
        if (!file) return;

        // 验证文件类型
        const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
        if (!validTypes.includes(file.type)) {
            showToast('请选择图片文件（JPG、PNG、GIF 或 WebP）', 'error');
            event.target.value = '';
            return;
        }

        // 验证文件大小（最大 5MB）
        const maxSize = 5 * 1024 * 1024;
        if (file.size > maxSize) {
            showToast('图片大小不能超过 5MB', 'error');
            event.target.value = '';
            return;
        }

        avatarFile = file;

        // 生成预览
        const reader = new FileReader();
        reader.onload = (e) => {
            avatarPreview = e.target.result;
        };
        reader.readAsDataURL(file);
    }

    // 移除已选择的头像
    function removeAvatar() {
        resetAvatarUpload();
    }

    // 删除现有头像
    function deleteExistingAvatar() {
        if (currentUser) {
            currentUser.avatar = '';
        }
    }

    async function saveUserHandler() {
        // 验证邮箱
        if (!currentUser.email || !currentUser.email.trim()) {
            showToast('邮箱不能为空', 'error');
            return;
        }

        // 验证邮箱格式
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(currentUser.email)) {
            showToast('邮箱格式不正确', 'error');
            return;
        }

        try {
            const submitData = {
                email: currentUser.email.trim(),
                name: currentUser.name?.trim() || '',
                emailVisibility: !!currentUser.emailVisibility,
                role: currentUser.role || []
            };

            if (currentUser.id) {
                // 编辑用户
                submitData.id = currentUser.id;

                // 只有填写了密码才更新密码
                if (currentUser.password && currentUser.password.trim()) {
                    if (currentUser.password !== currentUser.passwordConfirm) {
                        showToast('两次密码输入不一致', 'error');
                        return;
                    }
                    if (currentUser.password.length < 8) {
                        showToast('密码长度至少为8个字符', 'error');
                        return;
                    }
                    submitData.password = currentUser.password;
                    submitData.passwordConfirm = currentUser.passwordConfirm;
                }

                submitData.verified = !!currentUser.verified;

            } else {
                // 新增用户
                if (!currentUser.password || !currentUser.password.trim()) {
                    showToast('密码不能为空', 'error');
                    return;
                }
                if (currentUser.password !== currentUser.passwordConfirm) {
                    showToast('两次密码输入不一致', 'error');
                    return;
                }
                if (currentUser.password.length < 8) {
                    showToast('密码长度至少为8个字符', 'error');
                    return;
                }

                submitData.password = currentUser.password;
                submitData.passwordConfirm = currentUser.passwordConfirm;
            }

            // 处理头像上传
            if (avatarFile) {
                submitData.avatar = avatarFile;
            } else if (currentUser.avatar === '') {
                // 如果用户删除了头像，传递 null 或空字符串
                submitData.avatar = null;
            }

            console.log('提交数据:', submitData);

            await saveUser(submitData);
            showToast('保存成功', 'success');
            showDialog = false;
            resetAvatarUpload();
            await loadUsers({page});
        } catch (error) {
            console.error('Save user error:', error);

            let errorMsg = '操作失败';
            if (error?.data) {
                const errors = [];
                for (const [field, detail] of Object.entries(error.data)) {
                    if (typeof detail === 'object' && detail.message) {
                        errors.push(`${field}: ${detail.message}`);
                    }
                }
                if (errors.length > 0) {
                    errorMsg = errors.join('; ');
                } else if (error.message) {
                    errorMsg = error.message;
                }
            } else if (error?.message) {
                errorMsg = error.message;
            }

            showToast(errorMsg, 'error');
        }
    }

    async function handleDelete(event) {
        const user = event.detail;
        const ok = await confirmDialog(`确定要删除用户 "${user.name || user.email}" 吗？`);
        if (!ok) return;

        try {
            await deleteUser(user.id);
            showToast('删除成功', 'success');
            if (users.length === 1 && page > 1) {
                await loadUsers({page: page - 1});
            } else {
                await loadUsers({page});
            }
        } catch (error) {
            console.error('Delete user error:', error);
            showToast('删除失败: ' + (error?.message || '未知错误'), 'error');
        }
    }

    function removeRole(roleId) {
        currentUser.role = currentUser.role.filter(id => id !== roleId);
    }

    function addRole(roleId) {
        if (roleId && !currentUser.role.includes(roleId)) {
            currentUser.role = [...currentUser.role, roleId];
            roleSelectValue = '';
        }
    }

    onMount(async () => {
        try {
            await loadUsers();
            const res = await listRole(1, 500);
            allRoles = res.items || [];
        } catch (error) {
            console.error('初始化失败:', error);
            showToast('初始化失败', 'error');
        }
    });
</script>

<!-- 用户表格 -->
<DataTable
      data={users}
      {total}
      {columns}
      {page}
      searchPlaceholder="搜索邮箱或姓名..."
      addButtonText="添加用户"
      actions={{ add: "sys:user:add", edit: "sys:user:edit", delete: "sys:user:delete" }}
      on:add={handleAdd}
      on:edit={handleEdit}
      on:delete={handleDelete}
      on:search={handleSearch}
      on:sort={handleSort}
      on:pageChange={handlePageChange}
/>

<!-- 用户编辑/新增抽屉 -->
<Drawer show={showDialog} title={currentUser?.id ? '编辑用户' : '添加用户'} position="right"
        on:close={() => showDialog = false} width="550px">
    <form on:submit|preventDefault={saveUserHandler} class="space-y-4">

        <!-- 头像上传区域 -->
        <div>
            <label class="block text-sm font-medium mb-2">头像</label>
            <div class="flex items-start gap-4">
                <!-- 当前头像预览 -->
                <div class="flex-shrink-0">
                    {#if avatarPreview}
                        <!-- 新选择的头像预览 -->
                        <div class="relative group">
                            <img
                                  src={avatarPreview}
                                  alt="Avatar Preview"
                                  class="w-20 h-20 rounded-full object-cover border-2 border-blue-500"
                            />
                            <button
                                  type="button"
                                  on:click={removeAvatar}
                                  class="absolute -top-2 -right-2 bg-red-500 text-white rounded-full p-1 opacity-0 group-hover:opacity-100 transition-opacity"
                                  title="移除"
                            >
                                <X class="w-4 h-4"/>
                            </button>
                        </div>
                    {:else if currentAvatarUrl && currentUser.avatar}
                        <!-- 现有头像 -->
                        <div class="relative group">
                            <img
                                  src={currentAvatarUrl}
                                  alt="Current Avatar"
                                  class="w-20 h-20 rounded-full object-cover border-2 border-gray-300"
                            />
                            <button
                                  type="button"
                                  on:click={deleteExistingAvatar}
                                  class="absolute -top-2 -right-2 bg-red-500 text-white rounded-full p-1 opacity-0 group-hover:opacity-100 transition-opacity"
                                  title="删除头像"
                            >
                                <X class="w-4 h-4"/>
                            </button>
                        </div>
                    {:else}
                        <!-- 默认头像占位符 -->
                        <div class="w-20 h-20 rounded-full bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center text-white text-2xl font-semibold">
                            {currentUser?.name?.charAt(0)?.toUpperCase() || currentUser?.email?.charAt(0)?.toUpperCase() || 'U'}
                        </div>
                    {/if}
                </div>

                <!-- 上传按钮 -->
                <div class="flex-1">
                    <input
                          type="file"
                          bind:this={avatarInput}
                          on:change={handleAvatarSelect}
                          accept="image/jpeg,image/png,image/gif,image/webp"
                          class="hidden"
                          id="avatar-upload"
                    />
                    <label
                          for="avatar-upload"
                          class="inline-flex items-center gap-2 px-4 py-2 border border-gray-300 rounded-lg cursor-pointer hover:bg-gray-50 transition"
                    >
                        <Upload class="w-4 h-4"/>
                        <span>选择图片</span>
                    </label>
                    <p class="text-xs text-gray-500 mt-2">
                        支持 JPG、PNG、GIF、WebP 格式<br/>
                        最大 5MB
                    </p>
                </div>
            </div>
        </div>

        <!-- Email -->
        <div>
            <label class="block text-sm font-medium mb-1">
                邮箱 <span class="text-red-500">*</span>
            </label>
            <input
                  type="email"
                  bind:value={currentUser.email}
                  placeholder="请输入邮箱"
                  required
                  class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <p class="text-xs text-gray-500 mt-1">用于登录和接收通知</p>
        </div>

        <!-- Name -->
        <div>
            <label class="block text-sm font-medium mb-1">姓名</label>
            <input
                  type="text"
                  bind:value={currentUser.name}
                  placeholder="请输入姓名"
                  class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
        </div>

        <!-- Password -->
        <div class="space-y-3">
            <div>
                <label class="block text-sm font-medium mb-1">
                    密码
                    {#if !currentUser.id}
                        <span class="text-red-500">*</span>
                    {:else}
                        <span class="text-gray-500 text-xs">(留空表示不修改)</span>
                    {/if}
                </label>
                <input
                      type="password"
                      bind:value={currentUser.password}
                      placeholder={currentUser.id ? '留空表示不修改' : '至少8个字符'}
                      minlength="8"
                      class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
            </div>
            <div>
                <label class="block text-sm font-medium mb-1">
                    确认密码
                    {#if !currentUser.id}<span class="text-red-500">*</span>{/if}
                </label>
                <input
                      type="password"
                      bind:value={currentUser.passwordConfirm}
                      placeholder="请再次输入密码"
                      minlength="8"
                      class="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
            </div>
            {#if !currentUser.id}
                <p class="text-xs text-gray-500">密码至少8个字符</p>
            {/if}
        </div>

        <!-- Switches -->
        <div class="space-y-2">
            <label class="flex items-center gap-2 cursor-pointer">
                <input
                      type="checkbox"
                      bind:checked={currentUser.emailVisibility}
                      class="rounded border-gray-300 cursor-pointer"
                />
                <span class="text-sm">邮箱可见</span>
                <span class="text-xs text-gray-500">（其他用户是否可以看到此邮箱）</span>
            </label>

            {#if currentUser.id}
                <label class="flex items-center gap-2 cursor-pointer">
                    <input
                          type="checkbox"
                          bind:checked={currentUser.verified}
                          class="rounded border-gray-300 cursor-pointer"
                    />
                    <span class="text-sm">已验证</span>
                    <span class="text-xs text-gray-500">（邮箱验证状态）</span>
                </label>
            {:else}
                <p class="text-xs text-gray-500">新用户创建后需要验证邮箱</p>
            {/if}
        </div>

        <!-- Roles -->
        <div use:permission={"sys:user:role"} class="space-y-1">
            <label class="block text-sm font-medium mb-1">角色</label>
            <div class="border rounded px-3 py-2 min-h-[100px] focus-within:ring-2 focus-within:ring-blue-500">
                {#if selectedRoles.length > 0}
                    <div class="flex flex-wrap gap-1 mb-2">
                        {#each selectedRoles as r (r.id)}
                            <span class="flex items-center bg-blue-100 text-blue-700 text-xs px-2 py-1 rounded">
                                {r.roleName}
                                <button
                                      type="button"
                                      class="ml-1 hover:text-blue-900 font-bold"
                                      on:click={() => removeRole(r.id)}
                                >×</button>
                            </span>
                        {/each}
                    </div>
                {/if}

                {#if allRoles.length > 0}
                    <select
                          bind:value={roleSelectValue}
                          on:change={(e) => addRole(e.target.value)}
                          class="w-full border rounded px-2 py-1 focus:ring-0 focus:outline-none text-sm"
                    >
                        <option value="">选择角色...</option>
                        {#each allRoles as role (role.id)}
                            {#if !currentUser.role.includes(role.id)}
                                <option value={role.id}>{role.roleName}</option>
                            {/if}
                        {/each}
                    </select>
                {:else}
                    <p class="text-sm text-gray-500">暂无可用角色</p>
                {/if}
            </div>
        </div>

        <!-- 底部按钮 -->
        <div class="flex justify-end gap-2 mt-6 pt-4 border-t">
            <button
                  type="button"
                  on:click={() => { showDialog = false; resetAvatarUpload(); }}
                  class="px-4 py-2 border rounded hover:bg-gray-50 transition"
            >
                取消
            </button>
            <button
                  type="submit"
                  class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
            >
                保存
            </button>
        </div>
    </form>
</Drawer>