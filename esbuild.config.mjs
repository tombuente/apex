import * as esbuild from 'esbuild'
import cssModulesPlugin from 'esbuild-css-modules-plugin';

await esbuild.build({
	entryPoints: ['bundle/js/app.js'],
	outfile: 'static/js/app.js',
	bundle: true,
	minify: true,
})

await esbuild.build({
	entryPoints: ['bundle/css/app.css'],
	outfile: 'static/css/app.css',
	bundle: true,
	plugins: [cssModulesPlugin()],
	minify: true,
	loader: { '.css': 'css'},
});
