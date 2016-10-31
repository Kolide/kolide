require('es6-promise').polyfill();

var path = require('path');
var webpack = require('webpack');
var autoprefixer = require('autoprefixer');
var ExtractTextPlugin = require("extract-text-webpack-plugin");
var bourbon = require('node-bourbon').includePaths;

var plugins = [
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
    loaders: [
      {test: /\.(png|gif)$/, loader: 'url-loader?name=[name]@[hash].[ext]&limit=6000'},
      {test: /\.(pdf|ico|jpg|svg|eot|otf|woff|ttf|mp4|webm)$/, loader: 'file-loader?name=[name]@[hash].[ext]'},
      {test: /\.json$/, loader: 'raw-loader'},
      {test: /\.tsx?$/, loader: 'ts-loader'},
      {
        test: /\.css$/,
        loader: ExtractTextPlugin.extract("style-loader", "css-loader!autoprefixer-loader")
      },
      { test: /\.scss$/, loader: "style!css!sass?includePaths[]=" + bourbon + "!import-glob" },
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
    }
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
