import { git } from '../gitHelpers/git.js';
import { parseCommitMessage } from '../gitHelpers/parseCommitMessage.js';
import { Command } from './Command.js';
/**
 * Branch class to manage individual Git branches.
 *
 * It:
 * 1. Skips if using remote and the branch is from the wrong remote.
 * 2. Parses the last commit message to extract trailers.
 * 3. Instantiates and runs a Command based on the parsed trailers
 * 4. waits for the Command to complete.
 *
 * It has no clue about what the Command does internally.
 */
export class Branch {
    constructor({ config, branchName, isCurrentBranch }) {
        this.config = config;
        this.branchName = branchName;
        this.isCurrentBranch = isCurrentBranch;
    }

    async parseLastCommitMessage() {
        const log = await git.log([this.branchName, '-1']);
        const lastCommit = log.latest;
        return parseCommitMessage(lastCommit);
    }

    async run() {
        // If using remote branches, only process branches from the specified remote
        if (this.config.useRemote && !this.branchName.startsWith(`${this.config.useRemote}/`)) {
            return;
        }

        const {
            trailers,
            body,
        } = await this.parseLastCommitMessage();

        const command = new Command({
            config: this.config,
            branchName: this.branchName,
            isCurrentBranch: this.isCurrentBranch,
            trailers,
            body
        });
        await command.run();
    }
}