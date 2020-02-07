module.exports = {
  mode: "production",
  module: {
    rules: [
      {
        test: /\.ts(x)?$/,
        use: 'ts-loader',
        exclude: /node_modules/,
      },
      {
        test: /\.css$/,
        use: [
          'style-loader',
          'css-loader'
        ]
      }
    ],
  },
  resolve: {
    extensions: [ '.ts', '.js', '.tsx' ],
  },
  entry: "./src/index.tsx",
  externals: {
    "react": "React",
    "react-dom": "ReactDOM"
  }
};
