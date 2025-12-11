const Redis = require('ioredis');

const redis = new Redis({
  host: process.env.REDIS_HOST || 'redis',
  port: Number(process.env.REDIS_PORT) || 6379
});

const QUEUE_NAME = process.env.QUEUE_NAME || 'tasks_queue';

async function enqueueTask(payload) {
  const job = {
    id: Date.now(),
    payload,
    createdAt: new Date().toISOString()
  };

  await redis.lpush(QUEUE_NAME, JSON.stringify(job));
  return job;
}

async function dequeueTask() {
  // BRPOP = blocking pop from the right side of list
  const timeoutSeconds = 5;
  const res = await redis.brpop(QUEUE_NAME, timeoutSeconds);

  if (!res) return null; // timeout no job
  const [, jobStr] = res;
  return JSON.parse(jobStr);
}

module.exports = {
  enqueueTask,
  dequeueTask
};
