# tuipr - PR Lifecycle TUI Manager

Un gestor de Pull Requests en terminal inspirado en **LazyGit** y **tuicr**, con navegación estilo Neovim y colores Catppuccin Mocha.

## Características

- 🎨 **Colores Catppuccin Mocha** - Paleta de colores elegante
- ⌨️ **Navegación estilo LazyGit** - Paneles numerados (1, 2, 3)
- 📝 **Editor estilo Vim** - Modo Normal/Insert con `i` y `Esc`
- 🔗 **Integración con `gh`** - Lista, crea y mergea PRs directamente
- 🛠️ **Editor de conflictos** - Abre Neovim para resolver conflictos

## Instalación

```bash
# Clonar el repositorio
git clone https://github.com/ecortez/tuipr.git
cd tuipr

# Compilar
go build -o tuipr .

# Instalar (opcional)
mv tuipr /usr/local/bin/
```

## Uso

```bash
tuipr              # Abrir dashboard principal
tuipr -c           # Ir directo a Create PR
tuipr -m           # Ir a Merge (seleccionar PR)
tuipr -m 134       # Merge directo del PR #134
tuipr --help       # Mostrar ayuda
```

## Navegación

### Dashboard
| Tecla | Acción |
|-------|--------|
| `1` | Ir al panel de PRs |
| `2` | Ir al panel de Detalles |
| `3` | Ir al panel de Status/Conflicts |
| `j` / `k` | Navegar arriba/abajo |
| `c` | Crear nuevo PR |
| `m` | Merge PR seleccionado |
| `e` | Abrir nvim para conflictos |
| `r` | Refrescar lista |
| `q` | Salir |

### Create PR
| Tecla | Acción |
|-------|--------|
| `i` | Modo Insert |
| `Esc` | Volver a Normal |
| `j` / `k` | Navegar campos |
| `Tab` | Cambiar entre paneles |
| `Ctrl+s` | Crear PR |

### Merge PR
| Tecla | Acción |
|-------|--------|
| `1` | Merge Commit |
| `2` | Squash |
| `3` | Rebase |
| `d` | Toggle delete remote |
| `D` | Toggle delete local |
| `Enter` | Ejecutar merge |

## Diseño

```
┌──────────────────────────────────────────────────────────────┐
│  TUIPR - PR Lifecycle Manager | Branch: feature/auth        │
│  ─────────────────────────────────────────────────────────  │
│  ┌────────────┐  ┌──────────────┐  ┌─────────────────┐    │
│  │ (1) PRs    │  │ (2) Details  │  │ (3) Conflicts   │    │
│  │ ────────── │  │ ──────────── │  │ ─────────────── │    │
│  │ 🟢 #142    │  │ Fix auth     │  │ CI Status       │    │
│  │ ⚠️ #141    │  │ Author: @you │  │ ✓ Passed        │    │
│  │ 🟢 #140    │  │ Reviews: ✓   │  │ Conflicts       │    │
│  │            │  │              │  │ ✓ No conflicts  │    │
│  └────────────┘  └──────────────┘  └─────────────────┘    │
│  ─────────────────────────────────────────────────────────  │
│  [1] PRs [2] Details [3] Conflicts                          │
│  [j/k] Navigate  [c] Create  [m] Merge  [e] nvim  [r] Refresh │
└──────────────────────────────────────────────────────────────┘
```

## Requisitos

- Go 1.21+
- GitHub CLI (`gh`)
- Neovim (para resolver conflictos)
- Terminal con soporte para 256 colores

## Licencia

MIT
