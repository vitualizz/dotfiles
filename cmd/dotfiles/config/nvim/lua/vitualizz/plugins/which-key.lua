local gh = require('vitualizz.pack').gh

vim.pack.add { gh 'folke/which-key.nvim' }
require('which-key').setup {
  delay = 0,
  icons = { mappings = vim.g.have_nerd_font },
  spec = {
    { '<leader>s', group = '[S]earch', mode = { 'n', 'v' } },
    { '<leader>t', group = '[T]oggle' },
    { '<leader>h', group = 'Git [H]unk', mode = { 'n', 'v' } },
    { 'gr', group = 'LSP Actions', mode = { 'n' } },
    { '<leader>f', group = '[F]ind' },
    { '<leader>d', group = '[D]iagnostics (trouble)' },
    { '<leader>w', group = '[W]hich-key' },
    { 'gs', group = 'Surround', mode = { 'n', 'x' } },
  },
}

vim.keymap.set('n', '<leader>wK', '<cmd>WhichKey<CR>', { desc = 'Which-key: all keymaps' })
vim.keymap.set('n', '<leader>wk', function() vim.cmd('WhichKey ' .. vim.fn.input 'WhichKey: ') end, { desc = 'Which-key: query lookup' })
