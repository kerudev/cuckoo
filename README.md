# cuckoo

> [!WARNING]
> Use at your own risk! This library is experimental.

A crontab visualizer made in Go with Raylib bindings.

Made for learning purposes.

# Motivation

Imagine a process orchestrator (Airflow, Celery, etc.) that coordinates over
200 processes. Each process has a [cron](https://en.wikipedia.org/wiki/Cron)
that defines when it will be executed.

Now, think about how you would decide the cron for your new process. It's hard
to know if can't visualize that data right? That's where cuckoo comes in.

cuckoo lets you visualize crons as coordinates on a grid to:
- Get insights about your crons with just a glance.
- Know the next free spot where you can place a new cron.
- Identify periods where there is a work overload.

# Install & use

To install cuckoo:

```sh
go install github.com/kerudev/cuckoo
```

cuckoo takes a simple JSON file whit the following structure:

```json
{
    [CRON_NAME]: [CRON_VALUE],
    ...
}
```

To run cuckoo:

```sh
cuckoo -path path/to/data
```

# Dependencies

cuckoo's dependencies:
- [raylib-go](https://github.com/gen2brain/raylib-go): Go bindings for raylib.

dev dependencies:
- [air](https://github.com/air-verse/air): hot module reloading for Go.
