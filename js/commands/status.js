import fs from 'fs/promises';
import { constants } from 'fs';
import path from 'path';
import { git } from '../gitHelpers/git.js';
import { parseCommitMessage } from '../gitHelpers/parseCommitMessage.js';

function splitCommitMessage(fullMessage) {
  const normalized = fullMessage.replace(/\r\n/g, '\n');
  const lines = normalized.split('\n');
  const firstLine = lines.shift() || '';
  let body = lines.join('\n');
  body = body.replace(/^\n+/, '');
  return { firstLine, body };
}

function trailerValue(trailers, key) {
  const target = String(key).trim().toLowerCase();
  for (const [k, value] of Object.entries(trailers || {})) {
    if (String(k).trim().toLowerCase() !== target) {
      continue;
    }
    if (Array.isArray(value)) {
      return String(value[value.length - 1] || '').trim();
    }
    return String(value || '').trim();
  }
  return '';
}

function leaseStatusForState(state, leaseSecondsRaw, committerDate) {
  if (state !== 'working') {
    return 'n/a';
  }
  const leaseSeconds = Number.parseInt(leaseSecondsRaw, 10);
  if (!Number.isFinite(leaseSeconds)) {
    return 'unknown';
  }
  const lastCommitTime = new Date(committerDate).getTime();
  if (Number.isNaN(lastCommitTime)) {
    return 'unknown';
  }
  const expiresAt = lastCommitTime + leaseSeconds * 1000;
  return Date.now() > expiresAt ? 'expired' : 'active';
}

async function action(options) {
  try {
    const repoRoot = (await git.revparse(['--show-toplevel'])).trim();
    const branch = (await git.revparse(['--abbrev-ref', 'HEAD'])).trim();
    const headCommit = (await git.revparse(['HEAD'])).trim();
    const committerDate = (await git.raw(['log', '-1', '--format=%cI'])).trim();
    const fullMessage = await git.raw(['log', '-1', '--format=%B']);

    const { firstLine, body } = splitCommitMessage(fullMessage);
    const { trailers } = await parseCommitMessage({ message: firstLine, body });

    const state = trailerValue(trailers, 'aynig-state');
    const runId = trailerValue(trailers, 'aynig-run-id');
    const leaseSecondsRaw = trailerValue(trailers, 'aynig-lease-seconds');
    const originState = trailerValue(trailers, 'aynig-origin-state');
    const leaseStatus = leaseStatusForState(state, leaseSecondsRaw, committerDate);

    let commandStatus = 'missing';
    let commandPath = '';
    let commandState = state;
    let shouldResolveCommand = true;

    if (state === 'working' && originState) {
      commandState = originState;
    } else if (state === 'working') {
      shouldResolveCommand = false;
    }

  if (shouldResolveCommand && commandState && commandState !== 'working') {
    const roleName = String(process.env.AYNIG_ROLE || '').trim();
    if (roleName) {
      const rolePath = path.join(repoRoot, '.aynig', 'roles', roleName, 'command', commandState);
      try {
        await fs.access(rolePath, constants.X_OK);
        commandStatus = 'exists';
        commandPath = rolePath;
      } catch {
        commandStatus = 'missing';
        commandPath = rolePath;
      }
    }
    if (!commandPath) {
      commandPath = path.join(repoRoot, '.aynig', 'command', commandState);
      try {
        await fs.access(commandPath, constants.X_OK);
        commandStatus = 'exists';
      } catch {
        commandStatus = 'missing';
      }
    }
  } else if (!shouldResolveCommand) {
    commandStatus = 'lease';
  }

    console.log(`branch: ${branch}`);
    console.log(`head: ${headCommit}`);
    console.log(`aynig-state: ${state || 'n/a'}`);
    if (state === 'working' && originState) {
      console.log(`aynig-origin-state: ${originState}`);
    }
    console.log(`aynig-run-id: ${runId || 'n/a'}`);
    console.log(`lease: ${leaseStatus}`);
    console.log(`command: ${commandStatus}`);
    if (commandPath) {
      console.log(`command-path: ${commandPath}`);
    }
  } catch (error) {
    console.error('Error: Not a Git repository. Please run `git init` first.');
    process.exit(1);
  }
}

export function registerStatusCommand(program) {
  program
    .command('status')
    .description('Show the current AYNIG state for this branch')
    .action(action);
}
