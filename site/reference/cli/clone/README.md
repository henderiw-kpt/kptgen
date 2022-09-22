---
title: "`clone`"
linkTitle: "clone"
type: docs
description: >
  Clone copies the yaml files of a git repo to a target directory.
---

<!--mdtogo:Short
    Clone copies the yaml files of a git repo to a target directory.
-->

`clone` copies the yaml files of a git repo to a target directory.

### Synopsis

<!--mdtogo:Long-->

```
kptgen clone GIT_REPO_URL TARGET_DIR [flags]
```

#### Args

```
GIT_REPO:
  The source repo of manifests.

TARGET_DIR:
  The target diretcory the package would be created in
```

<!--mdtogo-->

### Examples

{{% hide %}}

<!-- @makeWorkplace @verifyExamples-->

```
# Set up workspace for the test.
TEST_HOME=$(mktemp -d)
cd $TEST_HOME
```

{{% /hide %}}

<!--mdtogo:Examples-->

<!-- @ @verifyStaleExamples-->

```shell
# Creates a new Kpt Package from kubebuilder manifest files.
$ kptgen clone https://github.com/henderiw-kpt/kptgen-templates ./blueprint/topology
```

<!--mdtogo-->
