## Building Elastifile's ECFS CSI provisioner

### Make targets

CSI ECFS plugin can be compiled in a form of a binary file or in a form of a Docker image.

When compiled as a binary file, the result is stored in `_output/` directory with the name `ecfsplugin`.

When compiled as an image, it's stored in the local Docker image registry.

#### Building binary
```bash
$ make binary
```

#### Building Docker image
This will build the binary and create docker image
```bash
$ make image
```

#### Pushing the image to Docker Hub
This will push an existing image
```bash
$ make push
```

#### One-stop shop
This will build the binary, create docker image and push it to Docker Hub
```bash
$ make all
```

#### Clean
This will remove the  built artifacts
```bash
$ make clean
```

### Changing the defaults
By default, make will tag the images as 'dev'.

If you want to use a different tag, set `PLUGIN_TAG` environment variable accordingly

Example:
```bash
$ PLUGIN_TAG=v0.0.1 make all
```

### Dependency management

This project uses [dep](https://github.com/golang/dep) as it dependency management tool.

After cloning the project, it is recommended to run `dep ensure`.

Another case where it's recommended to run `dep ensure` is after making changes to the project's imports.
