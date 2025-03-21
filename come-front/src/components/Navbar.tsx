// src/components/Navbar.tsx

import { AppBar, Toolbar, Typography, Button, Avatar, Menu, MenuItem, Tooltip } from '@mui/material';
import { FC, useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { getProfile, User } from '../api/user';
import { UserRole } from '../constants/roles';

const Navbar: FC = () => {
  const [user, setUser] = useState<User | null>(null);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      getProfile()
        .then((data) => {
          setUser(data);
        })
        .catch((_) => {
          localStorage.removeItem('token');
          setUser(null);
        });
    } else {
      setUser(null);
    }
  }, []);

  const handleAvatarClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    setUser(null);
    handleMenuClose();
    window.location.href = '/login';
  };

  const avatarSrc = user?.avatar ? `/${user.avatar}` : undefined;

  return (
    <AppBar position="static">
      <Toolbar>

        <Tooltip title="home">
          <Button color="inherit" component={Link} to="/" sx={{ textTransform: 'none', mr: 2 }}>
            🏠
          </Button>
        </Tooltip>
        <Tooltip title="chat">
          <Button color="inherit" component={Link} to="/chat" sx={{ textTransform: 'none', mr: 2}}>
            💬
          </Button>
        </Tooltip>
        <Typography variant="h6" sx={{ flexGrow: 1 }}>

        </Typography>
        { !user?.banned &&
          <Button color="inherit" component={Link} to="/post" sx={{ textTransform: 'none', mr: 2 }}>
            Post
          </Button>}
        {user ? (
          <>
            <Tooltip title={user.username}>
              <Avatar
                src={avatarSrc}
                sx={{
                  bgcolor: !avatarSrc? 'primary.main' : undefined,
                  width: 36,
                  height: 36,
                  cursor: 'pointer',
                  "&:hover": { border: "2px solid #1976d2" },
                }}
                onClick={handleAvatarClick}
              >
                {!user.avatar && user.username[0].toUpperCase()}
              </Avatar>
            </Tooltip>
            <Menu
              anchorEl={anchorEl}
              open={Boolean(anchorEl)}
              onClose={handleMenuClose}
              anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
              transformOrigin={{ vertical: 'top', horizontal: 'right' }}
              sx={{
                '& .MuiPaper-root': {
                  minWidth: 150,
                  bgcolor: '#ffffff',
                  boxShadow: '0 4px 8px rgba(0,0,0,0.1)',
                  borderRadius: 2,
                },
              }}
            >
              <MenuItem component={Link} to="/profile" onClick={handleMenuClose}>
                Profile
              </MenuItem>
              { Number(localStorage.getItem("role")) as UserRole === UserRole.Admin &&
                <MenuItem component={Link} to="/admin/dashboard" onClick={handleMenuClose}>
                  Dashboard
                </MenuItem>
              }
              <MenuItem onClick={handleLogout}>
                Logout
              </MenuItem>
            </Menu>
          </>
        ) : (
          <Button color="inherit" component={Link} to="/login" sx={{ textTransform: 'none' }}>
            Login
          </Button>
        )}
      </Toolbar>
    </AppBar>
  );
};

export default Navbar;
