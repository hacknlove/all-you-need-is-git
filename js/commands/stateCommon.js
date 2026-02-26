import fs from 'node:fs/promises';
import { git } from '../gitHelpers/git.js';
import { parseCommitMessage } from '../gitHelpers/parseCommitMessage.js';

export async function readHeadCommitMessage() {
  const raw = await git.raw(['show', '-s', '--format=%B', 'HEAD']);
  const normalized = raw.replace(/\r\n/g, '\n');
  const lines = normalized.split('\n');
  const firstLine = lines.shift() || '';
  let body = lines.join('\n');
  body = body.replace(/^\n+/, '');
  return parseCommitMessage({ message: firstLine, body });
}

export async function readPrompt({ prompt, promptFile, promptStdin }, fallback = '') {
  const sources = [
    Boolean(prompt && String(prompt).trim() !== ''),
    Boolean(promptFile && String(promptFile).trim() !== ''),
    Boolean(promptStdin)
  ].filter(Boolean).length;

  if (sources > 1) {
    throw new Error('Only one prompt source is allowed: --prompt | --prompt-file | --prompt-stdin');
  }

  if (prompt && String(prompt).trim() !== '') {
    return String(prompt);
  }
  if (promptFile && String(promptFile).trim() !== '') {
    return fs.readFile(promptFile, 'utf8');
  }
  if (promptStdin) {
    const chunks = [];
    for await (const chunk of process.stdin) {
      chunks.push(chunk);
    }
    return Buffer.concat(chunks.map(c => Buffer.isBuffer(c) ? c : Buffer.from(String(c)))).toString('utf8');
  }

  return fallback;
}

export function parseTrailerArg(raw) {
  const idx = String(raw).indexOf(':');
  if (idx <= 0) {
    throw new Error(`Invalid trailer format: "${raw}" (expected key:value)`);
  }
  const key = String(raw).slice(0, idx).trim();
  const value = String(raw).slice(idx + 1).trim();
  if (!key) {
    throw new Error(`Invalid trailer format: "${raw}" (empty key)`);
  }
  return { key, value };
}

export function trailerValue(trailers, key) {
  const target = String(key).toLowerCase();
  for (const [k, value] of Object.entries(trailers || {})) {
    if (k.toLowerCase() !== target) {
      continue;
    }
    if (Array.isArray(value)) {
      return String(value[value.length - 1] || '');
    }
    return String(value || '');
  }
  return '';
}

export function resolveAynigRemote(cliRemote, headTrailers) {
  if (cliRemote && String(cliRemote).trim() !== '') {
    return String(cliRemote).trim();
  }
  return trailerValue(headTrailers, 'aynig-remote').trim();
}

export async function validateRemoteExists(remote) {
  if (!remote) {
    return;
  }
  try {
    await git.raw(['remote', 'get-url', remote]);
  } catch {
    throw new Error(`Unknown remote "${remote}"`);
  }
}

export async function pushCurrentBranch(remote) {
  if (!remote) {
    return;
  }
  const branch = (await git.revparse(['--abbrev-ref', 'HEAD'])).trim();
  await git.push(remote, branch);
}
