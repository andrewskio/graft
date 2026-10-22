# Graft

A lightweight Linux CLI/TUI tool for setting up and managing servers in Docker Compose. 

![Demo image showing list and new commands](./docs/assets/demo.png)

## Why Graft?

- Lightweight: Single Go binary with Docker Compose being the only dependency. Just install and use.

- Stateless: Graft works directly with Compose configurations. You can add existing servers (WIP) or manually edit configuration without issue.

- Secure: No admin web interface/API by default, simply use SSH to prevent unnecessary exposure.

## Basic Usage

- `graft new` - Launch interactive TUI to create new server
- `graft ls` - List servers and respective container status
- `graft start/stop <server>` - Start/stop the container for a server
- `graft console <server>` - Send commands to the server console
- `graft log <server>` - Scrolling log viewer to view console logs
- `graft status <server>` - Print the current status of a server container

## Install From Source

1. Install go, docker, and docker compose. Add user to docker group.

2. Run `go install github.com/andrewskio/graft@latest`. Make sure the Go bin directory is added to path.

3. Run `graft help` to get started. Optionally, run `graft ac` and restart your terminal to add completions.


Server data is stored in `~/.local/share/graft/servers`. 
