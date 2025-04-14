#!/bin/zsh
shopt -s expand_aliases

changeAlias(){
  alias go=go1.23.8
}

while getopts ":yn" flag; do
  case ${flag} in
    y)
      changeAlias
      echo "Alias set: go -> go1.23.8"
      go mod tidy
      ;;
    n)
      unalias go 2>/dev/null
      echo "Alias removed for 'go'"
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      exit 1
      ;;
  esac
done