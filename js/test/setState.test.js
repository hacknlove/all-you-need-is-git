import { spawnSync } from 'node:child_process';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { test, expect } from 'vitest';
import { action } from '../commands/set-state.js';

const testDir = path.dirname(fileURLToPath(import.meta.url));
const projectDir = path.resolve(testDir, '..');

function expectedCommitMessage(prompt) {
  return `chore: set review\n\n${prompt}\n\ndwp-state: review\n\n`;
}

function runSetStateWithPrompt(prompt) {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'aynig-set-state-'));
  const run = (command, args) => spawnSync(command, args, {
    cwd: tempDir,
    encoding: 'utf8'
  });

  try {
    expect(run('git', ['init']).status).toBe(0);
    expect(run('git', ['config', 'user.name', 'AYNIG Test']).status).toBe(0);
    expect(run('git', ['config', 'user.email', 'aynig@example.com']).status).toBe(0);

    const result = spawnSync(process.execPath, [
      '--no-warnings',
      path.join(projectDir, 'index.js'),
      'set-state',
      '--dwp-state',
      'review',
      '--prompt',
      prompt
    ], {
      cwd: tempDir,
      encoding: 'utf8'
    });

    expect(result.status).toBe(0);
    expect(result.stderr).toBe('');

    const logResult = run('git', ['log', '-1', '--format=%B']);
    expect(logResult.status).toBe(0);
    return logResult.stdout;
  } finally {
    fs.rmSync(tempDir, { recursive: true, force: true });
  }
}

function runSetStateKeepingTrailers(prompt) {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'aynig-set-state-'));
  const run = (command, args) => spawnSync(command, args, {
    cwd: tempDir,
    encoding: 'utf8'
  });

  try {
    expect(run('git', ['init']).status).toBe(0);
    expect(run('git', ['config', 'user.name', 'AYNIG Test']).status).toBe(0);
    expect(run('git', ['config', 'user.email', 'aynig@example.com']).status).toBe(0);
    expect(run('git', ['commit', '--allow-empty', '-m', 'seed', '-m', [
      'body',
      '',
      'dwp-state: triage',
      'dwp-note: first',
      'dwp-note: second',
      'dwp-origin-state: queued',
      'custom: skip'
    ].join('\n')]).status).toBe(0);

    const result = spawnSync(process.execPath, [
      '--no-warnings',
      path.join(projectDir, 'index.js'),
      'set-state',
      '--dwp-state',
      'review',
      '--keep-trailers',
      '--prompt',
      prompt
    ], {
      cwd: tempDir,
      encoding: 'utf8'
    });

    expect(result.status).toBe(0);
    expect(result.stderr).toBe('');

    const logResult = run('git', ['log', '-1', '--format=%B']);
    expect(logResult.status).toBe(0);
    return logResult.stdout;
  } finally {
    fs.rmSync(tempDir, { recursive: true, force: true });
  }
}

test('set-state requires dwp-state', async () => {
  await expect(action({})).rejects.toThrow(/Missing required flag: --dwp-state/);
});

test('set-state rejects working', async () => {
  await expect(action({ dwpState: 'working' })).rejects.toThrow(/use aynig set-working/);
});

test('set-state preserves a one-line --prompt in the commit body', () => {
  const prompt = 'Una sola linea';
  expect(runSetStateWithPrompt(prompt)).toBe(expectedCommitMessage(prompt));
});

test('set-state preserves a multiline --prompt in the commit body', () => {
  const prompt = 'Primera linea\nSegunda linea\nTercera linea';
  expect(runSetStateWithPrompt(prompt)).toBe(expectedCommitMessage(prompt));
});

test('set-state preserves multiline --prompt lines that look like trailers', () => {
  const prompt = 'Contexto\nKEY: value\nMas contexto';
  expect(runSetStateWithPrompt(prompt)).toBe(expectedCommitMessage(prompt));
});

test('set-state preserves blank lines in multiline --prompt bodies', () => {
  const prompt = 'Primera linea\n\nSegunda linea\n\n';
  expect(runSetStateWithPrompt(prompt)).toBe(expectedCommitMessage(prompt));
});

test('set-state can keep existing aynig trailers from HEAD', () => {
  const prompt = 'Mantener trailers';
  expect(runSetStateKeepingTrailers(prompt)).toBe([
    'chore: set review',
    '',
    prompt,
    '',
    'dwp-state: review',
    'dwp-note: first',
    'dwp-note: second',
    'dwp-origin-state: queued',
    '',
    ''
  ].join('\n'));
});
