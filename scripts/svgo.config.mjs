export default {
  multipass: true,
  plugins: [
    {
      name: "preset-default",
      params: {
        overrides: {
          removeViewBox: false,
        },
      },
    },
    "convertStyleToAttrs",
    {
      name: "convertColors",
      params: {
        currentColor: true,
      },
    },
    {
      name: "removeAttrs",
      params: {
        attrs: "(stroke-opacity|fill-opacity|paint-order)",
      },
    },
    {
      name: "addAttributesToSVGElement",
      params: {
        attributes: [{ fill: "currentColor" }],
      },
    },
  ],
};
