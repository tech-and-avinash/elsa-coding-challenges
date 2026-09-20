import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { getSocket } from '../services/socket';

export const useLiveSessionStore = defineStore('liveSession', () => {
  const socket = getSocket();

  const code = ref(null);
  const sessionTitle = ref(null);
  const status = ref('LOBBY'); // LOBBY, LIVE, COMPLETED
  const currentQuestionIndex = ref(0);
  const questions = ref([
    { q: "What does “concise” mean?", a: ["Brief and clear", "Difficult to understand", "Very detailed", "Informal"], c: 0 },
    { q: "Which word means “to postpone”?", a: ["Confirm", "Defer", "Resolve", "Proceed"], c: 1 },
    { q: "Choose the best phrase for a meeting follow-up.", a: ["Please find below", "Let's circle back", "I am agree", "Do the needful"], c: 1 },
    { q: "What is the opposite of “expand”?", a: ["Increase", "Extend", "Reduce", "Explain"], c: 2 },
    { q: "Which word means “necessary or essential”?", a: ["Optional", "Relevant", "Mandatory", "Casual"], c: 2 }
  ]);

  const participantsMap = ref({});
  const distribution = ref([0, 0, 0, 0]);
  const answeredCount = ref(0);

  // Participant local state
  const myName = ref(localStorage.getItem('elsa_participant_name') || 'Alex');
  const myScore = ref(0);
  const myCorrectCount = ref(0);
  const hasAnsweredCurrent = ref(false);
  const selectedOptionIndex = ref(null);
  const lastAnswerResult = ref(null);

  // Connected participants array sorted by score
  const sortedParticipants = computed(() => {
    const list = Object.values(participantsMap.value);
    return list.sort((a, b) => b.score - a.score);
  });

  const totalParticipants = computed(() => Object.keys(participantsMap.value).length);

  const myRank = computed(() => {
    const idx = sortedParticipants.value.findIndex(p => p.name.toLowerCase() === myName.value.toLowerCase());
    return idx >= 0 ? idx + 1 : '—';
  });

  const currentQuestion = computed(() => {
    return questions.value[currentQuestionIndex.value] || questions.value[0];
  });

  // Setup Socket listeners
  const initListeners = () => {
    socket.off('room_state');
    socket.off('participants_updated');
    socket.off('session_started');
    socket.off('question_changed');
    socket.off('answer_result');
    socket.off('live_update');
    socket.off('session_ended');
    socket.off('session_created');

    socket.on('session_created', (data) => {
      code.value = data.code;
      sessionTitle.value = data.title;
      status.value = data.status;
      currentQuestionIndex.value = data.currentQuestionIndex || 0;
      questions.value = data.questions || questions.value;
      participantsMap.value = data.participants || {};
      // Trigger custom event for UI to listen to
      window.dispatchEvent(new CustomEvent('session-created', { detail: data }));
    });

    socket.on('room_state', (data) => {
      if (data.session) {
        code.value = data.session.code;
        sessionTitle.value = data.session.title;
        status.value = data.session.status;
        currentQuestionIndex.value = data.session.currentQuestionIndex || 0;
        questions.value = data.session.questions || questions.value;
        participantsMap.value = data.session.participants || {};
      }
      if (data.distribution) distribution.value = data.distribution;
      if (data.answeredCount !== undefined) answeredCount.value = data.answeredCount;
    });

    socket.on('participants_updated', (data) => {
      participantsMap.value = data.participants;
    });

    socket.on('session_started', (data) => {
      status.value = 'LIVE';
      currentQuestionIndex.value = 0;
      hasAnsweredCurrent.value = false;
      selectedOptionIndex.value = null;
      lastAnswerResult.value = null;
    });

    socket.on('question_changed', (data) => {
      currentQuestionIndex.value = data.currentQuestionIndex;
      if (data.distribution) distribution.value = data.distribution;
      if (data.answeredCount !== undefined) answeredCount.value = data.answeredCount;
      hasAnsweredCurrent.value = false;
      selectedOptionIndex.value = null;
      lastAnswerResult.value = null;
    });

    socket.on('answer_result', (data) => {
      lastAnswerResult.value = data;
      myScore.value = data.myScore;
      myCorrectCount.value = data.myCorrectCount;
    });

    socket.on('live_update', (data) => {
      if (data.participants) participantsMap.value = data.participants;
      if (data.distribution) distribution.value = data.distribution;
      if (data.answeredCount !== undefined) answeredCount.value = data.answeredCount;
    });

    socket.on('session_ended', (data) => {
      status.value = 'COMPLETED';
      if (data.participants) participantsMap.value = data.participants;
    });
  };

  // Actions
  const createSession = (quizId, customCode = null) => {
    initListeners();
    socket.emit('create_session', { quizId, customCode });
  };

  const joinRoom = (roomCode, name, isHost = false) => {
    initListeners();
    code.value = roomCode || '4821';
    if (name) {
      myName.value = name;
      localStorage.setItem('elsa_participant_name', name);
    }
    socket.emit('join_room', { code: code.value, name: myName.value, isHost });
  };

  const startSession = () => {
    socket.emit('start_session', { code: code.value });
  };

  const nextQuestion = () => {
    if (currentQuestionIndex.value < questions.value.length - 1) {
      socket.emit('next_question', {
        code: code.value,
        index: currentQuestionIndex.value + 1
      });
    } else {
      endSession();
    }
  };

  const submitAnswer = (optionIndex) => {
    if (hasAnsweredCurrent.value) return;
    hasAnsweredCurrent.value = true;
    selectedOptionIndex.value = optionIndex;
    socket.emit('submit_answer', {
      code: code.value,
      name: myName.value,
      questionIndex: currentQuestionIndex.value,
      selectedOption: optionIndex
    });
  };

  const endSession = () => {
    socket.emit('end_session', { code: code.value });
  };

  return {
    code,
    sessionTitle,
    status,
    currentQuestionIndex,
    questions,
    currentQuestion,
    participantsMap,
    sortedParticipants,
    totalParticipants,
    distribution,
    answeredCount,
    myName,
    myScore,
    myCorrectCount,
    myRank,
    hasAnsweredCurrent,
    selectedOptionIndex,
    lastAnswerResult,
    createSession,
    joinRoom,
    startSession,
    nextQuestion,
    submitAnswer,
    endSession
  };
});
