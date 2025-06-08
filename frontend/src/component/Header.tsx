import { useNavigate } from 'react-router-dom';
import './Header.css';

export const Header = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/');
  };

  return (
    <header className="header">
      <h1 className="logo">Xpense 💸</h1>
      <button className="logout-btn" onClick={handleLogout}>
        Logout
      </button>
    </header>
  );
};
