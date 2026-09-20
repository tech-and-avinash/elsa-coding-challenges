<template>
  <div>
    <!-- 1. Host Lobby State -->
    <div v-if="sessionStore.status === 'LOBBY'">
      <div class="toolbar">
        <div>
          <h1>Live session</h1>
          <div class="sub">{{ sessionStore.sessionTitle }} · Waiting room</div>
        </div>
        <div class="actions">
          <router-link to="/host/quizzes" class="btn">Back to quizzes</router-link>
          <button class="btn primary" @click="sessionStore.startSession">Start session</button>
        </div>
      </div>

      <div class="grid two">
        <div class="panel">
          <div class="panel-head">
            <h2>Join details</h2>
            <span class="badge ready">Ready</span>
          </div>
          <div class="panel-body">
            <div class="muted">Session code</div>
            <div class="session-code">{{ sessionStore.code }}</div>
            <p class="muted">Share this code or the link below with participants.</p>
            <div class="join-link">{{ joinUrl }}</div>
            <div style="height: 12px"></div>
            <button class="btn small" @click="copyJoinLink">
              {{ copied ? '✓ Copied!' : 'Copy join link' }}
            </button>
          </div>
        </div>

        <div class="panel">
          <div class="panel-head">
            <h2>Participants <span>({{ sessionStore.totalParticipants }})</span></h2>
          </div>
          <div class="panel-body">
            <div
              v-for="p in sessionStore.sortedParticipants"
              :key="p.name"
              style="padding: 8px 0; border-bottom: 1px solid var(--line);"
            >
              {{ p.name }}
            </div>
            <div v-if="sessionStore.totalParticipants === 0" class="muted">
              Waiting for participants to join...
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 2. Host Live Monitoring State -->
    <div v-else-if="sessionStore.status === 'LIVE'">
      <div class="toolbar">
        <div>
          <h1>{{ sessionStore.sessionTitle }}</h1>
          <div class="sub">Live monitoring · Session {{ sessionStore.code }}</div>
        </div>
        <div class="actions">
          <span class="badge live">● Live</span>
          <button class="btn primary" @click="sessionStore.nextQuestion">
            {{ isLastQuestion ? 'Finish session' : 'Next question →' }}
          </button>
          <button class="btn danger" @click="sessionStore.endSession">End session</button>
        </div>
      </div>

      <!-- KPIs Row -->
      <div class="kpis">
        <div class="kpi">
          <div class="label">Question</div>
          <div class="value">{{ sessionStore.currentQuestionIndex + 1 }} / {{ sessionStore.questions.length }}</div>
          <div class="detail">Live</div>
        </div>
        <div class="kpi">
          <div class="label">Responses</div>
          <div class="value">{{ sessionStore.answeredCount }} / {{ sessionStore.totalParticipants }}</div>
          <div class="detail">{{ responsePercentage }}% complete</div>
        </div>
        <div class="kpi">
          <div class="label">Accuracy</div>
          <div class="value">{{ currentAccuracy }}%</div>
          <div class="detail">Current question</div>
        </div>
        <div class="kpi">
          <div class="label">Players</div>
          <div class="value">{{ sessionStore.totalParticipants }}</div>
          <div class="detail">Connected</div>
        </div>
        <div class="kpi">
          <div class="label">Avg. score</div>
          <div class="value">{{ averageScore }}</div>
          <div class="detail">Across players</div>
        </div>
      </div>

      <!-- Question + Leaderboard Grid -->
      <div class="live-host-grid">
        <section class="panel">
          <div class="panel-head">
            <h2>Current question</h2>
            <span class="muted">Question {{ sessionStore.currentQuestionIndex + 1 }} of {{ sessionStore.questions.length }}</span>
          </div>
          <div class="panel-body">
            <div class="host-question-title">{{ sessionStore.currentQuestion.q }}</div>
            
            <div
              v-for="(option, idx) in sessionStore.currentQuestion.a"
              :key="idx"
              class="host-option"
              :class="{ correct: idx === sessionStore.currentQuestion.c }"
            >
              <span><strong>{{ String.fromCharCode(65 + idx) }}</strong> &nbsp;{{ option }}</span>
              <span v-if="idx === sessionStore.currentQuestion.c" class="tag">Correct</span>
            </div>

            <div style="margin-top: 14px;" class="notice">
              Host view · The question is shown here for context while you monitor the room.
            </div>
          </div>
        </section>

        <section class="panel">
          <div class="panel-head">
            <h2>Live leaderboard</h2>
            <span class="badge live">Updates live</span>
          </div>
          <div style="padding: 0;">
            <table class="table">
              <thead>
                <tr>
                  <th>#</th>
                  <th>Participant</th>
                  <th>Score</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(p, index) in sessionStore.sortedParticipants" :key="p.name">
                  <td class="rank">{{ index + 1 }}</td>
                  <td><strong>{{ p.name }}</strong></td>
                  <td>{{ p.score }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </div>

      <div style="height: 14px"></div>

      <!-- Response Overview -->
      <section class="panel">
        <div class="panel-head">
          <h2>Response overview</h2>
          <span class="muted">Collective answers for the current question</span>
        </div>
        <div class="panel-body">
          <div class="response-dist">
            <div
              v-for="(count, idx) in sessionStore.distribution"
              :key="idx"
              class="response-row"
            >
              <div class="response-label">{{ String.fromCharCode(65 + idx) }}</div>
              <div
                class="response-bar"
                :class="{ correct: idx === sessionStore.currentQuestion.c }"
              >
                <span :style="{ width: getBarWidth(count) + '%' }"></span>
              </div>
              <div class="response-count">{{ count }}</div>
            </div>
          </div>

          <div class="response-state">
            <span class="response-dot"></span>
            <strong>{{ sessionStore.answeredCount }} of {{ sessionStore.totalParticipants }} participants have answered</strong>
          </div>
        </div>
      </section>
    </div>

    <!-- 3. Host Completed State -->
    <div v-else-if="sessionStore.status === 'COMPLETED'">
      <div class="toolbar">
        <div>
          <h1>Session results</h1>
          <div class="sub">{{ sessionStore.sessionTitle }} · Session {{ sessionStore.code }} · Completed</div>
        </div>
        <div class="actions">
          <router-link to="/host/quizzes" class="btn">Back to quizzes</router-link>
          <router-link to="/host/reports" class="btn primary">View results</router-link>
        </div>
      </div>

      <div class="grid two">
        <div class="panel">
          <div class="panel-head">
            <h2>Final leaderboard</h2>
            <span class="badge done">Completed</span>
          </div>
          <div style="padding: 0;">
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
                <tr v-for="(p, index) in sessionStore.sortedParticipants" :key="p.name">
                  <td class="rank">{{ index + 1 }}</td>
                  <td><strong>{{ p.name }}</strong></td>
                  <td><strong>{{ p.score }}</strong></td>
                  <td>{{ Math.round((p.correctCount || 0) / sessionStore.questions.length * 100) }}%</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="panel">
          <div class="panel-head">
            <h2>Session summary</h2>
          </div>
          <div class="panel-body">
            <div class="field">
              <label>Participants</label>
              <div class="big-number">{{ sessionStore.totalParticipants }}</div>
            </div>
            <div class="field">
              <label>Questions</label>
              <div class="big-number">{{ sessionStore.questions.length }}</div>
            </div>
            <div class="field">
              <label>Average score</label>
              <div class="big-number">{{ averageScore }}</div>
            </div>
            <div class="notice">
              Final standings are available after the session. During play, participants see their own score and position while the host sees the full live leaderboard.
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useLiveSessionStore } from '../../stores/liveSession';

const route = useRoute();
const sessionStore = useLiveSessionStore();

const copied = ref(false);
const joinUrl = computed(() => {
  return `${window.location.origin}/join/${sessionStore.code}`;
});

const isLastQuestion = computed(() => {
  return sessionStore.currentQuestionIndex >= sessionStore.questions.length - 1;
});

const responsePercentage = computed(() => {
  if (sessionStore.totalParticipants === 0) return 0;
  return Math.round((sessionStore.answeredCount / sessionStore.totalParticipants) * 100);
});

const currentAccuracy = computed(() => {
  const correctIdx = sessionStore.currentQuestion.c;
  const correctAnswersCount = sessionStore.distribution[correctIdx] || 0;
  if (sessionStore.answeredCount === 0) return 0;
  return Math.round((correctAnswersCount / sessionStore.answeredCount) * 100);
});

const averageScore = computed(() => {
  if (sessionStore.totalParticipants === 0) return 0;
  const total = sessionStore.sortedParticipants.reduce((sum, p) => sum + p.score, 0);
  return Math.round(total / sessionStore.totalParticipants);
});

const getBarWidth = (count) => {
  const max = Math.max(...sessionStore.distribution, 1);
  return (count / max) * 100;
};

const copyJoinLink = () => {
  navigator.clipboard.writeText(joinUrl.value);
  copied.value = true;
  setTimeout(() => (copied.value = false), 2000);
};

onMounted(() => {
  const codeParam = route.params.code || '4821';
  sessionStore.joinRoom(codeParam, 'Host', true);
});
</script>
