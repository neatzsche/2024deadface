const mysql = require('mysql2/promise');
const bcrypt = require('bcryptjs');

async function initDB() {
  const connection = await mysql.createConnection({
    host: '127.0.0.1',
    user: 'root',
    password: 'your_mysql_root_password',
  });

  // Create database
  await connection.query('CREATE DATABASE IF NOT EXISTS my_app_db');
  await connection.query('USE my_app_db');

  // Create tables
  await connection.query(`
    CREATE TABLE IF NOT EXISTS users (
      id INT AUTO_INCREMENT PRIMARY KEY,
      username VARCHAR(255) UNIQUE NOT NULL,
      password VARCHAR(255) NOT NULL
    )
  `);

  await connection.query(`
    CREATE TABLE IF NOT EXISTS targets (
      id INT AUTO_INCREMENT PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      info TEXT NOT NULL
    )
  `);

  // Seed users
  const passwordHash = await bcrypt.hash('xpeanutxbutter', 10);
  const adminHash = await bcrypt.hash('sdfahkfhsakdjhfkjhsdakfhskdahfkhsakdhfkdsah', 10);

  await connection.query(
    'INSERT IGNORE INTO users (username, password) VALUES (?, ?), (?, ?)',
    ['admin', adminHash, 'user', passwordHash]
  );

  // Seed targets
  await connection.query(
    'INSERT INTO targets (name, info) VALUES (?, ?), (?, ?), (?, ?), (?,?), (?,?)',
    ['Alice', 'Works at Demonne Financial', 'Bob', 'Works at Lytton Labs', 'Charlie', 'Info about Charlie', 'Amber', 'May have password we need', 'zflag', 'flag{SQL-1nj3ct10n-thrU-x0r}']
  );

  // Create a read-only user and grant SELECT access
  await connection.query(`
    CREATE USER IF NOT EXISTS 'readonly'@'localhost' IDENTIFIED BY 'readonly_password'
  `);

  await connection.query(`
    GRANT SELECT ON my_app_db.* TO 'readonly'@'localhost'
  `);

  await connection.query(`
    FLUSH PRIVILEGES
  `);

  await connection.end();
}

initDB()
  .then(() => {
    console.log('Database initialized with readonly user');
    process.exit(0);
  })
  .catch((err) => {
    console.error('Error initializing database:', err);
    process.exit(1);
  });
