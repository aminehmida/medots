# medots - multiplatform config files manager

## Motivation

I want to be able to manage my config files (dot files) stored in a git repo across multiple OSs. The actual config files should be created as symlinks to the original files in the git repository.

## How does it works

me dots uses `dots.yaml` file saved in the same git repo as your dot files to understand how to "deploy" them correctly.

You will usually have to run `medots deploy` in the repo where the `dots.yaml` file is saved.

The best way to discover medots features is by building a real world `dots.yaml` config file.

## Features by example

### Symlinking files

```yaml
nvim:
  - source: ./nvim/init.vim
  - destination: ~/.config/nvim/
```

We start by defining the name of the app that we want to manage its config files. This is can be used to target that specific app if you don't want to deploy all your configs.

We can now run `medots deploy nvim`. The following sequence of actions will happen:

1. `~/.config/nvim/` directory will be recursively created if it doesn't exist
2. `~/.config/nvim/init.vim` will be renamed to `~/.config/nvim/init.vim.bak` if it already exists
3. `~/.config/nvim/init.vim` will be created as a symlink to `./nvim/init.vim`. The current directory dot `.` in the `source` and the tild `~` in the `destination` will be expanded when creating the symlink.

You can also set a different name for the destination file as bellow:

```yaml
nvim:
  - source: ./nvim/main_config.vim
    destination: ~/.config/nvim/init.vim
```

Because of this behavior make sure to always add a trailing `/` when the destination is a folder.

It's also possible to symlink multiple files using a glob expression:

```yaml
nvim:
  - source: ./nvim/*.vim
    destination: ~/.config/nvim/
```

Or by add multiple entries for you app:

```yaml
nvim:
  - source: ./nvim/main_config.vim
    destination: ~/.config/nvim/init.vim
  - source: ./nvim/plugins.vim
    destination: ~/.config/nvim/
```

### Targeting a specific OS

You can target a specific OS by adding `if_os`. The supported values are: `linux`, `darwin`, `windows`. Using me dots in Windows is not tested. Feel free  to report any issues.

```yaml
nvim:
  - source: ./nvim/init_darwin.vim
    destination: ~/.config/nvim/
    if_os: darwin
  - source: ./nvim/init_linux.vim
    destination: ~/.config/nvim/
    if_os: linux
```

### Running a command

Sometimes, a command need to be executed before or after the config file is symlinked. This can be done using `run` or `run_interactive`. `run` will capture the stdout/stderr of the command and show them at the end. `run_interactive` is useful if user input is required. Those commands will be executed by order of appearance.

```yaml
nvim:
  # Install vim-plug
  - run: curl -fLo "${XDG_DATA_HOME:-$HOME/.local/share}"/nvim/site/autoload/plug.vim --create-dirs https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
  # This will copy the config files then install the plugins
  - source: ./nvim/*
    destination: ~/.config/nvim/
    run_interactive: nvim +PlugInstall
  # This command will be executed last
  - run_interactive: nvim +UpdateRemotePlugins
```

### Full `dots.yaml` example

```yaml
nvim:
  # Install vim-plug
  - run: curl -fLo "${XDG_DATA_HOME:-$HOME/.local/share}"/nvim/site/autoload/plug.vim --create-dirs https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
  # This will copy the config files then install the plugins
  - source: ./nvim/*
    destination: ~/.config/nvim/
    run_interactive: nvim +PlugInstall
  # This command will be executed last
  - run_interactive: nvim +UpdateRemotePlugins

tmux:
  - source: tmux/.tmux.conf
    destination: ~/.tmux.conf
  - source: tmux/tmux.darwin.conf
    destination: ~/.config/tmux/tmux.custom.conf
    if_os: darwin
  - source: tmux/tmux.linux.conf
    destination: ~/.config/tmux.custom.conf
    if_os: linux

zsh:
  # Install ohmyzsh
  - run_interactive: 'sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"'
    if_os: darwin
  - run_nteractive: yay -S oh-my-zsh-git
    if_os: linux
  # Install zsh-histdb
  - run: mkdir -p $HOME/.oh-my-zsh/custom/plugins/; git clone https://github.com/larkery/zsh-histdb $HOME/.oh-my-zsh/custom/plugins/zsh-histdb
  - source: zsh/.zshrc
    destination: ~/.zshrc
  - source: zsh/*.zsh
    destination: ~/.config/zsh/
```

## Contact, bugs and feature requests

I am using `medots` for my personal needs and it's serving me well because I couldn't find any cross platform config file manager. If you find the project useful, please give it a Star ⭐️ that will make me happy and motivated to work more on it 🤗 also feel free to share it with anyone who could find it useful 😉

If you find a bug or you have a feature request, please open a github issue