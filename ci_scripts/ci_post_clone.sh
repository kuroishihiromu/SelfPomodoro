#!/bin/zsh
set -euo pipefail
mkdir -p ~/Library/org.swift.swiftpm/security/
cp "$(dirname "$0")/macros.json" ~/Library/org.swift.swiftpm/security/
echo "[ci_post_clone] Installed macros.json to ~/Library/org.swift.swiftpm/security/"
