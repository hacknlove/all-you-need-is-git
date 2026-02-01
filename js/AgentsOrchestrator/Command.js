import { git } from '../gitHelpers/git.js';
import { resolve } from 'path';
import { homedir } from 'os';
import { cwd } from 'process';

export class Command {
    constructor({
        config,
        branchName,
        command,
        body
    }) {
        this.config = config;
        this.branchName = branchName;
        this.command = command;
        this.body = body;
    }

    async run() {
        if (!this.command) {
            return;
        }

        // create worktree for the branch
        const worktreePath = resolve(this.config.dir, this.config.workTree || '.', `worktree-${this.branchName.replace(/\//g, '_')}`);
        await git.worktree(['add', worktreePath, this.branchName]);

        // Search for command in configured paths relative to repo, cwd, and home
        const searchRoots = [
            this.config.dir,           // repo directory
            cwd(),                      // current working directory
        ];

        // TODO: Complete command execution logic
    }
}