import { format } from 'util';

const levels = {
  debug: 0,
  info: 1,
  warn: 2,
  error: 3
};

function normalizeLevel(level) {
  const normalized = String(level || '').trim().toLowerCase();
  if (normalized === 'warning') {
    return 'warn';
  }
  return Object.prototype.hasOwnProperty.call(levels, normalized) ? normalized : 'error';
}

export class Logger {
  constructor(level) {
    this.level = normalizeLevel(level);
  }

  shouldLog(level) {
    return levels[this.level] <= levels[level];
  }

  debug(...args) {
    if (this.shouldLog('debug')) {
      console.log(`debug: ${format(...args)}`);
    }
  }

  info(...args) {
    if (this.shouldLog('info')) {
      console.log(`info: ${format(...args)}`);
    }
  }

  warn(...args) {
    if (this.shouldLog('warn')) {
      console.warn(`warn: ${format(...args)}`);
    }
  }

  error(...args) {
    if (this.shouldLog('error')) {
      console.error(`error: ${format(...args)}`);
    }
  }
}
