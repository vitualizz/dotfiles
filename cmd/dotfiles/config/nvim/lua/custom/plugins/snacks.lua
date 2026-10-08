-- snacks.nvim: a collection of small modules. Used here for the start screen,
-- file explorer, terminal, lazygit, indent guides and notifications.
-- https://github.com/folke/snacks.nvim

vim.pack.add { { src = 'https://github.com/folke/snacks.nvim', version = vim.version.range '2.*' } }

require('snacks').setup {
  dashboard = require 'custom.dashboard',
  explorer = { enabled = true, replace_netrw = true },
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
