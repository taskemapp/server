[![codecov](https://codecov.io/gh/taskemapp/server/graph/badge.svg?token=13F6QIQ717)](https://codecov.io/gh/taskemapp/server)
[![CodeFactor](https://www.codefactor.io/repository/github/taskemapp/server/badge)](https://www.codefactor.io/repository/github/taskemapp/server)

# Server for `taskem`

## Development

### Requirements

- libvips 8.3+ (8.8+ recommended)
- C compatible compiler such as gcc 4.6+ or clang 3.0+
- Go 1.3+

### Installation

1. [Download](https://www.msys2.org/) and install MSYS2.
2. After install open MSYS2 UCRT64
3. Update MSYS2
    ```bash
    pacman -Syuu
    ```
4. Install toolchain | gcc | pkg-config | vips
    ```bash
    pacman -S mingw-w64-ucrt-x86_64-toolchain
    pacman -S mingw-w64-ucrt-x86_64-gcc
    pacman -S pkg-config
    pacman -S mingw-w64-ucrt-x86_64-libvips
    ```
   After that add `C:/msys64/ucrt64/bin` to your `PATH`

## Coverage

![codecov graph](https://codecov.io/gh/taskemapp/server/graphs/sunburst.svg?token=13F6QIQ717)