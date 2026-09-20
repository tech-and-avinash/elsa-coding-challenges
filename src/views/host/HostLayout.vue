<template>
  <div class="host-app">
    <aside class="sidebar">
      <div class="brand" @click="router.push('/host/dashboard')">ELSA<span>•</span></div>
      <div class="workspace">Workspace</div>
      <nav class="nav">
        <router-link to="/host/dashboard">
          <span class="ico">⌂</span>Dashboard
        </router-link>
        <router-link to="/host/quizzes">
          <span class="ico">▤</span>Quizzes
        </router-link>
        <router-link to="/host/reports">
          <span class="ico">◉</span>Live sessions
        </router-link>
        <router-link to="/host/settings">
          <span class="ico">⚙</span>Settings
        </router-link>
      </nav>
      <div class="side-bottom">Live Quiz workspace<br>v1.0</div>
    </aside>

    <main class="main">
      <header class="topbar">
        <div class="top-title">{{ pageTitle }}</div>
        <div class="top-actions">
          <button class="btn small primary" @click="createNewQuiz">+ New quiz</button>
          <div class="avatar" style="cursor: pointer" title="Logged in as admin@demo.com" @click="handleLogout">
            {{ authStore.user?.initials || 'EL' }}
          </div>
        </div>
      </header>

      <div class="content">
        <div class="crumbs">
          <router-link to="/host/dashboard">Workspace</router-link>
          <span>›</span>
          <span class="current">{{ pageTitle }}</span>
        </div>

        <router-view />

        <div class="footer">ELSA Live Quiz · Workspace</div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import { useQuizStore } from '../../stores/quiz';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const quizStore = useQuizStore();

const pageTitle = computed(() => {
  switch (route.name) {
    case 'host-dashboard': return 'Dashboard';
    case 'host-quizzes': return 'Quizzes';
    case 'host-quiz-edit': return 'Quiz Editor';
    case 'host-session': return 'Live Session';
    case 'host-settings': return 'Settings';
    default: return 'Workspace';
  }
});

const createNewQuiz = () => {
  const newQuiz = quizStore.createQuiz();
  router.push(`/host/quiz/${newQuiz.id}/edit`);
};

const handleLogout = () => {
  if (confirm('Sign out of host workspace?')) {
    authStore.logout();
    router.push('/host');
  }
};
</script>
