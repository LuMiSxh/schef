package config

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"os"
	"path/filepath"
)

func ReadCache() *Cache {
	dir, userErr := os.UserHomeDir()
	if userErr != nil {
		color.Error.Block(userErr.Error())
		os.Exit(1)
	}
	conf, err := filepath.Abs(fmt.Sprintf("%s/.schef.json", dir))
	if err != nil {
		_, cErr := os.Create(fmt.Sprintf("%s/.schef.json", dir))
		color.Info.Prompt("Cache file could not be found and will be created at:  %s", string(conf))
		if cErr != nil {
			color.Error.Block(userErr.Error())
			os.Exit(1)
		}
	}
	conf, _ = filepath.Abs(fmt.Sprintf("%s/.schef.json", dir))

	file, err := os.ReadFile(conf)
	if err != nil {
		color.Error.Block(err.Error())
		os.Exit(1)
	}
	var c Cache
	if err := json.Unmarshal(file, &c); err != nil {
		color.Error.Block(err.Error())
		os.Exit(1)
	}

	return &c
}

func (c *Cache) WriteCache() {
	content, err := json.Marshal(c)
	if err != nil {
		color.Error.Block(err.Error())
		os.Exit(1)
	}

	dir, userErr := os.UserHomeDir()
	if userErr != nil {
		color.Error.Block(userErr.Error())
		os.Exit(1)
	}
	conf, err := filepath.Abs(fmt.Sprintf("%s/.schef.json", dir))
	if err != nil {
		color.Error.Block(userErr.Error())
		os.Exit(1)

	}

	if osErr := os.WriteFile(conf, content, 0644); osErr != nil {
		color.Warn.Block("%s\n%s", "The updated cache could not be saved", osErr.Error())
	}
}

func (c *Cache) IncludesWorkspace(workspace string) bool {
	for _, w := range c.Workspaces {
		if w.Name == workspace {
			return true
		}
	}
	return false
}

func (c *Cache) IncludesWorkspaceMember(workspace string, member string) bool {
	for _, w := range c.Workspaces {
		if w.Name == workspace {
			for _, m := range w.Members {
				if m.Name == member {
					return true
				}
			}
		}
	}
	return false
}

func (c *Cache) AddWorkspace(workspace string, members []CacheMember) {
	c.Workspaces = append(c.Workspaces, CacheWorkspace{Name: workspace, Members: members})
}

func (c *Cache) AddWorkspaceMember(workspace string, member CacheMember) {
	for i, w := range c.Workspaces {
		if w.Name == workspace {
			c.Workspaces[i].Members = append(c.Workspaces[i].Members, member)
			return
		}
	}
}

func (c *Cache) RemoveWorkspace(workspace string) {
	for i, w := range c.Workspaces {
		if w.Name == workspace {
			c.Workspaces[i] = c.Workspaces[len(workspace)-1]
			return
		}
	}
}

func (c *Cache) RemoveWorkspaceMember(workspace string, member string) {
	for i, w := range c.Workspaces {
		if w.Name == workspace {
			for j, m := range w.Members {
				if m.Name == member {
					c.Workspaces[i].Members[j] = c.Workspaces[i].Members[len(c.Workspaces[i].Members)-1]
					return
				}
			}
		}
	}
}

func (c *Cache) GetWorkspace(workspace string) *CacheWorkspace {
	for _, w := range c.Workspaces {
		if w.Name == workspace {
			return &w
		}
	}
	return nil
}

func (c *Cache) GetWorkspaceMember(workspace string, member string) *CacheMember {
	for _, w := range c.Workspaces {
		if w.Name == workspace {
			for _, m := range w.Members {
				if m.Name == member {
					return &m
				}
			}
		}
	}
	return nil
}

func (c *Cache) SetHash(workspace string, member string, hash string) {
	for i, w := range c.Workspaces {
		if w.Name == workspace {
			for j, m := range w.Members {
				if m.Name == member {
					c.Workspaces[i].Members[j] = CacheMember{Name: member, Hash: hash}
				}
			}
		}
	}
}
