# Cluster API

This is Giant Swarm's fork. See the upstream [cluster-api README](https://github.com/kubernetes-sigs/cluster-api/blob/main/README.md) for official documentation.

## How to work with this repo

Currently, we try to follow the upstream `release-X.Y` branch to always get the latest stable release and fixes, but not untested commits from `main`. Our only differences against upstream should be in this `README.md` and `.circleci/`. Other changes should be opened as PR for the upstream project first.

We release cluster-api versions with [cluster-api-app](https://github.com/giantswarm/cluster-api-app/). To provide the YAML manifests, we do not use GitHub releases as the upstream project, but package the files into the Docker image (see `.circleci/config.yml`). The app's scripts convert them into the final manifests.

### Repo setup

Since we follow upstream, add their Git repo as remote from which we merge commits:

```sh
git clone git@github.com:giantswarm/cluster-api.git
cd cluster-api
git remote add upstream https://github.com/kubernetes-sigs/cluster-api.git
```

### Pull/test/use latest release from upstream project

We want to use stable upstream release tags unless a hotfix is absolutely required ([decision](https://intranet.giantswarm.io/docs/product/pdr/010_fork_management/)).

First, merge the latest commits from the remote `release-X.Y` branch and push a feature branch:

```sh
git fetch upstream
git checkout release-X.Y
git pull
git checkout -b some-feature-branch
git merge upstream/release-X.Y
git push # new feature branch
```

Now, check the [CircleCI pipeline](https://app.circleci.com/pipelines/github/giantswarm/cluster-api) for the build results. If CircleCI builds fail, fix that first since pushing tags would otherwise lead to failed builds.

To test the changes in the app:

- Replace `.image.tag` in [the app's `values.yaml`](https://github.com/giantswarm/cluster-api-app/blob/master/helm/cluster-api/values.yaml) with your commit SHA (CircleCI pipeline output shows the image tag as well)
- Run `cd cluster-api-app && make generate` to update manifests
- Push a feature branch for the app
- Check the test app release thoroughly on a management cluster

Note that for testing changes to upstream, you probably better use the `main` branch and try your change together with the latest commits from upstream. This also avoids merge conflicts. The latest release branch is usually a bit behind `main`.

### Release

If you have a non-urgent fix, create an upstream PR and wait until it gets released. Then release like this:

- On our `release-*` branch, merge the latest upstream commit and build new tags

  ```sh
  git remote add upstream https://github.com/kubernetes-sigs/cluster-api.git
  git fetch upstream
  git checkout release-X.Y

  # Create a merge commit using upstream's desired release tag (the one we want
  # to upgrade to)
  git merge --no-ff vX.Y.Z

  # Since we want the combined content of our repo and the upstream Git tag,
  # we need to create our own tag on the merge commit
  git tag "vX.Y.Z-gs-$(git rev-parse --short HEAD)"

  git push

  # Push our own, single tag (assuming `origin` is the Giant Swarm fork)
  git push origin "vX.Y.Z-gs-$(git rev-parse --short HEAD)"
  ```

- Check that CircleCI pipeline succeeds for the desired Git tag in order to produce images. If the tag build fails, fix it.
- Replace `.image.tag` in [the app's `values.yaml`](https://github.com/giantswarm/cluster-api-app/blob/master/helm/cluster-api/values.yaml) with the new tag. If the build "just worked" earlier, use tag `vX.Y.Z`. But if you had to make build fixes, consider creating a custom tag `vX.Y.Z-gs-<short commit SHA>` (see hint below).
- Test as described above
- Open PR
- Once merged, manually bump the version in the respective collection to deploy it for one provider (e.g. [capa-app-collection](https://github.com/giantswarm/capa-app-collection/))

In the _rare_ case of an urgent hotfix, we can create an intermediate release tag like this and use it in `.image.tag` instead of a stable upstream release:

```sh
# Find the latest upstream release on which our `release-X.Y` branch is based.
# Get the short commit SHA with `git rev-parse --short HEAD`.
# Then add and push a Git tag to trigger image builds.
git tag vX.Y.Z-gs-SHORT_COMMIT_SHA
git push origin --tags TAG_FROM_ABOVE
```

### Keep fork customizations up to date

Only these files should differ between upstream and our fork:

```sh
git diff main..release-1.FILLME -- .circleci/ README.md
```

If this shows any output, please align the `main` branch with the release branches.
