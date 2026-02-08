import { git } from '../gitHelpers/git.js';
import { resolve } from 'path';
import { homedir } from 'os';
import { cwd } from 'process';

export class Command {
    constructor({
        config,
        branchName,
        isCurrentBranch,
        command,
        trailers,
        body
    }) {
        this.config = config;
        this.branchName = branchName;
        this.isCurrentBranch = isCurrentBranch;
        this.command = command;
        this.trailers = trailers;
        this.body = body;
    }

    async run() {
        if (!this.command) {
            return;
        }

        // Current branch runs in main workspace; other branches get a worktree
        let worktreePath;
        if (this.isCurrentBranch) {
            worktreePath = cwd();
        } else {
            worktreePath = resolve(this.config.dir, this.config.workTree || '.', `worktree-${this.branchName.replace(/\//g, '_')}`);
            await git.worktree(['add', worktreePath, this.branchName]);
        }

        // Search for command in configured paths relative to repo, cwd, and home
        const searchRoots = [
            this.config.dir,           // repo directory
            cwd(),                      // current working directory
        ];

        // TODO: Complete command execution logic
    }
}