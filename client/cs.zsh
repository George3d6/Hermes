#! /usr/bin/zsh
for scssFile in scss/*; do
    cat $scssFile >> scss/compresed.scss
done;
sassc scss/compresed.scss > static/bundle.css;
rm -rf scss/compresed.scss;
