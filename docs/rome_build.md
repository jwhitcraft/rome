## rome build

Build SugarCRM

### Synopsis


This will take a source version of Sugar and substitute out all the necessary build tags and create an
installable copy of Sugar for you to use and dev on.

By default this will ignore sugarcrm/node_modules, but build sugarcrm/sidecar/node_modules to save on time since the
node_modules are not required inside of SugarCRM but are for Sidecar.


```
rome build
```

### Examples

```
rome build -v 7.9.0.0 -f ent -d /tmp/sugar /path/to/mango/git/checkout
```

### Options

```
      --clean                  Remove Existing Build Before Building
      --clean-cache            Clears the cache before doing the build. This will only delete certain cache files before doing a build.
  -d, --destination string     Where should the built files be put
      --file-buffer-size int   Size of the file buffer before it gets reset (default 4096)
      --file-workers int       Number of Workers to start for processing files (default 80)
  -f, --flavor string          What Flavor of SugarCRM to build (default "ent")
  -l, --folder string          What folder should we build to on the server, if left empty, it will build to a folder named <version><flavor>
      --port string            What is the server port (default "47600")
  -s, --server string          What server should we build to
  -v, --version string         What Version is being built
```

### SEE ALSO
* [rome](rome.md)	 - A Tool for Building Sugar from source

###### Auto generated by spf13/cobra on 27-Feb-2017
