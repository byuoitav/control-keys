# control-keys
A lightweight HTTP service - the control keys service is part of the camera control system used with the Pi Control Processors.  The control keys are a security feature that limits who can control the camera in the local room by generating a random number which the Touchpanel-UI-Microservice uses displays in the UI.  

Designed for BYU AV environments, this service provides endpoints for retrieving mappings, refreshing room data, and basic health checks.

---

## 🚀 How It Works

The `control-keys` service initializes a code map on startup, builds a lookup table of control key <-> preset mappings, and exposes a small set of HTTP endpoints using the Echo web framework. These endpoints allow clients to:

- Look up a preset by control key
- Look up a control key by preset
- Refresh preset mappings for a room
- Perform a health check

Internally, the service uses a `CodeMap` object that manages these mappings and handles logic for lookups and refreshes.

---

## 🧾 Required Flags

There are **no required CLI flags** at this time. The service runs with the following hardcoded default:

- **Port:** `8029`

---

## 📌 Examples

Start the service:

```bash
go run main.go

## 🌐 HTTP Endpoints

| Method | Path                     | Description                                               | Response Type            |
| ------ | ------------------------ | --------------------------------------------------------- | ------------------------ |
| GET    | `/:controlKey/getPreset` | Gets the preset (room ID + preset name) for a control key | `Preset` JSON or 404     |
| GET    | `/:preset/getControlKey` | Gets the control key for a given preset string            | `ControlKey` JSON or 404 |
| GET    | `/:room/refresh`         | Refreshes mappings for a specific room                    | `204 No Content` or 404  |
| GET    | `/status`                | Health check endpoint                                     | `"Healthy!"`             |


