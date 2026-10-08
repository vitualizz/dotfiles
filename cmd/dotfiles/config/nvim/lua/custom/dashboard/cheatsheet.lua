-- Keymaps shown in the start screen's second column, grouped by topic.
-- This list is only documentation: keymaps are defined in init.lua and
-- lua/custom/plugins/*.lua. Keep both in sync when changing a keymap.

local groups = {
  {
    'Find',
    { '<leader>sf', 'Files' },
    { '<leader>sg', 'Text (grep)' },
    { '<leader>s.', 'Recent files' },
    { '<leader><leader>', 'Open buffers' },
    { '<leader>sk', 'Keymaps' },
  },
  {
    'Code',
    { 'grd', 'Go to definition' },
    { 'grr', 'References' },
    { 'grn', 'Rename' },
    { 'gra', 'Code action' },
    { 'K', 'Hover docs' },
    { '<leader>f', 'Format buffer' },
  },
  {
    'Tools',
    { '<C-n>', 'File explorer' },
    { '<A-i>', 'Floating terminal' },
    { '<A-h>', 'Bottom terminal' },
    { '<leader>gg', 'Lazygit' },
  },
  {
    'Git',
    { ']c / [c', 'Next / previous change' },
    { '<leader>hp', 'Preview change' },
    { '<leader>hb', 'Blame line' },
  },
}

local section = { pane = 2, padding = 1 }
for _, group in ipairs(groups) do
  table.insert(section, { text = { { group[1], hl = 'title' } } })
  for i = 2, #group do
    local key, desc = group[i][1], group[i][2]
    table.insert(section, {
      text = { { key, hl = 'key', width = 18 }, { desc, hl = 'desc' } },
      padding = i == #group and 1 or 0,
    })
  end
end

return section
