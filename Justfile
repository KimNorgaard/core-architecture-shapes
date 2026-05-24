# core-architecture-shapes Justfile

set shell := ["bash", "-c"]

# Export, optimize, and update gallery
all: export-icons optimize update-gallery

# Install dependencies (svgo)
install:
    npm install -g svgo

# Prepare the export file (manual copy if you want to perform manual flattening)
prepare-export:
    cp src/core-architecture-shapes-master.svg src/core-architecture-shapes-export.svg
    @echo "Export file created at src/core-architecture-shapes-export.svg"
    @echo "Open this file in Inkscape, perform manual Union/Cleanup, then run 'just all'"

# Export icons from the master (or export) file
export-icons:
    python3 scripts/export_icons.py

# Optimize exported SVGs and apply currentColor
optimize:
    npx svgo --config scripts/svgo.config.mjs -f dist -o dist/optimized

# Update the visual gallery in README.md
update-gallery:
    python3 scripts/generate_gallery.py

# Clean up exported and optimized assets
clean:
    rm -rf dist/*.svg dist/optimized/*.svg
