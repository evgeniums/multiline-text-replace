# multiline-text-replace
Simple Go tool for replacing text with multiple lines in files.

For usage description see help:

```
multiline-text-replace --help

Usage:
  multiline-text-replace.exe [OPTIONS]

Application Options:
  /p, /pattern-file:       File with pattern text to search.
  /s, /substitution-file:  File with substitution text to replace with.
  /f, /file:               Target file to replace text in. Either file or dir must be specified.
  /d, /dir:                Target directory containing files to replace text in. Either file or dir must be specified.
  /e, /ext:                Extensions of files in target directory to replace text in, can be multiple separated with
                           comma, e.g. .txt,.csv
  /r, /recursive           Process files recursively in the target directory.

Help Options:
  /?                       Show this help message
  /h, /help                Show this help message
```
