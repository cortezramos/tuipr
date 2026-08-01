package main

import (
	"testing"
)

func TestNewModel_Default(t *testing.T) {
	flags := Flags{}
	model := NewModel(flags)

	if model.ViewState != ViewDashboard {
		t.Errorf("ViewState should be ViewDashboard, got %v", model.ViewState)
	}
	if model.ActivePanel != PanelPRs {
		t.Errorf("ActivePanel should be PanelPRs, got %v", model.ActivePanel)
	}
	if model.SelectedPR != 0 {
		t.Errorf("SelectedPR should be 0, got %d", model.SelectedPR)
	}
	if model.MergeStrategy != StrategyMergeCommit {
		t.Errorf("MergeStrategy should be StrategyMergeCommit, got %v", model.MergeStrategy)
	}
	if !model.DeleteLocalBranch {
		t.Error("DeleteLocalBranch should be true by default")
	}
}

func TestNewModel_CreatePR(t *testing.T) {
	flags := Flags{CreatePR: true}
	model := NewModel(flags)

	if model.ViewState != ViewCreatePR {
		t.Errorf("ViewState should be ViewCreatePR, got %v", model.ViewState)
	}
	if model.CreatePanel != CreateFieldsPanel {
		t.Errorf("CreatePanel should be CreateFieldsPanel, got %v", model.CreatePanel)
	}
	if model.CreateField != FieldTarget {
		t.Errorf("CreateField should be FieldTarget, got %v", model.CreateField)
	}
	if model.CreateMode != modeNormal {
		t.Errorf("CreateMode should be modeNormal, got %s", model.CreateMode)
	}
}

func TestNewModel_MergePR(t *testing.T) {
	flags := Flags{MergePR: true}
	model := NewModel(flags)

	if model.ViewState != ViewMergePR {
		t.Errorf("ViewState should be ViewMergePR, got %v", model.ViewState)
	}
	if model.MergePanel != MergePanelStrategy {
		t.Errorf("MergePanel should be MergePanelStrategy, got %v", model.MergePanel)
	}
	if model.MergeMode != modeNormal {
		t.Errorf("MergeMode should be modeNormal, got %s", model.MergeMode)
	}
}

func TestNewModel_MergePRWithNumber(t *testing.T) {
	flags := Flags{MergePR: true, MergePRNum: 142}
	model := NewModel(flags)

	if model.ViewState != ViewMergePR {
		t.Errorf("ViewState should be ViewMergePR, got %v", model.ViewState)
	}
	if model.SelectedPR != 142 {
		t.Errorf("SelectedPR should be 142, got %d", model.SelectedPR)
	}
}

func TestViewStateConstants(t *testing.T) {
	if ViewDashboard != 0 {
		t.Errorf("ViewDashboard should be 0, got %d", ViewDashboard)
	}
	if ViewCreatePR != 1 {
		t.Errorf("ViewCreatePR should be 1, got %d", ViewCreatePR)
	}
	if ViewMergePR != 2 {
		t.Errorf("ViewMergePR should be 2, got %d", ViewMergePR)
	}
	if ViewBranchList != 3 {
		t.Errorf("ViewBranchList should be 3, got %d", ViewBranchList)
	}
}

func TestPanelIndexConstants(t *testing.T) {
	if PanelPRs != 0 {
		t.Errorf("PanelPRs should be 0, got %d", PanelPRs)
	}
	if PanelDetails != 1 {
		t.Errorf("PanelDetails should be 1, got %d", PanelDetails)
	}
	if PanelConflicts != 2 {
		t.Errorf("PanelConflicts should be 2, got %d", PanelConflicts)
	}
}

func TestMergeStrategyConstants(t *testing.T) {
	if StrategyMergeCommit != 0 {
		t.Errorf("StrategyMergeCommit should be 0, got %d", StrategyMergeCommit)
	}
	if StrategySquash != 1 {
		t.Errorf("StrategySquash should be 1, got %d", StrategySquash)
	}
	if StrategyRebase != 2 {
		t.Errorf("StrategyRebase should be 2, got %d", StrategyRebase)
	}
}

func TestModeConstants(t *testing.T) {
	if modeNormal != "normal" {
		t.Errorf("modeNormal should be 'normal', got '%s'", modeNormal)
	}
	if modeInsert != "insert" {
		t.Errorf("modeInsert should be 'insert', got '%s'", modeInsert)
	}
}

func TestErrorDefinitions(t *testing.T) {
	if ErrNoTitle == nil {
		t.Error("ErrNoTitle should not be nil")
	}
	if ErrNoPRSelected == nil {
		t.Error("ErrNoPRSelected should not be nil")
	}
	if ErrPRHasConflicts == nil {
		t.Error("ErrPRHasConflicts should not be nil")
	}
}
