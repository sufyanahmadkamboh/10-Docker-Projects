const express = require('express');
const db = require('./db');

const app = express();
const port = process.env.PORT || 3000;

app.use(express.json());

app.get('/', (req, res) => {
  res.json({ status: 'ok', message: 'Notes API is running' });
});

// Get all notes
app.get('/notes', async (req, res) => {
  try {
    const result = await db.query('SELECT id, title, body, created_at FROM notes ORDER BY created_at DESC');
    res.json({ notes: result.rows });
  } catch (err) {
    console.error('Error fetching notes:', err);
    res.status(500).json({ error: 'Internal server error' });
  }
});

// Create a note
app.post('/notes', async (req, res) => {
  try {
    const { title, body } = req.body;

    if (!title || !body) {
      return res.status(400).json({ error: 'title and body are required' });
    }

    const result = await db.query(
      'INSERT INTO notes (title, body) VALUES ($1, $2) RETURNING id, title, body, created_at',
      [title, body]
    );

    res.status(201).json({ note: result.rows[0] });
  } catch (err) {
    console.error('Error creating note:', err);
    res.status(500).json({ error: 'Internal server error' });
  }
});

app.listen(port, () => {
  console.log(`Notes API listening on port ${port}`);
});
