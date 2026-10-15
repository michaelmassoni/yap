package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestMapArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{"update", nil, []string{"-Syu"}},
		{"update", []string{"--devel"}, []string{"-Syu", "--devel"}},
		{"install", []string{"foo", "bar"}, []string{"-S", "foo", "bar"}},
		{"install", []string{"--needed", "foo"}, []string{"-S", "--needed", "foo"}},
		{"remove", []string{"foo"}, []string{"-Rs", "foo"}},
		{"purge", []string{"foo"}, []string{"-Rns", "foo"}},
		{"search", []string{"foo"}, []string{"-Ss", "foo"}},
		{"info", []string{"foo"}, []string{"-Si", "foo"}},
		{"info", []string{"--local", "foo"}, []string{"-Qi", "foo"}},
		{"clean", nil, []string{"-Sc"}},
		{"autoremove", nil, []string{"-Yc"}},
		{"list", nil, []string{"-Q"}},
		{"list", []string{"all"}, []string{"-Q"}},
		{"list", []string{"explicit"}, []string{"-Qe"}},
		{"list", []string{"native"}, []string{"-Qn"}},
		{"list", []string{"aur"}, []string{"-Qm"}},
		{"list", []string{"foo"}, []string{"-Q", "foo"}},
		{"install", []string{"--", "-Weird"}, []string{"-S", "--", "-Weird"}},
	}
	for _, tt := range tests {
		t.Run(tt.name+" "+strings.Join(tt.args, " "), func(t *testing.T) {
			got, err := mapArgs(tt.name, tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMapArgsErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{"install", nil, "requires at least one package name"},
		{"remove", nil, "requires at least one package name"},
		{"purge", nil, "requires at least one package name"},
		{"search", nil, "requires a query string"},
		{"info", nil, "requires a package name"},
		{"info", []string{"--local"}, "requires a package name"},
		{"frobnicate", nil, "unknown command: frobnicate"},
		{"info", []string{"-R", "foo"}, "changes the yay operation"},
		{"search", []string{"-Syu"}, "changes the yay operation"},
		{"install", []string{"--remove", "foo"}, "changes the yay operation"},
		{"list", []string{"--sync=foo"}, "changes the yay operation"},
	}
	for _, tt := range tests {
		t.Run(tt.name+" "+strings.Join(tt.args, " "), func(t *testing.T) {
			_, err := mapArgs(tt.name, tt.args)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("got error %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestExtractDryRun(t *testing.T) {
	tests := []struct {
		in     []string
		want   []string
		dryRun bool
	}{
		{[]string{"install", "foo"}, []string{"install", "foo"}, false},
		{[]string{"install", "foo", "--dry-run"}, []string{"install", "foo"}, true},
		{[]string{"--dry-run", "update"}, []string{"update"}, true},
		{[]string{"install", "--", "--dry-run"}, []string{"install", "--", "--dry-run"}, false},
	}
	for _, tt := range tests {
		got, dryRun := extractDryRun(tt.in)
		if !reflect.DeepEqual(got, tt.want) || dryRun != tt.dryRun {
			t.Errorf("extractDryRun(%q) = %q, %v; want %q, %v", tt.in, got, dryRun, tt.want, tt.dryRun)
		}
	}
}

func TestFormatCommand(t *testing.T) {
	got := formatCommand([]string{"-Ss", "foo bar", "it's"})
	want := `yay -Ss 'foo bar' 'it'\''s'`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

// Every command should appear in every completion script.
func TestCompletionsCoverCommands(t *testing.T) {
	for _, f := range []string{"completions/yap.bash", "completions/_yap", "completions/yap.fish"} {
		data, err := completions.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, c := range commands {
			if !strings.Contains(string(data), c.name) {
				t.Errorf("%s is missing command %q", f, c.name)
			}
		}
		for typ := range listTypes {
			if !strings.Contains(string(data), typ) {
				t.Errorf("%s is missing list type %q", f, typ)
			}
		}
	}
}
