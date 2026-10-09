local gh = require('vitualizz.pack').gh

vim.pack.add { { src = gh 'folke/trouble.nvim', version = vim.version.range '3.*' } }
require('trouble').setup {}

local function toggle(lhs, mode, desc) vim.keymap.set('n', lhs, '<cmd>Trouble ' .. mode .. ' toggle<cr>', { desc = desc }) end
toggle('<leader>xx', 'diagnostics', 'Diagnostics (project)')
toggle('<leader>xX', 'diagnostics filter.buf=0', 'Diagnostics (buffer)')
toggle('<leader>xs', 'symbols', 'Symbols')
toggle('<leader>xl', 'lsp', 'LSP definitions and references')
toggle('<leader>xq', 'qflist', 'Quickfix list')
toggle('<leader>xt', 'todo', 'TODOs')
