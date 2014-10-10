#!/usr/bin/env python

import os
import json
import functools

def run(cmd):
    rc = os.system(cmd)

    if rc != 0:
        raise Exception("Unexpected return code when running %s: %s" % (cmd, rc))

def set_contents(path, contents):
    with open(path, "w") as f:
        f.write(contents)

def set_json(path, contents):
    with open(path, "w") as f:
        json.dumps(contents, f)

def ensure_contents(path, expected_contents):
    with open(path, "r") as f:
        actual_contents = f.read()

    if actual_contents != expected_contents:
        raise Exception("Unexpected contents for %s: %s" % (path, actual_contents))

def ensure_gitignore(*expected_files):
    with open(".gitignore", "r") as f:
        actual_files = [line for line in f.read().split("\n") if line]

    if set(actual_files) != set(expected_files):
        raise Exception("Expected .gitignore entries:\n  %s\nBut got:\n  %s" % (", ".join(expected_files), ", ".join(actual_files)))

def ensure_json(path, expected_json):
    with open(path, "r") as f:
        actual_json = json.load(f)

    if actual_json != expected_json:
        raise Exception("Unexpected JSON for %s: %s" % (path, actual_json))

def ensure_dir(path, *files):
    expected_files = set(files)
    actual_files = set(os.listdir(path))

    if expected_files != actual_files:
        raise Exception("Expected files in %s:\n  %s\nBut got:\n  %s" % (path, ", ".join(expected_files), ", ".join(actual_files)))

ensure_spec = functools.partial(ensure_json, "git-fit.json")
ensure_emptyfile = lambda: ensure_contents("emptyfile", "")
ensure_helloworld = lambda contents: ensure_contents("helloworld", contents)
set_helloworld = functools.partial(set_contents, "helloworld")

def main():
    # Util methods
    check_full_gitignore = lambda: ensure_gitignore("/emptyfile", "/helloworld")
    check_full_spec = lambda helloworld_hash: ensure_spec(dict(version=1, files={"emptyfile": "da39a3ee5e6b4b0d3255bfef95601890afd80709", "helloworld": helloworld_hash}))

    print "# Initialize the repo"
    os.mkdir("integration")
    os.chdir("integration")
    run("git init")
    run("git fit init")
    ensure_gitignore()
    ensure_spec(dict(version=1, files={}))

    print "# Pull/push nothing - should be a no-op"
    run("git fit pull")
    run("git fit push")
    ensure_gitignore()
    ensure_spec(dict(version=1, files={}))

    print "# Push the files"
    set_contents("emptyfile", "")
    set_helloworld("hello world 1")
    run("git fit push emptyfile helloworld")
    check_full_gitignore()
    check_full_spec("96e58c52e52b5f3bcb307d3309264d420b60403c")
    ensure_emptyfile()
    ensure_helloworld("hello world 1")

    print "# Pull the existing files - should be no-ops"
    run("git fit pull")
    run("git fit pull emptyfile helloworld")
    check_full_gitignore()
    check_full_spec("96e58c52e52b5f3bcb307d3309264d420b60403c")
    ensure_emptyfile()
    ensure_helloworld("hello world 1")

    print "# Remove/replace the files"
    run("rm emptyfile")
    set_helloworld("hello world 2")

    print "# Pull the files - should replace emptyfile and skip helloworld"
    run("git fit pull")
    check_full_gitignore()
    check_full_spec("96e58c52e52b5f3bcb307d3309264d420b60403c")
    ensure_emptyfile()
    ensure_helloworld("hello world 2")

    print "# Pull helloworld explicitly - should still not replace it"
    run("git fit pull helloworld")
    check_full_gitignore()
    check_full_spec("96e58c52e52b5f3bcb307d3309264d420b60403c")
    ensure_emptyfile()
    ensure_helloworld("hello world 2")

    print "# Push the change"
    run("git fit push")
    check_full_gitignore()
    check_full_spec("42ad4ff8bdd0125e98eeaa23146d7899ee77577e")
    ensure_emptyfile()
    ensure_helloworld("hello world 2")

    print "# Remove/re-pull - should still be the new contents"
    run("rm helloworld")
    run("git fit pull")
    check_full_gitignore()
    check_full_spec("42ad4ff8bdd0125e98eeaa23146d7899ee77577e")
    ensure_emptyfile()
    ensure_helloworld("hello world 2")

    print "# Remove emptyfile from gf"
    run("git fit rm emptyfile")
    check_full_gitignore()
    ensure_spec(dict(version=1, files={"helloworld": "42ad4ff8bdd0125e98eeaa23146d7899ee77577e"}))
    ensure_emptyfile()
    ensure_helloworld("hello world 2")

    print "# Run gc and check before/after contents"
    ensure_dir(".git/fit", "da39a3ee5e6b4b0d3255bfef95601890afd80709", "96e58c52e52b5f3bcb307d3309264d420b60403c", "42ad4ff8bdd0125e98eeaa23146d7899ee77577e")
    run("git fit gc")
    ensure_dir(".git/fit", "42ad4ff8bdd0125e98eeaa23146d7899ee77577e")

if __name__ == "__main__":
    main()
