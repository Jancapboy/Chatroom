import { useState, useEffect } from 'react';
import { Lobby } from './pages/Lobby';
import { RoomPage } from './pages/Room';
import { CreateRoom } from './pages/CreateRoom';
import { Login } from './pages/Login';

type Page = 'login' | 'lobby' | 'room' | 'create';

interface NavigationState {
  page: Page;
  params?: Record<string, string>;
}

function App() {
  const [nav, setNav] = useState<NavigationState>({ page: 'lobby' });

  const navigate = (page: string, params?: Record<string, string>) => {
    setNav({ page: page as Page, params });
    // Update URL for shareability (optional)
    if (page === 'room' && params?.roomId) {
      window.history.replaceState(null, '', `/room/${params.roomId}`);
    } else if (page === 'lobby') {
      window.history.replaceState(null, '', '/');
    }
  };

  // Handle initial URL
  useEffect(() => {
    const path = window.location.pathname;
    if (path.startsWith('/room/')) {
      const roomId = path.split('/')[2];
      if (roomId) setNav({ page: 'room', params: { roomId } });
    }
  }, []);

  switch (nav.page) {
    case 'login':
      return <Login onNavigate={(p) => navigate(p)} />;
    case 'lobby':
      return <Lobby onNavigate={(p, params) => navigate(p, params)} />;
    case 'room':
      return <RoomPage onNavigate={(p) => navigate(p)} />;
    case 'create':
      return <CreateRoom onNavigate={(p) => navigate(p)} />;
    default:
      return <Lobby onNavigate={(p, params) => navigate(p, params)} />;
  }
}

export default App;
