Read `.dwp/CONTRACT.md`.

The current job is in state `stalled`.

Your goal is to decide the next state.

Instructions:
1) Read `.dwp/COMMANDS.md` to learn:
   - available states and their instructions
2) Find the `working` commit for this run:
   - use `dwp-stalled-run` to locate the last `dwp-state: working` commit for the run
   - read its `dwp-origin-state` trailer to recover the original intent
3) Find the initial commit of this job: the closest commit with `dwp-state: <origin-state>`.
4) Review all `working` commits since that start: commit messages and changes.
5) Assess the current state of the job:
   - Is it finished?
   - Can it continue?
   - Must it be restarted?
   - Is it not recoverable?

Decision: pick one action:
1. Retry from scratch: reset to the initial commit and re-emit the initial state, optionally with a refined prompt.
2. Continue: keep current history, re-emit the initial state with a new prompt to continue.
3. Switch state: set another available state if it fits better.
4. Fail: set `dwp-state: failed` and explain why in the commit body.

IMPORTANT RULES:
- The new commit MUST follow the contract.
- Do not change any files; only create a new commit with the new state trailers and body.
