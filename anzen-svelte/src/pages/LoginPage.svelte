<script>
    import {ArrowRight, Logo} from '../lib/icons';
    import {login} from "../api/sysApis.js";
    import {showToast} from '../stores/toastStore.js';
    import {setUser} from "../stores/userStore.js";
    import {push} from 'svelte-spa-router';

    let email = '';
    let password = '';
    let isLoading = false;

    const handleLogin = async () => {
        isLoading = true;
        try {
            const res = await login(email, password);
            setUser(res);
            showToast('login success！', 'success');
            setTimeout(() => {
                isLoading = false;
                push('/');
            }, 800);
        } catch (e) {
            console.log("login failed", e.message);
            showToast("login failed！", "error");
            isLoading = false;
        }
    };
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-100 p-6">
  <div class="bg-white rounded-xl p-10 w-full max-w-md shadow-lg">
    <div class="flex items-center justify-center gap-3 mb-10">
      <Logo size={48}/>
      <h1 class="text-3xl font-light text-gray-800">
        RUOYI-<span class="font-semibold">PB</span>
      </h1>
    </div>

    <h2 class="text-center text-xl font-medium text-gray-700 mb-8">Superuser login</h2>

    <form on:submit|preventDefault={handleLogin} class="flex flex-col gap-5">

      <div class="flex flex-col gap-1">
        <label for="email" class="text-sm font-medium text-gray-600">
          Email <span class="text-red-500">*</span>
        </label>
        <input
                id="email"
                type="email"
                bind:value={email}
                required
                class="w-full px-4 py-2.5 border border-gray-300 rounded-lg text-gray-800 bg-gray-50 focus:bg-white focus:border-blue-400 outline-none transition"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label for="password" class="text-sm font-medium text-gray-600">
          Password <span class="text-red-500">*</span>
        </label>
        <input
                id="password"
                type="password"
                bind:value={password}
                required
                class="w-full px-4 py-2.5 border border-gray-300 rounded-lg text-gray-800 bg-gray-50 focus:bg-white focus:border-blue-400 outline-none transition"
        />
      </div>

      <a href="#/forgot-password" class="text-sm text-gray-600 hover:text-gray-800 underline">
        Forgotten password?
      </a>

      <button type="submit"
              class="w-full py-3 bg-gray-800 text-white rounded-lg font-medium flex items-center justify-center gap-2 mt-2 hover:bg-gray-900 transition disabled:opacity-70 disabled:cursor-not-allowed"
              disabled={isLoading}
      >
        {isLoading ? 'Loading...' : 'Login'}
        {#if !isLoading}
          <ArrowRight size={16}/>
        {/if}
      </button>
    </form>
  </div>
</div>
