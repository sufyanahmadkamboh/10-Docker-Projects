const API_BASE = 'http://localhost:3000';

const catsCountEl = document.getElementById('cats-count');
const dogsCountEl = document.getElementById('dogs-count');
const totalCountEl = document.getElementById('total-count');

async function fetchResults() {
  try {
    const res = await fetch(`${API_BASE}/results`);
    const data = await res.json();
    const { cats, dogs } = data.votes;

    catsCountEl.textContent = cats;
    dogsCountEl.textContent = dogs;
    totalCountEl.textContent = data.total;
  } catch (err) {
    console.error('Error fetching results:', err);
  }
}

async function sendVote(option) {
  try {
    await fetch(`${API_BASE}/vote`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ option })
    });

    // Refresh results after voting
    fetchResults();
  } catch (err) {
    console.error('Error sending vote:', err);
  }
}

document.getElementById('vote-cats').addEventListener('click', () => {
  sendVote('cats');
});

document.getElementById('vote-dogs').addEventListener('click', () => {
  sendVote('dogs');
});

// Fetch results every 3 seconds
setInterval(fetchResults, 3000);
fetchResults();
