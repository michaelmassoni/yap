# bash completion for yap
_yap() {
    local cur=${COMP_WORDS[COMP_CWORD]}
    local cmds="update install remove purge search info clean autoremove list completion version help"

    if [[ $COMP_CWORD -eq 1 ]]; then
        COMPREPLY=($(compgen -W "$cmds" -- "$cur"))
        return
    fi

    case ${COMP_WORDS[1]} in
        install)
            COMPREPLY=($(compgen -W "--dry-run $(pacman -Slq 2>/dev/null)" -- "$cur")) ;;
        remove|purge)
            COMPREPLY=($(compgen -W "--dry-run $(pacman -Qq 2>/dev/null)" -- "$cur")) ;;
        info)
            COMPREPLY=($(compgen -W "--local --dry-run $(pacman -Slq 2>/dev/null)" -- "$cur")) ;;
        list)
            [[ $COMP_CWORD -eq 2 ]] && COMPREPLY=($(compgen -W "all explicit native aur" -- "$cur")) ;;
        completion)
            [[ $COMP_CWORD -eq 2 ]] && COMPREPLY=($(compgen -W "bash zsh fish" -- "$cur")) ;;
        *)
            COMPREPLY=($(compgen -W "--dry-run" -- "$cur")) ;;
    esac
}
complete -F _yap yap
