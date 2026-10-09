-- Powerline-style statusline (mini.statusline), matching the Starship prompt:
--  MODE  git · diff · diagnostics  filename        fileinfo  line:col
-- The mode block and both ends take the color of the current mode.
local theme = require 'vitualizz.plugins.colorscheme'
local c = theme.palette

local glyph = { right = '\u{e0b0}', left = '\u{e0b2}', cap_left = '\u{e0b6}', cap_right = '\u{e0b4}' }
local bg = { info = c.base02, file = c.base01, edge = (not theme.transparent) and c.base00 or nil }
local secondary = '#8f989b'

local modes = {
  Normal = c.base0D,
  Insert = c.base0B,
  Visual = c.base0E,
  Replace = c.base08,
  Command = c.base0A,
  Other = c.base0C,
}

local hl = vim.api.nvim_set_hl
for mode, color in pairs(modes) do
  local name = 'MiniStatuslineMode' .. mode
  hl(0, name, { fg = c.base00, bg = color, bold = true })
  hl(0, name .. 'ToInfo', { fg = color, bg = bg.info })
  hl(0, name .. 'ToFile', { fg = color, bg = bg.file })
  hl(0, name .. 'Edge', { fg = color, bg = bg.edge })
end
hl(0, 'MiniStatuslineInfoToFile', { fg = bg.info, bg = bg.file })
hl(0, 'MiniStatuslineDevinfo', { fg = secondary, bg = bg.info })
hl(0, 'MiniStatuslineFileinfo', { fg = secondary, bg = bg.info })
hl(0, 'MiniStatuslineFilename', { fg = c.base05, bg = bg.file })
hl(0, 'MiniStatuslineInactive', { fg = secondary, bg = bg.file })
hl(0, 'StatusLine', { fg = c.base05, bg = bg.edge })

local function sep(group, char) return '%#' .. group .. '#' .. char end

local function active()
  local s = MiniStatusline
  local mode, mode_hl = s.section_mode { trunc_width = 120 }
  -- section_diff shows "-" when the file has no changes; hide it instead.
  local diff = s.section_diff { trunc_width = 75 }
  if diff:match '%-$' then diff = '' end
  local info = vim.tbl_filter(function(x) return x ~= '' end, {
    s.section_git { trunc_width = 40 },
    diff,
    s.section_diagnostics { trunc_width = 75 },
    s.section_lsp { trunc_width = 75 },
  })
  local fileinfo = s.section_fileinfo { trunc_width = 120 }
  local position = { s.section_searchcount { trunc_width = 75 }, '%2l:%-2v' }

  return s.combine_groups {
    sep(mode_hl .. 'Edge', glyph.cap_left),
    { hl = mode_hl, strings = { mode } },
    #info > 0 and sep(mode_hl .. 'ToInfo', glyph.right) or sep(mode_hl .. 'ToFile', glyph.right),
    #info > 0 and { hl = 'MiniStatuslineDevinfo', strings = info } or '',
    #info > 0 and sep('MiniStatuslineInfoToFile', glyph.right) or '',
    '%<',
    { hl = 'MiniStatuslineFilename', strings = { s.section_filename { trunc_width = 140 } } },
    '%=',
    fileinfo ~= '' and sep('MiniStatuslineInfoToFile', glyph.left) or '',
    fileinfo ~= '' and { hl = 'MiniStatuslineFileinfo', strings = { fileinfo } } or '',
    fileinfo ~= '' and sep(mode_hl .. 'ToInfo', glyph.left) or sep(mode_hl .. 'ToFile', glyph.left),
    { hl = mode_hl, strings = position },
    sep(mode_hl .. 'Edge', glyph.cap_right),
  }
end

require('mini.statusline').setup {
  use_icons = vim.g.have_nerd_font,
  content = { active = active },
}
