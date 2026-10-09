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
    NormalFloat = true,
    WhichKeyFloat = true,
    SignColumn = true,
    FoldColumn = true,
    LineNr = true,
    LineNrAbove = true,
    LineNrBelow = true,
  }
  for name in pairs(vim.api.nvim_get_hl(0, {})) do
    if canvas[name] or name:match '^GitSigns' or name:match '^MiniDiffSign' or name:match '^DiagnosticFloating' then
      local hl = vim.api.nvim_get_hl(0, { name = name, link = false })
      hl.bg, hl.ctermbg = nil, nil
      vim.api.nvim_set_hl(0, name, hl)
    end
  end
end

return { palette = everblush, transparent = transparent }
