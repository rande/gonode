const merge = require('webpack-merge');

module.exports = merge(require('./webpack.common.js'), {
    mode: 'development',
    devtool: 'inline-source-map',
    devServer: {
        contentBase: './dist'
    }
});
