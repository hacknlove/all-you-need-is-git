# Documentation Suggestions

These are areas where I (as a new user/developer integrating AYNIG) needed clarification.

## 1. Commit Structure Visualization

The relationship between commit parts and AYNIG was unclear initially:
- Title → human-only summary
- Body → prompt delivered to the next command
- Trailers → structured metadata (dwp-state, etc.)

A simple diagram early in the docs would have helped:
```
commit message title (human summary)
---
body: instructions for the next command

---
dwp-state: <state>
dwp-origin-state: <state>
```

## 2. set-state vs set-working Purpose

The CLI reference shows syntax but not when to use each:
- `set-working` → lease commit (I'm actively working)
- `set-state <state>` → signal next command to run

A "Decision Flow" or "When to use X" section would clarify this.

## 3. Command Lifecycle Example

The authoring guide shows a simple hello-world example. A fuller example showing:
1. `dwp-state: review` triggers review command
2. Review command runs, decides to implement fixes
3. Creates `dwp-state: implement` with body containing fix instructions
4. AYNIG dispatches to implement command

This would show how state flows between commands.

## 4. Common Workflow Patterns

Some commands are "deciders" (like an idle command that picks next task), others are "doers" (implement fixes, write tests). A guide showing:
- Pattern: command that reads context and sets next state
- Pattern: command that always advances to a specific state
- Pattern: command that can loop (retry, blocked, etc.)

This helps developers structure their command logic.

## 5. Environment Variables Reference Location

The Environment Variables page exists but it's easy to miss. A prominent link in the authoring guide pointing to it would help, since commands need to know what context they receive (AYNIG_BODY, trailers, etc.).

## 6. Mention `AYNIG_BODY` Wherever `--prompt` Is Introduced

It took extra digging to confirm that the body passed with `aynig set-state --prompt ...` reaches the next command as `AYNIG_BODY`.

That connection is important for real command patterns such as an `idle` command that reads guidance from the triggering commit body to decide what to prioritize next.

The environment variables page documents `AYNIG_BODY`, but the same fact should also be called out directly in places where users are most likely to need it:
- command authoring docs
- `set-state` docs
- run/commit format docs

Without that cross-linking, `AYNIG_BODY` feels effectively undocumented unless the reader already knows where to look.
