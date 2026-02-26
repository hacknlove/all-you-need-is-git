import { git } from '../gitHelpers/git.js';
import simpleGit from 'simple-git';
import { isAbsolute, relative, resolve, sep } from 'path';
import { access, constants, mkdir } from 'fs/promises';
import { closeSync, openSync } from 'node:fs';
import { createHash, randomUUID } from 'crypto';
import { cwd } from 'process';
import { spawn } from 'child_process';
import { hostname } from 'os';
import { Logger } from '../utils/logger.js';

function branchHash(branchName) {
    return createHash('sha256').update(branchName).digest('hex').slice(0, 8);
}

async function getWorktreeList() {
    const output = await git.worktree(['list', '--porcelain']);

    const worktrees = [];
    let current = {};
    for (const line of output.split('\n')) {
        if (line === '') {
            if (current.path) {
                worktrees.push(current);
            }
            current = {};
            continue;
        }
        if (line.startsWith('worktree ')) {
            current.path = line.slice('worktree '.length);
        } else if (line.startsWith('branch ')) {
            current.branch = line.slice('branch '.length);
        }
    }
    if (current.path) {
        worktrees.push(current);
    }
    return worktrees;
}

export class Command {
    constructor({
        config,
        branchName,
        isCurrentBranch,
        trailers,
        body,
        commitDate,
        logger,
        logLevel,
        gitFactory,
        spawnImpl
    }) {
        this.config = config;
        this.branchName = branchName;
        this.isCurrentBranch = isCurrentBranch;
        const { state, error } = resolveStateTrailer(trailers);
        this.command = state;
        this.invalidState = error;
        this.trailers = trailers;
        this.body = body;
        this.commitDate = commitDate;
        this.gitFactory = gitFactory || simpleGit;
        this.spawnImpl = spawnImpl || spawn;
        this.logLevel = logLevel || config?.logLevel;
        this.logger = logger || config?.logger || new Logger(this.logLevel);
    }

    async findExistingWorktree() {
        const worktrees = await getWorktreeList();
        const ref = `refs/heads/${this.branchName}`;
        const match = worktrees.find(wt => wt.branch === ref);
        return match?.path;
    }

    async getWorkspace() {
        if (this.isCurrentBranch) {
            this.logger.debug('Using current working directory for %s', this.branchName);
            return cwd();
        }

        // Check if a worktree already exists for this branch
        const existing = await this.findExistingWorktree();
        if (existing) {
            this.logger.debug('Using existing worktree for %s', this.branchName);
            return existing;
        }

        // Build worktree path anchored to repo root, with hash to avoid collisions
        const safeName = this.branchName.replace(/\//g, '_');
        const hash = branchHash(this.branchName);
        const worktreePath = resolve(
            this.config.repoRoot,
            this.config.workTree || '.',
            `worktree-${safeName}-${hash}`
        );

        try {
            if (this.config.aynigRemote) {
                await git.worktree(['add', '-b', this.branchName, worktreePath, `${this.config.aynigRemote}/${this.branchName}`]);
            } else {
                await git.worktree(['add', worktreePath, this.branchName]);
            }
        } catch {
            this.logger.warn('Failed to create worktree for branch %s', this.branchName);
            return null;
        }

        this.logger.debug('Created worktree for %s at %s', this.branchName, worktreePath);

        return worktreePath;
    }

    async getCommandPath(worktreePath) {
        const baseDir = resolve(worktreePath, '.aynig', 'command');
        const commandName = this.command;

        if (!commandName) {
            return false;
        }

        const commandPath = resolve(baseDir, commandName);
        const rel = relative(baseDir, commandPath);
        if (!rel || isAbsolute(rel) || rel === '..' || rel.startsWith(`..${sep}`)) {
            return false;
        }

        try {
            await access(commandPath, constants.X_OK);
            return commandPath;
        } catch {
            return false;
        }
    }

    async checkWorking() {
        const trailers = this.trailers || {};

        const leaseSeconds = Number(trailers['aynig-lease-seconds']);
        if (!Number.isFinite(leaseSeconds) || leaseSeconds <= 0) {
            return;
        }

        const committedAt = new Date(this.commitDate).getTime();
        if (!Number.isFinite(committedAt)) {
            return;
        }

        const expired = Date.now() > committedAt + (leaseSeconds * 1000);
        if (!expired) {
            return;
        }

        this.logger.info('Lease expired for branch %s', this.branchName);

        const stalledRun = trailers['aynig-run-id']?.trim() || 'unknown';
        const worktreePath = await this.getWorkspace();
        if (!worktreePath) {
            return;
        }

        const worktreeGit = this.gitFactory(worktreePath);

        let stalledMessage = `chore: stalled

Lease expired

aynig-state: stalled
aynig-stalled-run: ${stalledRun}
`;
        if (this.config.aynigRemote) {
            stalledMessage += `aynig-remote: ${this.config.aynigRemote}\n`;
        }
        await worktreeGit.commit(stalledMessage, { '--allow-empty': null });

        if (this.config.aynigRemote) {
            try {
                this.logger.info('Pushing stalled state for %s to %s', this.branchName, this.config.aynigRemote);
                await worktreeGit.push(this.config.aynigRemote, this.branchName);
            } catch {
                return;
            }
        }
    }

    async run() {
        if (this.invalidState) {
            this.logger.warn('Skipping branch %s (%s)', this.branchName, this.invalidState);
            return;
        }
        if (!this.command) {
            this.logger.debug('Skipping branch %s (no aynig-state)', this.branchName);
            return;
        }

        if (this.command === 'working') {
            this.logger.debug('Checking lease on branch %s', this.branchName);
            return this.checkWorking();
        }

        this.logger.info('Running command %s on branch %s', this.command, this.branchName);

        const worktreePath = await this.getWorkspace();
        if (!worktreePath) {
            return;
        }

        const commandPath = await this.getCommandPath(worktreePath);
        if (!commandPath) {
            this.logger.debug('Command path not found for %s', this.command);
            return;
        }

        this.logger.debug('Command path: %s', commandPath);

        const leaseSeconds = this.config.leaseSeconds || 300;
        const runId = randomUUID();
        const runnerId = hostname();
        const originState = this.command;

        const worktreeGit = this.gitFactory(worktreePath);
        const currentCommitHash = (await worktreeGit.revparse(['HEAD'])).trim();

        let workingMessage = `chore: working

command ${this.command} takes control of the branch

aynig-state: working
aynig-origin-state: ${originState}
aynig-run-id: ${runId}
aynig-runner-id: ${runnerId}
aynig-lease-seconds: ${leaseSeconds}
`;
        if (this.config.aynigRemote) {
            workingMessage += `aynig-remote: ${this.config.aynigRemote}\n`;
        }
        await worktreeGit.commit(workingMessage, { '--allow-empty': null });

        if (this.config.aynigRemote) {
            try {
                this.logger.info('Pushing branch %s to %s', this.branchName, this.config.aynigRemote);
                await worktreeGit.push(this.config.aynigRemote, this.branchName);
            } catch {
                return;
            }
        }

        const env = {
            ...process.env,
            AYNIG_BODY: this.body,
            AYNIG_COMMIT_HASH: currentCommitHash
        };
        if (this.logLevel) {
            env.AYNIG_LOG_LEVEL = this.logLevel;
        }
        for (const [key, value] of Object.entries(this.trailers)) {
            const normalizedKey = key.replace(/-/g, '_').toUpperCase();
            const envValue = Array.isArray(value) ? value.join(',') : value;
            env[`AYNIG_TRAILER_${normalizedKey}`] = envValue;
        }

        const logsDir = resolve(worktreePath, '.aynig', 'logs');
        await mkdir(logsDir, { recursive: true });
        const logPath = resolve(logsDir, `${currentCommitHash}.log`);
        const logFd = openSync(logPath, 'w');

        let child;
        try {
            child = this.spawnImpl(commandPath, [], {
                cwd: worktreePath,
                env,
                detached: true,
                stdio: ['ignore', logFd, logFd]
            });
        } finally {
            closeSync(logFd);
        }
        child.unref();
        this.logger.info('Launched %s in %s', this.command, worktreePath);
        this.logger.debug('Command log: %s', logPath);
    }
}

function resolveStateTrailer(trailers = {}) {
    const raw = trailers['aynig-state'];
    if (raw === undefined || raw === null) {
        return { state: '', error: '' };
    }

    if (Array.isArray(raw)) {
        if (raw.length === 0) {
            return { state: '', error: '' };
        }
        if (raw.length > 1) {
            return { state: '', error: 'multiple aynig-state trailers' };
        }
        const state = String(raw[0]).trim().toLowerCase();
        if (!state) {
            return { state: '', error: 'empty aynig-state trailer' };
        }
        return { state, error: '' };
    }

    const state = String(raw).trim().toLowerCase();
    if (!state) {
        return { state: '', error: 'empty aynig-state trailer' };
    }
    return { state, error: '' };
}
