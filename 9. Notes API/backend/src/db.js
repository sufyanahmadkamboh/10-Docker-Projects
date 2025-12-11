const { Pool } = require('pg');

const pool = new Pool({
  host: process.env.DB_HOST || 'db',
  port: Number(process.env.DB_PORT) || 5432,
  user: process.env.DB_USER || 'notes_user',
  password: process.env.DB_PASSWORD || 'notes_password',
  database: process.env.DB_NAME || 'notes_db'
});

async function query(text, params) {
  const res = await pool.query(text, params);
  return res;
}

module.exports = {
  query
};
