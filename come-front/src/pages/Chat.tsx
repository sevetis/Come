// src/pages/Chat.tsx

import React, { useEffect, useState, useRef } from 'react';
import { Container, Typography, TextField, Button, List, ListItem, ListItemText, CircularProgress, Box, Chip } from '@mui/material';
import { getChatHistory, ChatMessage, getOnlineCount } from '../api/chat';
import { useNavigate } from 'react-router-dom';

const getWebSocketUrl = (token: string) => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const host = window.location.host;
  return `${protocol}//${host}:8080/api/chat?token=${token}`;
};

const Chat: React.FC = () => {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [onlineCount, setOnlineCount] = useState(0);
  const ws = useRef<WebSocket | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate();

  useEffect(() => {
    let mounted = true;

    const fetchOnlineCount = async () => {
      try {
        const count = await getOnlineCount()
        if (mounted) setOnlineCount(count);
      } catch (err) {
        console.error('Failed to fetch online count:', err);
      }
    };

    const initChat = async () => {
      try {
        setLoading(true);
        const history = await getChatHistory();
        if (mounted && history !== null) {
          setMessages(history.map(msg => ({
              ...msg,
              timestamp: msg.timestamp * 1000,
          })));
        }

        const token = localStorage.getItem('token') || '';
        ws.current = new WebSocket(getWebSocketUrl(token));

        ws.current.onopen = () => {
          console.log('Connected to chat');
          if (mounted) {
            setLoading(false);
            fetchOnlineCount();
          }
        };
        ws.current.onmessage = (event) => {
          const msg = JSON.parse(event.data);
          console.log('Received message:', msg);
          if (mounted) {
            setMessages((prev) => [...prev, {
              id: Date.now(),
              userId: msg.userId,
              username: msg.username,
              content: msg.content,
              timestamp: msg.timestamp * 1000,
              type: msg.type,
            }]);
            if (msg.type === 'join' || msg.type === 'leave')
              fetchOnlineCount();
          }
        };
        ws.current.onerror = (e) => {
          console.error('WebSocket error:', e);
          if (mounted) setError('Failed to connect to chat');
        };
        ws.current.onclose = (e) => {
          console.log('Disconnected from chat:', e.code, e.reason);
          if (mounted) setError('Chat connection closed');
        };
      } catch (err) {
        console.error('Init chat error:', err);
        if (mounted) setError('Failed to load chat history');
      } finally {
        if (mounted) setLoading(false);
      }
    };

    initChat();

    return () => {
      mounted = false;
      // ws.current?.close();
      if (ws.current) {
        ws.current.close();
        console.log("WebSocket closed on cleanup");
      }
    };
  }, []);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const sendMessage = () => {
    if (input.trim() && ws.current && ws.current.readyState === WebSocket.OPEN) {
      const msg = {"content": input};
      console.log('Sending message:', msg);
      ws.current.send(JSON.stringify({ content: input }));
      setInput('');
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  }

  if (loading) {
    return (
      <Container maxWidth="md" sx={{ mt: 4, textAlign: 'center' }}>
        <CircularProgress />
        <Typography>Loading chat...</Typography>
      </Container>
    );
  }

  if (error) {
    return (
      <Container maxWidth="md" sx={{ mt: 4 }}>
        <Typography color="error">{error}</Typography>
      </Container>
    );
  }

  return (
    <Container maxWidth="md" sx={{ mt: 4, mb: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h4">Chat Room</Typography>
          <Chip label={`online: ${onlineCount}`} color="primary" />
        </Box>
        <List sx={{ maxHeight: '400px', overflowY: 'auto', mb: 2, bgcolor: '#f5f5f5', borderRadius: 2, p: 1 }}>
          {messages.length > 0 ? (
            messages.map((msg) => (
              <ListItem key={msg.id} sx={{ bgcolor: msg.type === 'join' || msg.type === 'leave' ? '#e0f7fa' : 'inherit' }}>
                <ListItemText
                  primary={msg.type === 'message' ? `${msg.username}: ${msg.content}` : msg.content}
                  secondary={new Date(msg.timestamp).toLocaleTimeString()}
                  sx={{
                    color: msg.userId === 0 ? 'grey' : 'inherit',
                    fontStyle: msg.type !== 'message' ? 'italic' : 'normal',
                  }}
                />
              </ListItem>
            ))
          ) : (
            <Typography sx={{ p: 2 }}>No messages yet</Typography>
          )}
          <div ref={messagesEndRef} />
        </List>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
          <TextField
            fullWidth
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Type a message..."
            multiline
            minRows={2}
            maxRows={6}
            sx={{ bgcolor: 'white', borderRadius: 1 }}
          />
          <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
            <Button variant="outlined" onClick={() => {navigate(-1)}}>Back</Button>
            <Button variant="contained" onClick={sendMessage}>Send</Button>
          </Box>
        </Box>
      </Container>
  );
};

export default Chat;

