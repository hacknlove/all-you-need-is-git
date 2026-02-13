import { git } from '../gitHelpers/git.js';
import simpleGit from 'simple-git';
import { resolve } from 'path';
import { access, constants } from 'fs/promises';
import { createHash, randomUUID } from 'crypto';
import { cwd } from 'process';
import { spawn } from 'child_process';
import { hostname } from 'os';

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
        gitFactory,
        spawnImpl
    }) {
        this.config = config;
        this.branchName = branchName;
        this.isCurrentBranch = isCurrentBranch;
        this.command = trailers['aynig-state']?.trim().toLowerCase();
        this.trailers = trailers;
        this.body = body;
        this.commitDate = commitDate;
        this.gitFactory = gitFactory || simpleGit;
        this.spawnImpl = spawnImpl || spawn;
    }

    async findExistingWorktree() {
        const worktrees = await getWorktreeList();
        const ref = `refs/heads/${this.branchName}`;
        const match = worktrees.find(wt => wt.branch === ref);
        return match?.path;
    }

    async getWorkspace() {
        if (this.isCurrentBranch) {
            return cwd();
        }

        // Check if a worktree already exists for this branch
        const existing = await this.findExistingWorktree();
        if (existing) {
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
            if (this.config.useRemote) {
                await git.worktree(['add', '-b', this.branchName, worktreePath, `${this.config.useRemote}/${this.branchName}`]);
            } else {
                await git.worktree(['add', worktreePath, this.branchName]);
            }
        } catch {
            console.warn(`Failed to create worktree for branch ${this.branchName}`);
            return null;
        }

        return worktreePath;
    }

    async getCommandPath(worktreePath) {
        const baseDir = resolve(worktreePath, '.aynig', 'command');
        const commandName = this.command;

        if (!commandName) {
            return false;
        }

        const commandPath = resolve(baseDir, commandName);
        if (!commandPath.startsWith(`${baseDir}/`)) {
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

        const stalledRun = trailers['aynig-run-id']?.trim() || 'unknown';
        const worktreePath = await this.getWorkspace();
        if (!worktreePath) {
            return;
        }

        const worktreeGit = this.gitFactory(worktreePath);

        await worktreeGit.commit(`chore: stalled

Lease expired

aynig-state: stalled
aynig-stalled-run: ${stalledRun}
`, { '--allow-empty': null });

        if (this.config.useRemote) {
            try {
                await worktreeGit.push(this.config.useRemote, this.branchName);
            } catch {
                return;
            }
        }
    }

    async run() {
        if (!this.command) {
            return;
        }

        if (this.command === 'working') {
            return this.checkWorking();
        }

        const worktreePath = await this.getWorkspace();
        if (!worktreePath) {
            return;
        }

        const commandPath = await this.getCommandPath(worktreePath);
        if (!commandPath) {
            return;
        }

        const leaseSeconds = this.config.leaseSeconds || 300;
        const runId = randomUUID();
        const runnerId = hostname();

        const worktreeGit = this.gitFactory(worktreePath);
        const currentCommitHash = (await worktreeGit.revparse(['HEAD'])).trim();

        await worktreeGit.commit(`chore: working

command ${this.command} takes control of the branch

aynig-state: working
aynig-run-id: ${runId}
aynig-runner-id: ${runnerId}
aynig-lease-seconds: ${leaseSeconds}
`, { '--allow-empty': null });

        if (this.config.useRemote) {
            try {
                await worktreeGit.push(this.config.useRemote, this.branchName);
            } catch {
                return;
            }
        }

        const env = {
            ...process.env,
            AYNIG_BODY: this.body,
            AYNIG_COMMIT_HASH: currentCommitHash
        };
        for (const [key, value] of Object.entries(this.trailers)) {
            env[`AYNIG_TRAILER_${key.toUpperCase()}`] = value;
        }

        const child = this.spawnImpl(commandPath, [], {
            cwd: worktreePath,
            env,
            detached: true,
            stdio: 'ignore'
        });
        child.unref();
    }
}
