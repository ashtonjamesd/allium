# Allium

Allium is a simple CLI tool to convert Markdown source files into HTML. This tool is compliant with the [CommonMark](https://commonmark.org/) specification.


## It can:
- Convert CommonMark-compliant Markdown to HTML
- Sound cool on my CV


## Usage

```
go run ./src --convert=[tohtml | tomd] --path=<source> --output=*.[html | md]
```

<br/>

## Example

```
go run ./src --convert=tohtml --path=example.md
```