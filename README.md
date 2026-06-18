# tcell

Tcell is a Go module that provides a cell based view for text user interfaces.
It was inspired by termbox, and builds on the excellent work of the original
[gdamore/tcell](https://github.com/gdamore/tcell) project. It is hosted on
[Github](https://github.com/unstablebuild/tcell).

## About this fork

This is a fork of [gdamore/tcell](https://github.com/gdamore/tcell), tailored
for a more focused use-case. The upstream project supports a very broad range
of environments and features, including things like WebAssembly. This fork
trims the scope to keep the backend small and fast for our needs.

The main differences from upstream are:

- The widget and view packages have been removed (we use rune-go-sdk's
  [component](https://github.com/unstablebuild/rune-go-sdk/tree/main/component)
  and [handler](https://github.com/unstablebuild/rune-go-sdk/tree/main/handler)
  packages for those).
- The screen interface has been simplified and refactored to work better with
  event loops. The render path is now lock-free, with a focus on reducing
  allocations per render.

These changes make the backend lighter and faster for the kind of TUI
applications we build. The benchmarks below show the effect of these
optimizations on our integration suite:

```
benchmark                                                 old ns/op     new ns/op     delta
BenchmarkIntegration/tscreen_large-10                     22075794      1907189       -91.36%
BenchmarkIntegration/tscreen_medium-10                    3347864       88311         -97.36%
BenchmarkIntegration/tscreen_small-10                     446367        10658         -97.61%
BenchmarkIntegration/tscreen_large_with_resize-10         19868033      3121991       -84.29%
BenchmarkIntegration/tscreen_medium_with_resize-10        4043734       136762        -96.62%
BenchmarkIntegration/tscreen_small_with_resize-10         723866        22815         -96.85%

benchmark                                                 old allocs     new allocs     delta
BenchmarkIntegration/tscreen_large-10                     17570          0              -100.00%
BenchmarkIntegration/tscreen_medium-10                    2815           0              -100.00%
BenchmarkIntegration/tscreen_small-10                     398            0              -100.00%
BenchmarkIntegration/tscreen_large_with_resize-10         198972         0              -100.00%
BenchmarkIntegration/tscreen_medium_with_resize-10        10541          0              -100.00%
BenchmarkIntegration/tscreen_small_with_resize-10         1329           0              -100.00%

benchmark                                                 old bytes     new bytes     delta
BenchmarkIntegration/tscreen_large-10                     200315        3130          -98.44%
BenchmarkIntegration/tscreen_medium-10                    34947         145           -99.59%
BenchmarkIntegration/tscreen_small-10                     5556          20            -99.64%
BenchmarkIntegration/tscreen_large_with_resize-10         20368081      9569687       -53.02%
BenchmarkIntegration/tscreen_medium_with_resize-10        878374        388256        -55.80%
BenchmarkIntegration/tscreen_small_with_resize-10         109314        48380         -55.74%
```

If your project needs the full feature set, broad platform support, or the
widget/view packages, we recommend using the upstream
[gdamore/tcell](https://github.com/gdamore/tcell) directly. Many thanks to the
original authors for their work, which made this fork possible.

## Development

Clone the repository:

```sh
git clone git@github.com:unstablebuild/tcell.git
```

Install [pre-commit](https://pre-commit.com) if you haven't already. Then
install this repo's pre-commit hook `git/hooks/pre-commit` by running:

```sh
pre-commit install
```

Compile:

```sh
go build ./...
```

Run the test suite:

```sh
go test ./...
```

Cross-compile to many combinations of platform and OS:

```sh
./test_crosscompile.sh
```
