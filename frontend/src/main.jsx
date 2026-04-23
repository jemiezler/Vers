import React, { useState } from 'react';
import { createRoot } from 'react-dom/client';
import './styles.css';

function App() {
  const [filename, setFilename] = useState('go.mod');
  const [content, setContent] = useState('require github.com/gin-gonic/gin v1.10.0');
  const [result, setResult] = useState(null);
  const [error, setError] = useState('');

  async function submitReview(event) {
    event.preventDefault();
    setError('');
    setResult(null);

    const response = await fetch('/reviews', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ filename, content }),
    });
    const payload = await response.json();

    if (!response.ok) {
      setError(payload.error || 'Review failed');
      return;
    }

    setResult(payload);
  }

  return (
    <main className="shell">
      <section className="hero">
        <p className="eyebrow">Version-aware code review</p>
        <h1>Vers turns dependency manifests into review context.</h1>
        <p>
          Upload a manifest, retrieve version-specific docs, build an augmented prompt,
          and send it to a local review model.
        </p>
      </section>

      <form className="panel" onSubmit={submitReview}>
        <label>
          Manifest filename
          <input value={filename} onChange={(event) => setFilename(event.target.value)} />
        </label>
        <label>
          Manifest content
          <textarea value={content} onChange={(event) => setContent(event.target.value)} />
        </label>
        <button type="submit">Run scaffold review</button>
      </form>

      {error && <pre className="error">{error}</pre>}
      {result && <pre className="result">{JSON.stringify(result, null, 2)}</pre>}
    </main>
  );
}

createRoot(document.getElementById('root')).render(<App />);

