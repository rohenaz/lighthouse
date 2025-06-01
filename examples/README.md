# Lighthouse Examples

This directory contains example scripts and projects demonstrating how to use the Lighthouse BSV crowdfunding system.

## Getting Started

### Prerequisites
- Lighthouse CLI built and available (run `make cli` from the root directory)
- Basic understanding of BSV addresses and transactions

### Quick Start

1. **Create a simple project:**
   ```bash
   chmod +x create-project.sh
   ./create-project.sh
   ```

2. **Create multiple example projects:**
   ```bash
   chmod +x example-projects.sh
   ./example-projects.sh
   ```

## Example Projects

The `example-projects.sh` script creates five different types of crowdfunding projects:

### 1. BSV Wallet Library (10 BSV goal)
Open source software development project for creating comprehensive BSV wallet libraries.

### 2. BSV Developer Course (3.5 BSV goal)
Educational content creation for BSV development tutorials and documentation.

### 3. BSV Hardware Wallet (25 BSV goal)
Hardware development project for secure BSV storage devices.

### 4. BSV Conference 2024 (8 BSV goal)
Community event funding for developer conferences and meetups.

### 5. Scaling Research (15 BSV goal)
Research project investigating BSV network scaling and optimization.

## File Structure

After running the examples, you'll have:

```
examples/
├── README.md                    # This file
├── create-project.sh           # Simple project creation
├── example-projects.sh         # Multiple example projects
└── projects/                   # Generated project files
    ├── bsv-wallet-library.lighthouse
    ├── bsv-education.lighthouse
    ├── bsv-hardware-wallet.lighthouse
    ├── bsv-conference.lighthouse
    └── scaling-research.lighthouse
```

## CLI Usage Examples

### Project Management
```bash
# Create a new project
lighthouse project create "My Project" \
  --goal 5.0 \
  --address "1YourBSVAddress123..." \
  --description "Project description"

# View project details
lighthouse project view my-project.lighthouse

# Check funding status
lighthouse project status my-project.lighthouse

# Claim funds when goal is reached
lighthouse project claim my-project.lighthouse --broadcast
```

### Pledge Management
```bash
# Create a pledge (requires real UTXOs)
lighthouse pledge create my-project.lighthouse \
  --amount 0.5 \
  --wif "YourPrivateKeyInWIFFormat" \
  --utxo "txid:vout:satoshis" \
  --message "Supporting this great project!"

# View pledge details
lighthouse pledge view my-pledge.pledge

# Revoke a pledge (spend the UTXO elsewhere)
lighthouse pledge revoke my-pledge.pledge \
  --wif "YourPrivateKeyInWIFFormat"
```

## Understanding Assurance Contracts

Lighthouse uses **assurance contracts** - a trustless crowdfunding mechanism where:

1. **Pledgers** create partial transactions with `SIGHASH_ANYONECANPAY` signatures
2. **Pledges** can be combined when the funding goal is reached
3. **Funds** are only transferred when the goal is met
4. **Pledgers** can revoke their pledges anytime before claiming

### Key Benefits
- **Trustless**: No need to trust the project creator with funds
- **Revocable**: Pledgers can change their mind before goal is reached
- **Atomic**: Either the goal is reached and everyone pays, or no one pays
- **Minimal fees**: Only pay Bitcoin network fees, no platform fees

## Testing with Regtest

For testing purposes, you can use regtest (local Bitcoin network):

1. Set up a regtest environment
2. Generate test addresses and UTXOs
3. Use the generated UTXOs in pledge commands
4. Test the complete pledge → claim workflow

## Next Steps

- Try creating your own projects
- Test the pledge workflow with regtest
- Build a web interface using the CLI as a backend
- Integrate with BSV wallets and services

## Security Notes

- Never share your private keys (WIF format)
- Always verify addresses before sending funds
- Test with small amounts first
- Use proper UTXO management in production

## Support

For questions and issues, refer to the main Lighthouse documentation or check the CLI help:

```bash
lighthouse --help
lighthouse project --help
lighthouse pledge --help
```