# Atajos

`<leader>` es `Espacio`. `<leader>ch` abre esta hoja; `<leader>wK` muestra todos los atajos con which-key.

## 🔍 Buscar

| Atajo | Acción |
|-------|--------|
| `<leader>ff` · `<leader>sf` | Archivos |
| `<leader>fa` | Todos los archivos (ocultos e ignorados incluidos) |
| `<leader>fw` · `<leader>sg` | Texto en el proyecto (grep) |
| `<leader>sw` | La palabra bajo el cursor |
| `<leader>fz` | Buscar dentro del archivo actual |
| `<leader>s/` | Grep solo en archivos abiertos |
| `<leader>fo` · `<leader>s.` | Archivos recientes |
| `<leader>fb` · `<leader><leader>` | Buffers abiertos |
| `<leader>fh` · `<leader>sh` | Ayuda de Neovim |
| `<leader>ma` | Marcas |
| `<leader>sk` | Todos los atajos |
| `<leader>sc` | Comandos |
| `<leader>sn` | Archivos de tu config de Neovim |
| `<leader>sr` | Retomar la última búsqueda |
| `<leader>ss` | Todos los buscadores |
| `<leader>sR` | Buscar y reemplazar en el proyecto (o en la selección) |

## 💻 Código (LSP)

| Atajo | Acción |
|-------|--------|
| `grd` | Ir a la definición |
| `grr` | Referencias |
| `gri` | Implementaciones |
| `grt` | Definición del tipo |
| `grD` | Declaración |
| `grn` | Renombrar |
| `gra` | Code action |
| `gO` / `gW` | Símbolos del archivo / del proyecto |
| `K` | Documentación flotante |
| `<leader>fm` | Formatear |
| `<leader>th` | Mostrar u ocultar inlay hints |

## 🩺 Diagnósticos

| Atajo | Acción |
|-------|--------|
| `]d` / `[d` | Siguiente / anterior error |
| `<leader>dd` / `<leader>db` | Lista de diagnósticos del proyecto / del archivo (trouble) |
| `<leader>dS` | Símbolos (trouble) |
| `<leader>dl` | Definiciones y referencias del LSP (trouble) |
| `<leader>dt` | TODOs del proyecto (trouble) |
| `<leader>dq` | Quickfix (trouble) |
| `<leader>ds` · `<leader>q` | Diagnósticos en la location list |

## 📄 Buffers

| Atajo | Acción |
|-------|--------|
| `<Tab>` / `<S-Tab>` | Siguiente / anterior buffer |
| `<leader>b` | Buffer nuevo |
| `<leader>x` | Cerrar buffer (sin romper las ventanas) |
| `<leader>X` | Cerrar los demás buffers |
| `<C-s>` | Guardar |
| `<C-c>` | Copiar todo el archivo |

## 🛠️ Herramientas

| Atajo | Acción |
|-------|--------|
| `<C-n>` / `<leader>e` | Explorador: abrir o cerrar / enfocar el archivo actual |
| `<A-i>` | Terminal flotante |
| `<A-h>` / `<A-v>` | Terminal abajo / a la derecha |
| `<C-x>` · `<Esc><Esc>` | Salir del modo terminal |
| `<leader>gg` | Lazygit |
| `<leader>cm` / `<leader>gt` | Commits / estado de git |
| `<leader>n` / `<leader>rn` | Números de línea / números relativos |
| `<leader>wk` | Which-key: buscar un atajo |

## 🌿 Git (gitsigns)

| Atajo | Acción |
|-------|--------|
| `]c` / `[c` | Siguiente / anterior cambio |
| `<leader>hp` / `<leader>hi` | Vista previa del cambio (flotante / en línea) |
| `<leader>hs` / `<leader>hr` | Stage / reset del cambio |
| `<leader>hS` / `<leader>hR` | Stage / reset del archivo completo |
| `<leader>hb` | Blame de la línea |
| `<leader>tb` | Blame en vivo (activar o desactivar) |
| `<leader>hd` / `<leader>hD` | Diff contra el index / el último commit |
| `<leader>hq` / `<leader>hQ` | Cambios en quickfix (archivo / repo) |
| `<leader>tw` | Diff por palabra |
| `ih` | Seleccionar el cambio (objeto de texto) |

## ✏️ Edición

| Atajo | Acción |
|-------|--------|
| `gsa` / `gsd` / `gsr` | Agregar / borrar / reemplazar rodeado. Ej: `gsaiw"` pone comillas a una palabra |
| `va)` · `ci'` · `daq` | Objetos de texto mejorados (mini.ai); `aa` / `ii` = el siguiente |
| `<leader>/` · `gcc` / `gc` | Comentar línea / selección |

## 🦘 Movimiento

| Atajo | Acción |
|-------|--------|
| `s` | Saltar a cualquier lugar visible (flash) |
| `S` | Seleccionar un bloque de sintaxis (flash) |
| `<C-h/j/k/l>` | Moverte entre ventanas |
| `<Esc>` | Quitar el resaltado de búsqueda |

## ⌨️ Modo insert y autocompletado

| Atajo | Acción |
|-------|--------|
| `<Tab>` / `<S-Tab>` | Siguiente / anterior sugerencia (y saltar dentro de snippets) |
| `<CR>` | Aceptar la sugerencia |
| `<C-Space>` | Abrir sugerencias / documentación |
| `<C-e>` | Cerrar sugerencias; si no hay, ir al final de la línea |
| `<C-d>` / `<C-f>` | Subir / bajar la documentación |
| `<C-b>` | Ir al inicio de la línea |
| `<C-h/j/k/l>` | Mover el cursor |
