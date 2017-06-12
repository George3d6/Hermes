const webpack = require('webpack');
const path = require('path');

const babelOptions = {
    "presets": "es2017",
    plugins: [
  ['transform-runtime', {
    helpers: false,
    polyfill: false,
    regenerator: true, }],
  'transform-es2015-destructuring',
  'transform-async-to-generator',
  ],
}

module.exports = {
    entry: ['./ts/main'],
    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'static')
    },
    module: {
        rules: [
            {
                test: /\.ts(x?)$/,
                exclude: /node_modules/,
                use: [
                    {
                        loader: 'babel-loader',
                        options: babelOptions
                    }, {
                        loader: 'ts-loader'
                    }
                ]
            }, {
                test: /\.js$/,
                exclude: /node_modules/,
                use: [
                    {
                        loader: 'babel-loader',
                        options: babelOptions
                    }
                ]
            }
        ]
    },
    resolve: {
        extensions: ['.webpack.js', '.web.js', '.ts', '.tsx', '.js']
    },
    plugins: []
}

//transform-runtime
