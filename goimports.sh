#!/bin/bash
# Copyright (C) 2020 Reactive Markets Limited. All rights reserved.

set -euo pipefail

script=$(cat <<'EOF'
/^import [(]$/ { print $0; import=1; next }
/^[)]$/        { import=0 }
               { if (!import || length($0) > 0) { print $0 } }
EOF
)

git ls-files "*.go" | while read file; do
    if [[ "$file" != internal/app/sbetool/golang/template.go ]]; then
        awk "$script" <$file >tmp.$$ && mv tmp.$$ $file
    fi
done

gofmt -s -w .
goimports -w .
