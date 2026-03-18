# Orbit Attestation Command Reference

Verify, download, and inspect build provenance attestations using Sigstore bundles with in-toto attestation format and SLSA provenance predicates.

**Top-level command:** `orbit attestation` (alias: `attest`)

**Notes:**

- Attestation verification uses Sigstore bundle format.
- Provenance follows the in-toto attestation specification with SLSA predicates.
- All commands support `-o json` and `-o yaml` for structured output.

---

## Table of Contents

- [verify -- Verify artifact provenance](#verify)
- [download -- Download attestation bundle](#download)
- [inspect -- Inspect attestation bundle](#inspect)

---

## verify

Verify an artifact's build provenance against its Sigstore attestation bundle.

Computes the artifact digest (or accepts a pre-computed digest), loads the attestation bundle, parses the SLSA provenance predicate, and verifies the signer identity and source repository.

```
orbit attestation verify <artifact-path-or-digest> [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--bundle` | Path to local attestation bundle file | |
| `--owner` | GitHub org/owner to scope verification | |
| `--repo` | Specific repository (owner/repo) | |
| `--signer-identity` | Expected signer workflow identity | |
| `--digest-algorithm` | Hash algorithm: `sha256`, `sha512` | `sha256` |

**Examples:**

```bash
# Verify a local binary against a bundle
orbit attestation verify ./my-binary --bundle attestation.jsonl

# Verify with owner and signer identity checks
orbit attest verify ./artifact --bundle bundle.json --owner my-org --signer-identity "github.com/my-org/my-repo"

# Verify a pre-computed digest
orbit attestation verify abc123def456... --bundle att.json --digest-algorithm sha256

# Output as JSON
orbit attestation verify ./my-binary --bundle att.json -o json
```

**Output (table format):**

```
Verification: PASSED
Digest:       sha256:abc123def456...
Signer:       https://github.com/my-org/my-repo/.github/workflows/release.yml@refs/tags/v1.0.0
Builder:      https://github.com/actions/runner
Build Type:   https://slsa.dev/provenance/v1
Source:       git+https://github.com/my-org/my-repo@refs/tags/v1.0.0
Commit:       abc123def456
Materials:    3
```

---

## download

Download an attestation bundle for an artifact digest.

```
orbit attestation download <artifact-digest> [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--repo` | Repository (owner/repo) — **required** | |
| `--digest-algorithm` | Hash algorithm: `sha256`, `sha512` | `sha256` |

**Examples:**

```bash
# Download attestation bundle
orbit attestation download sha256:abc123... --repo owner/repo

# With explicit algorithm
orbit attest download abc123... --repo owner/repo --digest-algorithm sha256
```

---

## inspect

Display the contents of an attestation bundle, including SLSA provenance, signer identity, and build information.

```
orbit attestation inspect <bundle-file> [flags]
```

**Examples:**

```bash
# Inspect a bundle file
orbit attestation inspect attestation.jsonl

# Output as JSON for processing
orbit attest inspect bundle.json -o json
```

**Output (table format):**

```
Media Type:  application/vnd.dev.sigstore.bundle.v0.3+json
Payload:     application/vnd.in-toto+json
Signatures:  1
Signer:      https://github.com/my-org/my-repo/.github/workflows/release.yml@refs/tags/v1.0.0

Provenance:
  Builder:    https://github.com/actions/runner
  Build Type: https://slsa.dev/provenance/v1
  Source:     git+https://github.com/my-org/my-repo@refs/tags/v1.0.0
  Entry:      .github/workflows/release.yml
  sha1:       abc123def456
  Materials:
    - git+https://github.com/my-org/my-repo@refs/tags/v1.0.0
      sha1: abc123def456
```
