import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useQuizStore = defineStore('quiz', () => {
  const quizzes = ref([]);
  const sessions = ref([]);
  const activeQuiz = ref(null);
  const loading = ref(false);
  const quizSessions = ref({}); // Map quizId -> sessions array

  const fetchQuizzes = async () => {
    loading.value = true;
    try {
      const res = await fetch('/api/quizzes');
      if (res.ok) {
        const data = await res.json();
        if (Array.isArray(data) && data.length > 0) {
          quizzes.value = data;
        }
      }
    } catch (err) {
      console.warn('[QuizStore] Error fetching quizzes from Go API, using default fallback:', err);
    } finally {
      loading.value = false;
    }
  };

  const getQuizById = (id) => {
    return quizzes.value.find(q => q.id === id) || quizzes.value[0];
  };

  const createQuiz = async () => {
    const newId = `quiz-${Date.now()}`;
    const newQuiz = {
      id: newId,
      title: 'New Vocabulary Quiz',
      description: 'Custom vocabulary quiz for team practice.',
      category: 'Vocabulary · Custom',
      timePerQuestion: 20,
      scoringMode: 'Speed & Accuracy',
      status: 'Ready',
      updatedAt: 'Just now',
      questions: [
        { q: "What is the meaning of efficiency?", a: ["High productivity", "Low effort", "Slow work", "Random choice"], c: 0 }
      ]
    };
    
    quizzes.value.push(newQuiz);

    try {
      await fetch('/api/quizzes', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newQuiz)
      });
    } catch (err) {
      console.error('[QuizStore] Failed to create quiz via API:', err);
    }

    return newQuiz;
  };

  const saveQuiz = async (quizData) => {
    const idx = quizzes.value.findIndex(q => q.id === quizData.id);
    const updated = { ...quizData, updatedAt: 'Just now' };
    if (idx !== -1) {
      quizzes.value[idx] = updated;
    } else {
      quizzes.value.push(updated);
    }

    try {
      await fetch('/api/quizzes', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updated)
      });
    } catch (err) {
      console.error('[QuizStore] Failed to save quiz via API:', err);
    }
  };

  const fetchQuizSessions = async (quizId) => {
    try {
      const res = await fetch(`/api/quizzes/${quizId}/sessions`);
      if (res.ok) {
        const data = await res.json();
        quizSessions.value[quizId] = Array.isArray(data) ? data : [];
        return quizSessions.value[quizId];
      }
    } catch (err) {
      console.error('[QuizStore] Error fetching quiz sessions:', err);
    }
    quizSessions.value[quizId] = [];
    return [];
  };

  const fetchSessions = async () => {
    try {
      const res = await fetch('/api/sessions');
      if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
      }
      const data = await res.json();
      sessions.value = Array.isArray(data) ? data : [];
      return sessions.value;
    } catch (err) {
      console.error('[QuizStore] Error fetching sessions list:', err);
      sessions.value = [];
      return [];
    }
  };

  // Automatically fetch on store creation
  fetchQuizzes();

  return {
    quizzes,
    sessions,
    activeQuiz,
    loading,
    quizSessions,
    fetchQuizzes,
    fetchSessions,
    getQuizById,
    createQuiz,
    saveQuiz,
    fetchQuizSessions
  };
});
