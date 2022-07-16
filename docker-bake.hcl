variable "GO_VERSION" {
  default = "1.18"
}

target "_common" {
  args = {
    GO_VERSION = GO_VERSION
    BUILDKIT_CONTEXT_KEEP_GIT_DIR = 1
  }
}

group "default" {
  targets = ["binary"]
}

target "binary" {
  inherits = ["_common"]
  target = "binary"
  output = ["./bin"]
}

target "artifact" {
  inherits = ["_common"]
  target = "artifact"
  output = ["./dist"]
}

target "artifact-all" {
  inherits = ["artifact"]
  platforms = [
    "darwin/amd64",
    "darwin/arm64",
    "freebsd/amd64",
    "freebsd/386",
    "linux/386",
    "linux/amd64",
    "linux/arm/v5",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64",
    "windows/386",
    "windows/amd64",
    "windows/arm64"
  ]
}

target "vendor" {
  inherits = ["_common"]
  dockerfile = "./hack/vendor.Dockerfile"
  target = "update"
  output = ["."]
}

target "gomod-outdated" {
  inherits = ["_common"]
  dockerfile = "./hack/vendor.Dockerfile"
  target = "outdated"
  no-cache-filter = ["outdated"]
  output = ["type=cacheonly"]
}

group "validate" {
  targets = ["lint", "vendor-validate"]
}

target "lint" {
  inherits = ["_common"]
  dockerfile = "./hack/lint.Dockerfile"
  target = "lint"
  output = ["type=cacheonly"]
}

target "vendor-validate" {
  inherits = ["_common"]
  dockerfile = "./hack/vendor.Dockerfile"
  target = "validate"
  output = ["type=cacheonly"]
}
