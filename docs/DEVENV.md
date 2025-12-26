# melange2 Development Environment

To ease the development of melange2, we offer a script that builds
a [Wolfi](https://github.com/wolfi-dev/community)-based image with
the required tooling to compile and run melange2. When running the
devenv, you'll drop into a Wolfi shell, with your local fork directory
will mounted in `$PWD`.

## Launching the DevEnv

To launch the development environment, simply run the script from
the top of the melange2 environment:

```bash
# Clone the melange2 repo

git clone git@github.com:dlorenc/melange2.git

cd melange2

# Build and launch the devenv:

./hack/make-devenv.sh
                _                        
 _ __ ___   ___| | __ _ _ __   __ _  ___ 
| '_ ` _ \ / _ \ |/ _` | '_ \ / _` |/ _ \
| | | | | |  __/ | (_| | | | | (_| |  __/
|_| |_| |_|\___|_|\__,_|_| |_|\__, |\___|
                              |___/      

Welcome to the melange2 development environment!

To run melange from your local fork run:
        go run ./main.go


```

The image is built every time it is needed, it will be 
built using Chainguard's [apko image](https://github.com/chainguard-images/images/tree/main/images/apko).

## Requirements

The only requirement to run the melange development environment
is to haver docker running in you local system. The script will
take care of the rest.

# Reset/Delete the Development Environment

The script will only build the Wolfi image when it cannot find it
in the local docker daemon. If you need to rebuild it, simply 
delete it using docker:

```bash
docker rmi melange-inception:latest
```

Running the script again will force a clean build of the image.