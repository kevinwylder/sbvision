const path = require('path');

module.exports = {
  mode: "production",
  entry: {
  },
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
  output: {
    filename: '[name].js',
    path: path.resolve(__dirname, 'public/bundle'),
  },
  externals: {
    "react": "React",
    "react-dom": "ReactDOM"
  }
};
