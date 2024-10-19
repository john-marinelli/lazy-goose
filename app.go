package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

type RevFileResponse struct {
	Failed  []string `json:"failed"`
	Nodes   []Node   `json:"nodes"`
	Edges   []Edge   `json:"edges"`
	Success bool     `json:"success"`
}

type Node struct {
	Id       string       `json:"id"`
	Data     NodeData     `json:"data"`
	Position NodePosition `json:"position"`
}

type NodeData struct {
	Label string `json:"label"`
}

type NodePosition struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Edge struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type Revision struct {
	Id   string `json:"id"`
	Down string `json:"down"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetFiles() RevFileResponse {
	verDir := "/Users/johnm/dev/test-api/project/migrations/versions"
	d, err := os.ReadDir(verDir)
	if err != nil {
		panic(err)
	}

	var allRevs []Revision
	var nodes []Node
	var edges []Edge
	var failed []string

	for _, f := range d {
		if f.IsDir() {
			continue
		}

		revBytes, err := os.ReadFile(filepath.Join(verDir, f.Name()))
		if err != nil {
			panic(err)
		}
		varMap, err := getVarMap(
			string(revBytes),
		)
		if err != nil {
			continue
		}

		_, ok := varMap["revision"]
		if !ok {
			failed = append(failed, f.Name())
			continue
		}

		_, ok = varMap["down_revision"]
		if !ok {
			failed = append(failed, f.Name())
			continue
		}

		allRevs = append(
			allRevs,
			Revision{
				Id:   varMap["revision"],
				Down: varMap["down_revision"],
			},
		)

	}

	for _, r := range allRevs {
		nodes = append(
			nodes,
			Node{
				Id: r.Id,
				Data: NodeData{
					Label: r.Id,
				},
				Position: NodePosition{},
			},
		)
		if !strings.Contains(r.Down, "None") {
			edges = append(
				edges,
				Edge{
					Id:     strings.Join([]string{r.Down, r.Id}, "-"),
					Target: r.Id,
					Source: r.Down,
				},
			)
		}
	}

	return RevFileResponse{
		Success: true,
		Failed:  failed,
		Edges:   edges,
		Nodes:   nodes,
	}
}

func getVarMap(str string) (map[string]string, error) {
	rx := regexp.MustCompile(`(?m)^\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*(.+)$`)
	vm := make(map[string]string)

	// Find all matches
	matches := rx.FindAllStringSubmatch(str, -1)

	// Populate the map
	for _, match := range matches {
		vn := match[1]
		vv := cleanRevId(match[2])
		vm[vn] = vv
	}

	// Print the map
	return vm, nil
}

func cleanRevId(ri string) string {
	return strings.Trim(
		strings.Trim(
			strings.TrimSpace(
				ri,
			),
			"'",
		),
		"\"",
	)
}
