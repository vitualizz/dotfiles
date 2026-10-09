local gh = require('vitualizz.pack').gh

vim.pack.add { { src = gh 'folke/trouble.nvim', version = vim.version.range '3.*' } }
require('trouble').setup {}

local function toggle(lhs, mode, desc) vim.keymap.set('n', lhs, '<cmd>Trouble ' .. mode .. ' toggle<cr>', { desc = desc }) end
toggle('<leader>dd', 'diagnostics', 'Diagnostics (project)')
toggle('<leader>db', 'diagnostics filter.buf=0', 'Diagnostics (buffer)')
toggle('<leader>dS', 'symbols', 'Symbols')
toggle('<leader>dl', 'lsp', 'LSP definitions and references')
toggle('<leader>dq', 'qflist', 'Quickfix list')
toggle('<leader>dt', 'todo', 'TODOs')
