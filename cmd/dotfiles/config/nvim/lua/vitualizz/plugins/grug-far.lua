local gh = require('vitualizz.pack').gh

vim.pack.add { { src = gh 'MagicDuck/grug-far.nvim', version = vim.version.range '1.*' } }
require('grug-far').setup {}

-- <leader>sr stays "search resume"; R is for Replace.
vim.keymap.set('n', '<leader>sR', function() require('grug-far').open { transient = true } end, { desc = '[S]earch and [R]eplace in project' })
vim.keymap.set('x', '<leader>sR', function() require('grug-far').with_visual_selection { transient = true } end, { desc = '[S]earch and [R]eplace selection' })
