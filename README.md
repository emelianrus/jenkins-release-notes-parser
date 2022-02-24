# jenkins-release-notes-parser

parse jenkins plugins to get release notes based on git diff file (changed line in git diff)

[EXAMPLE_RELEASE_NOTES](./jenkins-plugins.txt_RELEASE_NOTES_EXAMPLE.md)

[EXAMPLE_GIT_DIFF_FILE](./diff.out)

```
git diff master jenkins-icecream-plugins.txt > diff.out

make go-run
```

Create exec:
```
make go-build
```


## NOTES

* some plugin names can be different to URL and location in git and some plugins doesn't have releases at all

* github allows to do 60 api calls for unauthenticated requests.

* do not support other then <name>:<version> pattern and probably will fail if in diff it will find text line change
