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

  let originState = trailerValue(headTrailers, 'dwp-origin-state').trim();
  if (!originState) {
    originState = trailerValue(headTrailers, 'dwp-state').trim();
  }
  if (!originState) {
    throw new Error('Missing required trailer: dwp-state');
  }

  const runId = trailerValue(headTrailers, 'dwp-run-id').trim() || `run-${randomUUID().replace(/-/g, '')}`;
  const leaseRaw = trailerValue(headTrailers, 'dwp-lease-seconds').trim();
  const leaseSeconds = Number.parseInt(leaseRaw, 10);
  const defaultLease = Number.isFinite(leaseSeconds) && leaseSeconds > 0 ? leaseSeconds : defaultConfig.leaseSeconds;
  const cliLease = Number.parseInt(options.leaseSeconds, 10);
  const lease = Number.isFinite(cliLease) && cliLease > 0 ? cliLease : defaultLease;

  const trailers = [
    { key: 'dwp-state', value: 'working' },
    { key: 'dwp-origin-state', value: originState },
    { key: 'dwp-run-id', value: runId },
    { key: 'dwp-runner-id', value: hostname() },
    { key: 'dwp-lease-seconds', value: String(lease) }
  ];
  if (remote) {
    trailers.push({ key: 'dwp-source', value: `git:${remote}` });
  }

  const reserved = new Set([
    'dwp-state',
    'dwp-origin-state',
    'dwp-run-id',
    'dwp-runner-id',
    'dwp-lease-seconds',
    'dwp-source'
  ]);
  const copied = Object.keys(headTrailers)
    .filter(key => key.toLowerCase().startsWith('dwp-'))
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
    .option('--lease-seconds <seconds>', 'Lease duration in seconds (overrides dwp-lease-seconds trailer)')
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
