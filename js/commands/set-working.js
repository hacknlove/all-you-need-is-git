import { randomUUID } from 'node:crypto';
import { hostname } from 'node:os';
import defaultConfig from '../defaultConfig.json' with { type: 'json' };
import { commitState } from '../AgentsOrchestrator/stateCommit.js';
import { git } from '../gitHelpers/git.js';
import {
  parseTrailerArg,
  pushCurrentBranch,
  readHeadCommitMessage,
  readPrompt,
  resolveAynigRemote,
  trailerValue,
  validateRemoteExists
} from './stateCommon.js';

export async function action(options) {
  const parsed = await readHeadCommitMessage();
  const headTrailers = parsed.trailers || {};

  const remote = resolveAynigRemote(options.aynigRemote, headTrailers);
  await validateRemoteExists(remote);

  const prompt = await readPrompt(options, 'Lease heartbeat');
  const subject = (options.subject || '').trim() || 'chore: working';

  let originState = trailerValue(headTrailers, 'aynig-origin-state').trim();
  if (!originState) {
    originState = trailerValue(headTrailers, 'aynig-state').trim();
  }
  if (!originState) {
    throw new Error('Missing required trailer: aynig-state');
  }

  const runId = trailerValue(headTrailers, 'aynig-run-id').trim() || `run-${randomUUID().replace(/-/g, '')}`;
  const leaseRaw = trailerValue(headTrailers, 'aynig-lease-seconds').trim();
  const leaseSeconds = Number.parseInt(leaseRaw, 10);
  const lease = Number.isFinite(leaseSeconds) && leaseSeconds > 0 ? leaseSeconds : defaultConfig.leaseSeconds;

  const trailers = [
    { key: 'aynig-state', value: 'working' },
    { key: 'aynig-origin-state', value: originState },
    { key: 'aynig-run-id', value: runId },
    { key: 'aynig-runner-id', value: hostname() },
    { key: 'aynig-lease-seconds', value: String(lease) }
  ];
  if (remote) {
    trailers.push({ key: 'aynig-remote', value: remote });
  }

  const reserved = new Set([
    'aynig-state',
    'aynig-origin-state',
    'aynig-run-id',
    'aynig-runner-id',
    'aynig-lease-seconds',
    'aynig-remote'
  ]);
  const copied = Object.keys(headTrailers)
    .filter(key => key.toLowerCase().startsWith('aynig-'))
    .filter(key => !reserved.has(key.toLowerCase()))
    .sort();
  for (const key of copied) {
    const value = trailerValue(headTrailers, key).trim();
    if (!value) {
      continue;
    }
    trailers.push({ key: key.toLowerCase(), value });
  }

  for (const raw of options.trailer || []) {
    trailers.push(parseTrailerArg(raw));
  }

  await commitState(git, subject, prompt, trailers);
  await pushCurrentBranch(remote);
}

export function registerSetWorkingCommand(program) {
  program
    .command('set-working')
    .description('Create or refresh a working lease commit')
    .option('--subject <text>', 'Commit title')
    .option('--prompt <text>', 'Commit prompt/body')
    .option('--prompt-file <path>', 'Path to file used as prompt/body')
    .option('--prompt-stdin', 'Read prompt/body from stdin')
    .option('--aynig-remote <name>', 'Remote name to push after commit')
    .option('--trailer <key:value>', 'Additional trailer (repeatable)', (value, acc) => {
      acc.push(value);
      return acc;
    }, [])
    .action(async (options) => {
      try {
        await action(options);
      } catch (error) {
        console.error(error.message || error);
        process.exit(1);
      }
    });
}
