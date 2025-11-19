<!-- src/components/LogTable.svelte -->
<script>
    import {createEventDispatcher} from 'svelte';
    import {ChevronRight} from "../lib/icons/index.js";

    export let logs = [];
    export let totalLogs = 0;
    export let currentPage = 1;
    export let isLoading = false;

    const dispatch = createEventDispatcher();

    let expandedRows = new Set();

    function toggleRow(id) {
        expandedRows = new Set(
              expandedRows.has(id)
                    ? [...expandedRows].filter(x => x !== id)
                    : [...expandedRows, id]
        );
    }

    function getLogLevelBadge(level) {
        const badges = {
            '-4': {label: 'DEBUG', class: 'bg-blue-100 text-blue-700 border-blue-300'},
            '0': {label: 'INFO', class: 'bg-green-100 text-green-700 border-green-300'},
            '4': {label: 'WARN', class: 'bg-yellow-100 text-yellow-700 border-yellow-300'},
            '8': {label: 'ERROR', class: 'bg-red-100 text-red-700 border-red-300'}
        };
        return badges[level] || {label: 'UNKNOWN', class: 'bg-gray-100 text-gray-700 border-gray-300'};
    }

    function formatDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
            hour12: false
        }) + ' UTC';
    }

    function parseLogMessage(message) {
        const parts = {};

        const methodMatch = message.match(/(GET|POST|PUT|DELETE|PATCH)\s+([^\s]+)/);
        if (methodMatch) {
            parts.method = methodMatch[1];
            parts.endpoint = methodMatch[2];
        }

        const statusMatch = message.match(/status:\s*(\d+)/);
        if (statusMatch) parts.status = statusMatch[1];

        const timeMatch = message.match(/execTime:\s*([\d.]+)ms/);
        if (timeMatch) parts.execTime = timeMatch[1];

        const authMatch = message.match(/auth:\s*(\w+)/);
        if (authMatch) parts.auth = authMatch[1];

        const ipMatch = message.match(/userIp:\s*([\d.]+)/);
        if (ipMatch) parts.userIp = ipMatch[1];

        const refererMatch = message.match(/referer:\s*([^\s]+)/);
        if (refererMatch) parts.referer = refererMatch[1];

        return parts;
    }

    function nextPage() {
        dispatch('pageChange', {page: currentPage + 1});
    }

    function prevPage() {
        if (currentPage > 1) dispatch('pageChange', {page: currentPage - 1});
    }
</script>

<div class="bg-white border border-gray-300 rounded overflow-hidden">
    <div class="overflow-x-auto">
        <table class="w-full text-sm">
            <thead>
            <tr class="bg-gray-50 text-gray-700 text-xs font-medium border-b border-gray-200">
                <th class="px-4 py-3 text-left w-10"></th>
                <th class="px-4 py-3 text-left w-24">Level</th>
                <th class="px-4 py-3 text-left">Message</th>
                <th class="px-4 py-3 text-left w-52">Created</th>
            </tr>
            </thead>

            <tbody>
            {#if isLoading}
                <tr>
                    <td colspan="4" class="px-4 py-8 text-center text-gray-500">Loading...</td>
                </tr>
            {:else if logs.length === 0}
                <tr>
                    <td colspan="4" class="px-4 py-8 text-center text-gray-500">No logs found</td>
                </tr>
            {:else}
                {#each logs as log (log.id)}
                    {@const badge = getLogLevelBadge(log.level)}

                    <tr class="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                        <!-- Expand icon -->
                        <td class="px-4 py-3">
                            <button
                                  on:click={() => toggleRow(log.id)}
                                  class="text-gray-400 hover:text-gray-600 transition-colors"
                                  aria-expanded={expandedRows.has(log.id)}
                            >
                               <span class="transition-transform"
                                     class:rotate-90={expandedRows.has(log.id)}>
    <ChevronRight size={14}/>
</span>

                            </button>
                        </td>

                        <!-- Level badge -->
                        <td class="px-4 py-3">
                            <span class={`inline-flex items-center gap-2 px-2 py-0.5 rounded text-xs font-medium border ${badge.class}`}>
                                <span class="w-1.5 h-1.5 rounded-full bg-current"></span>
                                {badge.label}({log.level})
                            </span>
                        </td>

                        <!-- Message + expanded details -->
                        <td class="px-4 py-3">
                            <div class="font-mono text-xs text-gray-900 break-all">
                                {log.message}
                            </div>

                            {#if expandedRows.has(log.id)}
                                {@const parsed = parseLogMessage(log.message)}

                                <div class="mt-2 p-3 bg-gray-50 rounded border border-gray-200 space-y-1">
                                    {#if parsed.method}
                                        <div class="flex gap-2">
                                            <span class="text-gray-500 w-20">Method:</span>
                                            <span class="font-semibold">{parsed.method}</span>
                                        </div>
                                    {/if}

                                    {#if parsed.endpoint}
                                        <div class="flex gap-2">
                                            <span class="text-gray-500 w-20">Endpoint:</span>
                                            <span class="font-mono text-xs text-gray-800">{parsed.endpoint}</span>
                                        </div>
                                    {/if}

                                    {#if parsed.status}
                                        <div class="flex gap-2">
                                            <span class="text-gray-500 w-20">Status:</span>
                                            <span class={`px-2 py-0.5 rounded text-xs ${
                                                parsed.status === '200'
                                                    ? 'bg-green-100 text-green-700'
                                                    : 'bg-red-100 text-red-700'
                                            }`}>
                                                {parsed.status}
                                            </span>
                                        </div>
                                    {/if}

                                    {#if parsed.execTime}
                                        <div class="flex gap-2">
                                            <span class="text-gray-500 w-20">Exec Time:</span>
                                            <span>{parsed.execTime}ms</span>
                                        </div>
                                    {/if}

                                    {#if parsed.auth}
                                        <div class="flex gap-2">
                                            <span class="text-gray-500 w-20">Auth:</span>
                                            <span>{parsed.auth}</span>
                                        </div>
                                    {/if}

                                    {#if parsed.userIp}
                                        <div class="flex gap-2">
                                            <span class="text-gray-500 w-20">User IP:</span>
                                            <span class="font-mono text-xs">{parsed.userIp}</span>
                                        </div>
                                    {/if}

                                    {#if parsed.referer}
                                        <div class="flex gap-2">
                                            <span class="text-gray-500 w-20">Referer:</span>
                                            <span class="font-mono text-xs text-gray-800 truncate">
                                                {parsed.referer}
                                            </span>
                                        </div>
                                    {/if}
                                </div>
                            {/if}
                        </td>

                        <!-- Created time -->
                        <td class="px-4 py-3 text-gray-600 text-xs">
                            {formatDate(log.created)}
                        </td>
                    </tr>
                {/each}
            {/if}
            </tbody>
        </table>
    </div>

    <!-- Pagination -->
    <div class="px-5 py-3 border-t border-gray-200 bg-gray-50 flex items-center justify-between text-xs">
        <div class="text-gray-600">
            Total: <strong class="text-gray-900 font-semibold">{totalLogs}</strong>
        </div>

        <div class="flex items-center gap-2">
            <button
                  on:click={prevPage}
                  disabled={currentPage === 1}
                  class="px-3 py-1.5 border border-gray-300 rounded bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed text-gray-700"
            >
                ← Prev
            </button>

            <span class="px-2">
                Page <strong>{currentPage}</strong>
            </span>

            <button
                  on:click={nextPage}
                  class="px-3 py-1.5 border border-gray-300 rounded bg-white hover:bg-gray-50 text-gray-700"
            >
                Next →
            </button>
        </div>
    </div>
</div>
