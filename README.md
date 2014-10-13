# git-fit #

## About ##

`git-fit` is a tool for efficiently managing your large repo assets outside of
git. Assets are stored in S3. Metadata is stored directly in git, so you can
have different versions of an asset across different commits.

## How It Works ##

`git-fit` doesn't use any special git techniques, hooks or features. All
metadata about files are stored in `git-fit.json` in the root of the
repository. This metadata is used to figure out what assets to pull from /
push to S3. The assets themselves are automatically added to `.gitignore` by
`git-fit`, so they're not at all stored on git.

`git-fit` maintains a cache of assets at `.git/fit`. Files are stored here and
in S3 by the SHA1 hash of their contents. This makes everything straight-
forward, and allows you to efficiently store multiple copes of the same asset
in different paths.

#### Why Not Use Smudge/Clean Filters? ####

There are a few tools out there
([git-media](https://github.com/schacon/git-media) being maybe the most
popular) that use smudge/clean filters to handle large assets. This allows the
tools to integrate into a normal git workflow, but there are some consequences
that bit us. Smudge/clean filters have a sometimes unintuitive execution
model, and can easily get your assets in a bad state. And tools using
this technique will execute frequently throughout the day - even when you'd
expect it not to (e.g. on `git diff`) - potentially slowing down your daily
workflow.

#### Why Not Use git-annex? ####

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
    pushd $GOROOT/src/github.com/dailymuse/git-fit
    make install
    popd

## Usage ##

Before using `git-fit` for the first time in a repo, run `git fit init`. This
will setup the repository to be able to use `git-fit` by adding configs and
creating a directory for storing cached assets.

### Pulling ###

Pulling will look at `git-fit.json` to see what versions of assets to pull.
If there's a cached copy, its contents will simply be copied to the asset's
location. If not, a copy will be fetched from S3 and cached.

You can do a partial pull by explicitly passing arguments; otherwise all
managed assets will be pulled. Pull will not overwrite existing files - this
is to prevent you from accidentally overwriting local changes that are
unsynced. To overwrite, remove the local copy first.

### Pushing ###

Pushing will hash the contents of the assets, store them in `git-fit.json`,
and push them off to S3 and the local cache if they aren't already stored.

As with pull, you can explicitly pass paths as arguments to push only certain
files. Otherwise all updated files will be pushed.

### GC ###

Every once in a while, it's a good idea to run `git fit gc` on a repo. This
will delete any cached assets that are not currently specified in
`git-fit.json`. Note that while this will free up space by removing old
versions of assets, it will also clear your cache, so future pulls may be
slow until the cache is warmed up again.
