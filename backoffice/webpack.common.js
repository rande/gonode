const path = require('path');

const HtmlWebpackPlugin = require('html-webpack-plugin');
const CleanWebpackPlugin = require('clean-webpack-plugin');
const { CheckerPlugin } = require('awesome-typescript-loader');


module.exports = {
    entry: './src/app.tsx',
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: 'awesome-typescript-loader',
                exclude: /(node_modules|(.*)test\.tsx?)/,
                include: path.resolve(__dirname, "src"),
            }
        ]
    },
    resolve: {
        extensions: [ '.tsx', '.ts', '.js' ],
        modules: [
            path.resolve(__dirname, "src"),
            path.resolve(__dirname, "src/fugue/packages"),
            path.resolve(__dirname, "node_modules")
        ]
    },
    output: {
        filename: 'counter.js',
        path: path.resolve(__dirname, 'dist'),
    },
    plugins: [
        new CheckerPlugin(),
        new CleanWebpackPlugin(['dist']),
        new HtmlWebpackPlugin({
            title: 'Admin Demo',
            template: './src/index.html'
        })
    ]
};
