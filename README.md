# volcabot

**volcabot** is a Discord bot created for and by the [Volca Heaven Discord server](https://discord.gg/EFvuAzr).

# Usage

In order to run Volcabot, you need a bot token which you can acquire through the [Discord Developer Portal](https://discord.com/developers/applications). Volcabot reads the token through the environment variable `VOLCABOT_TOKEN`.

Volcabot can be built and run with the standard Go toolchain:

```
$ go run .
```

## Environment variables

| Variable             | Description                                                                         |
| -------------------- | ----------------------------------------------------------------------------------- |
| `VOLCABOT_TOKEN`     | Discord bot API token                                                               |
| `VOLCABOT_LOG_LEVEL` | Log level (one of `panic`, `fatal`, `error`, `warning`, `info`, `debug`, or `trace` |

# Contributing

This repository has two main branches: `master` and `develop`. `master` is considered the stable branch and is the one the public bot is built from. `develop` is the development branch from which all other branches should be based, and into which all other branches are merged. If you want to submit a pull request, please fork `develop`!

This project keeps a fairly low barrier of entry in terms of code quality; feel free to submit a pull request even if it looks "hacky" or like it might not be up to par for a more professional project. We want anyone to be able to contribute!

Before adding a feature, I would consult the community in the Discord server to gauge actual interest. There are plenty of creative things you can do with a Discord bot, but it can be easy for a bot to start feeling pestiferous when it becomes bloated with gratuitous novelty features.

**Feel free to DM katie#3975 if you have any questions!**

**Please run `go fmt` before committing!**
