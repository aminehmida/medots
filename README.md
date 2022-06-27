# medots: Cross-platform config files manager
[![Build](https://github.com/aminehmida/medots/actions/workflows/build.yaml/badge.svg)](https://github.com/aminehmida/medots/actions/workflows/build.yaml) [![Go Report Card](https://goreportcard.com/badge/github.com/aminehmida/medots)](https://goreportcard.com/report/github.com/aminehmida/medots)

## Motivation

I use Linux, MacOS and WSL and I switch between them often. I was looking for a tool that can deploy my dot files in different environments quickly and easily. Not only symlinking those configs, I wanted to symlink OS specific files automatically and run commands before and/or after symlinking those file. I couldn't find such tool, so I created *medots*.

## How does it works

medots uses `dots.yaml` file saved in the same git repo as your dot files to understand how to "deploy" them correctly.

You will usually have to run `medots deploy` in the repo where the `dots.yaml` file is saved.

The best way to discover medots features is by building a real world `dots.yaml` config file.

After installing the tool, check the [Features by example](#features-by-example) section to see how to use the `dots.yaml` file.

## Install

### Using Homebrew

```shell
brew cask aminehmida/medots
brew install medots
```

### DEB package

Download the latest `.deb` file for your architecture from the [release page](https://github.com/aminehmida/medots/releases/latest)

### RPM package

Download the latest `.rpm` file for your architecture from the [release page](https://github.com/aminehmida/medots/releases/latest)

### From AUR for Arch Linux

```shell
yay -S medots
```

### APK package for Alpine Linux

Download the latest `.apk` file for your architecture from the [release page](https://github.com/aminehmida/medots/releases/latest).
The package is not signed. You can still verify the package by comparing its sha256 against values in `checksums_sha256.txt` file. After that you can use `--allow-insecure` flag:

```shell
sha256sum --check --ignore-missing checksums_sha256.txt
apk add --allow-insecure ./$apk_filename
```

### Binary archive

Download the latest `.tar.gz` file for your architecture from the [release page](https://github.com/aminehmida/medots/releases/latest) and extract it in your `$PATH`.

### From source

Make sure you have Go installed first.

```shell
go install github.com/aminehmida/medots
```

## Features by example

### Symlinking files

```yaml
nvim:
  - source: ./nvim/init.vim
    destination: ~/.config/nvim/
```

We start by defining the name of the app that we want to manage its config files. This is can be used to target that specific app if you don't want to deploy all your configs.

We can now run `medots deploy` to deploy all apps or `medots deploy nvim` to target nvim only.

The following sequence of actions will happen:

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

Or by adding multiple entries for your app:

```yaml
nvim:
  - source: ./nvim/main_config.vim
    destination: ~/.config/nvim/init.vim
  - source: ./nvim/plugins.vim
    destination: ~/.config/nvim/
```

### Targeting a specific OS

You can target a specific OS by adding `if_os`. The supported values are: `linux`, `darwin`, `windows`. Using medots in Windows is not tested. However, I have tested it in WSL. Feel free to report any issues related to Windows.

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
  # This will symlink the config files then install the plugins
  - source: ./nvim/*
    destination: ~/.config/nvim/
    run_interactive: nvim +PlugInstall
  # This command will be executed last
  - run_interactive: nvim +UpdateRemotePlugins

# You don't need to have a source and destination in your app block.
# This allows running any arbitrary command:
packages:
  - run_interactive: yay -S neovim tmux wget curl ripgrep jq fzf bat exa fd git-delta github-cli sops age awscli kubectl asdf-vm python-poetry
    if_os: linux
  - run_interactive: brew install neovim tmux wget curl ripgrep jq fzf bat exa fd git-delta gh sops age awscli kubectl asdf poetry
    if_os: darwin
```

### Full `dots.yaml` example

```yaml
# Install os specific packages
packages:
  - run_interactive: yay -S neovim tmux wget curl ripgrep jq fzf bat exa fd git-delta github-cli sops age awscli kubectl asdf-vm python-poetry
    if_os: linux
  - run_interactive: brew install neovim tmux wget curl ripgrep jq fzf bat exa fd git-delta gh sops age awscli kubectl asdf poetry
    if_os: darwin

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

I am using `medots` for my personal needs and it's serving me very well. I published this project because I couldn't find any other cross-platform config manager. If you find the it useful, please give it a Star ⭐️ that will make me happy and motivated to work more on it 🤗 also feel free to share it with anyone who could find it useful 😉

If you find a bug or you have a feature request, please open a github issue.
