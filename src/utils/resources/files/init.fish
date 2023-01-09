if not set -q FPM_CONFIG
    set -q XDG_CONFIG_HOME; or set -l XDG_CONFIG_HOME "$HOME/.config"
    set -q XDG_DATA_HOME; or set -l XDG_DATA_HOME "$HOME/.local/share"

    set -gx FPM_CONFIG "$XDG_CONFIG_HOME/fpm"
    set -gx FPM_PATH "$XDG_DATA_HOME/fpm"
end
test -f $FPM_CONFIG/before.init.fish and source $FPM_CONFIG/before.init.fish 2>/dev/null
test -f $FPM_CONFIG/theme
and read -l theme <$FPM_CONFIG/theme
or set -l theme default

# Require all packages
#emit perf:timer:start "Fish Package Manager init installed packages"
#require --path {$FPM_PATH,$FPM_CONFIG}/pkg/*
#emit perf:timer:finish "Fish Package Manager init installed packages"
#emit perf:timer:start "Oh My Fish init user config path"
#require --no-bundle --path $FPM_CONFIG
#emit perf:timer:finish "Oh My Fish init user config path"
# Load conf.d for current theme if exists
set -l theme_conf_path {$FPM_CONFIG,$FPM_PATH}/themes*/$theme/conf.d
for conf in $theme_conf_path/*.fish
    source $conf
end

function include_package
    set -l pkg $argv[1]
    set -l pkgVersion $argv[2]

    set -l package_path {$FPM_CONFIG,$FPM_PATH}/packages/$pkg/$pkgVersion

    set function_path $package_path/functions*
    set complete_path $package_path/completions*
    set init_path $package_path/init.fish*
    set conf_path $package_path/conf.d/*.fish

    # Autoload functions
    test -n "$function_path"
    and set fish_function_path $fish_function_path[1] \
        $function_path \
        $fish_function_path[2..-1]

    # Autoload completions
    test -n "$complete_path"
    and set fish_complete_path $fish_complete_path[1] \
        $complete_path \
        $fish_complete_path[2..-1]

    for init in $init_path
        source $init
    end

    for conf in $conf_path
        source $conf
    end

    return 0
end

# Load all packages
for pkg in $FPM_CONFIG/packages/* $FPM_PATH/packages/*
    test -d $pkg
    and include_package (basename $pkg) (basename $pkg)

    echo $pkg
end

source "$FPM_PATH/completions.fish"

# Add FPM to path
fish_add_path -g $FPM_PATH
