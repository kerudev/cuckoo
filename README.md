# cuckoo

A crontab visualizer made in Go with Raylib bindings.

Initially made for learning purposes, finished out of love <3 (or masochism).

<!-- Generated with: https://github.com/darkmavis1980/markdown-index-generator -->
<!-- index-start -->
## Index

- [Motivation](#motivation)
- [Install & use](#install--use)
- [Dependencies](#dependencies)
- [Build](#build)
- [Roadmap](#roadmap)
  - [CGO_ENABLED=0](#cgo_enabled0)
  - [WASM](#wasm)
<!-- index-end -->

## Motivation

Imagine a process orchestrator (Airflow, Celery, etc.) that coordinates over
200 processes. Each process has a [cron](https://en.wikipedia.org/wiki/Cron)
that defines when it will be executed.

Now, think about how you would decide the cron for your new process. This is
important so your new cron doesn't collide with heavy processes or workload
peaks, which would create overhead that could be better distributed.

It's hard if can't visualize the data, right? That's where `cuckoo` comes in.

`cuckoo` lets you visualize crons as coordinates on a grid to:
- Get insights about your crons with just a glance.
- Know the next free spot where you can place a new cron.
- Identify periods where there is a work overload.

## Install & use

To install cuckoo:

```sh
go install github.com/kerudev/cuckoo@latest
```

Takes a JSON file with the following structure:

```json
{
  "process_1": "*/5 * * * *",
  "process_2": "0 */2 * * *",
  "process_3": "15 3 * * *",
  // ...
}
```

Or a CSV file with two columns (header names are just a convention, but mind
the colons!):

```csv
name,cron
process_1,"*/5 * * * *"
process_2,"0 */2 * * *"
process_3,"15 3 * * *"
# ...
```

To run cuckoo:

```sh
cuckoo                      # Runs from the current directory
cuckoo -path path/to/data   # Runs from a directory / loads and parses a file
```

## Dependencies

cuckoo's dependencies:
- [raylib-go](https://github.com/gen2brain/raylib-go): Go bindings for raylib & raygui.

dev dependencies:
- [air](https://github.com/air-verse/air): hot module reloading for Go.

## Build

When running `go build` or `go run`, you may want to pass the `tags` argument.
This tells the bindings how to handle windows, and keyboard and mouse events.

Example:

```sh
go run -tags x11 . -path ...
```

I usually compile with the `noaudio` tag to save space, as there are no audio
involved in cuckoo:

```sh
go build -tags x11
du -sh cuckoo
7.8M    cuckoo

go build -tags noaudio,x11
du -sh cuckoo
6.5M    cuckoo
```

- List of tags: https://github.com/gen2brain/raylib-go/blob/master/README.md#build-tags
- Related discussion: https://github.com/gen2brain/raylib-go/discussions/554

## Roadmap

### CGO_ENABLED=0

Although cuckoo is a really small program, I feel like it spends too much time
on the CPU, probably due to not being able to compile without cgo.

The `debug` module has a `DebugServer()` that I use to measure performance.

I run cuckoo on one shell and this command on other shell:

```sh
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10

 flat  flat%   sum%        cum   cum%
7.96s 80.24% 80.24%      9.51s 95.87%  runtime.cgocall
0.29s  2.92% 83.17%      0.49s  4.94%  runtime.casgstatus
0.28s  2.82% 85.99%      0.76s  7.66%  runtime.reentersyscall
```

Seems like 95% of the time is spent on the CPU doing cgo stuff.

Right now, cuckoo can't compile with `CGO_ENABLED=0` because raygui doesn't
support it, but I plan for that to change.

### WASM

This project is currently being developed just for desktop. A WASM port is in
the works but not ready yet, as the project that ports raylib to WASM doesn't
support raygui yet
(see [BrownNPC/Raylib-Go-Wasm #10](https://github.com/BrownNPC/Raylib-Go-Wasm/pull/10)).

> [!CAUTION]
> These notes are subject to change, as they describe the current way of using
> the raygui bindings, which are under development on other project.

Complete guide: https://github.com/BrownNPC/Raylib-Go-Wasm/blob/master/README.md

To compile `cuckoo` to WASM:

1. Clone or create a symlink to the WASM bindings repository:

```sh
git clone https://github.com/BrownNPC/Raylib-Go-Wasm.git

ln -s path/to/Raylib-Go-Wasm/ Raylib-Go-Wasm  # I do this while developing the bindings
```

2. Checkout to the `raygui` branch.

3. Paste this into `go.mod` and run `go mod tidy`:

```go.mod
replace (
	github.com/gen2brain/raylib-go/raygui => ./Raylib-Go-Wasm/raygui
	github.com/gen2brain/raylib-go/raylib => ./Raylib-Go-Wasm/raylib
	github.com/BrownNPC/Raylib-Go-Wasm/wasm-runtime => ./Raylib-Go-Wasm/wasm-runtime
)
```

4. WASM doesn't support Raylib's `WindowShouldClose()`, so you need to change this:

```go
for !rl.WindowShouldClose() { ... }
```

To this:

```go
var update = func() { ... }

rl.SetMainLoop(update)
```

5. You might lots of errors in the code, but they are normal. Compile with:

```sh
GOOS=js GOARCH=wasm go build -o ./Raylib-Go-Wasm/index/main.wasm .
```

6. Run the server of your choice:
- On Python: `cd Raylib-Go-Wasm/index && python -m http.server`
- On Go: `go build ./Raylib-Go-Wasm/server/server.go` then `./server`
