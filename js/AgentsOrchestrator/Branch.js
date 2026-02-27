import { git } from '../gitHelpers/git.js';
import { parseCommitMessage } from '../gitHelpers/parseCommitMessage.js';
import { Command } from './Command.js';
import { Logger, BufferedLogger, resolveLevel } from '../utils/logger.js';
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
        const raw = await git.raw(['log', this.branchName, '-1', '--pretty=format:%s%x1f%b%x1f%cI']);
        const parts = raw.split('\x1f');
        const message = (parts[0] || '').replace(/\n+$/, '');
        const body = (parts[1] || '').replace(/\n+$/, '');
        const commitDate = (parts[2] || '').trim();
        const parsed = await parseCommitMessage({ message, body });
        return {
            ...parsed,
            commitDate
        };
    }

    async run() {
        const buffer = new BufferedLogger();
        buffer.debug('Inspecting branch %s', this.branchName);
        const baseLevel = resolveLevel({
            cliLevel: this.config.logLevel,
            cliSet: this.config.logLevelSet,
            trailerLevel: '',
            envLevel: process.env.AYNIG_LOG_LEVEL,
            defaultLevel: this.config.logLevel
        });
        const baseLogger = this.config.logger || new Logger(baseLevel);
        // If using remote branches, only process branches from the specified remote
        if (this.config.aynigRemote && !this.branchName.startsWith(`${this.config.aynigRemote}/`)) {
            buffer.debug('Skipping branch %s (not on remote %s)', this.branchName, this.config.aynigRemote);
            buffer.flush(baseLogger);
            return;
        }

        const {
            trailers,
            body,
            commitDate
        } = await this.parseLastCommitMessage();

        const trailerLevel = String(this.trailerValue(trailers, 'aynig-log-level')).trim();
        const resolvedLevel = resolveLevel({
            cliLevel: this.config.logLevel,
            cliSet: this.config.logLevelSet,
            trailerLevel,
            envLevel: process.env.AYNIG_LOG_LEVEL,
            defaultLevel: this.config.logLevel
        });
        const logger = new Logger(resolvedLevel);
        buffer.flush(logger);

        const command = new Command({
            config: this.config,
            branchName: this.branchName,
            isCurrentBranch: this.isCurrentBranch,
            trailers,
            body,
            commitDate,
            logger,
            logLevel: resolvedLevel
        });
        await command.run();
    }

    trailerValue(trailers, key) {
        const target = String(key).toLowerCase();
        for (const [k, value] of Object.entries(trailers || {})) {
            if (k.toLowerCase() !== target) {
                continue;
            }
            if (Array.isArray(value)) {
                return String(value[value.length - 1] || '');
            }
            return String(value || '');
        }
        return '';
    }
}
