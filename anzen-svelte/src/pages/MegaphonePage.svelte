<script>
    import {onDestroy, onMount} from 'svelte';
    import {pb} from '../api/sysApis';

    let messages = [];
    let newText = "";

    // Load initial data
    async function loadMessages() {
        messages = await pb.collection("message").getFullList({
            sort: "-created"
        });
    }

    // Add message
    async function addMessage() {
        if (!newText.trim()) return;
        await pb.collection("message").create({
            msg: newText
        });
        newText = "";
    }

    onMount(async () => {
        // 1. Optional: user login
        // await pb.collection('users').authWithPassword("xx@xx.com", "123456");
        await loadMessages();
        // 2. Subscribe to realtime events (* = all records)
        pb.collection("message").subscribe("*", e => {
            console.log("Realtime event:", e);
            if (e.action === "create") {
                messages = [e.record, ...messages];
            }
            if (e.action === "update") {
                messages = messages.map(m => m.id === e.record.id ? e.record : m);
            }
            if (e.action === "delete") {
                messages = messages.filter(m => m.id !== e.record.id);
            }
        });
    });
    // Cleanup
    onDestroy(() => {
        pb.collection("messages").unsubscribe("*");
    });
</script>

<!-- UI -->
<div class="p-6 max-w-xl mx-auto space-y-6">
    <h1 class="text-xl font-bold">PocketBase Live Messages (SSE Realtime)</h1>
    <!-- Input box -->
    <div class="flex gap-2">
        <input class="border rounded px-3 py-2 w-full" bind:value={newText} placeholder="Type something…"/>
        <button class="bg-blue-600 text-white px-4 rounded" on:click={addMessage}>
            Send
        </button>
    </div>
    <!-- Realtime messages -->
    <div class="space-y-3">
        {#each messages as msg (msg.id)}
            <div class="p-3 border rounded bg-gray-50">
                <div class="font-medium">{msg.msg}</div>
                <div class="text-xs text-gray-500">{msg.created}</div>
            </div>
        {/each}
    </div>
</div>
