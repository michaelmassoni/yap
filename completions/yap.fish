# fish completion for yap
set -l cmds update install remove purge search info clean autoremove list completion version help

complete -c yap -f
complete -c yap -n "not __fish_seen_subcommand_from $cmds" -a update -d 'Update system (yay -Syu)'
complete -c yap -n "not __fish_seen_subcommand_from $cmds" -a install -d 'Install packages (yay -S)'
complete -c yap -n "not __fish_seen_subcommand_from $cmds" -a remove -d 'Remove packages (yay -Rs)'
complete -c yap -n "not __fish_seen_subcommand_from $cmds" -a purge -d 'Remove packages without .pacsave backups (yay -Rns)'
complete -c yap -n "not __fish_seen_subcommand_from $cmds" -a search -d 'Search for packages (yay -Ss)'
complete -c yap -n "not __fish_seen_subcommand_from $cmds" -a info -d 'Show package information (yay -Si)'
complete -c yap -n "not __fish_seen_subcommand_from $cmds" -a clean -d 'Clean package cache (yay -Sc)'
complete -c yap -n "not __fish_seen_subcommand_from $cmds" -a autoremove -d 'Remove unneeded dependencies (yay -Yc)'
complete -c yap -n "not __fish_seen_subcommand_from $cmds" -a list -d 'List packages (yay -Q)'
complete -c yap -n "not __fish_seen_subcommand_from $cmds" -a completion -d 'Print shell completion script'
complete -c yap -n "not __fish_seen_subcommand_from $cmds" -a version -d 'Show version'
complete -c yap -n "not __fish_seen_subcommand_from $cmds" -a help -d 'Show help'

complete -c yap -n "__fish_seen_subcommand_from install info" -a "(pacman -Slq 2>/dev/null)"
complete -c yap -n "__fish_seen_subcommand_from remove purge" -a "(pacman -Qq 2>/dev/null)"
complete -c yap -n "__fish_seen_subcommand_from info" -l local -d 'Show installed package info (yay -Qi)'
complete -c yap -n "__fish_seen_subcommand_from list" -a "all explicit native aur"
complete -c yap -n "__fish_seen_subcommand_from completion" -a "bash zsh fish"
complete -c yap -l dry-run -d 'Print the yay command without running it'
