package config

import (
	"fmt"
	"github.com/gookit/color"
	"golang.org/x/mod/sumdb/dirhash"
	"os"
	"os/exec"
	"sync"
)

// Fetches cwd and generates Workspace and Member objects
func fetchDetails() *Workspace {
	var schef = GetSchef()

	workspace := schef.GenerateWorkspace()

	return workspace
}

// Goroutine to execute the commands provided
func executeCmd(wg *sync.WaitGroup, path string, cmd string, name string) {
	command := exec.Command("cmd.exe", "/c", cmd)
	command.Dir = fmt.Sprintf("%s", path)

	output, err := command.Output()
	if string(output) != "" {
		color.Info.Block("%s\n\t%s\n", name, string(output))
	}
	if err != nil {
		color.Warn.Block("%s\n\t%s\n", name, err.Error())
	}

	defer wg.Done()
}

// Glue Public method to bind the command executor to the cli
func Glue(member string, method string) {
	w := fetchDetails()
	c := ReadCache()
	wg := new(sync.WaitGroup)

	if member == "" {
		for _, wMember := range w.Members {
			inCache := false

			switch method {
			case "build":
				{
					if c.IncludesWorkspaceMember(w.Name, wMember.Name) && wMember.Cache {
						hash, err := dirhash.HashDir(fmt.Sprintf("%s\\%s", w.Cwd, wMember.Name), wMember.Name, dirhash.DefaultHash)
						if err != nil {
							color.Error.Block("%s\n\t%s\n", wMember.Name, err.Error())
							os.Exit(1)
						}

						if hash == c.GetWorkspaceMember(w.Name, wMember.Name).Hash {
							inCache = true
							color.Info.Block("%s\n\t%s\n", wMember.Name, "Using cached result")
						} else {
							c.SetHash(w.Name, wMember.Name, hash)
						}
					} else if !c.IncludesWorkspace(w.Name) {
						hash, err := dirhash.HashDir(fmt.Sprintf("%s\\%s", w.Cwd, wMember.Name), wMember.Name, dirhash.DefaultHash)
						if err != nil {
							color.Error.Block("%s\n\t%s\n", wMember.Name, err.Error())
							os.Exit(1)
						}

						members := make([]CacheMember, 0)
						members = append(members, CacheMember{Name: wMember.Name, Hash: hash})
						c.AddWorkspace(w.Name, members)
					} else if !c.IncludesWorkspaceMember(w.Name, wMember.Name) {
						hash, err := dirhash.HashDir(fmt.Sprintf("%s\\%s", w.Cwd, wMember.Name), wMember.Name, dirhash.DefaultHash)
						if err != nil {
							color.Error.Block("%s\n\t%s\n", wMember.Name, err.Error())
							os.Exit(1)
						}

						c.AddWorkspaceMember(w.Name, CacheMember{Name: wMember.Name, Hash: hash})
					}

					if !inCache {
						wg.Add(1)
						go executeCmd(wg, fmt.Sprintf("%s\\%s", w.Cwd, wMember.Name), wMember.BuildCmd, wMember.Name)
					}
				}
			case "dev":
				{
					if c.IncludesWorkspaceMember(w.Name, wMember.Name) && wMember.Cache {
						hash, err := dirhash.HashDir(fmt.Sprintf("%s\\%s", w.Cwd, wMember.Name), wMember.Name, dirhash.DefaultHash)
						if err != nil {
							color.Error.Block("%s\n\t%s\n", wMember.Name, err.Error())
							os.Exit(1)
						}

						if hash == c.GetWorkspaceMember(w.Name, wMember.Name).Hash {
							inCache = true
							color.Info.Block("%s\n\t%s\n", wMember.Name, "Using cached result")
						} else {
							c.SetHash(w.Name, wMember.Name, hash)
						}
					} else if !c.IncludesWorkspace(w.Name) {
						hash, err := dirhash.HashDir(fmt.Sprintf("%s\\%s", w.Cwd, wMember.Name), wMember.Name, dirhash.DefaultHash)
						if err != nil {
							color.Error.Block("%s\n\t%s\n", wMember.Name, err.Error())
							os.Exit(1)
						}

						members := make([]CacheMember, 0)
						members = append(members, CacheMember{Name: wMember.Name, Hash: hash})
						c.AddWorkspace(w.Name, members)
					} else if !c.IncludesWorkspaceMember(w.Name, wMember.Name) {
						hash, err := dirhash.HashDir(fmt.Sprintf("%s\\%s", w.Cwd, wMember.Name), wMember.Name, dirhash.DefaultHash)
						if err != nil {
							color.Error.Block("%s\n\t%s\n", wMember.Name, err.Error())
							os.Exit(1)
						}

						c.AddWorkspaceMember(w.Name, CacheMember{Name: wMember.Name, Hash: hash})
					}

					if !inCache {
						wg.Add(1)
						go executeCmd(wg, fmt.Sprintf("%s\\%s", w.Cwd, wMember.Name), wMember.DevCmd, wMember.Name)
					}
				}
			}
		}
	}

	for _, wMember := range w.Members {
		if wMember.Name == member {
			wg.Add(1)

			switch method {
			case "build":
				{
					go executeCmd(wg, fmt.Sprintf("%s\\%s", w.Cwd, wMember.Name), wMember.BuildCmd, wMember.Name)
				}
			case "dev":
				{
					go executeCmd(wg, fmt.Sprintf("%s\\%s", w.Cwd, wMember.Name), wMember.DevCmd, wMember.Name)
				}
			}
		}
	}
	defer c.WriteCache()
	wg.Wait()
}
