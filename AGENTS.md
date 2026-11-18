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

Pixel size is terminal resolution multiplied by constant font size.
For a 80x30 terminal, W = 80\*8 = 640, H = 480.
A box with `width: 8px` and `height: 16px` occupies exactly one cell on the screen.
Non-divisible sizes are rounded.

Text layout becomes simpler and faster, as the maximum width for the given text
is simply `len(text) * FontWidth`.

For background-only cells we can use the half-block element (U+2584: ▄) to
double our vertical color resolution.

# Non-goals

Don't want to spend too much energy on those.

- Optimizations. I already have a compiled language, and many aspects of the
  browser simply don't apply here. I don't have to keep a dozen GPU buffers,
  and font caches.
- Security. While JavaScript is already sandboxed, I am currently not willing
  to implement advanced security mechanisms, like process isolation.
- Implementing every Web API. Some just don't make sense in our context,
  some are too old, some are too obscured and not used by anyone.

# How we work on this

I like writing code. I will use your help with standards, architecture design
and debugging. You can propose changes but don't add any code without asking.

## Debugging

Dark can't run local files, so set you'll need an HTTP server to load a file.
