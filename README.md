# git-fit #

## About ##

`git-fit` is a tool for efficiently managing your large repo assets outside of
git. Assets are stored in S3, and tied to commits so you can have different
versions of an asset across different commits.

## How It Works ##

`git-fit` stores metadata about managed assets in `.git-config.json` in the
root of the repository. The assets themselves are gitignored, and `git-fit`
provides capabilities for pulling and pushing these assets to/from S3.

## Installation ##

    go get github.com/dailymuse/git-fit
    pushd $GOPATH/src/github.com/dailymuse/git-fit
    make install
    popd

## Alternatives ##

1. [git-media](https://github.com/schacon/git-media): Uses smudge/clean filters,
   which has has a sometimes unintuitive execution model and can easily get
   your assets in a bad state. You should expect significantly worse
   performance from git-media as well, as since it uses smudge/clean filters,
   it will execute frequently throughout the day - even when you'd expect it
   not to - e.g. on `git diff`.
2. [git-fat](https://github.com/jedbrown/git-fat): Uses smudge/clean filters
   as well. Better maintained than git-media, but only supports rsync.
3. [git-annex](https://git-annex.branchable.com/): By far the most flexible
   tool, but not very intuitive. Depending your needs this may be a better
   fit - especially if you're looking for something more than just fat asset
   management in S3.
