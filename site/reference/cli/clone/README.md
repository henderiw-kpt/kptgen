---
title: "`clone`"
linkTitle: "clone"
type: docs
description: >
  Clone copies some file from a git repo a target dir.
---

<!--mdtogo:Short
    Copy kubebuilder files to a package.
-->

`clone` copies files from a git repo to a kpt package.

### Synopsis

<!--mdtogo:Long-->

```
kptgen clone GIT_REPO TARGET_DIR [flags]
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
$ kptgen copy ./config ./blueprint/topology
```

<!--mdtogo-->
