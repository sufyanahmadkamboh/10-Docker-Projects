const express = require('express');
const cors = require('cors');
const Redis = require('ioredis');

const app = express();
const port = process.env.PORT || 3000;

app.use(cors());
app.use(express.json());

const redis = new Redis({
  host: process.env.REDIS_HOST || 'redis',
  port: Number(process.env.REDIS_PORT) || 6379
});

const OPTIONS = ['cats', 'dogs'];

app.get('/', (req, res) => {
  res.json({ status: 'ok', message: 'Voting API is running' });
});

app.post('/vote', async (req, res) => {
  try {
    const { option } = req.body;

    if (!OPTIONS.includes(option)) {
      return res.status(400).json({ error: 'Invalid option' });
    }

    await redis.incr(`votes:${option}`);

    const [cats, dogs] = await redis.mget('votes:cats', 'votes:dogs');
    return res.json({
      success: true,
      votes: {
        cats: Number(cats) || 0,
        dogs: Number(dogs) || 0
      }
    });
  } catch (err) {
    console.error('Error in /vote:', err);
    res.status(500).json({ error: 'Internal server error' });
  }
});

app.get('/results', async (req, res) => {
  try {
    const [cats, dogs] = await redis.mget('votes:cats', 'votes:dogs');
    const c = Number(cats) || 0;
    const d = Number(dogs) || 0;

    res.json({
      votes: {
        cats: c,
        dogs: d
      },
      total: c + d
    });
  } catch (err) {
    console.error('Error in /results:', err);
    res.status(500).json({ error: 'Internal server error' });
  }
});

app.listen(port, () => {
  console.log(`Voting API listening on port ${port}`);
});
