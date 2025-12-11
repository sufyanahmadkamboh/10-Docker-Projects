const express = require('express');
const { enqueueTask } = require('./queue');

const app = express();
const port = process.env.PORT || 4000;

app.use(express.json());

app.get('/', (req, res) => {
  res.json({
    status: 'ok',
    message: 'Task Queue Producer API is running'
  });
});

app.post('/job', async (req, res) => {
  try {
    const { type, data } = req.body;

    if (!type) {
      return res.status(400).json({ error: 'Field "type" is required' });
    }

    const job = await enqueueTask({ type, data });

    res.status(201).json({
      message: 'Job enqueued successfully',
      job
    });
  } catch (err) {
    console.error('Error enqueuing job:', err);
    res.status(500).json({ error: 'Internal server error' });
  }
});

app.listen(port, () => {
  console.log(`Producer API listening on port ${port}`);
});
