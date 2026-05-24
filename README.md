# core-architecture-shapes

A vector library for diagramming infrastructure topologies and network schematics.

## The Design Language

* Paths use a single, continuous 2px stroke profile.
* Shapes are aligned to a 48x48px coordinate plane.
* Boxed shapes represent network transit, firewalls, and middleware logic.
* Unboxed shapes represent data persistence, hardware primitives, and logical runtimes.

## Repository Structure

```text
core-architecture-shapes/
├── dist/                                   # Production ready assets
│   └── optimized/                          # Light/Dark-mode responsive, minified SVGs
├── src/                                    # Editable source
│   └── core-architecture-shapes-master.svg # The un-flattened master Inkscape workbook
└── scripts/                                # Optimization & pipeline configs
```

## Native Theme Support

All distribution shapes are compiled using `currentColor` for their vector fills.
They automatically inherit the text color of their parent container, meaning a single asset file supports both light and dark modes:

```css
/* Dark Mode implementation example */
@media (prefers-color-scheme: dark) {
  .topology-canvas {
    color: #ffffff; /* Shapes render white lines */
  }
}
```

## Development Workflow

This project uses `just` to manage the icon export and optimization pipeline.

### Prerequisites

* [Inkscape 1.2+](https://inkscape.org/)
* [Node.js](https://nodejs.org/)
* [Python 3](https://www.python.org/)

### Exporting Icons

To export icons from the master SVG and optimize them for production:

1. **Install dependencies:**

   ```bash
   just deps
   ```

2. **Export and Optimize:**

   ```bash
   just all
   ```

The icons will be exported to `dist/` and then optimized versions will be generated in `dist/optimized/`.

### Export Process Details

1. **Export:** The `scripts/export_icons.py` script parses the Inkscape master file, identifies all pages, and uses the Inkscape CLI to export each page as a standalone SVG. It automatically converts all objects to paths during export.
2. **Optimize:** SVGO is used with a custom configuration (`scripts/svgo.config.mjs`) to:
    * Minify the SVG data.
    * Remove Inkscape-specific metadata.
    * Convert all colors to `currentColor`.
    * Ensure a consistent `viewBox="0 0 48 48"`.

## License

MIT
