#!/bin/bash

# Create several example projects to demonstrate different use cases

set -e

echo "🚀 Creating example Lighthouse projects..."

# Make sure we're in the lighthouse directory
cd "$(dirname "$0")/.."

# Create examples directory for projects
mkdir -p examples/projects
cd examples/projects

# 1. Open Source Software Project
echo "📦 Creating open source software project..."
../../bin/lighthouse project create "BSV Wallet Library" \
    --goal 10.0 \
    --address "1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q" \
    --description "Fund development of a comprehensive BSV wallet library with full SPV support, advanced script templates, and easy-to-use APIs for developers." \
    --min-pledge 0.01 \
    --output "bsv-wallet-library.lighthouse"

# 2. Educational Content
echo "📚 Creating educational content project..."
../../bin/lighthouse project create "BSV Developer Course" \
    --goal 3.5 \
    --address "1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q" \
    --description "Create comprehensive video tutorials and documentation teaching BSV development, from basics to advanced topics including smart contracts and overlay networks." \
    --min-pledge 0.005 \
    --output "bsv-education.lighthouse"

# 3. Hardware Project
echo "🔧 Creating hardware project..."
../../bin/lighthouse project create "BSV Hardware Wallet" \
    --goal 25.0 \
    --address "1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q" \
    --description "Design and manufacture secure BSV hardware wallets with advanced features including multi-signature support, custom scripts, and easy recovery." \
    --min-pledge 0.1 \
    --output "bsv-hardware-wallet.lighthouse"

# 4. Community Event
echo "🎉 Creating community event project..."
../../bin/lighthouse project create "BSV Conference 2024" \
    --goal 8.0 \
    --address "1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q" \
    --description "Fund the annual BSV developer conference featuring workshops, presentations, and networking opportunities for the global BSV community." \
    --min-pledge 0.02 \
    --output "bsv-conference.lighthouse"

# 5. Research Project
echo "🔬 Creating research project..."
../../bin/lighthouse project create "Scaling Research" \
    --goal 15.0 \
    --address "1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q" \
    --description "Research project investigating BSV network scaling solutions, analyzing performance metrics, and developing optimization strategies for enterprise adoption." \
    --min-pledge 0.05 \
    --output "scaling-research.lighthouse"

echo ""
echo "✅ Created 5 example projects!"
echo ""
echo "📋 Project summaries:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for project in *.lighthouse; do
    echo ""
    ../../bin/lighthouse project view "$project"
    echo "Status: $(../../bin/lighthouse project status "$project" | grep "Status:" | cut -d: -f2 | xargs)"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
done

echo ""
echo "🎯 Usage examples:"
echo "• View any project: lighthouse project view examples/projects/[filename]"
echo "• Check status: lighthouse project status examples/projects/[filename]"
echo "• Create pledges: lighthouse pledge create examples/projects/[filename] --amount X --wif [your-key] --utxo [txid:vout:satoshis]"
echo ""
echo "💡 These projects demonstrate different funding goals and use cases for the BSV ecosystem!"