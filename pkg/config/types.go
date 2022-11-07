package config

type (
	Schef struct {
		Member map[string]Members
	}

	Options struct {
		Build string
		Dev   string
	}

	Members struct {
		Type       string
		Target     string
		EntryPoint string
		Cache      bool
		OutDir     string
		Options    Options
	}

	Workspace struct {
		Members []WorkspaceMember
		Cwd     string
		Name    string
	}

	WorkspaceMember struct {
		Name     string
		BuildCmd string
		DevCmd   string
		Cache    bool
	}

	Cache struct {
		Workspaces []CacheWorkspace
	}

	CacheWorkspace struct {
		Name    string
		Members []CacheMember
	}

	CacheMember struct {
		Name string
		Hash string
	}
)

// Rust

const RustBuild = "cargo build --release"
const RustDev = "cargo build"
const RustBuildWasm = "cargo build --release --target wasm32-unknown-unknown"
const RustDevWasm = "cargo build --target wasm32-unknown-unknown"

// Go

const GoBuild = "go build"
