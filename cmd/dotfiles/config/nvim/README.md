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
    colorscheme.lua         tokyonight
    editor.lua              guess-indent, todo-comments, mini (statusline, text objects, surround)
    git.lua                 gitsigns and its keymaps
    which-key.lua           keymap hints
    telescope.lua           fuzzy finder and its keymaps
    lsp.lua                 language servers (Mason) and LSP keymaps
    format.lua              conform.nvim
    completion.lua          blink.cmp and LuaSnip
    treesitter.lua          syntax parsers, highlighting, folds
    autopairs.lua           closing brackets and quotes
    snacks.lua              start screen, explorer, terminal, lazygit, indent, notifications
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
| `<leader><leader>` | Open buffers |
| `<C-n>` | File explorer |
| `<A-i>` / `<A-h>` | Floating / bottom terminal (also closes it from inside) |
| `<leader>gg` | Lazygit |
| `grd` `grr` `grn` `gra` `K` | Definition, references, rename, code action, hover (with an LSP) |
| `<leader>f` | Format buffer |

## Customizing

- **Start screen art**: replace the string in `lua/vitualizz/dashboard/header.lua` (up to ~60 columns).
- **Cheatsheet**: edit the groups in `lua/vitualizz/dashboard/cheatsheet.lua`. It is documentation only; when a keymap changes, update it there too.
- **New plugin**: create `lua/vitualizz/plugins/<name>.lua` with `local gh = require('vitualizz.pack').gh`, its `vim.pack.add { gh 'owner/repo' }` and `setup()`, then add `require 'vitualizz.plugins.<name>'` to `init.lua`. It installs on the next start without a prompt.
- **Formatting**: Lua files follow `.stylua.toml` (`stylua .`).
