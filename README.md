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
git clone https://github.com/LucasJezap/altimi-sync.git
cd altimi-sync
go build -o altimi-sync main.go
```

## 📌 Examples

### 🔄 Print help

```bash
./altimi-sync -h
./altimi-sync --help
```

### 🔄 Basic sync from source to target

```bash
./altimi-sync ./source ./target
```

### 🔄 Basic sync from source to target, delete files that do not exist in source

```bash
./altimi-sync -d ./source ./target
./altimi-sync --delete-missing ./source ./target
```

## 🧪 Tests

### 🔄 Print help

```bash
# normal 
go test ./...
# verbose
go test -v ./...
# with code coverage
go test -cover ./...
```
