## Contributing to Bucket

Bucket specialized framework for one-bucket usage of Couchbase. Our aim to provide simple usage of Couchbase over best practices like:

- Handle bucket as SQL's database
- No embedded fields over the JSON, only simple ones per document (this makes indexing lot more efficient)
- Full-text search capabilities
- Automatic document indexing based on tagged structs

#### Contributing

We are following the `git flow` principles.

- New changes comes as feature branches. Branch name looks like ex.: `feature/T52-short-name` where the T{num} represents the issue number
- We ONLY accepts pull requests into our develop branch.
- The most up-to-date branch is the develop, so new features should created from it.
- Releases happens frequently from develop to master with semantic versions.

#### Development environment.

Since it's a framework for Couchbase it's required to run a Couchbase. We can provide the necessary infrastructure by the following command. (NOTE: after the start it will setup a cluster and a bucket, so before running the tests it's necessary to wait a bit)

- `make dev`


