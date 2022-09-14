---
title: "`generate`"
linkTitle: "generate"
type: docs
description: >
  Generate a package.
---

<!--mdtogo:Short
    Generate a package.
-->

`generate` generates a kpt package from a kustomize config directory.

### Synopsis

<!--mdtogo:Long-->

```
kptgen generate SOURCE_DIR TARGET_DIR [flags]
```

#### Args

```
SOURCE_DIR:
  The source directory of the kustomize manifests.

TARGET_DIR:
  The target diretcory the package would be created in
```

#### Flags

```
--description
  Short description of the package. (default "sample description")

--keywords
  A list of keywords describing the package.

--site
  Link to page with information about the package.
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
# Creates a new Kpt Package from a kubebuilder manifest files.
$ kptgen generate ./config ./blueprint/topology
```

<!--mdtogo-->
