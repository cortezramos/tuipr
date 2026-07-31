package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/pelletier/go-toml/v2"
)

// ViewState represents the current view in the application.
type ViewState int

// View states.
const (
	ViewDashboard ViewState = iota
	ViewCreatePR
	ViewMergePR
	ViewBranchList
)

// PanelIndex represents panels in the dashboard view.
type PanelIndex int

// Dashboard panels.
const (
	PanelPRs PanelIndex = iota
	PanelDetails
	PanelConflicts
)

// CreatePanel represents panels in the Create PR view.
type CreatePanel int

// Create PR panels.
const (
	CreateFieldsPanel CreatePanel = iota
	CreateTitlePanel
	CreateDescPanel
)

// CreateField represents fields in the Create PR form.
type CreateField int

// Create PR form fields.
const (
	FieldSource CreateField = iota
	FieldTarget
)

// MergePanel represents panels in the Merge PR view.
type MergePanel int

// Merge PR panels.
const (
	MergePanelStrategy MergePanel = iota // Merge strategy selection.
	MergePanelOptions                   // Delete branch options.
	MergePanelChecklist                 // Information checklist.
	MergePanelCommit                    // Commit message input.
)

// MergeStrategy represents the git merge strategy.
type MergeStrategy int

// Merge strategies.
const (
	StrategyMergeCommit MergeStrategy = iota
	StrategySquash
	StrategyRebase
)

// PR represents a GitHub pull request.
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

// ReviewsInfo contains review statistics for a PR.
type ReviewsInfo struct {
	Approved int
	Total    int
	Required int
}

// Branch represents a git branch.
type Branch struct {
	Name      string
	IsCurrent bool
}

// Config represents application configuration.
type Config struct {
	DefaultTargetBranch string `toml:"default_target_branch"`
	TransparentPanels   bool   `toml:"transparent_panels"`
}

// Error definitions.
var (
	ErrNoTitle        = errors.New("title is required")
	ErrNoPRSelected   = errors.New("no PR selected")
	ErrPRHasConflicts = errors.New("cannot merge: has conflicts")
)

// Mode constants.
const (
	modeNormal = "normal"
	modeInsert = "insert"
)

// Model represents the application state.
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
	MergeStrategy       MergeStrategy
	DeleteRemoteBranch bool
	DeleteLocalBranch  bool
	CommitMessage      string
	MergePanel         MergePanel
	MergeMode          string
	ReviewsInfo        ReviewsInfo

	// Navigation cursors within panels.
	MergeCursor   int // 0-2: Merge Commit, Squash, Rebase.
	OptionsCursor int // 0-1: Delete remote, Delete local.

	// PR Status.
	PRMergeStatus struct {
		HasConflicts bool
		Conflicts    []string
	}

	// Flags
	Flags Flags

	// Config
	Config Config

	// Messages.
	Error     string
	StatusMsg string
}

// NewModel creates a new Model with the given flags.
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
		m.CreateMode = modeNormal
	} else if flags.MergePR {
		m.ViewState = ViewMergePR
		m.MergePanel = MergePanelStrategy
		m.MergeStrategy = StrategyMergeCommit
		m.MergeMode = modeNormal
		if flags.MergePRNum > 0 {
			m.SelectedPR = flags.MergePRNum
		}
	}

	return m
}

// loadConfig loads configuration from known paths.
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

// Init initializes the model with background commands.
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchCurrentBranch(),
		m.fetchBranches(),
		m.fetchPRs(),
	)
}

// fetchCurrentBranch retrieves the current git branch.
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

// fetchBranches retrieves all git branches.
func (m *Model) fetchBranches() tea.Cmd {
	return func() tea.Msg {
		exec.Command("git", "fetch", "--all").Run()

		cmd := exec.Command("git", "branch", "-a", "--format=%(refname:short)")
		out, err := cmd.Output()
		if err != nil {
			return BranchesMsg{}
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

// fetchPRs retrieves open pull requests from GitHub.
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

// fetchPRReviews retrieves review information for a PR.
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
			Reviews        []struct {
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

// BranchMsg indicates the current branch was updated.
type BranchMsg struct {
	Branch string
}

// BranchesMsg contains the list of git branches.
type BranchesMsg struct {
	Branches []Branch
}

// PRsMsg contains the list of pull requests.
type PRsMsg struct {
	PRs   []PR
	Error string
}

// ReviewsMsg contains review statistics for a PR.
type ReviewsMsg struct {
	Approved int
	Total    int
	Required int
}

// CreatePRResultMsg indicates the result of creating a PR.
type CreatePRResultMsg struct {
	Success  bool
	PRNumber int
	Error    string
}

// MergeResultMsg indicates the result of merging a PR.
type MergeResultMsg struct {
	Success bool
	Message string
	Error   string
}

// RefreshMsg triggers a data refresh.
type RefreshMsg struct{}
