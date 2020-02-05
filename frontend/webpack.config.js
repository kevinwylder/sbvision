module.exports = {
  mode: "production",
  module: {
    rules: [
      {
        test: /\.ts(x)?$/,
        use: 'ts-loader',
        exclude: /node_modules/,
      },
    ],
  },
  resolve: {
    extensions: [ '.ts', '.js', '.tsx' ],
  },
  entry: "./src/index.tsx"
};
