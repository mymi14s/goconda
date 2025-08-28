// src/plugins/socket.js
import { io } from 'socket.io-client';

// Create a socket connection
const socket = io('http://localhost:8000', {
  transports: ['websocket'],  // Force WebSocket only (optional)
//   auth: {
//     token: 'your-auth-token', // Send token during handshake
//   },
//   path: '/socket.io/',        // If your backend uses a custom path
  reconnection: true,         // Enable reconnection (default)
  reconnectionAttempts: 5,    // Try 5 times before giving up
  timeout: 10000,             // Connection timeout (ms)
  autoConnect: false,         // Don't connect immediately
//   query: {
//     userId: '12345'           // Additional query params
//   },
});

export default socket;
