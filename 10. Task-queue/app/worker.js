const { dequeueTask } = require('./queue');

console.log('Worker started. Waiting for jobs...');

async function processJob(job) {
  // Here you can do anything: send email, generate report, etc.
  console.log('Processing job:', job);

  // Simulate some work
  await new Promise(resolve => setTimeout(resolve, 1000));

  console.log(`Job ${job.id} done.`);
}

async function runWorkerLoop() {
  while (true) {
    try {
      const job = await dequeueTask();

      if (!job) {
        // No job received within timeout
        continue;
      }

      await processJob(job);
    } catch (err) {
      console.error('Error in worker loop:', err);
      // small delay to avoid tight error loop
      await new Promise(resolve => setTimeout(resolve, 1000));
    }
  }
}

runWorkerLoop();
