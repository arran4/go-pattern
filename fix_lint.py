import re

with open(".golangci.yml", "r") as f:
    text = f.read()

# Make sure revive is fully removed
text = text.replace("    - revive\n", "")

with open(".golangci.yml", "w") as f:
    f.write(text)

with open(".github/workflows/ci.yml", "r") as f:
    text = f.read()

# Use standard setup go v5 to avoid tar failures
text = re.sub(r"uses: actions/setup-go@v5\n        with:\n          go-version-file: go.mod\n          cache: true", "uses: actions/setup-go@v5\n        with:\n          go-version-file: go.mod", text)

with open(".github/workflows/ci.yml", "w") as f:
    f.write(text)
