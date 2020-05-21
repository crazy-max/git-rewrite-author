<p align="center"><a href="https://github.com/crazy-max/git-rewrite-author" target="_blank"><img width="100" src="https://raw.githubusercontent.com/crazy-max/git-rewrite-author/master/.github/git-rewrite-author.png"></a></p>

<p align="center">
  <a href="https://github.com/crazy-max/git-rewrite-author/releases/latest"><img src="https://img.shields.io/github/release/crazy-max/git-rewrite-author.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/crazy-max/git-rewrite-author/releases/latest"><img src="https://img.shields.io/github/downloads/crazy-max/git-rewrite-author/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://github.com/crazy-max/git-rewrite-author/actions?workflow=build"><img src="https://img.shields.io/github/workflow/status/crazy-max/git-rewrite-author/build?label=build&logo=github&style=flat-square" alt="Build Status"></a>
  <a href="https://goreportcard.com/report/github.com/crazy-max/git-rewrite-author"><img src="https://goreportcard.com/badge/github.com/crazy-max/git-rewrite-author?style=flat-square" alt="Go Report"></a>
  <a href="https://www.codacy.com/app/crazy-max/git-rewrite-author"><img src="https://img.shields.io/codacy/grade/356b78c4f48e4e2e9d286dd79be84d3f.svg?style=flat-square" alt="Code Quality"></a>
  <br /><a href="https://github.com/sponsors/crazy-max"><img src="https://img.shields.io/badge/sponsor-crazy--max-181717.svg?logo=github&style=flat-square" alt="Become a sponsor"></a>
  <a href="https://www.paypal.me/crazyws"><img src="https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal&style=flat-square" alt="Donate Paypal"></a>
</p>

## About

**git-rewrite-author** is a CLI application written in [Go](https://golang.org/) to rewrite one or several authors / committers history of a [Git](https://git-scm.com/) repository with ease. It was inspired by [this post on Github](https://help.github.com/articles/changing-author-info/).

___

* [Requirements](#requirements)
* [Download](#download)
* [Installation](#installation)
* [Usage](#usage)
* [How can I help?](#how-can-i-help)
* [License](#license)

## Requirements

You must have [Git](https://git-scm.com/) installed on your system and create a fresh, bare clone of your repository :

```
$ cd /tmp
$ git clone --bare https://github.com/user/repo.git
$ cd /tmp/repo.git
```

## Download

You can download the application matching your platform on the [**releases page**](https://github.com/crazy-max/git-rewrite-author/releases/latest).

## Installation

Place the executable in your Git repository. It is best to place it in your `PATH` so that you can use it anywhere in your system and also use it with the Git syntax `git rewrite-author`.

## Usage

```
Usage: git-rewrite-author <command>

Rewrite authors history of a Git repository with ease. More info:
https://github.com/crazy-max/git-rewrite-author

Flags:
  --help                Show context-sensitive help.
  --version
  --repo="."            Git repository path.
  --log-level="info"    Set log level.
  --log-caller          Add file:line of the caller to log output.

Commands:
  config-get
  config-set
  list
  rewrite
  rewrite-list

Run "git-rewrite-author <command> --help" for more information on a command.
```

You probably want to know the list of authors/committers for a repository before rewritting history:

```
$ git-rewrite-author list --repo /tmp/repo.git
ohcrap <ohcrap@bad.com>
GitHub <noreply@github.com>
root <root@localhost>
```

Then you can rewrite a single author/committer:

```
$ git-rewrite-author rewrite "ohcrap@bad.com" "John Smith <john.smith@domain.com>" --repo /tmp/repo.git

Following authors/committers will be rewritten:
- ohcrap@bad.com => John Smith <john.smith@domain.com

Rewrite 4b03c46d8f085f56014e5bee1e5597de86554139 (31/31) (22 seconds passed, remaining 0 predicted)
Ref 'refs/heads/master' was rewritten
Ref 'refs/remotes/origin/master' was rewritten
Ref 'refs/tags/0.15.1-1' was rewritten
Ref 'refs/tags/0.15.2-2' was rewritten
Ref 'refs/tags/0.15.310-3' was rewritten
Ref 'refs/tags/0.16.9-4' was rewritten
Ref 'refs/tags/0.17.13-5' was rewritten
Ref 'refs/tags/0.17.19-6' was rewritten
Ref 'refs/tags/0.18.14-7' was rewritten
Ref 'refs/tags/0.18.23-8' was rewritten
Ref 'refs/tags/0.18.23-9' was rewritten
Ref 'refs/tags/0.18.36-10' was rewritten
Ref 'refs/tags/0.19.48-11' was rewritten
Ref 'refs/tags/0.19.70-12' was rewritten
```

Or a list of authors/committers:

```
$ git-rewrite-author rewrite-list ../authors.json --repo /tmp/repo.git

Following authors/committers will be rewritten:
- root@localhost, noreply@github.com => John Smith <john.smith@domain.com>
- ohcrap@bad.com => Good Sir <goodsir@users.noreply.github.com>

Rewrite 4b03c46d8f085f56014e5bee1e5597de86554139 (31/31) (22 seconds passed, remaining 0 predicted)
Ref 'refs/heads/master' was rewritten
Ref 'refs/remotes/origin/master' was rewritten
Ref 'refs/tags/0.15.1-1' was rewritten
Ref 'refs/tags/0.15.2-2' was rewritten
Ref 'refs/tags/0.15.310-3' was rewritten
Ref 'refs/tags/0.16.9-4' was rewritten
Ref 'refs/tags/0.17.13-5' was rewritten
Ref 'refs/tags/0.17.19-6' was rewritten
Ref 'refs/tags/0.18.14-7' was rewritten
Ref 'refs/tags/0.18.23-8' was rewritten
Ref 'refs/tags/0.18.23-9' was rewritten
Ref 'refs/tags/0.18.36-10' was rewritten
Ref 'refs/tags/0.19.48-11' was rewritten
Ref 'refs/tags/0.19.70-12' was rewritten
```

Here the `authors.json` JSON file looks like this:

```json
[
    {
        "old": [ "root@localhost", "noreply@github.com" ],
        "correct_name": "John Smith",
        "correct_mail": "john.smith@domain.com"
    },
    {
        "old": [ "ohcrap@bad.com" ],
        "correct_name": "Good Sir",
        "correct_mail": "goodsir@users.noreply.github.com"
    }
]
```

Now confirm everything suits to you:

```
$ git-rewrite-author list --repo /tmp/repo.git
Good Sir <goodsir@users.noreply.github.com>
John Smith <john.smith@domain.com>
```

Review the new Git history for errors and push the corrected history to Git:

```
git push --force --all
```

## How can I help?

All kinds of contributions are welcome :raised_hands:! The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon: You can also support this project by [**becoming a sponsor on GitHub**](https://github.com/sponsors/crazy-max) :clap: or by making a [Paypal donation](https://www.paypal.me/crazyws) to ensure this journey continues indefinitely! :rocket:

Thanks again for your support, it is much appreciated! :pray:

## License

MIT. See `LICENSE` for more details.<br />
Icon credit to [ual Pharm](https://www.shareicon.net/author/ual-pharm).
