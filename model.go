package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/pelletier/go-toml/v2"
)

// ViewState
type ViewState int

const (
	ViewDashboard ViewState = iota
	ViewCreatePR
	ViewMergePR
	ViewBranchList
)

// PanelIndex para dashboard
type PanelIndex int

const (
	PanelPRs PanelIndex = iota
	PanelDetails
	PanelConflicts
)

// CreatePanel para Create PR
type CreatePanel int

const (
	CreateFieldsPanel CreatePanel = iota
	CreateTitlePanel
	CreateDescPanel
)

// CreateField para navegación en Fields
type CreateField int

const (
	FieldSource CreateField = iota
	FieldTarget
)

// MergePanel para navegación en Merge PR (4 paneles)
type MergePanel int

const (
	MergePanelStrategy MergePanel = iota // 0: Merge strategy
	MergePanelOptions                    // 1: Delete options
	MergePanelChecklist                  // 2: Checklist (solo informativo)
	MergePanelCommit                     // 3: Commit message
)

// Strategy dentro del panel de Merge
type MergeStrategy int

const (
	StrategyMergeCommit MergeStrategy = iota
	StrategySquash
	StrategyRebase
)

// PR
type PR struct {
	Number         int    `json:"number"`
	Title         string `json:"title"`
	HeadRefName   string `json:"headRefName"`
	BaseRefName   string `json:"baseRefName"`
	Mergeable     string `json:"mergeable"`
	ReviewDecision string `json:"reviewDecision"`
	State         string `json:"state"`
	Author        struct {
		Login string `json:"login"`
	} `json:"author"`
	URL string `json:"url"`
}

// Reviews info
type ReviewsInfo struct {
	Approved int
	Total    int
	Required int
}

// Branch
type Branch struct {
	Name      string
	IsCurrent bool
}

// Config
type Config struct {
	DefaultTargetBranch string `toml:"default_target_branch"`
	TransparentPanels   bool   `toml:"transparent_panels"`
}

// Model
type Model struct {
	ViewState ViewState

	// Dashboard
	ActivePanel PanelIndex
	PRs         []PR
	SelectedPR  int
	LoadingPRs  bool

	// Branch
	CurrentBranch string
	Branches      []Branch

	// Create PR
	SourceBranch string
	TargetBranch string
	TitleInput   string
	BodyInput    string
	CreatePanel  CreatePanel
	CreateField  CreateField
	CreateMode   string

	// Branch list
	BranchListField CreateField
	SelectedBranch  int

	// Merge PR
	MergeStrategy      MergeStrategy
	DeleteRemoteBranch bool
	DeleteLocalBranch bool
	CommitMessage     string
	MergePanel        MergePanel
	MergeMode         string
	ReviewsInfo       ReviewsInfo

	// Cursores de navegación dentro de cada panel
	MergeCursor       int // 0-2: Merge Commit, Squash, Rebase
	OptionsCursor     int // 0-1: Delete remote, Delete local

	// PR Status
	PRMergeStatus struct {
		HasConflicts bool
		Conflicts    []string
	}

	// Flags
	Flags Flags

	// Config
	Config Config

	// Mensajes
	Error     string
	StatusMsg string
}

// NewModel
func NewModel(flags Flags) *Model {
	config := loadConfig()

	m := &Model{
		ViewState:         ViewDashboard,
		ActivePanel:       PanelPRs,
		SelectedPR:        0,
		TargetBranch:      config.DefaultTargetBranch,
		MergeStrategy:     StrategyMergeCommit,
		DeleteLocalBranch: true,
		Flags:             flags,
		Config:            config,
	}

	if flags.CreatePR {
		m.ViewState = ViewCreatePR
		m.CreatePanel = CreateFieldsPanel
		m.CreateField = FieldTarget
		m.CreateMode = "normal"
	} else if flags.MergePR {
		m.ViewState = ViewMergePR
		m.MergePanel = MergePanelStrategy
		m.MergeStrategy = StrategyMergeCommit
		m.MergeMode = "normal"
		if flags.MergePRNum > 0 {
			m.SelectedPR = flags.MergePRNum
		}
	}

	return m
}

// loadConfig
func loadConfig() Config {
	config := Config{
		DefaultTargetBranch: "master",
		TransparentPanels:   true,
	}

	paths := []string{
		"tuipr.toml",
		"./tuipr.toml",
		".config/tuipr/config.toml",
		os.Getenv("HOME") + "/.config/tuipr/config.toml",
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err == nil {
			toml.Unmarshal(data, &config)
			break
		}
	}

	return config
}

// Init
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchCurrentBranch(),
		m.fetchBranches(),
		m.fetchPRs(),
	)
}

// fetchCurrentBranch
func (m *Model) fetchCurrentBranch() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("git", "branch", "--show-current")
		out, err := cmd.Output()
		branch := "unknown"
		if err == nil {
			branch = strings.TrimSpace(string(out))
		}
		return BranchMsg{Branch: branch}
	}
}

// fetchBranches
func (m *Model) fetchBranches() tea.Cmd {
	return func() tea.Msg {
		exec.Command("git", "fetch", "--all").Run()

		cmd := exec.Command("git", "branch", "-a", "--format=%(refname:short)")
		out, err := cmd.Output()
		if err != nil {
			return BranchesMsg{Branches: nil}
		}

		current := m.CurrentBranch
		var branches []Branch

		if current != "" && current != "unknown" {
			branches = append(branches, Branch{Name: current, IsCurrent: true})
		}

		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "remotes/origin/HEAD") {
				continue
			}
			line = strings.TrimPrefix(line, "remotes/origin/")
			if line == "" || line == current {
				continue
			}
			branches = append(branches, Branch{Name: line, IsCurrent: false})
		}

		return BranchesMsg{Branches: branches}
	}
}

// fetchPRs
func (m *Model) fetchPRs() tea.Cmd {
	return func() tea.Msg {
		m.LoadingPRs = true

		cmd := exec.Command("gh", "pr", "list",
			"--json", "number,title,headRefName,baseRefName,mergeable,reviewDecision,state,author,url",
			"--state", "open",
			"--limit", "50",
		)

		out, err := cmd.Output()
		if err != nil {
			m.LoadingPRs = false
			return PRsMsg{Error: fmt.Sprintf("gh error: %v", err)}
		}

		var prs []PR
		if err := json.Unmarshal(out, &prs); err != nil {
			m.LoadingPRs = false
			return PRsMsg{Error: fmt.Sprintf("parse error: %v", err)}
		}

		m.LoadingPRs = false
		return PRsMsg{PRs: prs}
	}
}

// fetchPRReviews obtiene reviews de un PR
func (m *Model) fetchPRReviews(prNumber int) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("gh", "pr", "view", fmt.Sprintf("%d", prNumber),
			"--json", "reviewDecision,reviews",
		)

		out, err := cmd.Output()
		if err != nil {
			return ReviewsMsg{}
		}

		var raw struct {
			ReviewDecision string `json:"reviewDecision"`
			Reviews       []struct {
				State string `json:"state"`
			} `json:"reviews"`
		}

		if err := json.Unmarshal(out, &raw); err != nil {
			return ReviewsMsg{}
		}

		approved := 0
		for _, r := range raw.Reviews {
			if r.State == "APPROVED" {
				approved++
			}
		}

		return ReviewsMsg{
			Approved: approved,
			Total:    len(raw.Reviews),
			Required: 1,
		}
	}
}

// Mensajes
type BranchMsg struct {
	Branch string
}

type BranchesMsg struct {
	Branches []Branch
}

type PRsMsg struct {
	PRs   []PR
	Error string
}

type ReviewsMsg struct {
	Approved int
	Total    int
	Required int
}

type CreatePRResultMsg struct {
	Success  bool
	PRNumber int
	Error    string
}

type MergeResultMsg struct {
	Success bool
	Message string
	Error   string
}

type RefreshMsg struct{}
