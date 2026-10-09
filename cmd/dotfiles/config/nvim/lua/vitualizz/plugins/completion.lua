local gh = require('vitualizz.pack').gh

vim.pack.add { { src = gh 'L3MON4D3/LuaSnip', version = vim.version.range '2.*' } }
require('luasnip').setup {}

vim.pack.add { { src = gh 'saghen/blink.cmp', version = vim.version.range '1.*' } }

local ui = require 'vitualizz.ui'
local theme = require 'vitualizz.plugins.colorscheme'
local c, shade = theme.palette, theme.ui

-- Each completion kind is a pill colored by category, like the statusline.
local kind_colors = {}
for color, kinds in pairs {
  [c.base0D] = { 'Function', 'Method', 'Constructor' },
  [c.base08] = { 'Variable', 'Field', 'Property' },
  [c.base0A] = { 'Class', 'Interface', 'Struct', 'Enum', 'Module', 'TypeParameter' },
  [c.base0E] = { 'Keyword', 'Operator', 'Snippet' },
  [c.base09] = { 'Constant', 'EnumMember', 'Value', 'Event' },
  [c.base0C] = { 'File', 'Folder', 'Reference', 'Unit', 'Color' },
  [c.base05] = { 'Text' },
} do
  for _, kind in ipairs(kinds) do
    kind_colors[kind] = color
    ui.pill('VitualizzKind' .. kind, color, shade.panel)
    vim.api.nvim_set_hl(0, 'BlinkCmpKind' .. kind, { fg = color })
  end
end

local function kind_pill_highlight(ctx)
  local group = 'VitualizzKind' .. (kind_colors[ctx.kind] and ctx.kind or 'Text')
  local left = #ui.cap.left
  local body = left + #ctx.kind_icon + 2
  return {
    { 0, left, group = group .. 'Cap', priority = 20000 },
    { left, body, group = group, priority = 20000 },
    { body, body + #ui.cap.right, group = group .. 'Cap', priority = 20000 },
  }
end

-- Menu, documentation and signature as cards (see colorscheme.lua floats).
local hl = vim.api.nvim_set_hl
for _, win in ipairs { 'BlinkCmpMenu', 'BlinkCmpDoc', 'BlinkCmpSignatureHelp' } do
  hl(0, win, { fg = c.base05, bg = shade.panel })
  hl(0, win .. 'Border', { fg = shade.panel, bg = shade.panel })
end
hl(0, 'BlinkCmpDocSeparator', { fg = shade.selection, bg = shade.panel })
hl(0, 'BlinkCmpMenuSelection', { bg = shade.selection, bold = true })
hl(0, 'BlinkCmpLabelMatch', { fg = c.base0D, bold = true })
hl(0, 'BlinkCmpLabelDescription', { fg = '#8f989b' })
hl(0, 'BlinkCmpLabelDetail', { fg = '#8f989b' })

require('blink.cmp').setup {
  keymap = {
    preset = 'default',
  },

  appearance = {
    nerd_font_variant = 'mono',
  },

  completion = {
    menu = {
      border = 'rounded',
      draw = {
        columns = { { 'kind_icon' }, { 'label', 'label_description', gap = 1 }, { 'kind' } },
        components = {
          kind_icon = {
            text = function(ctx) return ui.cap.left .. ' ' .. ctx.kind_icon .. ' ' .. ui.cap.right end,
            highlight = kind_pill_highlight,
          },
        },
      },
    },
    documentation = { auto_show = false, auto_show_delay_ms = 500, window = { border = 'rounded' } },
  },

  sources = {
    default = { 'lsp', 'path', 'snippets' },
  },

  snippets = { preset = 'luasnip' },

  fuzzy = { implementation = 'lua' },

  signature = { enabled = true, window = { border = 'rounded' } },
}
