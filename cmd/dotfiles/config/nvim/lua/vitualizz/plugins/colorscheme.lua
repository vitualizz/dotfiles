local gh = require('vitualizz.pack').gh

-- Let the terminal's background (and its opacity) show through.
local transparent = true

-- Everblush, with the base16 palette NvChad uses (base46 v3.0), applied by
-- mini.base16 (part of mini.nvim), which also colors snacks, blink, which-key,
-- the picker, gitsigns and mini.
local everblush = {
  base00 = '#141b1e',
  base01 = '#1e2528',
  base02 = '#282f32',
  base03 = '#2d3437',
  base04 = '#3c4346',
  base05 = '#dadada',
  base06 = '#e4e4e4',
  base07 = '#dadada',
  base08 = '#e57474',
  base09 = '#fcb163',
  base0A = '#e5c76b',
  base0B = '#8ccf7e',
  base0C = '#6cbfbf',
  base0D = '#67b0e8',
  base0E = '#c47fd5',
  base0F = '#ef7d7d',
}

vim.pack.add { gh 'nvim-mini/mini.nvim' }
require('mini.base16').setup { palette = everblush }
vim.g.colors_name = 'everblush'

-- Everblush's grays are below 3:1, so mini.base16 renders comments, line
-- numbers and other secondary text nearly invisible. Same hue, three tiers
-- (contrast on Everblush's and on the terminal's background):
--   main      #dadada  12:1   current line number
--   secondary #7e888c  4.7:1  comments, folds, which-key values, trouble locations
--   dim       #636c6f  3.2:1  line numbers, listchars, end-of-buffer
local readable = {
  [everblush.base05] = { 'CursorLineNr' },
  ['#7e888c'] = { 'Comment', 'Folded', 'WhichKeyValue', 'TroubleLocation' },
  ['#636c6f'] = { 'LineNr', 'LineNrAbove', 'LineNrBelow', 'NonText', 'Whitespace', 'SpecialKey', 'EndOfBuffer' },
}
for color, groups in pairs(readable) do
  for _, name in ipairs(groups) do
    local hl = vim.api.nvim_get_hl(0, { name = name, link = false })
    hl.fg = color
    vim.api.nvim_set_hl(0, name, hl)
  end
end

if transparent then
  local canvas = {
    Normal = true,
    NormalNC = true,
    SignColumn = true,
    FoldColumn = true,
    LineNr = true,
    LineNrAbove = true,
    LineNrBelow = true,
  }
  for name in pairs(vim.api.nvim_get_hl(0, {})) do
    if canvas[name] or name:match '^GitSigns' or name:match '^MiniDiffSign' then
      local hl = vim.api.nvim_get_hl(0, { name = name, link = false })
      hl.bg, hl.ctermbg = nil, nil
      vim.api.nvim_set_hl(0, name, hl)
    end
  end
end

-- Everblush's UI shades from NvChad (base_30), for panels and cards.
local ui = {
  panel = '#10171a',
  input = '#1a2124',
  selection = '#272e31',
}

-- Floating windows (hover, which-key, diagnostics, terminal, lazygit) as
-- cards: solid panel background, borders in the same color so they only add
-- padding, and titles as colored blocks like the statusline.
vim.api.nvim_set_hl(0, 'NormalFloat', { fg = everblush.base05, bg = ui.panel })
vim.api.nvim_set_hl(0, 'FloatBorder', { fg = ui.panel, bg = ui.panel })
vim.api.nvim_set_hl(0, 'FloatTitle', { fg = everblush.base00, bg = everblush.base0D, bold = true })
vim.api.nvim_set_hl(0, 'WhichKeyFloat', { bg = ui.panel })
for name in pairs(vim.api.nvim_get_hl(0, {})) do
  if name:match '^DiagnosticFloating' then
    local hl = vim.api.nvim_get_hl(0, { name = name, link = false })
    hl.bg = ui.panel
    vim.api.nvim_set_hl(0, name, hl)
  end
end

-- Diagnostic virtual text as tinted blocks: the level's color on a 15% tint
-- of it, so errors and warnings read as badges at the end of the line.
-- mini.base16 paints warnings purple; make them yellow everywhere (text,
-- signs, underlines, floats) so severities read red/yellow/cyan/blue.
local purple, yellow = tonumber(everblush.base0E:sub(2), 16), everblush.base0A
for name in pairs(vim.api.nvim_get_hl(0, {})) do
  if name:match 'Warn' then
    local hl = vim.api.nvim_get_hl(0, { name = name, link = false })
    if hl.fg == purple then hl.fg = yellow end
    if hl.sp == purple then hl.sp = yellow end
    vim.api.nvim_set_hl(0, name, hl)
  end
end

local blend = require('vitualizz.ui').blend
for _, level in ipairs { 'Error', 'Warn', 'Info', 'Hint', 'Ok' } do
  local fg = vim.api.nvim_get_hl(0, { name = 'Diagnostic' .. level, link = false }).fg
  if fg then
    local color = string.format('#%06x', fg)
    vim.api.nvim_set_hl(0, 'DiagnosticVirtualText' .. level, { fg = color, bg = blend(color, everblush.base00, 0.15) })
  end
end

return { palette = everblush, ui = ui, transparent = transparent }
