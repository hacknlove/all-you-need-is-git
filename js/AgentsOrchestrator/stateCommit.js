export function buildStateCommitMessage(subject, prompt, trailers) {
  const cleanSubject = String(subject || '').replace(/\n+$/, '');
  const cleanPrompt = String(prompt || '').replace(/\n+$/, '');

  const lines = [cleanSubject, '', cleanPrompt, ''];
  for (const trailer of trailers || []) {
    lines.push(`${trailer.key}: ${trailer.value}`);
  }

  return `${lines.join('\n')}\n`;
}

export async function commitState(worktreeGit, subject, prompt, trailers) {
  const message = buildStateCommitMessage(subject, prompt, trailers);
  await worktreeGit.commit(message, { '--allow-empty': null });
}
