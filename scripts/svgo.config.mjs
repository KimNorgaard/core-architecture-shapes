export default {
  multipass: true,
  plugins: [
    {
      name: 'preset-default',
      params: {
        overrides: {
          cleanupIds: true,
          convertShapeToPath: true,
          mergePaths: true,
        },
      },
    },
    {
      name: 'removeViewBox',
      active: false,
    },
    'convertStyleToAttrs',
    {
      name: 'convertColors',
      params: {
        currentColor: true,
      },
    },
    {
      name: 'removeAttrs',
      params: {
        attrs: '(stroke-opacity|fill-opacity|paint-order|id|class-name|stop-color|stop-opacity|overflow)',
      },
    },
    {
      name: 'addAttributesToSVGElement',
      params: {
        attributes: [
          { fill: 'currentColor' },
        ],
      },
    },
  ],
};
