# 💡 Lighthouse BSV

**Decentralized crowdfunding using Bitcoin assurance contracts**

A trustless crowdfunding platform where pledges are only collected when the funding goal is met. No central authority, no platform fees, just pure Bitcoin magic! ✨

Based on Mike Hearn's pioneering Lighthouse project, rebuilt for BSV blockchain with modern tooling.

---

## 🚀 **Quick Start**

### **Running the Complete System**

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/lighthouse
cd lighthouse

# 2. Build the CLI tool
make build

# 3. Start the web interface
cd web/lighthouse-web
bun install
bun dev

# 4. Open your browser
open http://localhost:7184
```

### **What You Get**

- 🖥️ **Web Interface**: Modern React app with Bitcoin authentication at `http://localhost:7184`
- ⚡ **CLI Tool**: Command-line interface for power users
- 🔐 **Bitcoin Auth**: No passwords, just cryptographic signatures
- 💰 **Assurance Contracts**: SIGHASH_ANYONECANPAY magic on BSV blockchain

---

## 🏗️ **Architecture**

### **System Components**

```
┌─────────────────────────────────────────────────────────────┐
│                   🌐 Web Interface                          │
│              Next.js + BigBlocks + BSV Auth                 │
│                   (Port 7184)                               │
└─────────────────┬───────────────────────────────────────────┘
                  │ HTTP API calls
┌─────────────────▼───────────────────────────────────────────┐
│                   🔧 Go CLI Backend                         │
│            Project & Pledge Management                      │
│              BSV SDK + Protocol Buffers                     │
└─────────────────┬───────────────────────────────────────────┘
                  │ BSV transactions
┌─────────────────▼───────────────────────────────────────────┐
│                  ⛓️  BSV Blockchain                         │
│              Assurance Contract Magic                       │
│          SIGHASH_ANYONECANPAY signatures                    │
└─────────────────────────────────────────────────────────────┘
```

### **Core Technologies**

- **Backend**: Go + BSV SDK + Protocol Buffers
- **Frontend**: Next.js 15 + TypeScript + BigBlocks
- **Blockchain**: BSV with assurance contracts
- **Authentication**: Bitcoin cryptographic signatures
- **Styling**: Tailwind CSS v4 + Original Lighthouse design

---

## 📦 **Installation & Setup**

### **Prerequisites**

- **Go 1.21+** - For building the CLI tool
- **Bun 1.0+** - For the web interface
- **BSV Wallet** - For testing pledges (can use testnet)

### **Step 1: Build the CLI**

```bash
# Clone and build
git clone https://github.com/yourusername/lighthouse
cd lighthouse

# Build the CLI tool
go mod tidy
go build -o bin/lighthouse cmd/lighthouse/main.go

# Test the CLI
./bin/lighthouse --help
```

### **Step 2: Setup Web Interface**

```bash
# Navigate to web interface
cd web/lighthouse-web

# Install dependencies
bun install

# Set environment variables
cp .env.example .env.local
# Edit .env.local with your settings

# Start development server
bun dev
```

### **Step 3: Create Example Projects**

```bash
# From the lighthouse root directory
./examples/example-projects.sh
```

---

## 🎯 **Usage Examples**

### **Web Interface (Recommended)**

1. **Browse Projects**: http://localhost:7184/projects
2. **Create Account**: Bitcoin-based authentication, no passwords!
3. **Create Project**: Set funding goal and description
4. **Make Pledges**: Support projects with SIGHASH_ANYONECANPAY
5. **Manage Pledges**: Revoke anytime before goal is reached

### **CLI Interface (Power Users)**

```bash
# Create a new crowdfunding project
./bin/lighthouse project create "Community Garden Project" \
  --goal 5.0 \
  --address "1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q" \
  --description "Help us build a beautiful community garden!" \
  --min-pledge 0.001

# View project details
./bin/lighthouse project view Community_Garden_Project.lighthouse

# Check funding status
./bin/lighthouse project status Community_Garden_Project.lighthouse

# Make a pledge (requires WIF and UTXO)
./bin/lighthouse pledge create Community_Garden_Project.lighthouse \
  --amount 0.5 \
  --wif "L1aW4aubDFB7yfras2S1mN3bqg9nwySY8nkoLmJebSLD5BWv3ENZ" \
  --utxo "txid:vout:satoshis"

# Claim funds when goal is reached
./bin/lighthouse project claim Community_Garden_Project.lighthouse
```

---

## 🛠️ **Development**

### **Project Structure**

```
lighthouse/
├── 📁 core/              # Go library - assurance contract logic
├── 📁 cmd/               # CLI implementation
├── 📁 proto/             # Protocol buffer definitions  
├── 📁 web/lighthouse-web/ # Next.js frontend application
├── 📁 examples/          # Example projects and scripts
├── 📁 docs/              # Documentation
├── 🔧 Makefile          # Build automation
├── 📋 go.mod            # Go dependencies
└── 📖 README.md         # This file
```

### **Development Commands**

```bash
# Build everything
make build

# Run tests
make test
go test ./...

# Clean build artifacts
make clean

# Generate protobuf files
make proto

# Web development
cd web/lighthouse-web
bun dev          # Development server
bun build        # Production build
bun test         # Run tests
```

### **Adding New Features**

1. **Backend**: Add CLI commands in `cmd/lighthouse/`
2. **Frontend**: Add React components in `web/lighthouse-web/src/`
3. **Protocol**: Extend protobuf definitions in `proto/`
4. **API**: Add API routes in `web/lighthouse-web/src/app/api/`

---

## ⚡ **How Assurance Contracts Work**

### **The Magic of SIGHASH_ANYONECANPAY**

```
1. 🎯 Project Created     → Set funding goal & BSV address
2. 💰 Pledges Made        → Partial transactions with special signatures
3. 🔄 Revocable Anytime   → Get your money back before goal reached  
4. 🎉 Goal Reached        → All pledges automatically combine and pay out
5. ✅ All or Nothing      → If goal not met, nobody pays anything
```

### **Key Benefits**

- **No Central Authority**: No platform holding your funds
- **Cryptographically Secure**: Bitcoin blockchain guarantees
- **Revocable Pledges**: Change your mind anytime before completion
- **No Platform Fees**: Only standard Bitcoin transaction fees
- **Instant Settlement**: Automatic payout when goal is reached

---

## 🌐 **API Reference**

### **Web API Endpoints**

```
GET    /api/projects          # List all projects
POST   /api/projects          # Create new project
GET    /api/projects/[id]     # Get project details
POST   /api/projects/[id]     # Pledge to project or claim funds

GET    /api/pledges           # List user's pledges  
DELETE /api/pledges           # Revoke a pledge

GET    /api/profile           # Get user profile
POST   /api/profile           # Update user profile
```

### **CLI Commands**

```bash
# Project management
lighthouse project create <title> [options]
lighthouse project view <file>
lighthouse project status <file>
lighthouse project claim <file>

# Pledge management  
lighthouse pledge create <project> [options]
lighthouse pledge view <file>
lighthouse pledge revoke <project> [options]

# Utility commands
lighthouse --help
lighthouse --version
```

---

## 🔧 **Configuration**

### **Environment Variables**

```bash
# Web Interface (.env.local)
NEXT_PUBLIC_APP_NAME="Lighthouse BSV"
NEXT_PUBLIC_APP_DESCRIPTION="BSV Crowdfunding with Assurance Contracts"
NEXT_PUBLIC_APP_URL="http://localhost:7184"

# Optional API keys for enhanced features
ANTHROPIC_API_KEY=your_key_here    # For AI assistance
REDIS_URL=redis://localhost:6379   # For production scaling
```

### **CLI Configuration**

```bash
# Set default network (testnet recommended for development)
export LIGHTHOUSE_NETWORK=testnet

# Set default data directory
export LIGHTHOUSE_DATA_DIR=~/.lighthouse
```

---

## 🧪 **Testing**

### **Run the Test Suite**

```bash
# Test the Go backend
go test ./...

# Test the web interface  
cd web/lighthouse-web
bun test

# Integration tests with example projects
./examples/example-projects.sh
```

### **Manual Testing Workflow**

1. **Start the web interface**: `bun dev`
2. **Create a Bitcoin identity**: Sign up with no password
3. **Create a test project**: Set a small funding goal
4. **Make a pledge**: Use testnet BSV
5. **Test revocation**: Revoke your pledge
6. **Test completion**: Make enough pledges to reach goal

---

## 🏛️ **Tribute to Original Lighthouse**

This project is inspired by **Mike Hearn's original Lighthouse** (2014-2016), the first implementation of Bitcoin assurance contracts. Key quotes from the original design doc:

> *"The goal of the Lighthouse design is to keep as much logic out of the server and in the fat client as possible... aiming for a highly decentralised design in which it's feasible for individuals with no sysadmin ability to create and run crowdfunding campaigns."*

### **Original vs BSV Implementation**

| Feature | Original (2014) | BSV Version (2024) |
|---------|----------------|-------------------|
| **Blockchain** | Bitcoin Core | BSV (unlimited scale) |
| **Interface** | JavaFX Desktop | Next.js Web + CLI |
| **Authentication** | Desktop wallet | Bitcoin cryptographic signatures |
| **Smart Contracts** | SIGHASH_ANYONECANPAY | Enhanced with BSV features |
| **File Format** | .lighthouse files | Compatible + protocol buffers |

---

## 📜 **License**

Apache 2.0 - See [LICENSE](LICENSE) file

## 🤝 **Contributing**

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🙏 **Acknowledgments**

- **Mike Hearn** - Original Lighthouse creator and Bitcoin pioneer
- **BSV Community** - For supporting unlimited blockchain scaling
- **BigBlocks Team** - For Bitcoin authentication components
- **Go BSV SDK** - For robust blockchain integration

---

## 🆘 **Support**

- **Documentation**: See `/docs` directory
- **Issues**: GitHub Issues for bugs and feature requests  
- **Discussions**: GitHub Discussions for questions
- **Email**: lighthouse-bsv@example.com

**Ready to revolutionize crowdfunding? Let's build the future together! 🚀**