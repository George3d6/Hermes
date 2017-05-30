let debug = process.env.NODE_ENV;
let webpack = require('webpack');
let path = require('path');

module.exports = {
    devServer: {
        port: 3000,
        historyApiFallback: true
    },
    context: __dirname,
    devtool: debug
        ? 'inline sourcemap'
        : null,
    entry: ['./js/main.js'],
    output: {
        path: __dirname,
        filename: 'static/bundle.js'
    },
    module: {
        loaders: [
            {
                loader: "babel-loader",
                exclude: /node_modules/,
                // Skip any files outside of your project's `src` directory
                include: [path.resolve(__dirname, "")],

                // Only run `.js` and `.jsx` files through Babel
                test: /\.jsx?$/,

                // Options to configure babel with
                query: {
                    plugins: [
                        'transform-runtime', "syntax-async-functions", "transform-regenerator"
                    ],
                    presets: ['es2015']
                }
            }
        ]
    },
    plugins: [
        new webpack.optimize.DedupePlugin(),
        new webpack.optimize.OccurenceOrderPlugin(),
        //new webpack.optimize.UglifyJsPlugin({mangle: false, sourcemap: false}),
        new webpack.DefinePlugin({
            'process.env': {
                NODE_ENV: JSON.stringify('production')
            }
        })
    ]
};
