# 🌀 Altimi Sync

**Altimi Sync** is a simple CLI tool written in Go for synchronizing files between two directories.

It supports:
- One-way sync from source to target
- File comparison by size, modification time, and checksum
- Optional deletion of files in the target that do not exist in the source
- Helpful logging and error handling (including permission issues)

---

## 📦 Installation

You can build the application from source:

```bash
git clone https://github.com/your-org/altimi-sync.git
cd altimi-sync
go build -o altimi-sync main.go
