# About this fork
This fork started because we did not agree on the latest changes in the original
https://github.com/gdamore/tcell repository. A lot of complexity has been added to support
esoteric use-cases such as WebAssembly.

All the widget/view packages have been removed (we use go-tui for that), the screen interface has been 
simplified and refactored to better support event loops (lock-free, reduced allocations per render),
resulting in improved performance.

From the start of the optimizations (7bf539846adad20240bd27c00d5d572b77a837c2) up until the
last optimization (c38ded0df481f83ebbc9fa0838fde95e3a11857c),
the impact is quite severe, which yields a much more lightweight TUI backend.

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
