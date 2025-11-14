<script>
    import {ArrowRight, Logo} from '../lib/icons';
    import {login} from "../api/userApis.js";
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
            setUser(res); // 存储核心用户信息
            showToast('login success！', 'success');
            setTimeout(() => {
                isLoading = false
                push('/')
            }, 800);
        } catch (e) {
            console.log("login failed", e.message);
            showToast("login failed！", "error");
            isLoading = false;
        }

    }

</script>
<div class="login-container">
    <div class="login-box">
        <div class="logo-container">
            <Logo size={48}/>
            <h1 class="logo-text">RUOYI-<span class="logo-bold">PB</span></h1>
        </div>
        <h2 class="login-title">Superuser login</h2>
        <form on:submit|preventDefault={handleLogin} class="login-form">
            <div class="form-group">
                <label for="email" class="form-label">
                    Email <span class="required">*</span>
                </label>
                <input
                        id="email"
                        type="email"
                        bind:value={email}
                        required
                        class="form-input"
                />
            </div>
            <div class="form-group">
                <label for="password" class="form-label">
                    Password <span class="required">*</span>
                </label>
                <input
                        id="password"
                        type="password"
                        bind:value={password}
                        required
                        class="form-input"
                />
            </div>
            <a href="#/forgot-password" class="forgot-link">Forgotten password?</a>
            <button type="submit" class="login-button" disabled={isLoading}>
                {isLoading ? 'Loading...' : 'Login'}
                {#if !isLoading}
                    <ArrowRight size={16}/>
                {/if}
            </button>
        </form>
    </div>
</div>
<style>
    .login-container {
        min-height: 100vh;
        display: flex;
        align-items: center;
        justify-content: center;
        background-color: #f7f7f7;
        padding: 20px;
    }

    .login-box {
        background: white;
        border-radius: 8px;
        padding: 48px 40px;
        width: 100%;
        max-width: 440px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    }

    .logo-container {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 32px;
        gap: 12px;
    }

    .logo-text {
        font-size: 28px;
        font-weight: 400;
        color: #2b3034;
        margin: 0;
    }

    .logo-bold {
        font-weight: 700;
    }

    .login-title {
        text-align: center;
        font-size: 20px;
        font-weight: 500;
        color: #2b3034;
        margin: 0 0 32px 0;
    }

    .login-form {
        display: flex;
        flex-direction: column;
        gap: 20px;
    }

    .form-group {
        display: flex;
        flex-direction: column;
        gap: 6px;
    }

    .form-label {
        font-size: 13px;
        color: #666;
        font-weight: 500;
    }

    .required {
        color: #e03f3f;
    }

    .form-input {
        width: 100%;
        padding: 10px 14px;
        border: 1px solid #d9d9d9;
        border-radius: 4px;
        font-size: 15px;
        color: #2b3034;
        background-color: #f0f4f8;
        transition: all 0.2s;
    }

    .form-input:focus {
        outline: none;
        border-color: #5a9fd4;
        background-color: white;
    }

    .forgot-link {
        font-size: 13px;
        color: #666;
        text-decoration: none;
        align-self: flex-start;
        margin-top: -8px;
    }

    .forgot-link:hover {
        color: #2b3034;
        text-decoration: underline;
    }

    .login-button {
        width: 100%;
        padding: 12px 24px;
        background-color: #2b3034;
        color: white;
        border: none;
        border-radius: 4px;
        font-size: 15px;
        font-weight: 500;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        margin-top: 12px;
        transition: background-color 0.2s;
    }

    .login-button:hover:not(:disabled) {
        background-color: #1a1d1f;
    }

    .login-button:disabled {
        opacity: 0.7;
        cursor: not-allowed;
    }
</style>