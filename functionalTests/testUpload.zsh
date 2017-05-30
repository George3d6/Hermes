#!/usr/bin/zsh
curl \
  -F "name=testUpload" \
  -F "ttl=59" \
  -F "file=@/home/george/Hermes/functionalTests/afile" \
  -F "compression=xz" \
  "localhost:3280/upload/"
