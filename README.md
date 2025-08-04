# GH Local Changes

GH Local Changes is a GitHub CLI extension that scans a given directory for Git repositories and reports on branches and changes that have not been pushed to the remote repository.

## Installation

To install GH Local Changes, you need to have the [GitHub CLI](https://cli.github.com/) installed. Then, you can install the extension with the following command:

```bash
gh extension install baruchiro/gh-local-changes
```

## Usage

To use GH Local Changes, run the following command:

```
gh local-changes
```

By default, GH Local Changes scans the current directory. You can specify a different directory by passing it as an argument:

```
gh local-changes ~/source
```

To print debug logs , use --debug flag :

```
gh local-changes --debug
gh local-changes --debug  ~/source
```

GH Local Changes will recursively scan the specified directory for Git repositories. For each repository, it will report the branches and changes that have not been pushed to the remote repository.

## Contributing

To develop GH Local Changes, you need to have [Go](https://golang.org) installed. You can then clone the repository and run the program with the following commands:

```bash
git clone https://github.com/baruchiro/gh-local-changes.git
cd gh-local-changes
go build
gh extension install .
```

To continuously check your changes, you can run the following command:

```bash
go build; gh local-changes ~/source
```

## License

GH Local Changes is released under the [MIT License](LICENSE).
