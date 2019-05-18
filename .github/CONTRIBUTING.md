# Contribution Guidelines

Thanks! Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

## Developer environment

Some development tools are needed to be ready.

```sh
git clone https://github.com/jesusprubio/online
cd online
cargo cmd deps
```

## Tests

We use [Clippy](https://github.com/rust-lang/rust-clippy) and [rustfmt](https://github.com/rust-lang/rustfmt) linters. Please run to be sure your code fits with them and the tests keep passing:

```sh
cargo cmd test
```

## Commit messages rules:
- It should be formed by a one-line subject, followed by one line of white space. Followed by one or more descriptive paragraphs, each separated by one￼￼￼￼ line of white space. All of them finished by a dot.
- If it fixes an issue, it should include a reference to the issue ID in the first line of the commit.
- It should provide enough information for a reviewer to understand the changes and their relation to the rest of the code.
