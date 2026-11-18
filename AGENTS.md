# dark

This is my attempt at writing a web browser for the terminal.

What makes it different from popular CLI browsers (lynx, w3m, etc.) is that it
executes JavaScript, meaning it can run modern websites and webapps.

Unlike similar apps (e.g. brow.sh) it does not have firefox or chromium underneath,
it is a custom browser engine written in Go. It is not in my goals to write everything
from scratch, so I use pre-built parts like the JavaScript engine, HTML parser and so on.

I have a couple goals/directions i'd like to keep in mind during development:

- Support modern JavaScript: this is the main one. What constitutes as modern?
  modules, latest EcmaScript built-ins (those we could add without patching the engine).
- Be standards-compliant: for a website, dark should behave like a regular browser.
  It does not know it's running inside a terminal.

I want to get to displaying real websites quicker by implementing the most essential
parts and APIs first.

# How much styling to support?

One thing goes right away, which is text scaling. Terminal is made of cells of
equal size, meaning every character, on every webpage ever will have the same size.

Other than than, I hope to support as much CSS as possible, especially:

- Colors: every modern terminal has truecolor
- Display `none`, `flexbox`, `grid`, `block`, `inline`
- Visibility `visible`, `hidden`
- All values for `position`
- Width, Height, Padding, Margin

# My take on rendering

A character terminal is fundamentally very different from a screen, so we should
treat it as such. The rendering output is a character buffer, which then goes
directly into the terminal.

Because the web was never really intended for character devices, there are gonna
be many compromises I am willing to accept, rather than try to engineer around them.

Pixel size is terminal resolution multiplied by cell size, determined via an
escape code, with a fallback to 8 as a default.
A box with `width: 8px` and `height: 16px` occupies exactly one cell on the screen.
Non-divisible sizes are rounded.

Text layout becomes simpler and faster, as the maximum width for the given text
is simply `len(text) * FontWidth`.

For background background-only cells we can use the half-block element (U+2584: ▄) to
double our vertical color resolution. With text on foreground, this becomes a single
vertical pixel per cell
This is the baseline for images.
On terminals that support kitty image protocol we would use that, as it allows
to show text over the image.

# Non-goals

Don't want to spend too much energy on those.

- Optimizations. I already have a compiled language, and many aspects of the
  browser simply don't apply here. I don't have to keep a dozen GPU buffers,
  and font caches.
- Security. While JavaScript is already sandboxed, I am currently not willing
  to implement advanced security mechanisms, like process isolation.
- Implementing every Web API. Some just don't make sense in our context,
  some are too old, some are too obscured and not used by anyone.

# The current state of things

A few sentences on the current progress with dark, so you don't have
to read half of the code to get a sense. You are welcome to change the lines below.
This is not a changelog, just a brief for yourself and future sessions.

Early stage: dark is a headless JS host for now. Nothing is rendered to the terminal yet.

`go run . <url>` fetches the page, parses it into our own DOM, prints the tree,
runs the page's scripts and then sits in the event loop forever.

What exists:

- **Navigation**: `Tab` holds a history stack and the active `Document`. One tab, no UI.
- **DOM**: `x/net/html` output is converted into our own tree (flat slice of nodes,
  addressed by `NodeHandle`). Not exposed to JS yet: `document` is an empty object.
- **Scripts**: classic and module scripts, inline and external. All of them run as
  deferred, after parsing finishes, because `x/net/html` builds the whole tree at once.
  `nomodule` scripts are skipped.
- **Modules**: static imports, dynamic `import()` (fetches synchronously),
  `import.meta.url` and `import.meta.resolve()`. Bare specifiers fail, import maps
  are recognised but not applied.
- **Web APIs**: `console.log`, timers, `queueMicrotask`, `window`/`self`,
  `navigator` (static values, user agent switchable with `-u`).
- **Event loop**: `Realm.MainLoop`. Timers run in a separate goroutine and send fired
  ids back over a channel, callbacks go through a task queue.

Not started: CSS, layout, rendering, input, DOM bindings, events, `fetch`/XHR, storage.
`terminal/` is an empty stub for cell size detection.

Tests: only `timers/` has Go tests. `tests/` holds HTML pages for manual testing
over an HTTP server.

# How we work on this

I like writing code. I will use your help with standards, architecture design
and debugging. You can propose changes but don't add any code without asking.

## Debugging

Dark can't run local files, so set you'll need an HTTP server to load a file.
