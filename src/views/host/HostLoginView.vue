<template>
  <div class="center" style="min-height: 100vh; display: flex; align-items: center; justify-content: center; padding: 20px;">
    <div class="join-box card" style="width: 100%; max-width: 440px;">
      <div class="card-body" style="padding: 30px;">
        <div style="text-align: center; margin-bottom: 24px;">
          <div class="logo" style="font-size: 26px; font-weight: 800; margin-bottom: 6px;">ELSA<span>•</span></div>
          <h1 style="font-size: 22px; margin: 0;">Host Workspace Sign In</h1>
          <p class="muted" style="font-size: 13px; margin-top: 4px;">Sign in to create quizzes and manage live sessions.</p>
        </div>

        <form @submit.prevent="handleLogin">
          <div class="field">
            <label>Email Address</label>
            <input v-model="email" type="email" required placeholder="admin@demo.com" />
          </div>
          <div class="field">
            <label>Password</label>
            <input v-model="password" type="password" required placeholder="••••••••" />
          </div>

          <div v-if="errorMessage" class="notice" style="border-color: #f1c1c1; background: #fff2f2; color: #c43d3d; margin-bottom: 14px;">
            {{ errorMessage }}
          </div>

          <button type="submit" class="btn primary full" style="padding: 11px;">Sign in as Host</button>
        </form>

        <div style="margin-top: 16px; text-align: center;">
          <button class="btn full" style="background: #f7f6ff; border-color: #ddd8ff; color: #5144c9;" @click="quickDemoLogin">
            ⚡ Quick Demo Sign In (admin@demo.com)
          </button>
        </div>

        <div style="margin-top: 20px; text-align: center; font-size: 12px;" class="muted">
          Looking to join a live quiz? <router-link to="/" style="color: var(--brand); font-weight: 600;">Join as Participant</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';

const router = useRouter();
const authStore = useAuthStore();

const email = ref('admin@demo.com');
const password = ref('admin123');
const errorMessage = ref('');

const handleLogin = () => {
  if (authStore.login(email.value, password.value)) {
    router.push('/host/dashboard');
  } else {
    errorMessage.value = 'Invalid login credentials';
  }
};

const quickDemoLogin = () => {
  email.value = 'admin@demo.com';
  password.value = 'admin123';
  handleLogin();
};
</script>
