import fs from 'fs/promises';
import fse from 'fs-extra';
import path from 'path';
import os from 'os';
import { simpleGit, pathspec } from 'simple-git';
import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

/**
 * Install command - installs AYNIG workflows from another repository
 */
async function action(repo, ref, subfolder) {
  const normalizedRepo = normalizeRepo(repo);
  const tmpDir = await fs.mkdtemp(path.join(os.tmpdir(), 'aynig-install-'));
  const git = simpleGit();

  try {
    // Clone the repository
    const cloneOptions = ['--depth', '1'];

    if (ref) {
      cloneOptions.push('--branch', ref);
      console.log(`Cloning ${normalizedRepo} at ${ref}...`);
    } else {
      console.log(`Cloning ${normalizedRepo} (default branch)...`);
    }

    await git.clone(normalizedRepo, tmpDir, cloneOptions);

    // Determine source path
    const sourcePath = path.join(tmpDir, subfolder || '.aynig');

    // Check if source .aynig exists
    try {
      await fs.access(sourcePath);
    } catch {
      const refInfo = ref ? ` at ${ref}` : '';
      const subfolderInfo = subfolder ? ` in ${subfolder}` : '';
      console.error(`Error: No .aynig directory found in ${normalizedRepo}${refInfo}${subfolderInfo}`);
      await fs.rm(tmpDir, { recursive: true, force: true });
      process.exit(1);
    }

    // Check for uncommitted changes in .aynig
    const destPath = '.aynig';
    const statusBefore = await git.status([pathspec(destPath)]);
    if (!statusBefore.isClean()) {
      console.error(`\nError: .aynig has uncommitted changes. Please commit or stash them first.`);
      await fs.rm(tmpDir, { recursive: true, force: true });
      process.exit(1);
    }

    // Copy entire .aynig directory recursively
    console.log('Copying workflows...');
    await fse.copy(sourcePath, destPath, {
      overwrite: true,
      filter: src => path.basename(src) !== 'README.md'
    });

    // Use git to detect what changed after the copy
    const statusAfter = await git.status([pathspec(destPath)]);

    // Overwritten files (excluding COMMANDS.md)
    const overwritten = statusAfter.modified
      .filter(f => f !== '.aynig/COMMANDS.md')
      .map(f => path.basename(f));

    // New files (excluding COMMANDS.md)
    const newFiles = statusAfter.not_added
      .filter(f => f !== '.aynig/COMMANDS.md')
      .map(f => path.basename(f));

    // Report what was installed
    if (newFiles.length > 0) {
      console.log(`✓ Installed new: ${newFiles.join(', ')}`);
    }

    if (overwritten.length > 0) {
      console.log(`⚠  Overwritten: ${overwritten.join(', ')}`);
    }

    // Handle COMMANDS.md if it was modified
    const commandsMdModified = statusAfter.modified.includes('.aynig/COMMANDS.md');

    if (commandsMdModified && overwritten.length > 0) {
      // Get the diff for COMMANDS.md
      const diff = await git.diff(['.aynig/COMMANDS.md']);

      const hasOpencode = await checkCommand('opencode');
      const hasClaude = await checkCommand('claude');

      if (hasOpencode || hasClaude) {
        const tool = hasOpencode ? 'opencode' : 'claude';
        console.log(`\n🤖 Calling ${tool} to merge COMMANDS.md...`);

        const prompt = `Some commands were just installed/overwritten. Please merge this COMMANDS.md file intelligently by looking at the git diff below. Keep the most accurate and valid description for each command, removing any duplicates.

Git diff:
${diff}`;

        try {
          const destCommandsPath = path.join(destPath, 'COMMANDS.md');
          if (hasOpencode) {
            await execAsync(`opencode "${destCommandsPath}" -p "${prompt}"`);
          } else {
            await execAsync(`claude -p "${prompt}" "${destCommandsPath}"`);
          }
          console.log(`✓ COMMANDS.md merged by ${tool}`);
        } catch (error) {
          console.error(`⚠  Error calling ${tool}: ${error.message}`);
          console.log(`⚠  Please manually review and merge COMMANDS.md`);
        }
      } else {
        console.log(`\n⚠  Warning: COMMANDS.md was modified. Please manually review the changes.`);
        console.log(`\nGit diff preview:\n${diff.split('\n').slice(0, 20).join('\n')}`);
        if (diff.split('\n').length > 20) {
          console.log('... (run `git diff .aynig/COMMANDS.md` to see full diff)');
        }
      }
    } else if (commandsMdModified) {
      console.log(`✓ COMMANDS.md updated`);
    }

    // Clean up temp directory
    await fs.rm(tmpDir, { recursive: true, force: true });

    console.log('\n✓ Workflows installed successfully!');
  } catch (error) {
    // Clean up on error
    try {
      await fs.rm(tmpDir, { recursive: true, force: true });
    } catch {}
    throw error;
  }
}

function normalizeRepo(repo) {
  if (isFullRepoUrl(repo)) {
    return repo;
  }

  if (isShorthandRepo(repo)) {
    return `https://github.com/${repo}.git`;
  }

  return repo;
}

function isFullRepoUrl(repo) {
  return /^(https?:\/\/|git@|ssh:\/\/)/i.test(repo);
}

function isShorthandRepo(repo) {
  return /^[A-Za-z0-9_.-]+\/[A-Za-z0-9_.-]+$/.test(repo);
}

/**
 * Check if a command exists in PATH
 */
async function checkCommand(command) {
  const checkCmd = process.platform === 'win32' ? 'where' : 'which';
  try {
    await execAsync(`${checkCmd} ${command}`);
    return true;
  } catch {
    return false;
  }
}

export function registerInstallCommand(program) {
  program
    .command('install <repo> [ref] [subfolder]')
    .description('Install AYNIG workflows from another repository')
    .action(action);
}
