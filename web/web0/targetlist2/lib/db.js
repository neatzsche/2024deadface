import mysql from 'mysql2/promise';

let pool;

export function getDBConnection() {
  if (!pool) {
    pool = mysql.createPool({
      host: '127.0.0.1',
      user: 'readonly',
      password: 'readonly_password',
      database: 'my_app_db',
      waitForConnections: true,
      connectionLimit: 10,
      queueLimit: 0,
    });
  }
  return pool;
}