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
- Run `make generate`
- Push a feature branch
- Check the test app release thoroughly on a management cluster
