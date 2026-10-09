local severity = vim.diagnostic.severity
local icons = {
  [severity.ERROR] = '\u{f057}',
  [severity.WARN] = '\u{f071}',
  [severity.INFO] = '\u{f05a}',
  [severity.HINT] = '\u{f0336}',
}

vim.diagnostic.config {
  update_in_insert = false,
  severity_sort = true,
  float = { border = 'rounded', source = 'if_many' },
  underline = { severity = { min = vim.diagnostic.severity.WARN } },

  signs = { text = icons },
  -- Rendered as a tinted badge (see colorscheme.lua): " <icon> message ".
  virtual_text = {
    spacing = 2,
    prefix = function(diagnostic) return ' ' .. icons[diagnostic.severity] end,
    format = function(diagnostic) return diagnostic.message .. ' ' end,
  },
  virtual_lines = false,

  jump = {
    on_jump = function(_, bufnr)
      vim.diagnostic.open_float {
        bufnr = bufnr,
        scope = 'cursor',
        focus = false,
      }
    end,
  },
}
