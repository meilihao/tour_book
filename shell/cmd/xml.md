# xml

## xmlstarlet
ref:
- [Modify multiple lines of an XML file using command line](https://unix.stackexchange.com/questions/309676/modify-multiple-lines-of-an-xml-file-using-command-line)

支持修改xml

### example
```bash
xmlstarlet sel -t -v "//element/@attribute" file.xml
```

## xmllint
`xmllint --format XML_FILE` 