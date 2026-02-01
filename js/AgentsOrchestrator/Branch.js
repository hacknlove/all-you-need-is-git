import { parseCommitMessage } from '../gitHelpers/parseCommitMessage.js';
import { Command } from './Command.js';
/**
 * Branch class to manage individual Git branches.
 * * It:
 * 1. Skips if the branch is from the wrong remote.
 * 2. Parses the last commit message to extract trailers.
 * 3. Instantiates and runs a Command based on the parsed trailers
 * 4. waits for the Command to complete.
 * 
 * It has no clue about what the Command does internally.
 */
export class Branch {
    constructor({ git, config, branchName }) {
        this.git = git;
        this.config = config;
        this.branchName = branchName;
    }

    async parseLastCommitMessage() {
        const log = await this.git.log([this.branchName, '-1']);
        const lastCommit = log.latest;
        return parseCommitMessage(lastCommit);
    }
    
    async run() {
        if (!this.branchName.startsWith(`${this.config.remote}/`)) {
            return;
        }

        const {
            trailers,
            body,
        } = await this.parseLastCommitMessage();

        const command = new Command({
            git: this.git,
            config: this.config,
            branchName: this.branchName,
            command: trailers.aynig?.trim().toLowerCase(),
            body
        });
        await command.run();
    }
}