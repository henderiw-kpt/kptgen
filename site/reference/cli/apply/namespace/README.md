---
title: "`namespace`"
linkTitle: "namespace"
type: docs
description: >
  Add a namespace to the package.
---

<!--mdtogo:Short
    Add a namespace to the package.
-->

`namespace` adds a namespace to the package.

### Synopsis

<!--mdtogo:Long-->

```
kptgen add namespace TARGET_DIR [flags]
```

#### Args

```
TARGET_DIR:
  The target directory of the package
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
# Add a namespace to the package
$ kptgen add namespace ./blueprint/admin 
```

<!--mdtogo-->
