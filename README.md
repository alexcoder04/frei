
# frei (FREe Improved)

![screenshot](./screenshot.png)

Since almost every basic command nowadays has a fancy Rust or Go rewrite, `free`
should not be an exception.

In case you didn't know, `free` shows your memory usage on Unix-based systems.
`frei` obtains memory data from `/proc/meminfo` and represents it in a colored
bar chart.

## Installation

### Arch Linux

frei is available on the [AUR](https://aur.archlinux.org/packages/frei).

### Pre-compiled binaries

Binaries for `i386`, `amd64` and `arm` are available on the [releases
page](https://github.com/alexcoder04/frei/releases/latest).

### Compiling from source

Make sure you have [Go](https://golang.org/doc/install.html) installed.

```sh
git clone https://github.com/alexcoder04/frei.git
cd frei
go build .
```

## Command-line options

| option     | description                                 |
|------------|---------------------------------------------|
| `-help`    | show list of options                        |
| `-h`       | human-readable numbers (implies `-numbers`) |
| `-key`     | display color key                           |
| `-numbers` | print numbers in addition to the chart      |
| `-version` | display version and exit                    |

