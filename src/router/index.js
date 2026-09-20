import { createRouter, createWebHistory } from 'vue-router';
import ParticipantApp from '../views/participant/ParticipantApp.vue';
import HostLoginView from '../views/host/HostLoginView.vue';
import HostLayout from '../views/host/HostLayout.vue';
import HostDashboardView from '../views/host/HostDashboardView.vue';
import HostQuizzesView from '../views/host/HostQuizzesView.vue';
import HostQuizEditorView from '../views/host/HostQuizEditorView.vue';
import HostSessionView from '../views/host/HostSessionView.vue';
import HostReportsView from '../views/host/HostReportsView.vue';
import HostSettingsView from '../views/host/HostSettingsView.vue';
import { useAuthStore } from '../stores/auth';

const routes = [
  {
    path: '/',
    name: 'participant',
    component: ParticipantApp
  },
  {
    path: '/join/:code?',
    name: 'participant-join',
    component: ParticipantApp
  },
  {
    path: '/host',
    name: 'host-login',
    component: HostLoginView
  },
  {
    path: '/host',
    component: HostLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'host-dashboard',
        component: HostDashboardView
      },
      {
        path: 'quizzes',
        name: 'host-quizzes',
        component: HostQuizzesView
      },
      {
        path: 'quiz/:id/edit',
        name: 'host-quiz-edit',
        component: HostQuizEditorView
      },
      {
        path: 'session/:code',
        name: 'host-session',
        component: HostSessionView
      },
      {
        path: 'reports',
        name: 'host-reports',
        component: HostReportsView
      },
      {
        path: 'settings',
        name: 'host-settings',
        component: HostSettingsView
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  if (to.meta.requiresAuth && !authStore.isAuthenticated()) {
    next('/host');
  } else {
    next();
  }
});

export default router;
