---
title: "`copy`"
linkTitle: "copy"
type: docs
description: >
  Copy some files from kubebuilder to a kpt package.
---

<!--mdtogo:Short
    Copy kubebuilder files to a package.
-->

`copy` copies files from a kubebuilder config directory to a kpt package.

### Synopsis

<!--mdtogo:Long-->

```
kptgen copy SOURCE_DIR TARGET_DIR [flags]
```

#### Args

```
SOURCE_DIR:
  The source directory of the kubebuilder manifests.

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
