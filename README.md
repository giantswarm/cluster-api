# Cluster API

This is Giant Swarm's fork. See the upstream [cluster-api README](https://github.com/kubernetes-sigs/cluster-api/blob/main/README.md) for official documentation.

## How to work with this repo

Currently, we try to follow the upstream `main` branch to always get the latest fixes. Our only differences against upstream should be in this README and `.circleci/`.

We release cluster-api versions with [cluster-api-app](https://github.com/giantswarm/cluster-api-app/). To provide the YAML manifests, we do not use GitHub releases as the upstream project, but package the files into the Docker image (see `.circleci/config.yml`). The app's scripts convert them into the final manifests.

### Repo setup

Since we follow upstream, add their Git repo as remote from which we merge commits:

```sh
git clone git@github.com:giantswarm/cluster-api.git
cd cluster-api
git remote add upstream https://github.com/kubernetes-sigs/cluster-api.git
```

### Pull and test latest changes from upstream project

First, merge the latest commits from the remote `main` branch and push a feature branch:

```sh
git fetch upstream
git checkout main
git pull
git checkout -b some-feature-branch
git merge upstream/main
git push
```

Now, check the [CircleCI pipeline](https://app.circleci.com/pipelines/github/giantswarm/cluster-api) for the build results.

To test the changes in the app:

- Replace `.image.tag` in [the app's `values.yaml`](https://github.com/giantswarm/cluster-api-app/blob/master/helm/cluster-api/values.yaml) with your commit SHA
- Run `cd cluster-api-app && make generate`
- Push a feature branch
- Check the test app release thoroughly on a management cluster

### Release

We want to use stable upstream release tags unless a hotfix is required ([decision](https://intranet.giantswarm.io/docs/product/pdr/010_fork_management/)).

So if you have a non-urgent fix, create an upstream PR and wait until it gets released. Then release like this:

- On our `release-*` branch, merge the latest upstream commit and build new tags

  ```sh
  git remote add upstream https://github.com/kubernetes-sigs/cluster-api.git
  git fetch upstream
  git checkout release-1.FILLME
  git merge vX.Y.Z # desired release of upstream
  git push && git push --tags
  ```

- Check that CircleCI pipeline succeeds for the desired Git tag in order to produce images
- Replace `.image.tag` in [the app's `values.yaml`](https://github.com/giantswarm/cluster-api-app/blob/master/helm/cluster-api/values.yaml) with the new tag
- Test as described above
- Open PR
- Once merged, bump the version in the respective collection to deploy it (e.g. [capa-app-collection](https://github.com/giantswarm/capa-app-collection/))

In the _rare_ case of an urgent hotfix, we can create an intermediate release tag like this and use it in `.image.tag` instead of a stable upstream release:

```sh
# Find the latest upstream release on which our `main` branch is based.
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

If this shows any output, please align the `main` branch with the others.
