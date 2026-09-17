package main

import (
	"testing"
)

func TestFlags_Struct(t *testing.T) {
	flags := Flags{}

	if flags.CreatePR {
		t.Error("CreatePR should be false by default")
	}
	if flags.MergePR {
		t.Error("MergePR should be false by default")
	}
	if flags.MergePRNum != 0 {
		t.Errorf("MergePRNum should be 0, got %d", flags.MergePRNum)
	}
	if flags.Version {
		t.Error("Version should be false by default")
	}
}

func TestFlags_CreatePR(t *testing.T) {
	flags := Flags{CreatePR: true}

	if !flags.CreatePR {
		t.Error("CreatePR should be true")
	}
}

func TestFlags_MergePR(t *testing.T) {
	flags := Flags{MergePR: true}

	if !flags.MergePR {
		t.Error("MergePR should be true")
	}
}

func TestFlags_MergePRWithNumber(t *testing.T) {
	flags := Flags{MergePR: true, MergePRNum: 134}

	if !flags.MergePR {
		t.Error("MergePR should be true")
	}
	if flags.MergePRNum != 134 {
		t.Errorf("MergePRNum should be 134, got %d", flags.MergePRNum)
	}
}

func TestFlags_DirectNumber(t *testing.T) {
	flags := Flags{MergePR: true, MergePRNum: 42}

	if !flags.MergePR {
		t.Error("MergePR should be true with direct number")
	}
	if flags.MergePRNum != 42 {
		t.Errorf("MergePRNum should be 42, got %d", flags.MergePRNum)
	}
}

func TestFlags_Version(t *testing.T) {
	flags := Flags{Version: true}

	if !flags.Version {
		t.Error("Version should be true")
	}
}

func TestFlags_VersionDefault(t *testing.T) {
	flags := Flags{}

	if flags.Version {
		t.Error("Version should be false by default")
	}
	if flags.Help {
		t.Error("Help should be false by default")
	}
}
