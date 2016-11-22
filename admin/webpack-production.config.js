const webpack               = require('webpack');
const path                  = require('path');
const buildPath             = path.resolve(__dirname, 'build');
const nodeModulesPath       = path.resolve(__dirname, 'node_modules');
const TransferWebpackPlugin = require('transfer-webpack-plugin');

const config = {
    entry:   ['whatwg-fetch', path.join(__dirname, '/src/app.js')],
    // Render source-map file for final build
    // devtool: 'source-map',
    // output config
    output:  {
        path:     buildPath, // Path of output file
        filename: 'app.js', // Name of output file
    },
    plugins: [
        // Minify the bundle
        new webpack.optimize.UglifyJsPlugin({
            compress: {
                // supresses warnings, usually from module minification
                warnings: false,
            },
        }),
        new webpack.DefinePlugin({
            'process.env': {
                'NODE_ENV': JSON.stringify('production')
            }
        }),
        // Allows error warnings but does not stop compiling.
        new webpack.NoErrorsPlugin(),
        // Transfer Files
        new TransferWebpackPlugin([
            {from: 'static'},
        ], path.resolve(__dirname, 'src')),
    ],
    module:  {
        loaders: [
            {
                test:    /\.js$/, // All .js files
                loaders: ['babel-loader'], // react-hot is like browser sync and babel loads jsx and es6-7
                exclude: [nodeModulesPath],
            },
            {
                test: /\.json$/,
                loader: 'json'
            }
        ],
    },
};

module.exports = config;
