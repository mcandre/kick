# CONFIGURATION

The kick configuration is specified with environment variables.

# `KICK_NONCE`

When set to `1`, enables updating a `.kick` nonce file with a timestamp (default: `0`).

Useful for generating commits when a repository is otherwise unchanged.

# `KICK_MESSAGE`

Customize the git commit message (default: `"up"`).

Blank messages trigger git's configured `core.editor` to prompt for a dynamically chosen message.
