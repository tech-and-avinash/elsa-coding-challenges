<template>
  <div>
    <div class="toolbar">
      <div>
        <h1>Dashboard</h1>
        <div class="sub">Create quizzes, open live sessions, and review results.</div>
      </div>
      <div class="actions">
        <button class="btn primary" @click="handleCreateQuiz">Create quiz</button>
      </div>
    </div>

    <div class="grid three">
      <div class="panel stat">
        <div class="label">Quizzes</div>
        <div class="value">{{ quizStore.quizzes.length }}</div>
        <div class="meta">3 updated this week</div>
      </div>
      <div class="panel stat">
        <div class="label">Live sessions</div>
        <div class="value">2</div>
        <div class="meta">1 currently running</div>
      </div>
      <div class="panel stat">
        <div class="label">Participants</div>
        <div class="value">184</div>
        <div class="meta">Across recent sessions</div>
      </div>
    </div>

    <div style="height: 14px"></div>

    <div class="panel">
      <div class="panel-head">
        <h2>Recent sessions</h2>
        <router-link to="/host/reports" class="btn small">View all</router-link>
      </div>
      <table class="table">
        <thead>
          <tr>
            <th>Quiz</th>
            <th>Status</th>
            <th>Participants</th>
            <th>Updated</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>Business English Basics</td>
            <td><span class="badge live">Live</span></td>
            <td>18</td>
            <td>Just now</td>
            <td>
              <router-link to="/host/session/4821" class="btn small primary">Open</router-link>
            </td>
          </tr>
          <tr>
            <td>Travel Vocabulary</td>
            <td><span class="badge done">Completed</span></td>
            <td>32</td>
            <td>Yesterday</td>
            <td>
              <router-link to="/host/session/3910" class="btn small">View result</router-link>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';
import { useQuizStore } from '../../stores/quiz';

const router = useRouter();
const quizStore = useQuizStore();

const handleCreateQuiz = () => {
  const newQ = quizStore.createQuiz();
  router.push(`/host/quiz/${newQ.id}/edit`);
};
</script>
