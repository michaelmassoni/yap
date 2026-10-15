package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"syscall"
)

// version is set at build time with -ldflags "-X main.version=...".
var version = "dev"

//go:embed completions/yap.bash completions/_yap completions/yap.fish
var completions embed.FS

// command maps a yap verb to a yay operation.
type command struct {
	name    string
	flag    string
	usage   string
	desc    string
	needArg string // if set, the command requires arguments; used in the error message
}

var commands = []command{
	{"update", "-Syu", "update", "Update system", ""},
	{"install", "-S", "install <pkg...>", "Install packages", "at least one package name"},
	{"remove", "-Rs", "remove <pkg...>", "Remove packages and unneeded deps", "at least one package name"},
	{"purge", "-Rns", "purge <pkg...>", "Like remove, but skip .pacsave config backups", "at least one package name"},
	{"search", "-Ss", "search <query>", "Search for packages", "a query string"},
	{"info", "-Si", "info [--local] <pkg>", "Show package info; --local uses -Qi", "a package name"},
	{"clean", "-Sc", "clean", "Clean package cache", ""},
	{"autoremove", "-Yc", "autoremove", "Remove unneeded dependencies", ""},
	{"list", "-Q", "list [type]", "List packages; type: all, explicit, native, aur", ""},
}

var listTypes = map[string]string{
	"all":      "-Q",
	"explicit": "-Qe",
	"native":   "-Qn",
	"aur":      "-Qm",
}

// Long options that select a different pacman/yay operation.
var operationLongOpts = map[string]bool{
	"--sync": true, "--remove": true, "--query": true, "--upgrade": true,
	"--database": true, "--files": true, "--deptest": true, "--version": true,
	"--yay": true, "--show": true, "--getpkgbuild": true, "--web": true, "--build": true,
}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(argv []string) int {
	argv, dryRun := extractDryRun(argv)

	if len(argv) == 0 {
		printUsage(os.Stderr)
		return 1
	}

	name, args := argv[0], argv[1:]

	switch name {
	case "help", "--help", "-h":
		printUsage(os.Stdout)
		return 0
	case "version", "--version", "-v":
		fmt.Printf("yap %s\n", buildVersion())
		return 0
	case "completion":
		return printCompletion(args)
	}

	yayArgs, err := mapArgs(name, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		var unknown unknownCommandError
		if errors.As(err, &unknown) {
			printUsage(os.Stderr)
		}
		return 1
	}

	if dryRun {
		fmt.Println(formatCommand(yayArgs))
		return 0
	}
	return execYay(yayArgs)
}

type unknownCommandError string

func (e unknownCommandError) Error() string { return fmt.Sprintf("unknown command: %s", string(e)) }

// mapArgs translates a yap command and its arguments into yay arguments.
func mapArgs(name string, args []string) ([]string, error) {
	var cmd *command
	for i := range commands {
		if commands[i].name == name {
			cmd = &commands[i]
			break
		}
	}
	if cmd == nil {
		return nil, unknownCommandError(name)
	}

	flag := cmd.flag
	switch name {
	case "list":
		if len(args) > 0 {
			if f, ok := listTypes[args[0]]; ok {
				flag = f
				args = args[1:]
			}
		}
	case "info":
		if len(args) > 0 && args[0] == "--local" {
			flag = "-Qi"
			args = args[1:]
		}
	}

	if cmd.needArg != "" && len(args) == 0 {
		return nil, fmt.Errorf("'%s' requires %s", name, cmd.needArg)
	}
	if err := checkOperationFlags(name, args); err != nil {
		return nil, err
	}

	return append([]string{flag}, args...), nil
}

// checkOperationFlags rejects arguments that would switch yay to a different
// operation (e.g. 'yap info -R pkg'), so each verb only does what it says.
// pacman/yay operations are the uppercase short flags; options are lowercase.
// Anything after "--" is treated as a target and passed through untouched.
func checkOperationFlags(name string, args []string) error {
	for _, a := range args {
		if a == "--" {
			return nil
		}
		isOp := false
		if strings.HasPrefix(a, "--") {
			opt, _, _ := strings.Cut(a, "=")
			isOp = operationLongOpts[opt]
		} else if strings.HasPrefix(a, "-") {
			isOp = strings.ContainsAny(a[1:], "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		}
		if isOp {
			return fmt.Errorf("'%s' changes the yay operation and can't be used with 'yap %s'; run yay directly instead", a, name)
		}
	}
	return nil
}

// extractDryRun removes --dry-run from the arguments (before any "--").
func extractDryRun(argv []string) ([]string, bool) {
	out := make([]string, 0, len(argv))
	dryRun := false
	for i, a := range argv {
		if a == "--" {
			out = append(out, argv[i:]...)
			break
		}
		if a == "--dry-run" {
			dryRun = true
			continue
		}
		out = append(out, a)
	}
	return out, dryRun
}

// formatCommand renders the yay command line, quoting arguments where needed.
func formatCommand(args []string) string {
	parts := []string{"yay"}
	for _, a := range args {
		if a == "" || strings.ContainsAny(a, " \t\n'\"\\$`*?[]{}()<>|&;#~") {
			a = "'" + strings.ReplaceAll(a, "'", `'\''`) + "'"
		}
		parts = append(parts, a)
	}
	return strings.Join(parts, " ")
}

func execYay(args []string) int {
	path, err := exec.LookPath("yay")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: yay not found in PATH: %v\n", err)
		return 1
	}

	// Status goes to stderr so piped output (e.g. yap list) stays clean.
	fmt.Fprintf(os.Stderr, "-> Executing: %s\n", formatCommand(args))

	// Replace this process with yay so signals, exit codes and the terminal
	// behave exactly as if yay had been run directly.
	err = syscall.Exec(path, append([]string{"yay"}, args...), os.Environ())
	fmt.Fprintf(os.Stderr, "Error executing yay: %v\n", err)
	return 1
}

func buildVersion() string {
	if version != "dev" {
		return version
	}
	// Fall back to the module version for 'go install' builds.
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}

func printCompletion(args []string) int {
	files := map[string]string{
		"bash": "completions/yap.bash",
		"zsh":  "completions/_yap",
		"fish": "completions/yap.fish",
	}
	if len(args) != 1 || files[args[0]] == "" {
		fmt.Fprintln(os.Stderr, "Usage: yap completion <bash|zsh|fish>")
		return 1
	}
	data, err := completions.ReadFile(files[args[0]])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	os.Stdout.Write(data)
	return 0
}

func printUsage(w *os.File) {
	fmt.Fprintln(w, "yap - Your Arch Helper Helper")
	fmt.Fprintln(w, "\nUsage:")
	fmt.Fprintln(w, "  yap <command> [arguments] [--dry-run]")
	fmt.Fprintln(w, "\nCommands:")
	for _, c := range commands {
		fmt.Fprintf(w, "  %-22s %s (yay %s)\n", c.usage, c.desc, c.flag)
	}
	fmt.Fprintf(w, "  %-22s %s\n", "completion <shell>", "Print shell completion script (bash, zsh, fish)")
	fmt.Fprintf(w, "  %-22s %s\n", "version", "Show version")
	fmt.Fprintln(w, "\nOptions:")
	fmt.Fprintf(w, "  %-22s %s\n", "--dry-run", "Print the yay command without running it")
}
