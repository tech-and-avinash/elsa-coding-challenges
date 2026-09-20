import express from 'express';
import { createServer } from 'http';
import { Server } from 'socket.io';

const app = express();
const httpServer = createServer(app);
const io = new Server(httpServer, {
  cors: {
    origin: '*',
    methods: ['GET', 'POST']
  }
});

app.use(express.json());

// Initial quiz templates
const defaultQuizzes = [
  {
    id: 'quiz-1',
    title: 'Business English Basics',
    description: 'Everyday vocabulary for meetings, email, and workplace communication.',
    category: 'Vocabulary · General',
    timePerQuestion: 20,
    scoringMode: 'Speed & Accuracy',
    status: 'Ready',
    updatedAt: 'Today',
    questions: [
      { q: "What does “concise” mean?", a: ["Brief and clear", "Difficult to understand", "Very detailed", "Informal"], c: 0 },
      { q: "Which word means “to postpone”?", a: ["Confirm", "Defer", "Resolve", "Proceed"], c: 1 },
      { q: "Choose the best phrase for a meeting follow-up.", a: ["Please find below", "Let's circle back", "I am agree", "Do the needful"], c: 1 },
      { q: "What is the opposite of “expand”?", a: ["Increase", "Extend", "Reduce", "Explain"], c: 2 },
      { q: "Which word means “necessary or essential”?", a: ["Optional", "Relevant", "Mandatory", "Casual"], c: 2 }
    ]
  },
  {
    id: 'quiz-2',
    title: 'Travel Vocabulary',
    description: 'Essential phrases for airports, hotels, and directions abroad.',
    category: 'Vocabulary · B1',
    timePerQuestion: 20,
    scoringMode: 'Speed & Accuracy',
    status: 'Ready',
    updatedAt: 'Yesterday',
    questions: [
      { q: "What is an itinerary?", a: ["A travel plan", "A luggage bag", "A passport visa", "A flight delay"], c: 0 },
      { q: "Which phrase is used at hotel check-in?", a: ["I have a reservation", "I want a refund", "Where is the exit", "Boarding pass please"], c: 0 },
      { q: "What does 'layover' mean?", a: ["A stopover between flights", "Overbooked seat", "Lost luggage", "Direct flight"], c: 0 }
    ]
  }
];

// Active sessions memory store
// Session structure:
// { code: '4821', quizId: 'quiz-1', title: '...', questions: [], currentQuestionIndex: 0, status: 'LOBBY' | 'LIVE' | 'COMPLETED', participants: {} }
const sessions = {};

// Create initial demo session 4821 so it's always ready for demo!
sessions['4821'] = {
  code: '4821',
  quizId: 'quiz-1',
  title: 'Business English Basics',
  questions: defaultQuizzes[0].questions,
  timePerQuestion: 20,
  currentQuestionIndex: 0,
  status: 'LOBBY',
  participants: {
    'p-priya': { name: 'Priya', score: 980, correctCount: 5, answers: {} },
    'p-daniel': { name: 'Daniel', score: 920, correctCount: 5, answers: {} },
    'p-maya': { name: 'Maya', score: 860, correctCount: 4, answers: {} },
    'p-noah': { name: 'Noah', score: 780, correctCount: 4, answers: {} },
    'p-sara': { name: 'Sara', score: 720, correctCount: 4, answers: {} },
    'p-liam': { name: 'Liam', score: 680, correctCount: 3, answers: {} },
    'p-emma': { name: 'Emma', score: 640, correctCount: 3, answers: {} },
    'p-ravi': { name: 'Ravi', score: 600, correctCount: 3, answers: {} }
  }
};

// REST endpoints
app.get('/api/quizzes', (req, res) => {
  res.json(defaultQuizzes);
});

app.get('/api/sessions', (req, res) => {
  const list = Object.values(sessions).map(s => ({
    code: s.code,
    title: s.title,
    status: s.status,
    participantCount: Object.keys(s.participants).length,
    updatedAt: s.status === 'LIVE' ? 'Just now' : 'Yesterday'
  }));
  res.json(list);
});

// Helper: Calculate answer distribution for a question index
function getQuestionDistribution(session, qIndex) {
  const dist = [0, 0, 0, 0];
  let answeredCount = 0;
  Object.values(session.participants).forEach(p => {
    if (p.answers && p.answers[qIndex] !== undefined) {
      const optionIdx = p.answers[qIndex].selectedOption;
      if (optionIdx >= 0 && optionIdx < 4) {
        dist[optionIdx]++;
      }
      answeredCount++;
    }
  });
  return { dist, answeredCount };
}

// Socket.io real-time handling
io.on('connection', (socket) => {
  console.log(`[Socket] Client connected: ${socket.id}`);

  // Create or reset host session
  socket.on('create_session', ({ quizId, customCode }) => {
    const quiz = defaultQuizzes.find(q => q.id === quizId) || defaultQuizzes[0];
    const code = customCode || Math.floor(1000 + Math.random() * 9000).toString();

    sessions[code] = {
      code,
      quizId: quiz.id,
      title: quiz.title,
      questions: quiz.questions,
      timePerQuestion: quiz.timePerQuestion || 20,
      currentQuestionIndex: 0,
      status: 'LOBBY',
      participants: sessions[code]?.participants || {
        'p-priya': { name: 'Priya', score: 980, correctCount: 5, answers: {} },
        'p-daniel': { name: 'Daniel', score: 920, correctCount: 5, answers: {} },
        'p-maya': { name: 'Maya', score: 860, correctCount: 4, answers: {} }
      }
    };

    socket.join(`room_${code}`);
    socket.emit('session_created', sessions[code]);
  });

  // Participant or Host joins room
  socket.on('join_room', ({ code, name, isHost }) => {
    let session = sessions[code];
    if (!session) {
      // Auto-create room for demo if missing
      const quiz = defaultQuizzes[0];
      session = {
        code: code || '4821',
        quizId: quiz.id,
        title: quiz.title,
        questions: quiz.questions,
        timePerQuestion: 20,
        currentQuestionIndex: 0,
        status: 'LOBBY',
        participants: {}
      };
      sessions[session.code] = session;
    }

    socket.join(`room_${session.code}`);

    let participantId = null;
    if (!isHost && name) {
      // Find existing by name or create new participant
      const existingKey = Object.keys(session.participants).find(
        k => session.participants[k].name.toLowerCase() === name.toLowerCase()
      );

      if (existingKey) {
        participantId = existingKey;
        session.participants[existingKey].socketId = socket.id;
      } else {
        participantId = `p_${socket.id.substring(0, 6)}`;
        session.participants[participantId] = {
          name,
          score: 0,
          correctCount: 0,
          socketId: socket.id,
          answers: {}
        };
      }
    }

    const distInfo = getQuestionDistribution(session, session.currentQuestionIndex);

    // Emit updated state to joining socket
    socket.emit('room_state', {
      session,
      participantId,
      distribution: distInfo.dist,
      answeredCount: distInfo.answeredCount
    });

    // Broadcast participant update to host and others in room
    io.to(`room_${session.code}`).emit('participants_updated', {
      participants: session.participants,
      count: Object.keys(session.participants).length
    });
  });

  // Host starts session
  socket.on('start_session', ({ code }) => {
    const session = sessions[code];
    if (session) {
      session.status = 'LIVE';
      session.currentQuestionIndex = 0;
      io.to(`room_${code}`).emit('session_started', {
        status: 'LIVE',
        currentQuestionIndex: 0
      });
    }
  });

  // Host moves to next question
  socket.on('next_question', ({ code, index }) => {
    const session = sessions[code];
    if (session && index < session.questions.length) {
      session.currentQuestionIndex = index;
      const distInfo = getQuestionDistribution(session, index);
      io.to(`room_${code}`).emit('question_changed', {
        currentQuestionIndex: index,
        distribution: distInfo.dist,
        answeredCount: distInfo.answeredCount
      });
    }
  });

  // Participant submits an answer
  socket.on('submit_answer', ({ code, name, questionIndex, selectedOption }) => {
    const session = sessions[code];
    if (!session) return;

    const pKey = Object.keys(session.participants).find(
      k => session.participants[k].name.toLowerCase() === name.toLowerCase()
    );

    if (pKey) {
      const participant = session.participants[pKey];
      const q = session.questions[questionIndex];
      const isCorrect = selectedOption === q.c;

      if (!participant.answers) participant.answers = {};
      
      // Calculate score if not already answered this question
      if (participant.answers[questionIndex] === undefined) {
        const points = isCorrect ? 1000 : 0;
        participant.score += points;
        if (isCorrect) participant.correctCount++;
        participant.answers[questionIndex] = { selectedOption, isCorrect, points };
      }

      const distInfo = getQuestionDistribution(session, questionIndex);

      // Send live feedback to this participant
      socket.emit('answer_result', {
        questionIndex,
        isCorrect,
        correctOption: q.c,
        myScore: participant.score,
        myCorrectCount: participant.correctCount
      });

      // Broadcast distribution & leaderboard update to host and room
      io.to(`room_${code}`).emit('live_update', {
        participants: session.participants,
        distribution: distInfo.dist,
        answeredCount: distInfo.answeredCount
      });
    }
  });

  // Host ends session
  socket.on('end_session', ({ code }) => {
    const session = sessions[code];
    if (session) {
      session.status = 'COMPLETED';
      io.to(`room_${code}`).emit('session_ended', {
        status: 'COMPLETED',
        participants: session.participants
      });
    }
  });

  socket.on('disconnect', () => {
    console.log(`[Socket] Client disconnected: ${socket.id}`);
  });
});

const PORT = process.env.PORT || 3001;
httpServer.listen(PORT, () => {
  console.log(`[Server] ELSA Quiz Real-Time Server running on http://localhost:${PORT}`);
});
