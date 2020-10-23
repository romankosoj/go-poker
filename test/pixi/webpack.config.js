const path = require("path");
const { SourceMapDevToolPlugin } = require("webpack");

module.exports = {
    module: {
        rules: [
            {
                test: /\.js&/,
                enforce: "pre",
                use: ["source-map-loader"],
            },
        ],
    },
    plugins: [
        new SourceMapDevToolPlugin({
            filename: "[file].map"
        }),
    ],
    output: {
        path: path.join(__dirname, "dist"),
        filename: "main.js",
        sourceMapFilename: "main.js.map",
    },
    devtool: "source-map",
    devServer: {
        contentBase: path.join(__dirname, "dist"),
        compress: true,
        port: 5000,
    },
}