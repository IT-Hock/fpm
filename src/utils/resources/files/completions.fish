set -l commands install search remove list github

complete -c fpm -f

complete -c fpm -l yes -o y -d 'Answer all prompts with yes'

complete -c fpm -l no -o n -d 'Answer all prompts with no'

complete -c fpm -s h -l help -d 'Print a short help text and exit'

complete -c fpm -n "not __fish_seen_subcommand_from $commands" \
     -a 'install' -d 'Install a package'

#complete -c fpm -n "__fish_seen_subcommand_from install" \
     #-a '(__fish_complete_packages)' -d 'Package to install'

complete -c fpm -n "not __fish_seen_subcommand_from $commands" \
     -a 'search' -d 'Search for a package'

complete -c fpm -n "not __fish_seen_subcommand_from $commands" \
     -a 'remove' -d 'Remove a package'

complete -c fpm -n "not __fish_seen_subcommand_from $commands" \
     -a 'list' -d 'List installed packages'

complete -c fpm -n "not __fish_seen_subcommand_from $commands" \
     -a 'github' -d 'Manage GitHub token'

complete -c fpm -n "__fish_seen_subcommand_from github" \
     -a 'login' -d 'Login to GitHub'

complete -c fpm -n "__fish_seen_subcommand_from github" \
     -a 'logout' -d 'Logout from GitHub'

complete -c fpm -n "__fish_seen_subcommand_from github" \
     -a 'token' -d 'Get GitHub token'
