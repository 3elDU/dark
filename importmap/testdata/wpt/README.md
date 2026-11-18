# WPT import map tests

Data-driven tests from [web-platform-tests](https://github.com/web-platform-tests/wpt/tree/master/import-maps/data-driven),
copied from commit `7f8d6b3544da6ce717dc0eb4e63754e5fbc1a3f0`. Licensed under 3-Clause BSD, see `LICENSE.md`.

The format is described in the upstream [README](https://github.com/web-platform-tests/wpt/blob/master/import-maps/data-driven/README.md).
`parsing-*.json` files check the parsed import map (`expectedParsedImportMap`), the rest check resolution (`expectedResults`).

To update, replace `SHA` with a newer commit and re-download every `.json` file:

```sh
SHA=7f8d6b3544da6ce717dc0eb4e63754e5fbc1a3f0
for f in *.json; do
  curl -sSf -o "$f" "https://raw.githubusercontent.com/web-platform-tests/wpt/$SHA/import-maps/data-driven/resources/$f"
done
```
