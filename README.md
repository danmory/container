# Container

## How to run

1. Install Latest version of Go - <https://go.dev/dl/>

2. Install netsetgo and move it to /usr/local/bin

    ```bash
        wget "https://github.com/teddyking/netsetgo/releases/download/0.0.1/netsetgo"

        sudo mv netsetgo /usr/local/bin/

        sudo chown root:root /usr/local/bin/netsetgo

        sudo chmod 4755 /usr/local/bin/netsetgo
    ```

3. Install rootfs - <https://www.alpinelinux.org/downloads/> and move it to root "/" directory

4. Compile container

    `` $ go build container.go ``

5. Run container

    `` $ sudo ./container run /bin/sh ``
