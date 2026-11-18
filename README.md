# Dark

This is a web browser written in Go that runs in the terminal.

Spiritual successor of [coox](https://github.com/3elDU/coox).

The main idea behind it is, what if a CLI web browser could run JS? Suddenly
you can run all those fancy web apps in the terminal.

It does not use chrome or firefox underdeath, this is my own browser engine. Although I did make a clear boundary: i don't want to write my own parsers/interpreters, so I use [sober](https://pkg.go.dev/github.com/grafana/sobek) for running JavaScript and [x/net/html](https://pkg.go.dev/golang.org/x/net/html) for parsing HTML.

This is a recreational programming project.

Some of the high level goals:

- Pass captchas like Cloudflare and Anubis
- Run React, Vue, Svelte
- Be modern (don't spend time supporting ancient web standards)
- Support enough APIs to run modern web apps.

Non-goals:

- Be fast: i prefer simplicity over optimizations
- Implement every Web API: some just don't make sense for the terminal

# AI Disclosure

[![AI:R READ](https://img.shields.io/badge/AI-READ-grey?style=flat&logoSize=auto&labelColor=13171f)](https://0rwx.com)

This is not a vibecoded project, but I do use limited AI assistance.
My main uses are debugging, getting a second opinion on my work, and some feature planning.
I limit AI code to small bugfixes and I don't copy-paste any code directly, so
95% of the code is written by myself.
