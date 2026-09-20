<template>
  <div v-if="quiz">
    <div class="toolbar">
      <div>
        <h1>{{ quiz.title || 'Untitled Quiz' }}</h1>
        <div class="sub">Edit the question set before opening a live session.</div>
      </div>
      <div class="actions">
        <button class="btn" @click="router.push('/host/quizzes')">Cancel</button>
        <button class="btn primary" @click="handleSave">Save quiz</button>
      </div>
    </div>

    <div class="grid two">
      <div class="panel">
        <div class="panel-head">
          <h2>Questions</h2>
          <button class="btn small" @click="addQuestion">+ Add question</button>
        </div>
        <div class="panel-body">
          <div v-for="(q, i) in quiz.questions" :key="i" class="q-item">
            <div class="q-head">
              <div class="q-num">{{ i + 1 }}</div>
              <div class="q-text">{{ q.q || 'Untitled question' }}</div>
              <div class="q-controls">
                <button class="iconbtn" @click="moveQuestion(i, -1)">↑</button>
                <button class="iconbtn" @click="moveQuestion(i, 1)">↓</button>
                <button class="iconbtn" @click="removeQuestion(i)">×</button>
              </div>
            </div>
            <div class="q-body">
              <div class="field">
                <label>Question text</label>
                <input v-model="q.q" placeholder="Enter question..." />
              </div>
              <div v-for="(opt, j) in q.a" :key="j" class="option">
                <input
                  type="radio"
                  :name="`correct_${i}`"
                  class="correct"
                  :checked="q.c === j"
                  @change="q.c = j"
                />
                <input v-model="q.a[j]" :placeholder="`Answer option ${j + 1}`" />
                <button class="iconbtn" @click="removeOption(i, j)">×</button>
              </div>
              <button class="btn small" style="margin-top: 6px;" @click="addOption(i)">+ Add answer option</button>
            </div>
          </div>
        </div>
      </div>

      <div class="panel">
        <div class="panel-head">
          <h2>Quiz settings</h2>
        </div>
        <div class="panel-body">
          <div class="field">
            <label>Title</label>
            <input v-model="quiz.title" placeholder="Quiz Title" />
          </div>
          <div class="field">
            <label>Description</label>
            <textarea v-model="quiz.description" placeholder="Short summary of this quiz"></textarea>
          </div>
          <div class="field">
            <label>Time per question</label>
            <select v-model="quiz.timePerQuestion">
              <option :value="20">20 seconds</option>
              <option :value="30">30 seconds</option>
              <option :value="45">45 seconds</option>
            </select>
          </div>
          <div class="field">
            <label>Scoring</label>
            <select v-model="quiz.scoringMode">
              <option value="Speed & Accuracy">Correctness + response speed</option>
              <option value="Accuracy Only">Correctness only</option>
            </select>
          </div>
          <div class="notice">
            The host sees the full live leaderboard. Participants see their own score and current position.
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useQuizStore } from '../../stores/quiz';

const route = useRoute();
const router = useRouter();
const quizStore = useQuizStore();

const quiz = ref(null);

onMounted(() => {
  const targetId = route.params.id;
  const existing = quizStore.getQuizById(targetId);
  quiz.value = JSON.parse(JSON.stringify(existing));
});

const addQuestion = () => {
  quiz.value.questions.push({
    q: '',
    a: ['', '', '', ''],
    c: 0
  });
};

const removeQuestion = (index) => {
  if (quiz.value.questions.length > 1) {
    quiz.value.questions.splice(index, 1);
  }
};

const moveQuestion = (index, direction) => {
  const newIndex = index + direction;
  if (newIndex >= 0 && newIndex < quiz.value.questions.length) {
    const item = quiz.value.questions.splice(index, 1)[0];
    quiz.value.questions.splice(newIndex, 0, item);
  }
};

const addOption = (qIdx) => {
  quiz.value.questions[qIdx].a.push('');
};

const removeOption = (qIdx, optIdx) => {
  if (quiz.value.questions[qIdx].a.length > 2) {
    quiz.value.questions[qIdx].a.splice(optIdx, 1);
    if (quiz.value.questions[qIdx].c === optIdx) {
      quiz.value.questions[qIdx].c = 0;
    } else if (quiz.value.questions[qIdx].c > optIdx) {
      quiz.value.questions[qIdx].c--;
    }
  }
};

const handleSave = () => {
  quizStore.saveQuiz(quiz.value);
  router.push('/host/quizzes');
};
</script>
