# cuckoo

> [!NOTE]
> This program is on early development.

A crontab visualizer made in Go with Raylib bindings.

Made for learning purposes.

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
go install github.com/kerudev/cuckoo
```

cuckoo takes a simple JSON file whit the following structure:

```json
{
    "CRON_NAME": "CRON_VALUE",
    ...
}
```

To run cuckoo:

```sh
cuckoo -path path/to/data
```

## Dependencies

cuckoo's dependencies:
- [raylib-go](https://github.com/gen2brain/raylib-go): Go bindings for raylib.

dev dependencies:
- [air](https://github.com/air-verse/air): hot module reloading for Go.

## About WASM

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
