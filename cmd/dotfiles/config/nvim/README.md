# Neovim config

Started from [kickstart.nvim](https://github.com/nvim-lua/kickstart.nvim) (upstream `8e35dbe`, MIT, see `LICENSE.md`) and split into modules. Requires Neovim 0.12+: plugins are managed with the built-in `vim.pack`, pinned in `nvim-pack-lock.json`.

## Layout

```
init.lua                    loads every module below, in order
lua/vitualizz/
  options.lua               editor options and leader key
  pack.lua                  vim.pack setup: no install prompt, build hooks, `gh()` helper
  keymaps.lua               general keymaps (windows, terminal, search highlight)
  diagnostics.lua           how LSP diagnostics are displayed
  autocmds.lua              autocommands (highlight on yank)
  health.lua                `:checkhealth vitualizz`
  plugins/
    colorscheme.lua         Everblush (NvChad base16 palette via mini.base16)
    editor.lua              guess-indent, todo-comments, mini (text objects, surround)
    statusline.lua          powerline-style statusline colored by mode
    git.lua                 gitsigns and its keymaps
    which-key.lua           keymap hints
    picker.lua              fuzzy finder (snacks.picker) and LSP pickers
    lsp.lua                 language servers (Mason) and LSP keymaps
    format.lua              conform.nvim
    completion.lua          blink.cmp and LuaSnip
    treesitter.lua          syntax parsers, highlighting, folds
    autopairs.lua           closing brackets and quotes
    flash.lua               jump anywhere with `s`
    trouble.lua             diagnostics, symbols and TODO lists
    grug-far.lua            project-wide search and replace
    snacks.lua              start screen, picker, explorer, terminal, lazygit, bufdelete, indent, notifications
  dashboard/
    init.lua                start screen layout and quick actions
    header.lua              ASCII art at the top
    cheatsheet.lua          keymaps listed on the start screen
```

One file per concern: options and keymaps that belong to a plugin live in that plugin's file. Order matters only in `init.lua` (options and `pack` first, then plugins).

## Keymaps

`<leader>` is `Space`. Press it and wait to see every keymap (which-key).

| Key | Action |
|-----|--------|
| `<leader>sf` / `<leader>sg` | Search files / text |
| `<leader>sR` | Search and replace in the project (grug-far) |
| `<leader><leader>` / `<leader>bd` | Open buffers / close buffer without closing its window |
| `s` / `S` | Jump anywhere / select a syntax node (flash) |
| `gsa` `gsd` `gsr` | Surround: add / delete / replace (`gs` prefix keeps `s` for flash) |
| `<leader>xx` / `<leader>xt` | Diagnostics / TODOs list (trouble) |
| `<C-n>` | File explorer |
| `<A-i>` / `<A-h>` | Floating / bottom terminal (also closes it from inside) |
| `<leader>gg` | Lazygit |
| `grd` `grr` `grn` `gra` `K` | Definition, references, rename, code action, hover (with an LSP) |
| `<leader>f` | Format buffer |

## Languages

| Language | LSP | Formatter | Requires |
|----------|-----|-----------|----------|
| Lua | `lua_ls` | `stylua` | nothing (Mason) |
| TypeScript / JavaScript | `ts_ls`, `eslint` | `prettier` (the project's own if installed) | Node on PATH, e.g. `mise use -g node@lts` (Mason installs the servers with npm) |
| Ruby | `ruby_lsp` | RuboCop, through ruby-lsp, when the project uses it | `gem install ruby-lsp` for each Ruby version (e.g. from mise). Not installed by Mason; enabled only if `ruby-lsp` is on PATH |

Servers and formatters are listed in `plugins/lsp.lua` and `plugins/format.lua`; syntax parsers in `plugins/treesitter.lua`. Formatting runs with `<leader>f` (not on save).

## Customizing

- **Start screen art**: replace the string in `lua/vitualizz/dashboard/header.lua` (up to ~60 columns).
- **Cheatsheet**: edit the groups in `lua/vitualizz/dashboard/cheatsheet.lua`. It is documentation only; when a keymap changes, update it there too.
- **New plugin**: create `lua/vitualizz/plugins/<name>.lua` with `local gh = require('vitualizz.pack').gh`, its `vim.pack.add { gh 'owner/repo' }` and `setup()`, then add `require 'vitualizz.plugins.<name>'` to `init.lua`. It installs on the next start without a prompt.
- **Formatting**: Lua files follow `.stylua.toml` (`stylua .`).
