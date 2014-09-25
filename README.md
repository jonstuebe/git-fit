# git-fit #

## About ##

`git-fit` is a tool for efficiently managing your large repo assets outside of
git. Assets are stored in S3, and tied to commits so you can have different
versions of an asset across different commits.

## How It Works ##

`git-fit` doesn't use any special git techniques, hooks or features. All
metadata about files are stored in `.git-fit.json` in the root of the
repository. This metadata is used to figure out what assets to pull from /
push to S3. The assets themselves are automatically added to `.gitignore` by
`git-fit`, so they're not at all stored on git.

### Why Not Use Smudge/Clean Filters? ###

There are a few tools out there
([git-media](https://github.com/schacon/git-media) being maybe the most
popular) that use smudge/clean filters to handle large assets. This allows the
tools to integrate into a normal git workflow, but there are some consequences
that bit us. Smudge/clean filters have a sometimes unintuitive execution
model, and can easily get your assets in a bad state. Furthermore, tools using
this technique will execute frequently throughout the day - even when you'd
expect it not to (e.g. on `git diff`) - potentially slowing down your daily
workflow.

### Why Not Use git-annex? ###

[git-annex](https://git-annex.branchable.com/) is another popular tool for
large asset management in git - it effectively invents its own git protocol
for managing these assets.

We wanted something simpler, where the execution model was very
straight-forward in order to prevent mistakes, and did not feel confident we
could achieve that with git-annex. But depending on your needs, this may be a
better fit - especially if you're looking for something more than just fast
asset management in git.

## Installation ##

    go get github.com/dailymuse/git-fit
    pushd $GOPATH/src/github.com/dailymuse/git-fit
    make install
    popd
