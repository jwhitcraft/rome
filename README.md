# rome
A new build system for SugarCRM

## Installing
```shell
curl -L http://h2ik.co/rome/rome-`uname -s`-`uname -m` -o /usr/local/bin/rome; chmod +x /usr/local/bin/rome
```

## Updating
To update rome, just run `rome self-update`

## Usage

### Build
This command will build all the files inside of the source directory

`rome build --version=7.9.0.0 --flavor=ent --destination=/tmp/sugar-build /path/to/mango/checkout/sugarcrm`

### Watch
This command will keep a process running and build each file as it's created or modified (experimental!)

`rome watch --version=7.9.0.0 --flavor=ent --destinations=/tmp/sugar-build /path/to/mango/checkout/sugarcrm`

## Help
`rome help build`

## Building Rome
Make sure that you have golang installed

`make clean; make -e VERSION=2.0.0AlphaX`
