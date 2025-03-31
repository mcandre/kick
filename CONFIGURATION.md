# CONFIGURATION

The kick configuration is specified with environment variables.

# `KICK_MESSAGE`

Customize the git commit message (default: `"up"`).

Blank messages trigger git's configured `core.editor` to prompt for a dynamically chosen message.

# `KICK_NONCE`

When set to `1`, enables updating a `.kick` nonce file with a timestamp (default: `0`).

Useful for generating commits when a repository is otherwise unchanged.

# `KICK_FETCH_ALL`

When set to `1`, enables fetching (tags) from all remotes (default: `1`).

# `KICK_PULL_ALL`

When set to `1`, enables pulling from all remotes (default: `1`).

# `KICK_PUSH_ALL`

When set to `1`, enables pushing to all remotes (default: `1`).

# `KICK_SYNC_TAGS`

When set to `1`, enables pushing and pulling tags (default: `1`).
