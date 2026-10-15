package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	var yayArgs []string

	switch command {
	case "update":
		// yap update -> yay -Syu
		yayArgs = []string{"-Syu"}
		// Append any extra args (like --noconfirm) just in case
		yayArgs = append(yayArgs, args...)

	case "install":
		// yap install pkg -> yay -S pkg
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: 'install' requires at least one package name.")
			os.Exit(1)
		}
		yayArgs = []string{"-S"}
		yayArgs = append(yayArgs, args...)

	case "remove":
		// yap remove pkg -> yay -Rs pkg (keep config)
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: 'remove' requires at least one package name.")
			os.Exit(1)
		}
		yayArgs = []string{"-Rs"}
		yayArgs = append(yayArgs, args...)

	case "purge":
		// yap purge pkg -> yay -Rns pkg (remove config)
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: 'purge' requires at least one package name.")
			os.Exit(1)
		}
		yayArgs = []string{"-Rns"}
		yayArgs = append(yayArgs, args...)

	case "search":
		// yap search query -> yay -Ss query
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: 'search' requires a query string.")
			os.Exit(1)
		}
		yayArgs = []string{"-Ss"}
		yayArgs = append(yayArgs, args...)

	case "info":
		// yap info pkg -> yay -Qi pkg (local) or fallback?
		// Usually users want to see remote info if not installed, but let's map to -Si for consistency with "search" intent or -Qi?
		// Standard managers usually distinguish. Let's map to -Qi for now as that's "query info".
		// Actually, yay -Si is "Sync Info" (repository). yay -Qi is "Query Info" (local).
		// Let's default to -Si (Remote) because usually you want to find out about a package.
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: 'info' requires a package name.")
			os.Exit(1)
		}
		yayArgs = []string{"-Si"}
		yayArgs = append(yayArgs, args...)

	case "clean":
		// yap clean -> yay -Sc (clean cache)
		yayArgs = []string{"-Sc"}
		yayArgs = append(yayArgs, args...)

	case "autoremove":
		// yap autoremove -> yay -Yc (clean dependencies)
		yayArgs = []string{"-Yc"}
		yayArgs = append(yayArgs, args...)

	case "list":
		// yap list -> yay -Q
		yayArgs = []string{"-Q"}
		if len(args) > 0 {
			switch args[0] {
			case "all":
				args = args[1:]
			case "explicit":
				yayArgs = []string{"-Qe"}
				args = args[1:]
			case "native":
				yayArgs = []string{"-Qn"}
				args = args[1:]
			case "aur":
				yayArgs = []string{"-Qm"}
				args = args[1:]
			}
		}
		yayArgs = append(yayArgs, args...)

	case "help", "--help", "-h":
		printUsage()
		return

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}

	runYay(yayArgs)
}

func runYay(args []string) {
	path, err := exec.LookPath("yay")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: yay not found in PATH: %v\n", err)
		os.Exit(1)
	}

	// Status goes to stderr so piped output (e.g. yap list) stays clean.
	fmt.Fprintf(os.Stderr, "-> Executing: yay %s\n", strings.Join(args, " "))

	// Replace this process with yay so signals, exit codes and the terminal
	// behave exactly as if yay had been run directly.
	err = syscall.Exec(path, append([]string{"yay"}, args...), os.Environ())
	fmt.Fprintf(os.Stderr, "Error executing yay: %v\n", err)
	os.Exit(1)
}

func printUsage() {
	fmt.Println("yap - Your Arch Helper Helper")
	fmt.Println("\nUsage:")
	fmt.Println("  yap <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  update              Update system (yay -Syu)")
	fmt.Println("  install <pkg...>    Install packages (yay -S <pkg>)")
	fmt.Println("  remove <pkg...>     Remove packages (yay -Rs <pkg>)")
	fmt.Println("  purge <pkg...>      Remove packages and config (yay -Rns <pkg>)")
	fmt.Println("  search <query>      Search for packages (yay -Ss <query>)")
	fmt.Println("  info <pkg>          Show package information (yay -Si <pkg>)")
	fmt.Println("  clean               Clean package cache (yay -Sc)")
	fmt.Println("  autoremove          Remove unneeded dependencies (yay -Yc)")
	fmt.Println("  list [type]         List packages (all, explicit, native, aur)")
}
