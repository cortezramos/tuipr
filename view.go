package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	switch m.ViewState {
	case ViewDashboard:
		return m.dashboardView()
	case ViewCreatePR:
		return m.createPRView()
	case ViewMergePR:
		return m.mergePRView()
	case ViewBranchList:
		return m.branchListView()
	default:
		return m.dashboardView()
	}
}

// ============================================
// DASHBOARD
// ============================================

func (m *Model) dashboardView() string {
	var out strings.Builder
	out.WriteString(m.dashboardHeader())
	out.WriteString("\n")

	p1 := m.renderPRsPanel()
	p2 := m.renderDetailsPanel()
	p3 := m.renderConflictsPanel()
	right := lipgloss.JoinVertical(lipgloss.Left, p2, p3)
	main := lipgloss.JoinHorizontal(lipgloss.Top, p1, right)
	out.WriteString(main)
	out.WriteString("\n")
	out.WriteString(m.dashboardFooter())

	if m.Error != "" {
		out.WriteString("\n" + ErrorStyle.Render("Error: "+m.Error))
	}
	if m.StatusMsg != "" {
		out.WriteString("\n" + SuccessStyle.Render("OK: "+m.StatusMsg))
	}
	return out.String()
}

func (m *Model) dashboardHeader() string {
	return fmt.Sprintf(" %s | Branch: %s | %s ",
		TitleStyle.Render("TUIPR"),
		BranchStyle.Render(m.CurrentBranch),
		MutedStyle.Render(fmt.Sprintf("%d PRs", len(m.PRs))))
}

func (m *Model) renderPRsPanel() string {
	var items []string
	items = append(items, panelTitle(1, "PRs", m.ActivePanel == PanelPRs))
	items = append(items, sep())

	if len(m.PRs) == 0 {
		items = append(items, MutedStyle.Render("No hay PRs abiertos"))
		items = append(items, "", KeybindStyle.Render("[c] Crear PR"))
	} else {
		for i, pr := range m.PRs {
			items = append(items, m.renderPRItem(pr, i == m.SelectedPR))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	style := ActivePanelStyle
	if m.ActivePanel != PanelPRs {
		style = InactivePanelStyle
	}
	return style.Width(30).Height(24).Render(content)
}

func (m *Model) renderPRItem(pr PR, sel bool) string {
	icon, iconStyle := getMergeStatusIcon(pr.Mergeable)
	prNum := PRNumberStyle.Render(fmt.Sprintf("#%d", pr.Number))
	pref := "  "
	if sel {
		pref = "-> "
	}
	title := truncate(pr.Title, 22)
	if sel {
		title = SelectedItemStyle.Render(title)
	} else {
		title = NormalItemStyle.Render(title)
	}
	return fmt.Sprintf("%s%s %s %s %s %s",
		pref, iconStyle.Render(icon), prNum, iconStyle.Render("|"), title, getReviewStatus(pr.ReviewDecision))
}

func (m *Model) renderDetailsPanel() string {
	var items []string
	items = append(items, panelTitle(2, "Details", m.ActivePanel == PanelDetails))
	items = append(items, sep())

	if len(m.PRs) > 0 && m.SelectedPR < len(m.PRs) {
		pr := m.PRs[m.SelectedPR]
		items = append(items, TitleTextStyle.Render(pr.Title))
		items = append(items, "", MutedStyle.Render(fmt.Sprintf("Author: %s", pr.Author.Login)))
		items = append(items, fmt.Sprintf("Branch: %s -> %s",
			BranchStyle.Render(pr.HeadRefName), BranchStyle.Render(pr.BaseRefName)))
		items = append(items, "")
		s, st := getMergeStatusText(pr.Mergeable)
		items = append(items, fmt.Sprintf("Merge: %s", st.Render(s)))
		s, st = getReviewStatusText(pr.ReviewDecision)
		items = append(items, fmt.Sprintf("Reviews: %s", st.Render(s)))
	} else {
		items = append(items, MutedStyle.Render("Selecciona un PR"))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	style := ActivePanelStyle
	if m.ActivePanel != PanelDetails {
		style = InactivePanelStyle
	}
	return style.Width(93).Height(12).Render(content)
}

func (m *Model) renderConflictsPanel() string {
	var items []string
	items = append(items, panelTitle(3, "Status", m.ActivePanel == PanelConflicts))
	items = append(items, sep())

	if len(m.PRs) > 0 && m.SelectedPR < len(m.PRs) {
		pr := m.PRs[m.SelectedPR]
		items = append(items, "", TitleTextStyle.Render("CI"))
		items = append(items, GreenStyle.Render("OK Checks Passed"))
		items = append(items, "", TitleTextStyle.Render("Conflicts"))
		if pr.Mergeable == "CONFLICTING" {
			items = append(items, RedStyle.Render("W CONFLICTS"))
			items = append(items, KeybindStyle.Render("[e] nvim"))
		} else {
			items = append(items, GreenStyle.Render("OK None"))
		}
	} else {
		items = append(items, MutedStyle.Render("Selecciona un PR"))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	style := ActivePanelStyle
	if m.ActivePanel != PanelConflicts {
		style = InactivePanelStyle
	}
	return style.Width(93).Height(11).Render(content)
}

func (m *Model) dashboardFooter() string {
	sep := SeparatorStyle.Render(strings.Repeat("-", 125))
	panels := fmt.Sprintf("%s %s %s",
		panelIndicator(1, m.ActivePanel == PanelPRs),
		panelIndicator(2, m.ActivePanel == PanelDetails),
		panelIndicator(3, m.ActivePanel == PanelConflicts))
	keys := "[j/k] Nav  [c] Create  [m] Merge  [q] Quit"
	return fmt.Sprintf("%s\n%s  %s", sep, panels, keys)
}

// ============================================
// CREATE PR
// ============================================

func (m *Model) createPRView() string {
	var out strings.Builder
	out.WriteString(m.createPRHeader())
	out.WriteString("\n")

	p1 := m.renderCreateFieldsPanel()
	p2a := m.renderCreateTitlePanel()
	p2b := m.renderCreateDescPanel()
	p2 := lipgloss.JoinVertical(lipgloss.Left, p2a, p2b)
	main := lipgloss.JoinHorizontal(lipgloss.Top, p1, p2)
	out.WriteString(main)
	out.WriteString("\n" + m.createPRFooter())

	if m.Error != "" {
		out.WriteString("\n" + ErrorStyle.Render("Error: "+m.Error))
	}
	return out.String()
}

func (m *Model) createPRHeader() string {
	return fmt.Sprintf(" %s | Branch: %s ", TitleStyle.Render("TUIPR - Create PR"), BranchStyle.Render(m.SourceBranch))
}

func (m *Model) renderCreateFieldsPanel() string {
	var items []string
	items = append(items, panelTitle(1, "Fields", m.CreatePanel == CreateFieldsPanel))
	items = append(items, sep())
	items = append(items, "")
	items = append(items, m.renderSourceField())
	items = append(items, m.renderTargetField())

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	style := ActivePanelStyle
	if m.CreatePanel != CreateFieldsPanel {
		style = InactivePanelStyle
	}
	return style.Width(28).Height(24).Render(content)
}

func (m *Model) renderSourceField() string {
	pref := "  "
	if m.CreatePanel == CreateFieldsPanel && m.CreateField == FieldSource {
		pref = "-> "
	}
	return fmt.Sprintf("%s %s %s %s", pref, MutedStyle.Render("Source"), MutedStyle.Render("|"), BranchStyle.Render(m.SourceBranch))
}

func (m *Model) renderTargetField() string {
	active := m.CreatePanel == CreateFieldsPanel && m.CreateField == FieldTarget
	pref := "  "
	if active {
		pref = "-> "
	}
	val := NormalItemStyle.Render(m.TargetBranch)
	if active && m.CreateMode == "insert" {
		val = InsertInputStyle.Render(m.TargetBranch) + MutedStyle.Render("_")
	}
	return fmt.Sprintf("%s %s %s %s %s", pref, MutedStyle.Render("Target"), MutedStyle.Render("|"), val, KeybindStyle.Render("[Enter]"))
}

func (m *Model) renderCreateTitlePanel() string {
	var items []string
	items = append(items, panelTitle(2, "Title", m.CreatePanel == CreateTitlePanel))
	items = append(items, sep())

	active := m.CreatePanel == CreateTitlePanel
	if m.TitleInput == "" {
		if active && m.CreateMode == "insert" {
			items = append(items, InsertInputStyle.Render("")+MutedStyle.Render("_"))
		} else {
			items = append(items, MutedStyle.Render("Enter title..."))
		}
	} else {
		if active && m.CreateMode == "insert" {
			items = append(items, InsertInputStyle.Render(m.TitleInput)+MutedStyle.Render("_"))
		} else {
			items = append(items, NormalItemStyle.Render(m.TitleInput))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	style := ActivePanelStyle
	if m.CreatePanel != CreateTitlePanel {
		style = InactivePanelStyle
	}
	return style.Width(95).Height(6).Render(content)
}

func (m *Model) renderCreateDescPanel() string {
	var items []string
	items = append(items, panelTitle(3, "Description", m.CreatePanel == CreateDescPanel))
	items = append(items, sep())

	active := m.CreatePanel == CreateDescPanel
	if m.BodyInput == "" {
		if active && m.CreateMode == "insert" {
			items = append(items, InsertInputStyle.Render("")+MutedStyle.Render("_"))
		} else {
			items = append(items, MutedStyle.Render("Write your PR description..."))
		}
	} else {
		lines := strings.Split(m.BodyInput, "\n")
		if len(lines) > 14 {
			lines = lines[:14]
			lines = append(lines, MutedStyle.Render("..."))
		}
		for i, l := range lines {
			if active && m.CreateMode == "insert" && i == len(lines)-1 {
				items = append(items, InsertInputStyle.Render(l)+MutedStyle.Render("_"))
			} else {
				items = append(items, NormalItemStyle.Render(l))
			}
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	style := ActivePanelStyle
	if m.CreatePanel != CreateDescPanel {
		style = InactivePanelStyle
	}
	return style.Width(95).Height(17).Render(content)
}

func (m *Model) createPRFooter() string {
	sep := SeparatorStyle.Render(strings.Repeat("-", 125))
	mode := NormalModeStyle.Render("NORMAL")
	if m.CreateMode == "insert" {
		mode = InsertModeStyle.Render("INSERT")
	}
	panels := fmt.Sprintf("%s %s %s",
		panelIndicator(1, m.CreatePanel == CreateFieldsPanel),
		panelIndicator(2, m.CreatePanel == CreateTitlePanel),
		panelIndicator(3, m.CreatePanel == CreateDescPanel))
	keys := "[j/k] Nav  [1-3] Panel  [Enter] Select  [i] Edit  [Ctrl+s] Submit"
	if m.CreateMode == "insert" {
		keys = "[Esc] Normal  [Tab] Panel  [Ctrl+s] Submit"
	}
	return fmt.Sprintf("%s\n%s | %s  %s", sep, mode, panels, keys)
}

// ============================================
// MERGE PR - 4 Paneles en 2x2
// ============================================

func (m *Model) mergePRView() string {
	var out strings.Builder
	out.WriteString(m.mergePRHeader())
	out.WriteString("\n")

	// Left column: Merge + Options
	p1 := m.renderMergeStrategyPanel()
	p2 := m.renderMergeOptionsPanel()
	left := lipgloss.JoinVertical(lipgloss.Left, p1, p2)

	// Right column: Checklist + Commit
	p3 := m.renderMergeChecklistPanel()
	p4 := m.renderMergeCommitPanel()
	right := lipgloss.JoinVertical(lipgloss.Left, p3, p4)

	main := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	out.WriteString(main)
	out.WriteString("\n" + m.mergePRFooter())

	if m.Error != "" {
		out.WriteString("\n" + ErrorStyle.Render("Error: "+m.Error))
	}
	return out.String()
}

func (m *Model) mergePRHeader() string {
	info := ""
	if len(m.PRs) > 0 && m.SelectedPR < len(m.PRs) {
		pr := m.PRs[m.SelectedPR]
		info = fmt.Sprintf(" | #%d: %s", pr.Number, pr.Title)
	} else if m.Flags.MergePRNum > 0 {
		info = fmt.Sprintf(" | #%d", m.Flags.MergePRNum)
	}
	return fmt.Sprintf(" %s%s ", TitleStyle.Render("TUIPR - Merge PR"), info)
}

func (m *Model) renderMergeStrategyPanel() string {
	var items []string
	items = append(items, m.mergePanelTitle(1, "Merge"))
	
	// Solo mostrar [x] donde está la selección guardada (MergeStrategy)
	items = append(items, m.mergeCursorItem(0, m.MergeCursor == 0, "Merge Commit", m.MergeStrategy == StrategyMergeCommit))
	items = append(items, m.mergeCursorItem(1, m.MergeCursor == 1, "Squash", m.MergeStrategy == StrategySquash))
	items = append(items, m.mergeCursorItem(2, m.MergeCursor == 2, "Rebase", m.MergeStrategy == StrategyRebase))

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	style := ActivePanelStyle
	if m.MergePanel != MergePanelStrategy {
		style = InactivePanelStyle
	}
	return style.Width(50).Height(10).Render(content)
}

func (m *Model) renderMergeOptionsPanel() string {
	var items []string
	items = append(items, m.mergePanelTitle(2, "Options"))
	
	items = append(items, m.mergeCursorItem(0, m.OptionsCursor == 0, "Delete remote", m.DeleteRemoteBranch))
	items = append(items, m.mergeCursorItem(1, m.OptionsCursor == 1, "Delete local", m.DeleteLocalBranch))

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	style := ActivePanelStyle
	if m.MergePanel != MergePanelOptions {
		style = InactivePanelStyle
	}
	return style.Width(50).Height(10).Render(content)
}

func (m *Model) renderMergeChecklistPanel() string {
	var items []string
	items = append(items, m.mergePanelTitle(3, "Checklist"))
	items = append(items, "")

	// CI (info)
	items = append(items, infoItem("CI Checks", true))

	// Reviews
	approved := m.ReviewsInfo.Approved
	required := m.ReviewsInfo.Required
	items = append(items, infoItem(fmt.Sprintf("Reviews (%d/%d)", approved, required), approved >= required))

	// Conflicts
	hasConflicts := len(m.PRs) > 0 && m.SelectedPR < len(m.PRs) && m.PRs[m.SelectedPR].Mergeable == "CONFLICTING"
	items = append(items, infoItem("Conflicts", !hasConflicts))
	if hasConflicts {
		items = append(items, "  "+RedStyle.Render("[e] nvim"))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	style := ActivePanelStyle
	if m.MergePanel != MergePanelChecklist {
		style = InactivePanelStyle
	}
	return style.Width(50).Height(12).Render(content)
}

func (m *Model) renderMergeCommitPanel() string {
	var items []string
	items = append(items, m.mergePanelTitle(4, "Commit Message"))
	items = append(items, "")
	items = append(items, commitMessageField(m.MergePanel == MergePanelCommit, m.MergeMode == "insert", m.CommitMessage))

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	style := ActivePanelStyle
	if m.MergePanel != MergePanelCommit {
		style = InactivePanelStyle
	}
	return style.Width(50).Height(12).Render(content)
}

func (m *Model) mergePanelTitle(num int, name string) string {
	active := false
	switch num {
	case 1:
		active = (m.MergePanel == MergePanelStrategy)
	case 2:
		active = (m.MergePanel == MergePanelOptions)
	case 3:
		active = (m.MergePanel == MergePanelChecklist)
	case 4:
		active = (m.MergePanel == MergePanelCommit)
	}
	if active {
		return PanelActiveStyle.Render("-> ") + TitleTextStyle.Render(name)
	}
	return PanelInactiveStyle.Render("   ") + TitleTextStyle.Render(name)
}

func (m *Model) mergeCursorItem(cursor int, hasCursor bool, name string, selected bool) string {
	// Muestra -> y [x]/[ ] según estado
	prefix := "  "
	arrow := "  "
	if hasCursor {
		prefix = ""
		arrow = "-> "
	}
	
	// Mostrar [x] o [ ] según si está seleccionado
	selectedStr := "[ ]"
	selectedStyle := MutedStyle
	if selected {
		selectedStr = "[x]"
		selectedStyle = GreenStyle
	}
	
	// Si tiene cursor, el nombre también se resalta
	nameStyle := MutedStyle
	if hasCursor && selected {
		nameStyle = GreenStyle
	}
	
	return fmt.Sprintf("%s%s %s %s", prefix, arrow, selectedStyle.Render(selectedStr), nameStyle.Render(name))
}

func infoItem(name string, ok bool) string {
	icon := "[ ]"
	style := MutedStyle
	if ok {
		icon = "[x]"
		style = GreenStyle
	} else {
		style = RedStyle
	}
	return fmt.Sprintf("  %s %s", MutedStyle.Render(icon), style.Render(name))
}

func commitMessageField(active, editing bool, msg string) string {
	pref := "  "
	arrow := "  "
	if active {
		pref = ""
		arrow = "-> "
	}
	if msg == "" {
		if active && editing {
			return pref + arrow + InsertInputStyle.Render("") + MutedStyle.Render("_")
		}
		return pref + arrow + MutedStyle.Render("(optional)")
	}
	if active && editing {
		return pref + arrow + InsertInputStyle.Render(msg) + MutedStyle.Render("_")
	}
	return pref + arrow + NormalItemStyle.Render(truncate(msg, 35))
}

func (m *Model) mergePRFooter() string {
	sep := SeparatorStyle.Render(strings.Repeat("-", 125))

	p1 := PanelInactiveStyle.Render("[1] Merge")
	p2 := PanelInactiveStyle.Render("[2] Options")
	p3 := PanelInactiveStyle.Render("[3] Checklist")
	p4 := PanelInactiveStyle.Render("[4] Commit")

	// Highlight active panel
	switch m.MergePanel {
	case MergePanelStrategy:
		p1 = PanelActiveStyle.Render("[1] Merge")
	case MergePanelOptions:
		p2 = PanelActiveStyle.Render("[2] Options")
	case MergePanelChecklist:
		p3 = PanelActiveStyle.Render("[3] Checklist")
	case MergePanelCommit:
		p4 = PanelActiveStyle.Render("[4] Commit")
	}

	keys := "[Tab] Next  [Shift+Tab] Prev  [j/k] Nav  [Space] Select  [i] Edit  [Esc] Back"
	return fmt.Sprintf("%s\n%s %s %s %s  %s", sep, p1, p2, p3, p4, keys)
}

// ============================================
// BRANCH LIST
// ============================================

func (m *Model) branchListView() string {
	var out strings.Builder
	out.WriteString(sep() + "\n")
	out.WriteString(centerText("SELECT TARGET BRANCH", 88) + "\n")
	out.WriteString(sep() + "\n\n")

	if len(m.Branches) == 0 {
		out.WriteString("  " + MutedStyle.Render("No branches found") + "\n")
	} else {
		for i, b := range m.Branches {
			sel := i == m.SelectedBranch
			pref := "  "
			if sel {
				pref = "-> "
			}
			name := b.Name
			if b.IsCurrent {
				name += " " + MutedStyle.Render("(current)")
			}
			if sel {
				out.WriteString(fmt.Sprintf("  %s%s %s\n", pref, GreenStyle.Render("*"), SelectedItemStyle.Render(name)))
			} else {
				out.WriteString(fmt.Sprintf("  %s  %s\n", pref, NormalItemStyle.Render(name)))
			}
		}
	}

	out.WriteString("\n" + sep() + "\n")
	out.WriteString(KeybindStyle.Render("[j/k] Navigate  [Enter] Select  [Esc] Cancel"))
	return out.String()
}

// ============================================
// HELPERS
// ============================================

func panelTitle(num int, name string, active bool) string {
	if active {
		return PanelActiveStyle.Render("-> ") + TitleTextStyle.Render(name)
	}
	return PanelInactiveStyle.Render("   ") + TitleTextStyle.Render(name)
}

func panelIndicator(num int, active bool) string {
	if active {
		return PanelActiveStyle.Render(fmt.Sprintf("[%d]", num))
	}
	return PanelInactiveStyle.Render(fmt.Sprintf("[%d]", num))
}

func sep() string {
	return SeparatorStyle.Render("--------------------------------")
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func centerText(s string, w int) string {
	if len(s) >= w {
		return s
	}
	pad := (w - len(s)) / 2
	return strings.Repeat(" ", pad) + s
}

func getMergeStatusIcon(m string) (string, lipgloss.Style) {
	switch m {
	case "MERGEABLE":
		return "*", GreenStyle
	case "CONFLICTING":
		return "W", RedStyle
	default:
		return "o", MutedStyle
	}
}

func getReviewStatus(d string) string {
	switch d {
	case "APPROVED":
		return GreenStyle.Render(" OK")
	case "CHANGES_REQUESTED":
		return RedStyle.Render(" NO")
	case "REVIEW_REQUIRED":
		return YellowStyle.Render(" ...")
	default:
		return MutedStyle.Render(" ---")
	}
}

func getMergeStatusText(m string) (string, lipgloss.Style) {
	switch m {
	case "MERGEABLE":
		return "Ready OK", GreenStyle
	case "CONFLICTING":
		return "Conflicts W", RedStyle
	default:
		return "Unknown", YellowStyle
	}
}

func getReviewStatusText(d string) (string, lipgloss.Style) {
	switch d {
	case "APPROVED":
		return "Approved OK", GreenStyle
	case "CHANGES_REQUESTED":
		return "Changes NO", RedStyle
	case "REVIEW_REQUIRED":
		return "Pending ...", YellowStyle
	default:
		return "No Reviews", MutedStyle
	}
}
