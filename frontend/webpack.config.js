const fs = require('fs');

function writeApiURL(url) {
  fs.writeFileSync("./src/api/url.ts", `
export const API_URL = "${url}";
  `);
}

module.exports = env => {

  let isDev = env && env.dev;
  if (isDev) {
    writeApiURL(`http://localhost:1080`);
  } else {
    writeApiURL(`https://sbvision.kwylder.com`);
  }

  return {
    mode: (isDev) ? "development" : "production",
    watch: isDev,
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
};
