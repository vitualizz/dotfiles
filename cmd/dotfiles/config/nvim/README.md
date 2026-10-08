# Neovim config

Based on [kickstart.nvim](https://github.com/nvim-lua/kickstart.nvim) (upstream `8e35dbe`, MIT, see `LICENSE.md`). Requires Neovim 0.12+: plugins are managed with the built-in `vim.pack`, pinned in `nvim-pack-lock.json`.

## Layout

```
init.lua                    kickstart base, read it top to bottom. Section 10 loads our modules
lua/
  kickstart/plugins/        optional modules shipped by kickstart (autopairs is enabled)
  custom/                   everything specific to this config
    plugins/snacks.lua      snacks.nvim: start screen, explorer, terminal, lazygit, indent, notifications
    dashboard/
      init.lua              start screen layout and quick actions
      header.lua            ASCII art at the top
      cheatsheet.lua        keymaps listed on the start screen
```

Rule of thumb: kickstart's files stay close to upstream so they can be compared; new plugins and settings go under `lua/custom/` and are required from section 10 of `init.lua`, in order.

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

- **Start screen art**: replace the string in `lua/custom/dashboard/header.lua` (up to ~60 columns).
- **Cheatsheet**: edit the groups in `lua/custom/dashboard/cheatsheet.lua`. It is documentation only; when a keymap changes, update it there too.
- **New plugin**: create `lua/custom/plugins/<name>.lua` with its `vim.pack.add { ... }` and `setup()`, then `require 'custom.plugins.<name>'` in section 10 of `init.lua`. New plugins install on the next start without a prompt.
- **Formatting**: Lua files follow `.stylua.toml` (`stylua .`).
