#!/usr/bin/env bash

sqlite3 bd.sqlite "SELECT recipe_title FROM stats;" | while IFS= read -r b64; do
  dec=$(echo "$b64" | base64 -d 2>/dev/null)
  printf 'ENC: %s - DEC: %s\n' "$b64" "$dec"
done