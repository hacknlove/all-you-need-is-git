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

async function action(options) {
  try {
    const historyEnabled = Boolean(options.history);
    let limit = historyEnabled ? Number.parseInt(options.limit, 10) : 1;
    if (!Number.isFinite(limit) || limit < 1) {
      limit = 1;
    }

    const raw = await git.raw([
      'log',
      '-n',
      String(limit),
      '--format=%H%x1f%cI%x1f%B%x1e'
    ]);

    const records = raw.split('\x1e').filter(Boolean);
    const events = [];

    for (const record of records) {
      const parts = record.split('\x1f');
      const commit = parts[0] || '';
      const date = parts[1] || '';
      const message = parts.slice(2).join('\x1f');
      const { firstLine, body } = splitCommitMessage(message);
      const { trailers } = await parseCommitMessage({ message: firstLine, body });
      const trailersLower = toLowerCaseKeys(trailers);
      const state = normalizeTrailerValue(trailersLower['aynig-state']);
      const runId = normalizeTrailerValue(trailersLower['aynig-run-id']);
      const originState = normalizeTrailerValue(trailersLower['aynig-origin-state']);

      events.push({
        commit,
        date,
        message: firstLine,
        state: state || null,
        runId: runId || null,
        originState: originState || null,
        trailers
      });
    }

    if (options.json) {
      console.log(JSON.stringify(events, null, 2));
      return;
    }

    if (events.length === 0) {
      console.log('No commits found.');
      return;
    }

    for (const event of events) {
      const shortCommit = event.commit ? event.commit.slice(0, 7) : 'unknown';
      const state = event.state || 'n/a';
      const runId = event.runId || 'n/a';
      const origin = event.originState ? ` origin=${event.originState}` : '';
      console.log(`${shortCommit} ${event.date} state=${state} run=${runId}${origin} ${event.message}`);
    }
  } catch (error) {
    console.error('Error: Not a Git repository. Please run `git init` first.');
    process.exit(1);
  }
}

export function registerEventsCommand(program) {
  program
    .command('events')
    .description('Show recent AYNIG events from commit trailers')
    .option('--history', 'Allow scanning recent history')
    .option('-n, --limit <number>', 'Number of commits to inspect (requires --history)', '10')
    .option('--json', 'Output JSON')
    .action(action);
}
