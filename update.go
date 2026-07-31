package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

// Update maneja los mensajes/eventos
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

// handleKeyMsg
func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	inInsertMode := (m.ViewState == ViewCreatePR && m.CreateMode == "insert") ||
		(m.ViewState == ViewMergePR && m.MergeMode == "insert")

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "q":
		if inInsertMode {
			m.handleCharInput(msg)
			return m, nil
		}
		if m.ViewState == ViewDashboard || m.ViewState == ViewBranchList {
			return m, tea.Quit
		}
		if m.ViewState == ViewCreatePR || m.ViewState == ViewMergePR {
			return m, tea.Quit
		}
	
	case "tab", "shift+tab":
		// Solo manejar Tab en Merge PR
		if m.ViewState == ViewMergePR && m.MergeMode != "insert" {
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
		m.CreateMode = "normal"

	case "m":
		if len(m.PRs) > 0 && m.SelectedPR < len(m.PRs) {
			m.ViewState = ViewMergePR
			m.MergePanel = MergePanelStrategy
			m.MergeMode = "normal"
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
		if m.CreateMode == "insert" {
			m.CreateMode = "normal"
		} else {
			m.ViewState = ViewDashboard
		}
		return m, nil
	}

	if m.CreateMode == "normal" {
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
			if m.CreatePanel == CreateFieldsPanel {
				if m.CreateField == FieldSource {
					m.CreateField = FieldTarget
				}
			}

		case "k":
			if m.CreatePanel == CreateFieldsPanel {
				if m.CreateField == FieldTarget {
					m.CreateField = FieldSource
				}
			}

		case "i":
			m.CreateMode = "insert"

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

	// Modo Insert
	if m.CreateMode == "insert" {
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
				m.CreateMode = "normal"
				m.CreatePanel = CreateDescPanel
			case CreateDescPanel:
				m.CreateMode = "normal"
				m.CreatePanel = CreateFieldsPanel
			default:
				m.CreateMode = "normal"
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
	// Modo Insert para commit message
	if m.MergeMode == "insert" {
		switch msg.String() {
		case "esc":
			m.MergeMode = "normal"
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

	// Modo Normal
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
		// Barra espaciadora - seleccionar/opción
		switch m.MergePanel {
		case MergePanelStrategy:
			// Rotar entre estrategias según cursor
			m.MergeStrategy = MergeStrategy(m.MergeCursor)
		case MergePanelOptions:
			// Toggle según cursor (0=remote, 1=local)
			if m.OptionsCursor == 0 {
				m.DeleteRemoteBranch = !m.DeleteRemoteBranch
			} else {
				m.DeleteLocalBranch = !m.DeleteLocalBranch
			}
		case MergePanelChecklist:
			// Solo informativo, no hace nada
		case MergePanelCommit:
			m.MergeMode = "insert"
		}
		return m, nil

	case "j":
		// Navegar hacia abajo dentro del panel
		m.navigateMergeDown()
	case "k":
		// Navegar hacia arriba dentro del panel
		m.navigateMergeUp()

	case "i":
		if m.MergePanel == MergePanelCommit {
			m.MergeMode = "insert"
		}

	case "enter":
		if m.MergePanel == MergePanelCommit {
			m.MergeMode = "insert"
		}
	}

	return m, nil
}

func (m *Model) navigateMergeDown() {
	switch m.MergePanel {
	case MergePanelStrategy:
		// Solo mover la flecha, NO selecciona automáticamente
		if m.MergeCursor < 2 {
			m.MergeCursor++
		}
	case MergePanelOptions:
		if m.OptionsCursor < 1 {
			m.OptionsCursor++
		}
	case MergePanelChecklist:
		// Solo informativo
	case MergePanelCommit:
		// Solo informativo
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
	case MergePanelChecklist:
		// Solo informativo
	case MergePanelCommit:
		// Solo informativo
	}
}

// handleMergeTabKeys maneja Tab y Shift+Tab en Merge PR
func (m *Model) handleMergeTabKeys(key string) (tea.Model, tea.Cmd) {
	if key == "tab" {
		// Avanzar panel
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
		// Retroceder panel
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
// ACCIONES
// ============================================

func (m *Model) createPR() tea.Cmd {
	return func() tea.Msg {
		if m.TitleInput == "" {
			return CreatePRResultMsg{Success: false, Error: "Title is required"}
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
			strconv.Atoi(re.FindStringSubmatch(string(out))[1])
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
			return MergeResultMsg{Success: false, Error: "No PR selected"}
		}

		if len(m.PRs) > 0 && m.SelectedPR < len(m.PRs) {
			if m.PRs[m.SelectedPR].Mergeable == "CONFLICTING" {
				return MergeResultMsg{
					Success: false,
					Error:   "Cannot merge: has conflicts. Press [e] to resolve.",
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
