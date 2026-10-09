-- snacks.nvim: a collection of small modules. Used here for the start screen,
-- file explorer, terminal, lazygit, indent guides and notifications.
-- https://github.com/folke/snacks.nvim

local gh = require('vitualizz.pack').gh

vim.pack.add { { src = gh 'folke/snacks.nvim', version = vim.version.range '2.*' } }

require('snacks').setup {
  dashboard = require 'vitualizz.dashboard',
  explorer = { enabled = true, replace_netrw = true },
  picker = { enabled = true },
  indent = { enabled = true },
  notifier = { enabled = true },
}

-- Terminals are told apart by `count`, so the floating and the bottom one
-- keep separate sessions.
local function toggle_terminal(count, position)
  return function() Snacks.terminal.toggle(nil, { count = count, win = { position = position } }) end
end

vim.keymap.set('n', '<C-n>', function() Snacks.explorer() end, { desc = 'Toggle file explorer' })
vim.keymap.set({ 'n', 't' }, '<A-i>', toggle_terminal(1, 'float'), { desc = 'Toggle floating terminal' })
vim.keymap.set({ 'n', 't' }, '<A-h>', toggle_terminal(2, 'bottom'), { desc = 'Toggle bottom terminal' })
vim.keymap.set('n', '<leader>gg', function() Snacks.lazygit() end, { desc = '[G]it: lazy[g]it' })
vim.keymap.set('n', '<leader>bd', function() Snacks.bufdelete() end, { desc = '[B]uffer [D]elete (keeps the window layout)' })
vim.keymap.set('n', '<leader>bo', function() Snacks.bufdelete.other() end, { desc = '[B]uffer: delete [O]thers' })
