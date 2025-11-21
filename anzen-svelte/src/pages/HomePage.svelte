<script>
    import {onMount} from "svelte";
    import {pb} from "../api/sysApis";
    import {Identifi, ShieldCheck, Users} from "../lib/icons/index.js";

    let loading = true;
    let stats = {users: 0, roles: 0, menus: 0};

    async function loadStats() {
        const users = await pb.collection("users").getList(1, 1);
        const roles = await pb.collection("sys_role").getList(1, 1);
        const menus = await pb.collection("sys_menu").getList(1, 1);

        stats.users = users.totalItems;
        stats.roles = roles.totalItems;
        stats.menus = menus.totalItems;
    }

    function subscribeRealtime() {
        pb.collection("users").subscribe("*", () => loadStats());
        pb.collection("sys_role").subscribe("*", () => loadStats());
        pb.collection("sys_menu").subscribe("*", () => loadStats());
    }

    onMount(async () => {
        await loadStats();
        loading = false;
        subscribeRealtime();
    });

    // --- Chart Data ---
    $: lineData = stats.users ? [stats.users, stats.users + 2, stats.users + 4, stats.users + 7, stats.users + 9] : [];
    $: roleData = stats.roles ? [stats.roles, stats.roles + 1, stats.roles + 1, stats.roles + 2, stats.roles + 3] : [];
    $: menuData = stats.menus ? [stats.menus, stats.menus + 1, stats.menus + 1, stats.menus + 2, stats.menus + 3] : [];

</script>

<div class="bg-[#F9FAFB] space-y-4">
    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">

        <!-- Users -->
        <div class="bg-white border border-gray-200 rounded-xl p-5 shadow hover:shadow-lg transition flex flex-col items-center gap-3 w-full">
            <Users size="{26}"/>
            <p class="text-sm text-gray-500">Total Users</p>
            {#if loading}
                <div class="animate-pulse h-7 w-20 bg-gray-200 rounded"></div>
            {:else}
                <p class="text-2xl font-bold text-gray-800">{stats.users}</p>
            {/if}
        </div>

        <!-- Roles -->
        <div class="bg-white border border-gray-200 rounded-xl p-5 shadow hover:shadow-lg transition flex flex-col items-center gap-3 w-full">
            <Identifi size="{26}"/>
            <p class="text-sm text-gray-500">Total Roles</p>
            {#if loading}
                <div class="animate-pulse h-7 w-20 bg-gray-200 rounded"></div>
            {:else}
                <p class="text-2xl font-bold text-gray-800">{stats.roles}</p>
            {/if}
        </div>

        <!-- Menus -->
        <div class="bg-white border border-gray-200 rounded-xl p-5 shadow hover:shadow-lg transition flex flex-col items-center gap-3 w-full">
            <ShieldCheck size="{26}"/>
            <p class="text-sm text-gray-500">Total Menus</p>
            {#if loading}
                <div class="animate-pulse h-7 w-20 bg-gray-200 rounded"></div>
            {:else}
                <p class="text-2xl font-bold text-gray-800">{stats.menus}</p>
            {/if}
        </div>
    </div>

    <!-- Growth Line Chart -->
    <div class="bg-white border border-gray-200 rounded-xl p-6 shadow">
        <p class="text-gray-600 text-sm mb-4">Growth Trend (Real-time)</p>

        {#if loading}
            <div class="animate-pulse h-58 bg-gray-200 rounded"></div>
        {:else}
            <div class="relative h-58 pl-10 pb-10 border-l border-b border-gray-300">
                <!-- Grid lines + Y-axis labels -->
                {#each [0, 25, 50, 75, 100, 125, 150, 175, 200] as y}
                    <div class="absolute left-0 w-full border-t border-gray-200 text-gray-400 text-xs -mt-1"
                         style="bottom: {y}px">{y}</div>
                {/each}

                <!-- Users Line -->
                <svg class="absolute inset-0 w-full h-full">
                    <polyline fill="none" stroke="#3B82F6" stroke-width="3" stroke-linecap="round"
                              points="{lineData.map((v,i)=>`${i*90},${160-v}`).join(' ')}"/>
                </svg>

                <!-- Roles Line -->
                <svg class="absolute inset-0 w-full h-full">
                    <polyline fill="none" stroke="#10B981" stroke-width="3" stroke-linecap="round"
                              points="{roleData.map((v,i)=>`${i*90},${160-v}`).join(' ')}"/>
                </svg>

                <!-- Menus Line -->
                <svg class="absolute inset-0 w-full h-full">
                    <polyline fill="none" stroke="#F59E0B" stroke-width="3" stroke-linecap="round"
                              points="{menuData.map((v,i)=>`${i*90},${160-v}`).join(' ')}"/>
                </svg>
            </div>

            <!-- Legend -->
            <div class="flex gap-8 mt-4 text-sm text-gray-600 justify-center">
                <div class="flex items-center gap-2"><span class="w-3 h-3 bg-blue-500 rounded-sm"></span> Users</div>
                <div class="flex items-center gap-2"><span class="w-3 h-3 bg-green-500 rounded-sm"></span> Roles</div>
                <div class="flex items-center gap-2"><span class="w-3 h-3 bg-yellow-500 rounded-sm"></span> Menus</div>
            </div>
        {/if}
    </div>

</div>
