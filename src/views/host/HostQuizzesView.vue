<template>
  <div>
    <div class="toolbar">
      <div>
        <h1>Quizzes</h1>
        <div class="sub">Manage question sets and launch them into a live session.</div>
      </div>
      <div class="actions">
        <button class="btn primary" @click="handleCreateQuiz">+ New quiz</button>
      </div>
    </div>

    <div class="panel">
      <table class="table">
        <thead>
          <tr>
            <th>Quiz</th>
            <th>Questions</th>
            <th>Sessions</th>
            <th>Last updated</th>
            <th>Status</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="q in quizStore.quizzes" :key="q.id">
            <td>
              <strong>{{ q.title }}</strong>
              <div class="muted">{{ q.category }}</div>
            </td>
            <td>{{ q.questions.length }}</td>
            <td>
              <button class="btn small" @click="toggleSessions(q.id)">
                {{ (quizStore.quizSessions[q.id] || []).length }} sessions
              </button>
            </td>
            <td>{{ q.updatedAt }}</td>
            <td><span class="badge ready">Ready</span></td>
            <td>
              <div class="actions">
                <router-link :to="`/host/quiz/${q.id}/edit`" class="btn small">Edit</router-link>
                <button class="btn small primary" @click="launchSession(q.id)">Host live</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Session History Panel -->
    <div v-if="expandedQuizId" class="panel" style="margin-top: 16px;">
      <div class="panel-head">
        <h2>Session History</h2>
        <span class="muted">{{ quizStore.quizzes.find(q => q.id === expandedQuizId)?.title }}</span>
      </div>
      <div class="panel-body">
        <table class="table">
          <thead>
            <tr>
              <th>Session Code</th>
              <th>Status</th>
              <th>Participants</th>
              <th>Last Updated</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="session in quizStore.quizSessions[expandedQuizId] || []" :key="session.code">
              <td><strong>{{ session.code }}</strong></td>
              <td><span class="badge" :class="session.status.toLowerCase()">{{ session.status }}</span></td>
              <td>{{ session.participantCount }}</td>
              <td>{{ session.updatedAt }}</td>
              <td>
                <router-link :to="`/host/session/${session.code}`" class="btn small">View</router-link>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-if="!quizStore.quizSessions[expandedQuizId] || quizStore.quizSessions[expandedQuizId].length === 0" class="muted">
          No sessions found for this quiz.
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useQuizStore } from '../../stores/quiz';
import { useLiveSessionStore } from '../../stores/liveSession';

const router = useRouter();
const quizStore = useQuizStore();
const sessionStore = useLiveSessionStore();

const expandedQuizId = ref(null);

const handleCreateQuiz = () => {
  const newQ = quizStore.createQuiz();
  router.push(`/host/quiz/${newQ.id}/edit`);
};

const launchSession = (quizId) => {
  const oldCode = sessionStore.code;

  // Listen for the custom session-created event
  const handleSessionCreated = (event) => {
    if (event.detail.code && event.detail.code !== oldCode) {
      window.removeEventListener('session-created', handleSessionCreated);
      router.push(`/host/session/${event.detail.code}`);
    }
  };

  window.addEventListener('session-created', handleSessionCreated);

  // Fallback timeout in case WebSocket event doesn't fire
  setTimeout(() => {
    window.removeEventListener('session-created', handleSessionCreated);
    if (sessionStore.code && sessionStore.code !== oldCode) {
      router.push(`/host/session/${sessionStore.code}`);
    }
  }, 2000);

  sessionStore.createSession(quizId);
};

const toggleSessions = async (quizId) => {
  if (expandedQuizId.value === quizId) {
    expandedQuizId.value = null;
  } else {
    expandedQuizId.value = quizId;
    await quizStore.fetchQuizSessions(quizId);
  }
};
</script>
