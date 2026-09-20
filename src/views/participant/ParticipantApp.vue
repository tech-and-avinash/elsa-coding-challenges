<template>
  <div class="participant-app">
    <header class="participant-header">
      <div class="logo">ELSA<span>•</span></div>
      <div class="participant-session">
        Live session <strong>{{ sessionStore.code }}</strong>
      </div>
      <div class="user">
        <span>{{ displayName }}</span>
        <span class="avatar">{{ displayName.charAt(0).toUpperCase() }}</span>
      </div>
    </header>

    <main class="participant-main">
      <!-- 1. Join View -->
      <section v-if="currentStage === 'join'" class="center">
        <div class="join-box card">
          <div class="card-body">
            <h1>Join live quiz</h1>
            <p class="muted">Enter the session code and the name you want other participants to see.</p>

            <form @submit.prevent="handleJoin">
              <div class="field">
                <label>Session code</label>
                <input v-model="inputCode" placeholder="e.g. 4821" required />
              </div>
              <div class="field">
                <label>Your name</label>
                <input v-model="inputName" placeholder="e.g. Alex" required />
              </div>
              <button type="submit" class="btn primary full">Join session</button>
            </form>
          </div>
        </div>
      </section>

      <!-- 2. Waiting Room View -->
      <section v-else-if="currentStage === 'waiting'" class="center">
        <div class="wait">
          <div class="pulse">
            <span class="dot"></span> Connected
          </div>
          <h1>You're in</h1>
          <p class="muted">{{ sessionStore.sessionTitle }}</p>
          <div class="code">{{ sessionStore.code }}</div>
          <div class="card">
            <div class="card-body">
              <strong>{{ displayName }}</strong>
              <p class="muted">Waiting for the host to start the quiz.</p>
              <div class="notice">
                Keep this page open. The first question will appear automatically when the session starts.
              </div>
            </div>
          </div>
          <div style="height: 14px"></div>
          <button class="btn" @click="currentStage = 'join'">Change details</button>
        </div>
      </section>

      <!-- 3. Live Question View -->
      <section v-else-if="currentStage === 'live'">
        <div class="progress">
          <span :style="{ width: progressPercent + '%' }"></span>
        </div>
        <div class="participant-layout">
          <div class="card">
            <div class="card-body">
              <div class="question-meta">
                <span>Question <strong>{{ sessionStore.currentQuestionIndex + 1 }}</strong> of <strong>{{ sessionStore.questions.length }}</strong></span>
                <span class="timer">{{ timerCount }} sec</span>
              </div>
              <div class="question">{{ sessionStore.currentQuestion.q }}</div>
              
              <div class="answers">
                <button
                  v-for="(option, idx) in sessionStore.currentQuestion.a"
                  :key="idx"
                  class="answer"
                  :class="getOptionClass(idx)"
                  :disabled="sessionStore.hasAnsweredCurrent"
                  @click="handleAnswer(idx)"
                >
                  <span class="letter">{{ String.fromCharCode(65 + idx) }}</span>
                  <span>{{ option }}</span>
                </button>
              </div>

              <div v-if="sessionStore.hasAnsweredCurrent" class="feedback-container">
                <div v-if="isSelectedCorrect" class="feedback correct">
                  Correct · +1,000 points
                </div>
                <div v-else class="feedback wrong">
                  Not quite · 0 points
                </div>
              </div>
            </div>
          </div>

          <aside class="card side-card">
            <div class="side-block">
              <div class="side-title">My score</div>
              <div class="score">{{ sessionStore.myScore }}</div>
              <div class="muted">points</div>
            </div>
            <div class="side-block">
              <div class="side-title">My position</div>
              <div class="rank">{{ sessionStore.myRank }}</div>
              <div class="muted">in this session</div>
            </div>
            <div class="side-block">
              <div class="side-title">Your progress</div>
              <div class="muted">
                {{ sessionStore.currentQuestionIndex + (sessionStore.hasAnsweredCurrent ? 1 : 0) }} of {{ sessionStore.questions.length }} answered
              </div>
            </div>
            <div class="side-block">
              <div class="notice">
                The host can see the full live leaderboard. Your screen shows your score and current position.
              </div>
            </div>
          </aside>
        </div>
      </section>

      <!-- 4. Final View -->
      <section v-else-if="currentStage === 'final'">
        <div class="final">
          <h1>Quiz complete</h1>
          <div class="final-sub">{{ sessionStore.sessionTitle }} · Session {{ sessionStore.code }}</div>
          
          <div class="summary">
            <div class="card">
              <div class="card-body">
                <div class="muted">Your score</div>
                <div class="num">{{ sessionStore.myScore }}</div>
              </div>
            </div>
            <div class="card">
              <div class="card-body">
                <div class="muted">Your position</div>
                <div class="num">#{{ sessionStore.myRank }}</div>
              </div>
            </div>
          </div>

          <div class="card">
            <div class="card-head">
              <strong>Final leaderboard</strong>
              <span class="muted">Session completed</span>
            </div>
            <table class="table">
              <thead>
                <tr>
                  <th>#</th>
                  <th>Participant</th>
                  <th>Score</th>
                  <th>Accuracy</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(p, index) in sessionStore.sortedParticipants"
                  :key="index"
                  :class="{ me: p.name.toLowerCase() === displayName.toLowerCase() }"
                >
                  <td class="rank">{{ index + 1 }}</td>
                  <td>
                    <strong>{{ p.name }}</strong>
                    <span v-if="p.name.toLowerCase() === displayName.toLowerCase()" class="me-tag">you</span>
                  </td>
                  <td><strong>{{ p.score }}</strong></td>
                  <td>{{ Math.round((p.correctCount || 0) / sessionStore.questions.length * 100) }}%</td>
                </tr>
              </tbody>
            </table>
          </div>

          <div style="height: 14px"></div>
          <div class="notice">
            Final standings are available after the session. During the quiz, your screen only exposes your own live score and position.
          </div>
        </div>
      </section>
    </main>

    <footer class="footer">ELSA Live Quiz · Participant</footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { useLiveSessionStore } from '../../stores/liveSession';

const route = useRoute();
const sessionStore = useLiveSessionStore();

const inputCode = ref(route.params.code || '4821');
const inputName = ref(localStorage.getItem('elsa_participant_name') || 'Alex');
const displayName = computed(() => sessionStore.myName || inputName.value);

const currentStage = ref('join'); // join, waiting, live, final
const timerCount = ref(20);
let timerInterval = null;

const progressPercent = computed(() => {
  return (sessionStore.currentQuestionIndex / sessionStore.questions.length) * 100;
});

const isSelectedCorrect = computed(() => {
  if (sessionStore.lastAnswerResult) {
    return sessionStore.lastAnswerResult.isCorrect;
  }
  return sessionStore.selectedOptionIndex === sessionStore.currentQuestion.c;
});

const getOptionClass = (idx) => {
  if (!sessionStore.hasAnsweredCurrent) return {};
  const isSelected = sessionStore.selectedOptionIndex === idx;
  const isCorrect = idx === sessionStore.currentQuestion.c;
  return {
    selected: isSelected,
    correct: isCorrect,
    wrong: isSelected && !isCorrect
  };
};

const startTimer = () => {
  clearInterval(timerInterval);
  timerCount.value = 20;
  timerInterval = setInterval(() => {
    if (timerCount.value > 0) {
      timerCount.value--;
    } else {
      clearInterval(timerInterval);
    }
  }, 1000);
};

const handleJoin = () => {
  sessionStore.joinRoom(inputCode.value, inputName.value, false);
  currentStage.value = 'waiting';
};

const handleAnswer = (optionIdx) => {
  sessionStore.submitAnswer(optionIdx);
};

// React to session state changes from WebSocket
watch(() => sessionStore.status, (newStatus) => {
  if (newStatus === 'LIVE') {
    currentStage.value = 'live';
    startTimer();
  } else if (newStatus === 'COMPLETED') {
    currentStage.value = 'final';
    clearInterval(timerInterval);
  }
});

watch(() => sessionStore.currentQuestionIndex, () => {
  if (currentStage.value === 'live') {
    startTimer();
  }
});

onMounted(() => {
  if (route.params.code) {
    inputCode.value = route.params.code;
  }
});
</script>
