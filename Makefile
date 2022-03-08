run:
	@go run *.go

sage-assets:
	npm install
	rm -R assets/static/sage-assets/fonts
	rm -R assets/static/sage-assets/main.css
	cp -R node_modules/@kajabi/sage-assets/dist/fonts assets/static/sage-assets
	cp -R node_modules/@kajabi/sage-assets/dist/main.css assets/static/sage-assets

sage-system:
	npm install
	# rm assets/static/sage-react/main.js
	# cp -R node_modules/@kajabi/sage-react/dist/main.js assets/static/sage-assets
	npx webpack build ./node_modules/@kajabi/sage-react/dist/main.js -o assets/static/sage-assets
