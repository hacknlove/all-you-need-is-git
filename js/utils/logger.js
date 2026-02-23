import { format } from 'util';

const levels = {
  debug: 0,
  info: 1,
  warn: 2,
  error: 3
};

export function normalizeLevel(level) {
  const normalized = String(level || '').trim().toLowerCase();
  if (normalized === 'warning') {
    return 'warn';
  }
  return Object.prototype.hasOwnProperty.call(levels, normalized) ? normalized : null;
}

export function resolveLevel({
  cliLevel,
  cliSet,
  trailerLevel,
  envLevel,
  defaultLevel
}) {
  if (cliSet) {
    const normalized = normalizeLevel(cliLevel);
    if (normalized) {
      return normalized;
    }
  }

  if (trailerLevel) {
    const normalized = normalizeLevel(trailerLevel);
    if (normalized) {
      return normalized;
    }
  }

  if (envLevel) {
    const normalized = normalizeLevel(envLevel);
    if (normalized) {
      return normalized;
    }
  }

  return normalizeLevel(defaultLevel) || 'error';
}

export class Logger {
  constructor(level) {
    this.level = normalizeLevel(level) || 'error';
  }

  shouldLog(level) {
    return levels[this.level] <= levels[level];
  }

  log(level, ...args) {
    if (!this.shouldLog(level)) {
      return;
    }
    const message = format(...args);
    if (level === 'warn') {
      console.warn(`warn: ${message}`);
      return;
    }
    if (level === 'error') {
      console.error(`error: ${message}`);
      return;
    }
    console.log(`${level}: ${message}`);
  }

  debug(...args) {
    this.log('debug', ...args);
  }

  info(...args) {
    this.log('info', ...args);
  }

  warn(...args) {
    this.log('warn', ...args);
  }

  error(...args) {
    this.log('error', ...args);
  }
}

export class BufferedLogger {
  constructor() {
    this.entries = [];
  }

  debug(...args) {
    this.entries.push({ level: 'debug', args });
  }

  info(...args) {
    this.entries.push({ level: 'info', args });
  }

  warn(...args) {
    this.entries.push({ level: 'warn', args });
  }

  error(...args) {
    this.entries.push({ level: 'error', args });
  }

  flush(logger) {
    for (const entry of this.entries) {
      logger.log(entry.level, ...entry.args);
    }
    this.entries = [];
  }
}
