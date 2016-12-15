require('es6-promise').polyfill();

var path = require('path');
var webpack = require('webpack');
var autoprefixer = require('autoprefixer');
var ExtractTextPlugin = require("extract-text-webpack-plugin");
var bourbon = require('node-bourbon').includePaths;
var WebpackBuildNotifierPlugin = require('webpack-build-notifier');
var BitBarWebpackProgressPlugin = require("bitbar-webpack-progress-plugin");

var plugins = [
  new BitBarWebpackProgressPlugin(),
  new webpack.NoErrorsPlugin(),
  new webpack.optimize.DedupePlugin(),
  new ExtractTextPlugin("bundle.css", {allChunks: false})
];

if (process.env.NODE_ENV === 'production') {
  plugins = plugins.concat([
    new webpack.optimize.UglifyJsPlugin({
      output: {comments: false},
      test: /bundle\.js?$/
    }),
    new webpack.DefinePlugin({
      'process.env': {NODE_ENV: JSON.stringify('production')}
    })
  ]);
};

if (process.argv.indexOf('--notify') > -1) {
  plugins = plugins.concat([
    new WebpackBuildNotifierPlugin({
      title: "Kolide",
      logo: path.resolve("./assets/images/kolide-logo.svg"),
      suppressWarning: true,
      suppressSuccess: true,
      sound: false
    })
  ])
};

var repo = __dirname

var config  = {
  entry: {
    bundle: path.join(repo, 'frontend/index.jsx')
  },
  output: {
    path: path.join(repo, 'assets/'),
    publicPath: "/assets/",
    filename: '[name].js'
  },
  plugins: plugins,
  module: {
    // The following noParse suppresses the warning about sqlite-parser being a
    // pre-compiled JS file. See https://goo.gl/N4s6bB.
    noParse: /node_modules\/sqlite-parser\/dist\/sqlite-parser-min.js/,
    loaders: [
      {test: /\.(png|gif)$/, loader: 'url-loader?name=[name]@[hash].[ext]&limit=6000'},
      {test: /\.(pdf|ico|jpg|svg|eot|otf|woff|ttf|mp4|webm)$/, loader: 'file-loader?name=[name]@[hash].[ext]'},
      {test: /\.json$/, loader: 'raw-loader'},
      {test: /\.tsx?$/, exclude: /node_modules/, loader: 'ts-loader'},
      {
        test: /\.scss$/,
        exclude: /node_modules/,
        loader: ExtractTextPlugin.extract('style-loader', 'css!sass?sourceMap=true&includePaths[]=' + bourbon + '!import-glob')
      },
      {
        test: /\.css$/,
        loader: ExtractTextPlugin.extract('style-loader', 'css-loader!autoprefixer-loader')
      },
      {
        test: /\.jsx?$/,
        include: path.join(repo, 'frontend'),
        loaders: ['babel']
      }
    ]
  },
  resolve: {
    extensions: ['', '.js', '.jsx', '.json'],
    alias: {
      '#app': path.join(repo, 'frontend'),
      '#components': path.join(repo, 'frontend/components'),
    },
    root: [
      path.resolve(path.join(repo, './frontend'))
    ]
  },
  svgo1: {
    multipass: true,
    plugins: [
      // by default enabled
      {mergePaths: false},
      {convertTransform: false},
      {convertShapeToPath: false},
      {cleanupIDs: false},
      {collapseGroups: false},
      {transformsWithOnePath: false},
      {cleanupNumericValues: false},
      {convertPathData: false},
      {moveGroupAttrsToElems: false},
      // by default disabled
      {removeTitle: true},
      {removeDesc: true}
    ]
  }
};

module.exports = config;
