# Graft

A lightweight CLI/TUI tool for setting up and managing servers in Docker Compose.

**Why Graft?**

    - Lightweight: Single Go binary with Docker Compose being the only dependency. Just install and use.
    - Stateless: Graft works directly with Compose configurations. You can add existing servers or manually edit configuration without issue.
    - Secure: No admin web interface/API by default, simply use SSH to prevent unnecessary exposure.