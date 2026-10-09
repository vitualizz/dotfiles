local gh = require('vitualizz.pack').gh

-- Everblush, with the base16 palette NvChad uses (base46 v3.0), applied by
-- mini.base16 (part of mini.nvim), which also colors snacks, blink, which-key,
-- telescope, gitsigns and mini.
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

-- Everblush's grays are below 3:1 on its background, so comments are hard to
-- read. Same hue, lifted to 4.5:1 (comments) and 3:1 (line numbers).
vim.api.nvim_set_hl(0, 'Comment', { fg = '#798488' })
vim.api.nvim_set_hl(0, 'LineNr', { fg = '#5e676a' })
