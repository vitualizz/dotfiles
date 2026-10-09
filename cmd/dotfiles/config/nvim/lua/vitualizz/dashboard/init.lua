-- Start screen (snacks.nvim dashboard), shown when nvim opens without files.
--   header.lua      ASCII art at the top
--   cheatsheet.lua  keymaps listed in the second column
-- Options: https://github.com/folke/snacks.nvim/blob/main/docs/dashboard.md

local ui = require 'vitualizz.ui'
ui.pill('VitualizzDashboardKey', require('vitualizz.plugins.colorscheme').palette.base09)

return {
  enabled = true,
  -- Quick-action keys as pills.
  formats = {
    key = function(item) return ui.pill_chunks('VitualizzDashboardKey', item.key) end,
  },
  preset = {
    header = require 'vitualizz.dashboard.header',
    keys = {
      { icon = ' ', key = 'f', desc = 'Find file', action = ":lua Snacks.dashboard.pick('files')" },
      { icon = ' ', key = 'g', desc = 'Find text', action = ":lua Snacks.dashboard.pick('live_grep')" },
      { icon = ' ', key = 'r', desc = 'Recent files', action = ":lua Snacks.dashboard.pick('oldfiles')" },
      { icon = ' ', key = 'n', desc = 'New file', action = ':ene | startinsert' },
      { icon = ' ', key = 'c', desc = 'Config', action = ":lua Snacks.dashboard.pick('files', { cwd = vim.fn.stdpath('config') })" },
      { icon = ' ', key = 'q', desc = 'Quit', action = ':qa' },
    },
  },
  sections = {
    { section = 'header' },
    { section = 'keys', gap = 1, padding = 1 },
    require 'vitualizz.dashboard.cheatsheet',
  },
}
