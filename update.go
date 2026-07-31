package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming messages and returns the updated model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case BranchMsg:
		m.CurrentBranch = msg.Branch
		m.SourceBranch = msg.Branch
		return m, m.fetchBranches()

	case BranchesMsg:
		m.Branches = msg.Branches
		return m, nil

	case PRsMsg:
		m.LoadingPRs = false
		if msg.Error != "" {
			m.Error = msg.Error
		}
		m.PRs = msg.PRs
		return m, nil

	case ReviewsMsg:
		m.ReviewsInfo.Approved = msg.Approved
		m.ReviewsInfo.Total = msg.Total
		m.ReviewsInfo.Required = msg.Required
		return m, nil

	case CreatePRResultMsg:
		if msg.Success {
			m.StatusMsg = fmt.Sprintf("PR #%d created!", msg.PRNumber)
			m.ViewState = ViewDashboard
			return m, m.fetchPRs()
		}
		m.Error = msg.Error
		return m, nil

	case MergeResultMsg:
		if msg.Success {
			m.StatusMsg = msg.Message
			m.ViewState = ViewDashboard
			return m, m.fetchPRs()
		}
		m.Error = msg.Error
		return m, nil

	case RefreshMsg:
		cmds := []tea.Cmd{m.fetchPRs(), m.fetchCurrentBranch(), m.fetchBranches()}
		if len(m.PRs) > 0 && m.SelectedPR < len(m.PRs) {
			cmds = append(cmds, m.fetchPRReviews(m.PRs[m.SelectedPR].Number))
		}
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	}

	return m, nil
}

// handleKeyMsg processes keyboard input.
func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	inInsertMode := (m.ViewState == ViewCreatePR && m.CreateMode == modeInsert) ||
		(m.ViewState == ViewMergePR && m.MergeMode == modeInsert)

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "q":
		if inInsertMode {
			m.handleCharInput(msg)
			return m, nil
		}
		if m.ViewState == ViewDashboard || m.ViewState == ViewBranchList ||
			m.ViewState == ViewCreatePR || m.ViewState == ViewMergePR {
			return m, tea.Quit
		}

	case "tab", "shift+tab":
		// Only handle Tab in Merge PR view.
		if m.ViewState == ViewMergePR && m.MergeMode != modeInsert {
			return m.handleMergeTabKeys(msg.String())
		}
	}

	switch m.ViewState {
	case ViewDashboard:
		return m.handleDashboardKeys(msg)
	case ViewCreatePR:
		return m.handleCreatePRKeys(msg)
	case ViewMergePR:
		return m.handleMergePRKeys(msg)
	case ViewBranchList:
		return m.handleBranchListKeys(msg)
	}

	return m, nil
}

// ============================================
// DASHBOARD KEYS
// ============================================

func (m *Model) handleDashboardKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return m, nil

	case "1":
		m.ActivePanel = PanelPRs
	case "2":
		m.ActivePanel = PanelDetails
	case "3":
		m.ActivePanel = PanelConflicts

	case "j":
		if m.ActivePanel == PanelPRs && m.SelectedPR < len(m.PRs)-1 {
			m.SelectedPR++
		}

	case "k":
		if m.ActivePanel == PanelPRs && m.SelectedPR > 0 {
			m.SelectedPR--
		}

	case "c":
		m.ViewState = ViewCreatePR
		m.CreatePanel = CreateFieldsPanel
		m.CreateField = FieldTarget
		m.CreateMode = modeNormal

	case "m":
		if len(m.PRs) > 0 && m.SelectedPR < len(m.PRs) {
			m.ViewState = ViewMergePR
			m.MergePanel = MergePanelStrategy
			m.MergeMode = modeNormal
			return m, m.fetchPRReviews(m.PRs[m.SelectedPR].Number)
		}

	case "r":
		cmds := []tea.Cmd{m.fetchPRs(), m.fetchCurrentBranch(), m.fetchBranches()}
		if len(m.PRs) > 0 && m.SelectedPR < len(m.PRs) {
			cmds = append(cmds, m.fetchPRReviews(m.PRs[m.SelectedPR].Number))
		}
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

// ============================================
// CREATE PR KEYS
// ============================================

func (m *Model) handleCreatePRKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" {
		if m.CreateMode == modeInsert {
			m.CreateMode = modeNormal
		} else {
			m.ViewState = ViewDashboard
		}
		return m, nil
	}

	if m.CreateMode == modeNormal {
		switch msg.String() {
		case "1":
			m.CreatePanel = CreateFieldsPanel
		case "2":
			m.CreatePanel = CreateTitlePanel
		case "3":
			m.CreatePanel = CreateDescPanel

		case "tab":
			switch m.CreatePanel {
			case CreateFieldsPanel:
				m.CreatePanel = CreateTitlePanel
			case CreateTitlePanel:
				m.CreatePanel = CreateDescPanel
			case CreateDescPanel:
				m.CreatePanel = CreateFieldsPanel
			}

		case "j":
			if m.CreatePanel == CreateFieldsPanel && m.CreateField == FieldSource {
				m.CreateField = FieldTarget
			}

		case "k":
			if m.CreatePanel == CreateFieldsPanel && m.CreateField == FieldTarget {
				m.CreateField = FieldSource
			}

		case "i":
			m.CreateMode = modeInsert

		case "enter":
			if m.CreatePanel == CreateFieldsPanel && m.CreateField == FieldTarget {
				m.ViewState = ViewBranchList
				m.BranchListField = FieldTarget
				m.SelectedBranch = 0
			}

		case "ctrl+s":
			return m, m.createPR()
		}
		return m, nil
	}

	// Insert mode.
	if m.CreateMode == modeInsert {
		switch msg.String() {
		case "enter":
			if m.CreatePanel == CreateDescPanel {
				m.BodyInput += "\n"
			}

		case "backspace":
			m.handleBackspace()

		case "tab":
			switch m.CreatePanel {
			case CreateTitlePanel:
				m.CreateMode = modeNormal
				m.CreatePanel = CreateDescPanel
			case CreateDescPanel:
				m.CreateMode = modeNormal
				m.CreatePanel = CreateFieldsPanel
			default:
				m.CreateMode = modeNormal
			}

		case "ctrl+s":
			return m, m.createPR()

		default:
			m.handleCharInput(msg)
		}
	}

	return m, nil
}

func (m *Model) handleBackspace() {
	if m.ViewState == ViewCreatePR {
		switch m.CreatePanel {
		case CreateFieldsPanel:
			if m.CreateField == FieldTarget && len(m.TargetBranch) > 0 {
				m.TargetBranch = m.TargetBranch[:len(m.TargetBranch)-1]
			}
		case CreateTitlePanel:
			if len(m.TitleInput) > 0 {
				m.TitleInput = m.TitleInput[:len(m.TitleInput)-1]
			}
		case CreateDescPanel:
			if len(m.BodyInput) > 0 {
				m.BodyInput = m.BodyInput[:len(m.BodyInput)-1]
			}
		}
	}
	if m.ViewState == ViewMergePR {
		if len(m.CommitMessage) > 0 {
			m.CommitMessage = m.CommitMessage[:len(m.CommitMessage)-1]
		}
	}
}

func (m *Model) handleCharInput(msg tea.KeyMsg) {
	char := msg.String()
	if len(char) != 1 {
		return
	}

	if m.ViewState == ViewCreatePR {
		switch m.CreatePanel {
		case CreateFieldsPanel:
			if m.CreateField == FieldTarget {
				m.TargetBranch += char
			}
		case CreateTitlePanel:
			m.TitleInput += char
		case CreateDescPanel:
			m.BodyInput += char
		}
	}
	if m.ViewState == ViewMergePR {
		m.CommitMessage += char
	}
}

// ============================================
// BRANCH LIST KEYS
// ============================================

func (m *Model) handleBranchListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.ViewState = ViewCreatePR
		return m, nil

	case "j":
		if m.SelectedBranch < len(m.Branches)-1 {
			m.SelectedBranch++
		}

	case "k":
		if m.SelectedBranch > 0 {
			m.SelectedBranch--
		}

	case "enter", "l":
		if m.SelectedBranch < len(m.Branches) {
			m.TargetBranch = m.Branches[m.SelectedBranch].Name
		}
		m.ViewState = ViewCreatePR
		return m, nil
	}

	return m, nil
}

// ============================================
// MERGE PR KEYS
// ============================================

func (m *Model) handleMergePRKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Insert mode for commit message.
	if m.MergeMode == modeInsert {
		switch msg.String() {
		case "esc":
			m.MergeMode = modeNormal
		case "backspace":
			m.handleBackspace()
		case "enter":
			m.CommitMessage += "\n"
		case "ctrl+s":
			return m, m.mergePR()
		default:
			m.handleCharInput(msg)
		}
		return m, nil
	}

	// Normal mode.
	switch msg.String() {
	case "esc":
		m.ViewState = ViewDashboard
		return m, nil

	case "1":
		m.MergePanel = MergePanelStrategy
		m.MergeCursor = 0
	case "2":
		m.MergePanel = MergePanelOptions
		m.OptionsCursor = 0
	case "3":
		m.MergePanel = MergePanelChecklist
	case "4":
		m.MergePanel = MergePanelCommit

	case " ":
		switch m.MergePanel {
		case MergePanelStrategy:
			m.MergeStrategy = MergeStrategy(m.MergeCursor)
		case MergePanelOptions:
			if m.OptionsCursor == 0 {
				m.DeleteRemoteBranch = !m.DeleteRemoteBranch
			} else {
				m.DeleteLocalBranch = !m.DeleteLocalBranch
			}
		case MergePanelCommit:
			m.MergeMode = modeInsert
		}
		return m, nil

	case "j":
		m.navigateMergeDown()
	case "k":
		m.navigateMergeUp()

	case "i", "enter":
		if m.MergePanel == MergePanelCommit {
			m.MergeMode = modeInsert
		}
	}

	return m, nil
}

func (m *Model) navigateMergeDown() {
	switch m.MergePanel {
	case MergePanelStrategy:
		if m.MergeCursor < 2 {
			m.MergeCursor++
		}
	case MergePanelOptions:
		if m.OptionsCursor < 1 {
			m.OptionsCursor++
		}
	}
}

func (m *Model) navigateMergeUp() {
	switch m.MergePanel {
	case MergePanelStrategy:
		if m.MergeCursor > 0 {
			m.MergeCursor--
		}
	case MergePanelOptions:
		if m.OptionsCursor > 0 {
			m.OptionsCursor--
		}
	}
}

// handleMergeTabKeys handles Tab and Shift+Tab in Merge PR view.
func (m *Model) handleMergeTabKeys(key string) (tea.Model, tea.Cmd) {
	if key == "tab" {
		switch m.MergePanel {
		case MergePanelStrategy:
			m.MergePanel = MergePanelOptions
		case MergePanelOptions:
			m.MergePanel = MergePanelChecklist
		case MergePanelChecklist:
			m.MergePanel = MergePanelCommit
		case MergePanelCommit:
			m.MergePanel = MergePanelStrategy
		}
	} else if key == "shift+tab" {
		switch m.MergePanel {
		case MergePanelStrategy:
			m.MergePanel = MergePanelCommit
		case MergePanelOptions:
			m.MergePanel = MergePanelStrategy
		case MergePanelChecklist:
			m.MergePanel = MergePanelOptions
		case MergePanelCommit:
			m.MergePanel = MergePanelChecklist
		}
	}
	return m, nil
}

// ============================================
// ACTIONS
// ============================================

func (m *Model) createPR() tea.Cmd {
	return func() tea.Msg {
		if m.TitleInput == "" {
			return CreatePRResultMsg{Success: false, Error: ErrNoTitle.Error()}
		}

		args := []string{
			"pr", "create",
			"--base", m.TargetBranch,
			"--head", m.SourceBranch,
			"--title", m.TitleInput,
		}

		if m.BodyInput != "" {
			args = append(args, "--body", m.BodyInput)
		}

		cmd := exec.Command("gh", args...)
		out, err := cmd.CombinedOutput()

		if err != nil {
			return CreatePRResultMsg{
				Success: false,
				Error:   fmt.Sprintf("Error: %s", string(out)),
			}
		}

		prNum := 0
		if re := regexp.MustCompile(`pull/(\d+)`); re.Match(out) {
			prNum, _ = strconv.Atoi(re.FindStringSubmatch(string(out))[1])
		}

		return CreatePRResultMsg{
			Success:  true,
			PRNumber: prNum,
		}
	}
}

func (m *Model) mergePR() tea.Cmd {
	return func() tea.Msg {
		prNum := m.Flags.MergePRNum
		if len(m.PRs) > 0 && m.SelectedPR < len(m.PRs) {
			prNum = m.PRs[m.SelectedPR].Number
		}

		if prNum == 0 {
			return MergeResultMsg{Success: false, Error: ErrNoPRSelected.Error()}
		}

		if len(m.PRs) > 0 && m.SelectedPR < len(m.PRs) {
			if m.PRs[m.SelectedPR].Mergeable == "CONFLICTING" {
				return MergeResultMsg{
					Success: false,
					Error:   ErrPRHasConflicts.Error(),
				}
			}
		}

		var strategy string
		switch m.MergeStrategy {
		case StrategyMergeCommit:
			strategy = "--admin"
		case StrategySquash:
			strategy = "--squash"
		case StrategyRebase:
			strategy = "--rebase"
		}

		args := []string{"pr", "merge", strconv.Itoa(prNum), strategy}
		if m.DeleteLocalBranch {
			args = append(args, "-d")
		}

		cmd := exec.Command("gh", args...)
		out, err := cmd.CombinedOutput()

		if err != nil {
			return MergeResultMsg{
				Success: false,
				Error:   fmt.Sprintf("Error: %s", string(out)),
			}
		}

		return MergeResultMsg{
			Success: true,
			Message: fmt.Sprintf("PR #%d merged!", prNum),
		}
	}
}
