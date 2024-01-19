# dtsync
`dtsync` is a flexible file synchronization tool designed to keep your directories in sync. 
It provides a simple command-line interface that allows you to specify source and destination directories, and offers options to remove or replace files. Whether you're managing backups, mirroring directories, or simply moving files around, dtsync can support you.

## Usage
Currently supported arguments can be listed as follow.
```bash
$ ./dtsync -help
Usage of dtsync:
  -src string
        The source root path (required)
  -dst string
        The destination root path (required)
  -remove
        Remove files and directories in dst not included in src
  -replace
        Replace file on dst when different
```

### Default Case
```bash
$ ./dtsync -src /a -dst /b
```
<img alter="Default Sync" src=".media/default_sync.png" width="350">

### Case With Replace
```bash
$ ./dtsync -src /a -dst /b -replace
```
<img alter="Replace Sync" src=".media/replace_only_sync.png" width="350">

### Case With Remove
```bash
$ ./dtsync -src /a -dst /b -remove
```
<img alter="Remove Sync" src=".media/remove_only_sync.png" width="350">

### Case With Replace And Remove
```bash
$ ./dtsync -src /a -dst /b -replace -remove
```
<img alter="Replace And Remove Sync" src=".media/full_sync.png" width="350">

### Disclaimer
`dtsync` is provided "as is", without warranty of any kind. 
The authors or copyright holders will not be liable for any damage, data loss, or any other issue that may occur as a result of using this tool. 
Use at your own risk.