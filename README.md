# core-architecture-shapes

A vector library for diagramming infrastructure topologies and network schematics.

## The Design Language

* Paths use a single, continuous 2px stroke profile.
* Shapes are aligned to an absolute 48x48px coordinate plane.
* Boxed shapes represent network transit, firewalls, and middleware logic.
* Unboxed shapes represent data persistence, hardware primitives, and logical runtimes.

## Repository Structure

```text
core-architecture-shapes/
├── dist/                                   # Production ready assets (tracked in git)
│   └── optimized/                          # Light/Dark-mode responsive, minified SVGs
├── src/                                    # Editable source
│   └── core-architecture-shapes-master.svg # The un-flattened master Inkscape workbook
└── scripts/                                # Optimization & pipeline configs
```

## Native Theme Support

All distribution shapes are compiled using `currentColor` for their vector fills.
They inherit the text color of their parent container, meaning a single asset file supports both light and dark modes:

```css
/* Dark Mode implementation example */
@media (prefers-color-scheme: dark) {
  .topology-canvas {
    color: #ffffff; /* Shapes render white lines */
  }
}
```

## Icon Gallery

| Icon | Name | Icon | Name | Icon | Name |
| :---: | :--- | :---: | :--- | :---: | :--- |
| <img src='./dist/optimized/bridge-l2.svg' width='48' height='48' /> | Bridge L2 | <img src='./dist/optimized/cloud.svg' width='48' height='48' /> | Cloud | <img src='./dist/optimized/container-group.svg' width='48' height='48' /> | Container Group | 
| <img src='./dist/optimized/container.svg' width='48' height='48' /> | Container | <img src='./dist/optimized/cpu.svg' width='48' height='48' /> | Cpu | <img src='./dist/optimized/data-volume.svg' width='48' height='48' /> | Data Volume | 
| <img src='./dist/optimized/database.svg' width='48' height='48' /> | Database | <img src='./dist/optimized/datacenter.svg' width='48' height='48' /> | Datacenter | <img src='./dist/optimized/endpoint.svg' width='48' height='48' /> | Endpoint | 
| <img src='./dist/optimized/firewall-virtual.svg' width='48' height='48' /> | Firewall Virtual | <img src='./dist/optimized/firewall.svg' width='48' height='48' /> | Firewall | <img src='./dist/optimized/gateway.svg' width='48' height='48' /> | Gateway | 
| <img src='./dist/optimized/hard-drive.svg' width='48' height='48' /> | Hard Drive | <img src='./dist/optimized/laptop.svg' width='48' height='48' /> | Laptop | <img src='./dist/optimized/load-balancer.svg' width='48' height='48' /> | Load Balancer | 
| <img src='./dist/optimized/memory.svg' width='48' height='48' /> | Memory | <img src='./dist/optimized/object-storage-appliance.svg' width='48' height='48' /> | Object Storage Appliance | <img src='./dist/optimized/object-storage-bucket.svg' width='48' height='48' /> | Object Storage Bucket | 
| <img src='./dist/optimized/queue.svg' width='48' height='48' /> | Queue | <img src='./dist/optimized/rack.svg' width='48' height='48' /> | Rack | <img src='./dist/optimized/router-virtual.svg' width='48' height='48' /> | Router Virtual | 
| <img src='./dist/optimized/router.svg' width='48' height='48' /> | Router | <img src='./dist/optimized/server-group.svg' width='48' height='48' /> | Server Group | <img src='./dist/optimized/server.svg' width='48' height='48' /> | Server | 
| <img src='./dist/optimized/site.svg' width='48' height='48' /> | Site | <img src='./dist/optimized/ssd.svg' width='48' height='48' /> | Ssd | <img src='./dist/optimized/storage-appliance.svg' width='48' height='48' /> | Storage Appliance | 
| <img src='./dist/optimized/switch-l2-virtual.svg' width='48' height='48' /> | Switch L2 Virtual | <img src='./dist/optimized/switch-l2.svg' width='48' height='48' /> | Switch L2 | <img src='./dist/optimized/switch-l3-virtual.svg' width='48' height='48' /> | Switch L3 Virtual | 
| <img src='./dist/optimized/switch-l3.svg' width='48' height='48' /> | Switch L3 | <img src='./dist/optimized/virtual-machine.svg' width='48' height='48' /> | Virtual Machine | <img src='./dist/optimized/vlan.svg' width='48' height='48' /> | Vlan | 
| <img src='./dist/optimized/vrf.svg' width='48' height='48' /> | Vrf | <img src='./dist/optimized/vtep.svg' width='48' height='48' /> | Vtep | <img src='./dist/optimized/vxlan.svg' width='48' height='48' /> | Vxlan | 

## Development Workflow

This project uses `just` to manage the icon export and optimization pipeline.

### Prerequisites

* [Inkscape 1.2+](https://inkscape.org/)
* [Node.js](https://nodejs.org/)
* [Python 3](https://www.python.org/)

### Standard Export Pipeline

To export icons from the master SVG and optimize them for production:

1. **Install dependencies:**

    ```bash
    just install
    ```

2. **Run full pipeline:**

    ```bash
    just all
    ```

This will export pages from `src/core-architecture-shapes-master.svg` to `dist/`, then generate optimized versions in `dist/optimized/`.

### Advanced Workflow (Manual Flattening)

If an icon requires complex manual geometry (e.g., specific boolean operations that the automated union cannot handle):

1. **Prepare a flattened source:**

    ```bash
    just prepare-export
    ```

2. Open `src/core-architecture-shapes-export.svg` in Inkscape.
3. Perform your manual `Path -> Union` or other adjustments.
4. Run `just all`. The script will prioritize the `*-export.svg` file if it exists.

### Export Process Details

1. **Automated Flattening:** The `scripts/export_icons.py` script uses the Inkscape CLI to perform `Stroke to Path` and `Union` operations automatically on every page before export.

2. **Optimization:** SVGO is used with a custom configuration (`scripts/svgo.config.mjs`) to:
    * Minify the SVG data and remove Inkscape metadata.
    * Convert all colors to `currentColor`.
    * Ensure a consistent `viewBox="0 0 48 48"`.

## License

MIT
