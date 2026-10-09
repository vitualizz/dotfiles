-- Fuzzy finder: snacks.picker (enabled in snacks.lua), with the keymaps
-- kickstart had for telescope. It also replaces vim.ui.select.
local function pick(source, opts)
  return function() Snacks.picker[source](opts) end
end

local map = vim.keymap.set
map('n', '<leader>sh', pick 'help', { desc = '[S]earch [H]elp' })
map('n', '<leader>sk', pick 'keymaps', { desc = '[S]earch [K]eymaps' })
map('n', '<leader>sf', pick 'files', { desc = '[S]earch [F]iles' })
map('n', '<leader>ss', pick 'pickers', { desc = '[S]earch [S]elect picker' })
map({ 'n', 'v' }, '<leader>sw', pick 'grep_word', { desc = '[S]earch current [W]ord' })
map('n', '<leader>sg', pick 'grep', { desc = '[S]earch by [G]rep' })
map('n', '<leader>sd', pick 'diagnostics', { desc = '[S]earch [D]iagnostics' })
map('n', '<leader>sr', pick 'resume', { desc = '[S]earch [R]esume' })
map('n', '<leader>s.', pick 'recent', { desc = '[S]earch Recent Files ("." for repeat)' })
map('n', '<leader>sc', pick 'commands', { desc = '[S]earch [C]ommands' })
map('n', '<leader><leader>', pick 'buffers', { desc = '[ ] Find existing buffers' })
map('n', '<leader>/', pick 'lines', { desc = '[/] Fuzzily search in current buffer' })
map('n', '<leader>s/', pick 'grep_buffers', { desc = '[S]earch [/] in Open Files' })
map('n', '<leader>sn', pick('files', { cwd = vim.fn.stdpath 'config' }), { desc = '[S]earch [N]eovim files' })

vim.api.nvim_create_autocmd('LspAttach', {
  group = vim.api.nvim_create_augroup('vitualizz-picker-lsp-attach', { clear = true }),
  callback = function(event)
    local function lsp(lhs, source, desc) map('n', lhs, pick(source), { buffer = event.buf, desc = desc }) end
    lsp('grr', 'lsp_references', '[G]oto [R]eferences')
    lsp('gri', 'lsp_implementations', '[G]oto [I]mplementation')
    lsp('grd', 'lsp_definitions', '[G]oto [D]efinition')
    lsp('gO', 'lsp_symbols', 'Open Document Symbols')
    lsp('gW', 'lsp_workspace_symbols', 'Open Workspace Symbols')
    lsp('grt', 'lsp_type_definitions', '[G]oto [T]ype Definition')
  end,
})
