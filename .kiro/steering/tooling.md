# Tooling & code exploration

## Always use codegraph to parse and explore code

This repository is indexed by codegraph (see the `.codegraph/` folder). When you
need to understand, locate, or trace code, use the codegraph tools FIRST — before
falling back to plain file reads, grep, or find.

- To answer "how does X work", "where is X", "what calls X", or to survey an area
  before editing, call `codegraph_explore` with the relevant symbol/file names or a
  natural-language question. It returns the verbatim, line-numbered source of the
  relevant symbols grouped by file, plus the call flow among them — in one call.
- Treat source returned by codegraph as already read: do NOT re-open those files
  with the read tool afterward.
- Before editing a symbol, check its "blast radius" (callers) from the codegraph
  output so changes stay consistent across the codebase.
- For a specific symbol whose body wasn't shown, run another `codegraph_explore`
  (or `codegraph_node`) with its exact name rather than reading the whole file.

Reserve raw `grep`/`find`/full-file reads for cases codegraph can't serve well:
non-code files (TMX maps, JSON, YAML, build workflows, assets), exact-string
searches across config, or when codegraph returns nothing for a query.

## Verification

- After any code change, run `go build ./...` and, for changed packages,
  `go vet ./<pkg>/...`.
- There is no automated test suite; confirm behavior by building. Long-running
  processes (the game itself) should be launched by the user, not blocked on here.

## Housekeeping

- Do not read or commit the generated binaries: `gamejam.exe`, `gamejam.wasm`, and
  the `__debug_bin*.exe` Delve artifacts.
