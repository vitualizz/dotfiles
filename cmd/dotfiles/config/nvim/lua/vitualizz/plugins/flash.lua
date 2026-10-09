local gh = require('vitualizz.pack').gh

vim.pack.add { { src = gh 'folke/flash.nvim', version = vim.version.range '2.*' } }

-- f/t/F/T keep Vim's behavior; flash only adds its own keys. Surround uses
-- the `gs` prefix (editor.lua) so it doesn't wait on `s`.
require('flash').setup { modes = { char = { enabled = false } } }

-- Flash's labels link to Substitute, which Everblush colors like Search, so
-- labels and matches looked the same. Labels get their own color.
local c = require('vitualizz.plugins.colorscheme').palette
vim.api.nvim_set_hl(0, 'FlashLabel', { fg = c.base00, bg = c.base0E, bold = true })

local flash = require 'flash'
vim.keymap.set({ 'n', 'x', 'o' }, 's', flash.jump, { desc = 'Flash jump' })
vim.keymap.set({ 'n', 'x', 'o' }, 'S', flash.treesitter, { desc = 'Flash treesitter selection' })
vim.keymap.set('o', 'r', flash.remote, { desc = 'Flash remote' })
vim.keymap.set({ 'o', 'x' }, 'R', flash.treesitter_search, { desc = 'Flash treesitter search' })
vim.keymap.set('c', '<C-s>', flash.toggle, { desc = 'Toggle flash in search' })
