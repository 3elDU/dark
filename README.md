# Dark

This is a web browser written in Go that runs in the terminal.

Spiritual successor of [coox](https://github.com/3elDU/coox).

The main idea behind it is, what if a CLI web browser could run JS? Suddenly
you can run all those fancy web apps in the terminal.

It does not use chrome or firefox underdeath, this is my own browser engine. Although I did make a clear boundary: i don't want to write my own parsers/interpreters, so I use [sober](https://pkg.go.dev/github.com/grafana/sobek) for running JavaScript and [x/net/html](https://pkg.go.dev/golang.org/x/net/html) for parsing HTML.

This is a recreational programming project.

Goals:

- Be modern (don't spend time supporting ancient web standards)
- Support enough APIs to run modern web apps.

Non-goals:

- Be fast: i prefer simplicity over optimizations
- Implement every Web API: some just don't make sense for the terminal
