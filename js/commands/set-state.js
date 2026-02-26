import { commitState } from '../AgentsOrchestrator/stateCommit.js';
import { git } from '../gitHelpers/git.js';
import {
  parseTrailerArg,
  pushCurrentBranch,
  readHeadCommitMessage,
  readPrompt,
  resolveAynigRemote,
  validateRemoteExists
} from './stateCommon.js';

export async function action(options) {
  const state = String(options.aynigState || '').trim().toLowerCase();
  if (!state) {
    throw new Error('Missing required flag: --aynig-state');
  }
  if (state === 'working') {
    throw new Error('Invalid aynig-state: working (use aynig set-working)');
  }

  const parsed = await readHeadCommitMessage();
  const headTrailers = parsed.trailers || {};

  const remote = resolveAynigRemote(options.aynigRemote, headTrailers);
  await validateRemoteExists(remote);

  const prompt = await readPrompt(options, '');
  const subject = (options.subject || '').trim() || `chore: set ${state}`;

  const trailers = [{ key: 'aynig-state', value: state }];
  if (remote) {
    trailers.push({ key: 'aynig-remote', value: remote });
  }

  for (const raw of options.trailer || []) {
    const parsedTrailer = parseTrailerArg(raw);
    if (parsedTrailer.key.trim().toLowerCase() === 'aynig-state') {
      throw new Error('Invalid trailer: aynig-state is managed by --aynig-state');
    }
    trailers.push(parsedTrailer);
  }

  await commitState(git, subject, prompt, trailers);
  await pushCurrentBranch(remote);
}

export function registerSetStateCommand(program) {
  program
    .command('set-state')
    .description('Create a commit with a new non-working aynig-state')
    .requiredOption('--aynig-state <state>', 'Next state (must not be working)')
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
