# bigcrock

[![Versión de Go](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Licencia: MIT](https://img.shields.io/badge/Licencia-MIT-yellow.svg)](LICENSE)
[![Releases](https://img.shields.io/github/v/release/Ars-byte/bigcrock?color=brightgreen)](https://github.com/Ars-byte/bigcrock/releases)

> **Idioma / Language**: **Español** | [English](README.md)

**bigcrock** es un analizador de uso de espacio en disco y ranking de directorios en árbol para la terminal, rápido, ligero e inspirado en el diseño de Ranger, desarrollado en Go utilizando [Bubble Tea](https://github.com/charmbracelet/bubbletea) y [Lip Gloss](https://github.com/charmbracelet/lipgloss).

Escanea directorios y los clasifica según su peso, mostrando de forma visual e inmediata qué carpetas y directorios ocultos (como `.local`, `.cache`, etc.) consumen más espacio.

---

## Características

- **Columnas Miller al estilo Ranger**:
  - **Columna izquierda**: Directorio padre para mantener el contexto de navegación.
  - **Columna central**: Lista interactiva clasificada por peso con etiquetas de tamaño y barras de progreso visuales.
  - **Columna derecha**: Inspector detallado con ranking de subdirectorios internos, archivos más pesados, metadatos y previsualización.
- **Ranking de los mayores consumidores de espacio**: Muestra de inmediato qué carpetas ocupan más (`#1`, `#2`, `#3`...).
- **Directorios ocultos incluidos**: Explora fácilmente carpetas ocultas del sistema como `.local`, `.config` y `.cache` (conmutable con `.` o `H`).
- **Navegación interactiva en árbol**: Entra a cualquier directorio o subdirectorio (`Enter` / `l` / `→`) y sube de nivel (`←` / `h` / `Backspace`).
- **Escaneo asíncrono en segundo plano**: Cálculo concurrente sin bloqueos para mantener la interfaz a 60 FPS incluso con árboles de miles de archivos.
- **Caché en memoria inteligente**: Respuesta instantánea al volver a carpetas previamente analizadas.
- **Soporte bilingüe integrado**: Alterna entre español e inglés al instante pulsando `Tab` o `t`.
- **Múltiples modos de ordenación**: Por tamaño descendente/ascendente, orden alfabético o cantidad de elementos (`s`).
- **Pantalla Acerca de**: Vista limpia y minimalista (`?` o `a`).
- **Arquitectura modular**: Estructura organizada siguiendo los patrones del proyecto Crayoz.

---

## Atajos de teclado

| Tecla | Acción |
| --- | --- |
| `↑` / `k` | Mover cursor hacia arriba |
| `↓` / `j` | Mover cursor hacia abajo |
| `Enter` / `→` / `l` | Abrir / entrar a la carpeta seleccionada |
| `←` / `h` / `Backspace` | Subir al directorio padre |
| `Home` / `g` | Ir al primer elemento (#1 más pesado) |
| `End` / `G` | Ir al último elemento |
| `PgUp` / `PgDn` | Desplazar una página arriba / abajo |
| `s` | Alternar orden (Peso ↓ / Peso ↑ / Nombre A-Z / Nombre Z-A / Elementos) |
| `.` / `H` | Mostrar / ocultar archivos y carpetas ocultas |
| `r` | Recalcular y re-escanear directorio actual |
| `Tab` / `t` | **Cambiar idioma (Español ↔ English)** |
| `?` / `a` | **Ver pantalla Acerca de** |
| `Cualquier tecla` | Salir de pantalla Acerca de |
| `q` / `Ctrl+C` | Salir de la aplicación |

---

## Instalación y Compilación

### Requisitos
- [Go](https://go.dev/) 1.22 o superior.

### Compilar desde el código fuente
```bash
git clone https://github.com/Ars-byte/bigcrock.git
cd bigcrock
go mod download
go build -ldflags="-s -w" -o bigcrock .
./bigcrock
```

### Uso
```bash
# Analizar directorio actual
./bigcrock

# Analizar la carpeta personal (Home)
./bigcrock ~

# Analizar una ruta específica
./bigcrock /var/log
```

---

## Licencia

Este proyecto está bajo la Licencia MIT. Consulta el archivo [LICENSE](LICENSE) para más detalles.
