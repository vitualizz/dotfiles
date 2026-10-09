vim.keymap.set('n', '<Esc>', '<cmd>nohlsearch<CR>')
vim.keymap.set('n', '<leader>q', vim.diagnostic.setloclist, { desc = 'Open diagnostic [Q]uickfix list' })

vim.keymap.set('t', '<Esc><Esc>', '<C-\\><C-n>', { desc = 'Exit terminal mode' })

vim.keymap.set('n', '<C-h>', '<C-w><C-h>', { desc = 'Move focus to the left window' })
vim.keymap.set('n', '<C-l>', '<C-w><C-l>', { desc = 'Move focus to the right window' })
vim.keymap.set('n', '<C-j>', '<C-w><C-j>', { desc = 'Move focus to the lower window' })
vim.keymap.set('n', '<C-k>', '<C-w><C-k>', { desc = 'Move focus to the upper window' })

-- NvChad keymaps (v2.5), kept for muscle memory.
local map = vim.keymap.set
map('n', '<C-s>', '<cmd>w<CR>', { desc = 'Save file' })
map('n', '<C-c>', '<cmd>%y+<CR>', { desc = 'Copy whole file' })
map('n', '<leader>n', '<cmd>set nu!<CR>', { desc = 'Toggle line [N]umbers' })
map('n', '<leader>rn', '<cmd>set rnu!<CR>', { desc = 'Toggle [R]elative [N]umbers' })
map('n', '<leader>ds', vim.diagnostic.setloclist, { desc = '[D]iagnostics loclist' })
map('n', '<leader>b', '<cmd>enew<CR>', { desc = 'New [B]uffer' })
map('n', '<Tab>', '<cmd>bnext<CR>', { desc = 'Next buffer' })
map('n', '<S-Tab>', '<cmd>bprevious<CR>', { desc = 'Previous buffer' })
map('n', '<leader>/', 'gcc', { remap = true, desc = 'Toggle comment' })
map('v', '<leader>/', 'gc', { remap = true, desc = 'Toggle comment' })
map('t', '<C-x>', '<C-\\><C-N>', { desc = 'Exit terminal mode' })
map('i', '<C-b>', '<ESC>^i', { desc = 'Move to beginning of line' })
map('i', '<C-e>', '<End>', { desc = 'Move to end of line' })
map('i', '<C-h>', '<Left>', { desc = 'Move left' })
map('i', '<C-l>', '<Right>', { desc = 'Move right' })
map('i', '<C-j>', '<Down>', { desc = 'Move down' })
map('i', '<C-k>', '<Up>', { desc = 'Move up' })
