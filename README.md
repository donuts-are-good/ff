
![donuts-are-good's followers](https://img.shields.io/github/followers/donuts-are-good?&color=555&style=for-the-badge&label=followers) ![donuts-are-good's stars](https://img.shields.io/github/stars/donuts-are-good?affiliations=OWNER%2CCOLLABORATOR&color=555&style=for-the-badge) ![donuts-are-good's visitors](https://komarev.com/ghpvc/?username=donuts-are-good&color=555555&style=for-the-badge&label=visitors)

# ff
Welcome to `ff`, a fast file finder written in Go. This tool allows you to search for files by name or by contents.

## 📖 Usage

To use `ff`, you need to specify either a name or contents keyword to search for. Optionally, you can also specify a path to search in. If no path is specified, `ff` will search from the root directory.

Here's the basic syntax:

```bash
ff --name "keyword" [path]
ff --contents "keyword" [path]
```

You can also include hidden files in your search with the `--hidden` flag:

```bash
ff --name "keyword" --hidden [path]
ff --contents "keyword" --hidden [path]
```

## 📝 Examples

Search for files with "test" in the name:

```bash
ff --name "test"
```

Search for files with "hello world" in the contents:

```bash
ff --contents "hello world"
```

Search for hidden files with "config" in the name in the home directory:

```bash
ff --name "config" --hidden ~/
```

