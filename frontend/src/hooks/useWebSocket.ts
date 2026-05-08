import { useEffect, useRef, useCallback, useState } from 'react';
import type { WSClientMessage, WSServerMessage } from '../types/ws';

const WS_BASE_URL = import.meta.env.VITE_WS_URL || `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.hostname}:4001`;
const RECONNECT_INTERVAL = 3000;
const HEARTBEAT_INTERVAL = 30000;

export function useWebSocket(roomId: string | null) {
  const [isConnected, setIsConnected] = useState(false);
  const [lastMessage, setLastMessage] = useState<WSServerMessage | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const heartbeatTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const maxReconnectAttempts = 10;

  const connect = useCallback(() => {
    if (!roomId) return;
    if (wsRef.current?.readyState === WebSocket.OPEN) return;

    try {
      const ws = new WebSocket(`${WS_BASE_URL}/ws/rooms/${roomId}`);
      wsRef.current = ws;

      ws.onopen = () => {
        setIsConnected(true);
        reconnectAttemptsRef.current = 0;
        // Start heartbeat
        heartbeatTimerRef.current = setInterval(() => {
          if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ type: 'ping' }));
          }
        }, HEARTBEAT_INTERVAL);
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data) as WSServerMessage;
          setLastMessage(data);
        } catch {
          console.warn('Invalid WS message:', event.data);
        }
      };

      ws.onclose = () => {
        setIsConnected(false);
        if (heartbeatTimerRef.current) {
          clearInterval(heartbeatTimerRef.current);
        }

        // Auto reconnect
        if (reconnectAttemptsRef.current < maxReconnectAttempts) {
          reconnectAttemptsRef.current++;
          reconnectTimerRef.current = setTimeout(() => {
            connect();
          }, RECONNECT_INTERVAL);
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
    } catch (error) {
      console.error('Failed to create WebSocket:', error);
    }
  }, [roomId]);

  const disconnect = useCallback(() => {
    if (reconnectTimerRef.current) {
      clearTimeout(reconnectTimerRef.current);
    }
    if (heartbeatTimerRef.current) {
      clearInterval(heartbeatTimerRef.current);
    }
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    setIsConnected(false);
  }, []);

  const sendMessage = useCallback((msg: WSClientMessage) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(msg));
      return true;
    }
    return false;
  }, []);

  const sendUserMessage = useCallback((content: string) => {
    return sendMessage({
      type: 'user_message',
      payload: { content },
    });
  }, [sendMessage]);

  const sendCommand = useCallback((command: 'pause' | 'resume' | 'next_phase' | 'fork') => {
    return sendMessage({
      type: 'command',
      payload: { command },
    });
  }, [sendMessage]);

  useEffect(() => {
    if (roomId) {
      connect();
    }
    return () => {
      disconnect();
    };
  }, [roomId, connect, disconnect]);

  return {
    isConnected,
    lastMessage,
    sendUserMessage,
    sendCommand,
    connect,
    disconnect,
  };
}
