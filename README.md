# Simple TAR creator/extractor

## Install
```
go get -u github.com/Warashi/star
```
or get binary from [releases](https://github.com/Warashi/star/releases)

## Usage
### create archive
```
star -c [target directory] > output.tar
```
If you want to compress archive such as gzip, you can pipe stdout to gzip.
i.e. `star -c target | gzip > output.tar.gz`

### extract archive
```
star -x [destination] < input.tar
```
When you don't specify `destination`, current directory is used.
If you want to compressed archive such as tar.gz, you must decompress with another tool.
i.e. `gzip -dc input.tar.gz | star -x dest`

## Special Thanks
Code mainly taken from [here](https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07).
Thanks!!
