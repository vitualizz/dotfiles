-- Building blocks for the statusline's visual language, shared by the rest
-- of the UI: rounded "pills" (dark bold text on a color, with rounded ends)
-- and color blending for tinted backgrounds.
local M = {}

M.cap = { left = '\u{e0b6}', right = '\u{e0b4}' }

-- Defines `name` (the pill body) and `name .. 'Cap'` (its rounded ends).
-- `cap_bg` is the background behind the pill; nil for the transparent canvas.
function M.pill(name, color, cap_bg)
  local dark = require('vitualizz.plugins.colorscheme').palette.base00
  vim.api.nvim_set_hl(0, name, { fg = dark, bg = color, bold = true })
  vim.api.nvim_set_hl(0, name .. 'Cap', { fg = color, bg = cap_bg })
end

-- Text chunks { text, hl } for a pill, in the format snacks.dashboard uses.
function M.pill_chunks(name, text)
  return {
    { M.cap.left, hl = name .. 'Cap' },
    { ' ' .. text .. ' ', hl = name },
    { M.cap.right, hl = name .. 'Cap' },
  }
end

-- Mixes two '#rrggbb' colors; alpha is the share of `fg` (0..1).
function M.blend(fg, bg, alpha)
  local function channel(hex, i) return tonumber(hex:sub(i, i + 1), 16) end
  local function mix(i) return math.floor(channel(fg, i) * alpha + channel(bg, i) * (1 - alpha) + 0.5) end
  return string.format('#%02x%02x%02x', mix(2), mix(4), mix(6))
end

return M
