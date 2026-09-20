class WsClient {
  constructor() {
    this.ws = null;
    this.listeners = new Map();
    this.connected = false;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.pendingEmits = [];
    this.connect();
  }

  getWsUrl() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    return `${protocol}//${host}/ws`;
  }

  connect() {
    if (this.ws && (this.ws.readyState === WebSocket.OPEN || this.ws.readyState === WebSocket.CONNECTING)) {
      return;
    }

    const url = this.getWsUrl();
    console.log('[WebSocket Service] Connecting to Go backend:', url);
    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      console.log('[WebSocket Service] Connected successfully');
      this.connected = true;
      this.reconnectAttempts = 0;
      this.trigger('connect');

      // Flush queued emit events
      while (this.pendingEmits.length > 0) {
        const { event, data } = this.pendingEmits.shift();
        this.emit(event, data);
      }
    };

    this.ws.onmessage = (messageEvent) => {
      try {
        const msg = JSON.parse(messageEvent.data);
        if (msg.event) {
          this.trigger(msg.event, msg.payload);
        }
      } catch (err) {
        console.error('[WebSocket Service] Error parsing message payload:', err);
      }
    };

    this.ws.onerror = (err) => {
      console.error('[WebSocket Service] Socket error:', err);
    };

    this.ws.onclose = () => {
      console.log('[WebSocket Service] Disconnected from Go backend');
      this.connected = false;
      this.trigger('disconnect');

      if (this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++;
        setTimeout(() => this.connect(), 1000 * Math.min(this.reconnectAttempts, 5));
      }
    };
  }

  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, []);
    }
    this.listeners.get(event).push(callback);
  }

  off(event, callback) {
    if (!this.listeners.has(event)) return;
    if (!callback) {
      this.listeners.delete(event);
    } else {
      const callbacks = this.listeners.get(event).filter(cb => cb !== callback);
      this.listeners.set(event, callbacks);
    }
  }

  trigger(event, payload) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).forEach(cb => cb(payload));
    }
  }

  emit(event, data = {}) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ event, data }));
    } else {
      this.pendingEmits.push({ event, data });
      if (!this.ws || this.ws.readyState === WebSocket.CLOSED) {
        this.connect();
      }
    }
  }
}

let socketInstance = null;

export function getSocket() {
  if (!socketInstance) {
    socketInstance = new WsClient();
  }
  return socketInstance;
}
