import fs from 'fs/promises';
import path from 'path';
import simpleGit from 'simple-git';

/**
 * Init command - initializes AYNIG in the current repository
 */
async function action() {
  const git = simpleGit();

  // Check if we're in a git repository
  try {
    await git.status();
  } catch (error) {
    console.error('Error: Not a Git repository. Please run `git init` first.');
    process.exit(1);
  }

  // Create .aynig/ directory
  const aynigDir = '.aynig';
  let aynigCreated = false;
  try {
    await fs.mkdir(aynigDir, { recursive: false });
    console.log(`✓ Created ${aynigDir}/`);
    aynigCreated = true;
  } catch (error) {
    if (error.code === 'EEXIST') {
      console.log(`⊘ ${aynigDir}/ already exists, skipping`);
    } else {
      throw error;
    }
  }

  // Only create files if directory was just created
  if (aynigCreated) {
    // Create COMMANDS.md
    const commandsContent = `# Available Commands

clean: Clean up the worktree and mark workflow as done
`;

    await fs.writeFile(path.join(aynigDir, 'COMMANDS.md'), commandsContent);
    console.log(`✓ Created COMMANDS.md`);

    // Create clean script
    const cleanScript = `#!/bin/bash
# Clean workflow - removes worktree and marks workflow as done

# Commit completion to current branch
git commit --allow-empty -m "chore: aynig is done"

# Remove the worktree
rm -rf "$AYNIG_WORKTREE_PATH"
`;

    const cleanPath = path.join(aynigDir, 'clean');
    await fs.writeFile(cleanPath, cleanScript);
    await fs.chmod(cleanPath, 0o755);
    console.log(`✓ Created clean script`);
  }

  // Create .worktrees/ directory
  const worktreesDir = '.worktrees';
  try {
    await fs.mkdir(worktreesDir, { recursive: false });
    console.log(`✓ Created ${worktreesDir}/`);
  } catch (error) {
    if (error.code === 'EEXIST') {
      console.log(`⊘ ${worktreesDir}/ already exists, skipping`);
    } else {
      throw error;
    }
  }

  // Add .worktrees/ to .gitignore
  const gitignorePath = '.gitignore';
  const worktreesEntry = '.worktrees/';

  try {
    let gitignoreContent = '';
    try {
      gitignoreContent = await fs.readFile(gitignorePath, 'utf-8');
    } catch (error) {
      if (error.code !== 'ENOENT') throw error;
    }

    if (gitignoreContent.split('\n').some(line => line.trim() === worktreesEntry)) {
      console.log(`⊘ ${worktreesEntry} already in .gitignore, skipping`);
    } else {
      const newContent = gitignoreContent
        ? `${gitignoreContent.trimEnd()}\n${worktreesEntry}\n`
        : `${worktreesEntry}\n`;
      await fs.writeFile(gitignorePath, newContent);
      console.log(`✓ Added ${worktreesEntry} to .gitignore`);
    }
  } catch (error) {
    console.error(`Error updating .gitignore: ${error.message}`);
    process.exit(1);
  }

  console.log('\nAYNIG initialized successfully!');
}

export function registerInitCommand(program) {
  program
    .command('init')
    .description('Initialize AYNIG in the current repository')
    .action(action);
}
