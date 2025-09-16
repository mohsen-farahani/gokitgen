# 🚀 GoKitGen — Go Kit Code Generator with Style

> ✨ A powerful, interactive CLI tool to generate production-ready Go Kit microservices — with models, gRPC/HTTP APIs, repositories, tests, and Protobuf — all in seconds.

<div align="center">
  <img src="./docs/Screenshot from 2025-09-16 11-05-22.png" alt="lazyssh logo" width="100%" height="600"/>
</div>

---
> [!WARNING]
> This project is being tested.

## ✨ Features

- 🧩 **Interactive Wizard** — Beautiful TUI with Bubble Tea
- 📦 **Generate Models** — With GORM, Enums, Relations, Validation
- 🌐 **HTTP & gRPC APIs** — Fully generated with Transport, Endpoints, Routes
- 🧪 **Auto-generated Tests** — For both HTTP and gRPC transports
- 📜 **Protobuf Support** — Auto-generate `.proto` files for gRPC
- 🗃️ **Repository Layer** — With `CommonBehaviorRepository` pattern
- 🧱 **Project Structure** — Clean, scalable, Go Kit standard
- 🛠️ **Installable CLI** — Use `gokitgen` anywhere after `go install`

---

## 🚀 Installation

### Prerequisites

- Go 1.24+
- `protoc` (optional, for gRPC codegen)

### Install via Go

```bash
go install github.com/mohsen-farahani/gokitgen/cmd/gokitgen@latest
```

---

## 🎯 Usage

```bash
gokitgen

Or generate a model directly:

gokitgen model
```

### Example: Generate an Order Service

- Run gokitgen
- Select 📦 Generate Model
- Enter model name: Order
- Add enum: OrderStatus → values: PENDING, CANCELLED
- Add field: Status → type: OrderStatus
- Add field: Amount → type: uint
- Add relation: Market → type: Ref:Market
- Select transport: HTTP + gRPC
- Generate tests: Yes

✅ Output: Fully generated Go Kit service in ./generated/

---

### 🤝 Contributing

Contributions are welcome! Please open an issue or PR.

- Fork the repo
- Create your feature branch (git checkout -b feature/AmazingFeature)
- Commit your changes (git commit -m 'Add some AmazingFeature')
- Push to the branch (git push origin feature/AmazingFeature)
- Open a Pull Request

---

### 📜 License
MIT License

---

### 🙌 Credits
- Go Kit
- Bubble Tea
- GORM
- Testify
