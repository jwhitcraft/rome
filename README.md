# rome
A new build system for SugarCRM

# docker

Build the docker image:
```
docker build -t registry.sugarcrm.net/rome:1.0.0 .
```

Generate a build using the docker image and volume mounts:
```
docker run -it --rm -v /path/to/mango/sugarcrm:/sugarcrm -v /tmp/sugarcrm-build:/sugarcrm-build --name rome registry.sugarcrm.net/rome:1.0.0 go-wrapper run build /sugarcrm -d /sugarcrm-build -f ent -v 7.8.0.0
```
