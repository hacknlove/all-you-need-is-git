import { spawn } from 'node:child_process';
import { parseGitTrailersStrict } from '../gitHelpers/parseTrailers.js';

export function buildStateCommitMessage(subject, prompt, trailers) {
  const cleanSubject = String(subject || '').replace(/\n+$/, '');
  const cleanPrompt = String(prompt || '').replace(/\n+$/, '');

  const lines = [cleanSubject, '', cleanPrompt, ''];
  for (const trailer of trailers || []) {
    lines.push(`${trailer.key}: ${trailer.value}`);
  }

  return `${lines.join('\n')}\n`;
}

export async function commitState(worktreeGit, subject, prompt, trailers) {
  const message = buildStateCommitMessage(subject, prompt, trailers);
  await worktreeGit.commit(message, { '--allow-empty': null });

  const headBody = await worktreeGit.raw(['show', '-s', '--format=%B', 'HEAD']);
  const trailersRaw = await interpretTrailers(headBody);
  const parsed = parseGitTrailersStrict(trailersRaw);
  const stateRaw = parsed['aynig-state'];
  const states = Array.isArray(stateRaw) ? stateRaw : stateRaw ? [stateRaw] : [];
  if (states.length === 0) {
    throw new Error('Invalid trailer block: trailers must be contiguous at end of message');
  }
  if (states.length > 1) {
    throw new Error('Invalid trailer block: multiple aynig-state trailers');
  }
  if (!String(states[0]).trim()) {
    throw new Error('Invalid trailer block: empty aynig-state trailer');
  }
}

function interpretTrailers(body) {
  return new Promise((resolve, reject) => {
    const child = spawn('git', ['interpret-trailers', '--parse', '--only-trailers']);
    let stdout = '';
    let stderr = '';

    child.stdout.on('data', (chunk) => {
      stdout += chunk;
    });
    child.stderr.on('data', (chunk) => {
      stderr += chunk;
    });

    child.on('error', reject);
    child.on('close', (code) => {
      if (code === 0) {
        resolve(stdout);
        return;
      }
      const suffix = stderr.trim() ? `: ${stderr.trim()}` : '';
      reject(new Error(`git interpret-trailers failed with code ${code}${suffix}`));
    });

    child.stdin.on('error', (error) => {
      if (error?.code !== 'EPIPE') {
        reject(error);
      }
    });

    child.stdin.end(body);
  });
}
