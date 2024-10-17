// pages/api/login.js

import { getDBConnection } from '../../lib/db';
import bcrypt from 'bcryptjs';
import { serialize } from 'cookie';

export default async function handler(req, res) {
  if (req.method === 'POST') {
    const { username, password } = req.body;

    try {
      const db = getDBConnection();
      const [users] = await db.query('SELECT * FROM users WHERE username = ?', [username]);

      if (users.length > 0) {
        const user = users[0];
        const isValid = await bcrypt.compare(password, user.password);

        if (isValid) {
          // Set a cookie
          res.setHeader(
            'Set-Cookie',
            serialize('auth', user.username + ':secretkey729374jhjh98279', {
              path: '/',
              httpOnly: true,
              sameSite: 'strict',
              maxAge: 60 * 60 * 24, // 1 day
            })
          );
          res.status(200).json({ message: 'Logged in' });
          return;
        }
      }

      res.status(401).json({ message: 'Invalid credentials' });
    } catch (error) {
      console.error('Login error:', error);
      res.status(500).json({ message: 'Internal server error' });
    }
  } else {
    res.status(405).end(); // Method Not Allowed
  }
}
