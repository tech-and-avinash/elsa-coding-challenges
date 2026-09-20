<template>
  <div>
    <div class="toolbar">
      <div>
        <h1>Live sessions</h1>
        <div class="sub">Review completed session results and open each result in detail.</div>
      </div>
    </div>

    <div class="panel">
      <table class="table">
        <thead>
          <tr>
            <th>Session</th>
            <th>Quiz</th>
            <th>Status</th>
            <th>Participants</th>
            <th>Updated</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="session in quizStore.sessions" :key="session.code">
            <td><strong>{{ session.code }}</strong></td>
            <td>{{ session.title }}</td>
            <td><span class="badge" :class="session.status?.toLowerCase()">{{ session.status }}</span></td>
            <td>{{ session.participantCount ?? 0 }}</td>
            <td>{{ session.updatedAt || 'Recently' }}</td>
            <td>
              <router-link :to="`/host/session/${session.code}`" class="btn small">View result</router-link>
            </td>
          </tr>
          <tr v-if="!quizStore.sessions.length">
            <td colspan="6" class="muted">No sessions found.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import { useQuizStore } from '../../stores/quiz';

const quizStore = useQuizStore();

onMounted(async () => {
  await quizStore.fetchSessions();
});
</script>
