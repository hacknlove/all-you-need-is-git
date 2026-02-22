import { git } from '../gitHelpers/git.js';
import { parseCommitMessage } from '../gitHelpers/parseCommitMessage.js';
import { Command } from './Command.js';
import { Logger } from '../utils/logger.js';
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
        const parsed = await parseCommitMessage(lastCommit);
        return {
            ...parsed,
            commitDate: lastCommit?.date
        };
    }

    async run() {
        const logger = this.config.logger || new Logger(this.config.logLevel);
        // If using remote branches, only process branches from the specified remote
        if (this.config.useRemote && !this.branchName.startsWith(`${this.config.useRemote}/`)) {
            logger.debug('Skipping branch %s (not on remote %s)', this.branchName, this.config.useRemote);
            return;
        }

        const {
            trailers,
            body,
            commitDate
        } = await this.parseLastCommitMessage();
        logger.debug('Inspecting branch %s', this.branchName);

        const command = new Command({
            config: this.config,
            branchName: this.branchName,
            isCurrentBranch: this.isCurrentBranch,
            trailers,
            body,
            commitDate
        });
        await command.run();
    }
}
