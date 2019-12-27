cd DivProxyReact && npm run-script build && cd ..
cp -r DivProxyReact/build resources/
sed -i 's/^/<script>require("popper.js");<\/script><script>require("bootstrap");<\/script><script>window\.\$ = window\.jQuery = require("jquery");<\/script><script src="\.\.\/app\/js\/index\.js"><\/script><script src="\.\.\/app\/js\/connectBackend\.js"><\/script>/g' resources/build/index.html
astilectron-bundler -v
cp -r resources/ output/linux-amd64/