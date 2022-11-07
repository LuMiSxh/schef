package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gookit/color"
	"os"
	"strings"
)

// GetSchef returns the workspace (Schef) in which the program got invoked
func GetSchef() *Schef {
	file, err := os.ReadFile("./schef.toml")
	if err != nil {
		color.Error.Block(err.Error())
		os.Exit(1)
	}

	var workspace Schef

	if err := toml.Unmarshal(file, &workspace); err != nil {
		color.Error.Block(err.Error())
		os.Exit(1)
	}

	return &workspace
}

// GenerateWorkspace generates the Workspace struct and generates all commands required
func (s *Schef) GenerateWorkspace() *Workspace {
	workspace := new(Workspace)
	workspace.Cwd, _ = os.Getwd()

	idx := strings.LastIndex(workspace.Cwd, "\\")
	workspace.Name = workspace.Cwd[idx+1 : len(workspace.Cwd)-1]

	for k, member := range s.Member {
		buildCmd := ""
		devCmd := ""

		switch member.Type {
		case "rust":
			{
				if strings.ToLower(member.Target) == "wasm" {
					buildCmd += RustBuildWasm
					devCmd += RustDevWasm
				} else {
					buildCmd += RustBuild
					devCmd += RustDev
				}

				if member.OutDir != "" {
					buildCmd += fmt.Sprintf(" --target-dir ../%s", member.OutDir)
					devCmd += fmt.Sprintf(" --target-dir ../%s/dev", member.OutDir)
				}
			}
		case "go":
			{
				buildCmd += GoBuild
				devCmd += GoBuild

				if member.EntryPoint != "" {
					buildCmd += fmt.Sprintf(" %s", member.EntryPoint)
					devCmd += fmt.Sprintf(" %s", member.EntryPoint)
				}

				if member.OutDir != "" {
					buildCmd += fmt.Sprintf(" --o ../%s", member.OutDir)
					devCmd += fmt.Sprintf(" --o ../%s/dev", member.OutDir)
				}
			}
		default:
			{
				buildCmd = member.Options.Build
				devCmd = member.Options.Dev
			}
		}

		if member.Options.Dev != "" {
			devCmd = member.Options.Dev
		}

		if member.Options.Build != "" {
			buildCmd = member.Options.Build
		}

		workspace.Members = append(workspace.Members, WorkspaceMember{Name: k, BuildCmd: buildCmd, DevCmd: devCmd, Cache: member.Cache})
	}

	return workspace
}
