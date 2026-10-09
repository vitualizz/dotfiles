-- Keymaps shown in the start screen's second column, grouped by topic.
-- This list is only documentation: keymaps are defined in init.lua and
-- lua/vitualizz/plugins/*.lua. Keep both in sync when changing a keymap.

local groups = {
  {
    'Find',
    { '<leader>sf', 'Files' },
    { '<leader>sg', 'Text (grep)' },
    { '<leader>s.', 'Recent files' },
    { '<leader><leader>', 'Open buffers' },
    { '<leader>sR', 'Search and replace' },
  },
  {
    'Move',
    { 's', 'Jump (flash)' },
    { 'S', 'Select syntax node' },
    { 'gsa / gsd / gsr', 'Surround add/del/replace' },
  },
  {
    'Code',
    { 'grd', 'Go to definition' },
    { 'grr', 'References' },
    { 'grn', 'Rename' },
    { 'gra', 'Code action' },
    { 'K', 'Hover docs' },
    { '<leader>f', 'Format buffer' },
    { '<leader>xx', 'Diagnostics list' },
  },
  {
    'Tools',
    { '<C-n>', 'File explorer' },
    { '<A-i>', 'Floating terminal' },
    { '<A-h>', 'Bottom terminal' },
    { '<leader>gg', 'Lazygit' },
    { '<leader>bd', 'Close buffer' },
  },
  {
    'Git',
    { ']c / [c', 'Next / previous change' },
    { '<leader>hp', 'Preview change' },
    { '<leader>hb', 'Blame line' },
  },
}

-- Group titles are pills, one color per topic, like the statusline modes.
local ui = require 'vitualizz.ui'
local c = require('vitualizz.plugins.colorscheme').palette
local title_colors = { Find = c.base0D, Move = c.base0E, Code = c.base0B, Tools = c.base0A, Git = c.base08 }
for title, color in pairs(title_colors) do
  ui.pill('VitualizzCheatsheet' .. title, color)
end

local section = { pane = 2, padding = 1 }
for _, group in ipairs(groups) do
  table.insert(section, { text = ui.pill_chunks('VitualizzCheatsheet' .. group[1], group[1]) })
  for i = 2, #group do
    local key, desc = group[i][1], group[i][2]
    table.insert(section, {
      text = { { key, hl = 'key', width = 18 }, { desc, hl = 'desc' } },
      padding = i == #group and 1 or 0,
    })
  end
end

return section
