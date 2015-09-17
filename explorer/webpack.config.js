var ExtractTextPlugin = require('extract-text-webpack-plugin');
var HtmlWebpackPlugin = require('html-webpack-plugin');
var webpack = require('webpack');
var path = require('path');
var WebpackDevServer = require('webpack-dev-server');

process.env.NODE_ENV = process.env.NODE_ENV || 'development';

console.log("NODE_ENV: " + process.env.NODE_ENV);

var config = {
    context: path.join(__dirname, 'src'),
    entry: {
        app: ["./index.jsx", 'webpack/hot/only-dev-server'],
        client: 'webpack-dev-server/client?http://0.0.0.0:3000' // WebpackDevServer host and port
    },
    output: {
        path: path.join(__dirname, 'dist'),
        filename: "js/[name]-[hash:8].js",
        publicPath: "/" // require to make the specific fonts work (absolute path)
    },
    resolve: {
        fallback: [
            path.join(__dirname, 'src')
        ],
        modulesDirectories: ['node_modules']
    },
    module: {
        loaders: [{
            test: /\.jsx?$/,
            exclude: /(node_modules)/,
            include: path.join(__dirname, 'src'),
            loaders: process.env.NODE_ENV == 'production' ? ['babel'] : ['react-hot', 'babel']
        }, {
            test: /\.scss$/,
            // style is require for inline loading
            // resolve-url allow to resolve url() please keep it mind to set the output.publicPath
            //   to use absolute path. This is usefull is you have fonts and css folders
            loader: 'style!css!resolve-url!sass'
        }, {
            test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/,
            loader: "file?name=fonts/[name]-[hash:8].[ext]"
        },  {
            test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/,
            loader: "file?name=fonts/[name]-[hash:8].[ext]"
        }, {
            test: /\.png$/,
            loader: "url-loader?limit=100000"
        }, {
            test: /\.jpg$/,
            loader: "file-loader"
        }]
    },
    plugins: [
        new webpack.BannerPlugin("This code is part of the GoNode Explorer file", {
            raw: false,
            entryOnly: false
        }),
        new ExtractTextPlugin("css/[name]-[id]-[contenthash:8].css", {
            allChunks: true
        }),
        new HtmlWebpackPlugin({
            inject: true,
            title: 'GoNode Explorer',
            template: 'src/index.html',
            filename: 'index.html'
        }),
        new webpack.optimize.CommonsChunkPlugin({
            name: "commons",
            filename: "js/[name]-[hash:8].js",
            chunks: [] // entries to use
        }),
        new webpack.DefinePlugin({
            VERSION: JSON.stringify("0.0.1-DEV"),
            'process.env': {
                'NODE_ENV': JSON.stringify(process.env.NODE_ENV)
            }
        })
    ]
};

if ('development' === process.env.NODE_ENV) {
    config.devtool = 'eval';
    config.plugins.push(new webpack.HotModuleReplacementPlugin());
}

if ('production' === process.env.NODE_ENV) {
    config.plugins.push(new webpack.optimize.UglifyJsPlugin({
        compress: {

        },
        mangle: { // default values
            except: ['$super', '$', 'exports', 'require']
        }
    }));
}

module.exports = config;

if (require.main === module) { // if call from node directly we start the dev webserver
    var server = new WebpackDevServer(webpack(config), {
        publicPath: config.output.publicPath,
        hot: true,
        historyApiFallback: true,
        stats: {
            colors: true
        },
        progress: true
    });

    server.listen(3000, 'localhost', function (err, result) {
        if (err) {
            console.log(err);
        }

        console.log('Listening at localhost:3000');
    });
}