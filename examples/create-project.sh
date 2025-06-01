#!/bin/bash

# Example script to create a sample crowdfunding project
# This demonstrates how to use the Lighthouse CLI

set -e

echo "🏗️  Creating a sample crowdfunding project..."

# Create a community garden project
./bin/lighthouse project create "Community Garden Project" \
    --goal 5.0 \
    --address "1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q" \
    --description "Help us build a beautiful community garden in our neighborhood! This space will provide fresh vegetables, a place for kids to learn about nature, and bring our community together." \
    --min-pledge 0.001

echo ""
echo "✅ Project created! Check the generated .lighthouse file"

# View the project details
echo ""
echo "📋 Project details:"
./bin/lighthouse project view Community_Garden_Project.lighthouse

echo ""
echo "📊 Project status:"
./bin/lighthouse project status Community_Garden_Project.lighthouse

echo ""
echo "🎯 Next steps:"
echo "1. Share the Community_Garden_Project.lighthouse file with potential supporters"
echo "2. Supporters can create pledges using: lighthouse pledge create"
echo "3. When funding goal is reached, claim funds with: lighthouse project claim"