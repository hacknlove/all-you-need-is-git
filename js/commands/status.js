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

function normalizeTrailerValue(value) {
  if (Array.isArray(value)) {
    return value[value.length - 1];
  }
  return value;
}

function toLowerCaseKeys(trailers) {
  const out = {};
  for (const [key, value] of Object.entries(trailers || {})) {
    out[key.toLowerCase()] = value;
  }
  return out;
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
    const trailersLower = toLowerCaseKeys(trailers);

    const state = normalizeTrailerValue(trailersLower['aynig-state']);
    const runId = normalizeTrailerValue(trailersLower['aynig-run-id']);
    const leaseSecondsRaw = normalizeTrailerValue(trailersLower['aynig-lease-seconds']);
    const originState = normalizeTrailerValue(trailersLower['aynig-origin-state']);
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
      commandPath = path.join(repoRoot, '.aynig', 'command', commandState);
      try {
        await fs.access(commandPath, constants.X_OK);
        commandStatus = 'exists';
      } catch {
        commandStatus = 'missing';
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
