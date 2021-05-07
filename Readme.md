#  More-Ports

More Ports is a proxy service to establish all web-based applications on different ports on the server-side over a well known TCP port.

It is good to use while your CDN does not allow few ports or inside the restricted network such as school.

## How it works ?

Firstly More-Ports runs as server mode on the server-side. Server applications start to listen to port 80 for HTTP and port 443 for HTTPS request. When the client tries to connect to non-regular HTTPS(80) and HTTP(443) ports, more-ports forward their request to localhost:8080. If the client requests different ports such as 3000, this request is handled by the more-ports client and forwarded to the more-ports server and the server redirects this request to the indicated port, for this example which is localhost:3000. 

![Work Logic](/docs/work-logic.png)

## Server-Side

By default, the system starts a web server on port 80 and 443 for HTTP and HTTPS. If the request does not have a port information such as http://example.com or https://example.com, the request forwarded to http://localhost:8080 by default. If it has a port information, request forwarded to indicated port on localhost.

**Note:** server mode is only supported to Linux.

### Configuration

Firstly start this service in host network mode to access shared ports if you use a container. If you don’t provide any certificate for https server, the program will create a self signed certificate.   
To use your own certificates, you can mount your certificates inside a container with the docker mount option and provide the certificate locations with server-cert and server-key arguments.

If you want to use this software behind the CDN solution and authorize actions taken at edge, you can control whether a request comes from CDN or not with client-certificate check for HTTPS. Also don’t forget to set the server-name option, to prevent different domains from hitting your server from the same CDN and redirect all http requests to https with ` -https-redirect` flag.  
Cloudflare's 'Authenticated Origin Pulls’ cert is included by default in the container, to use add `--client-cert /config/client-cloudflare.pem` arg on your docker run command.

```bash
  -client-cert string
        Client certificate location to authorize the client.
  -default-port string
        Port not defined schemes redirect to this port. (default "8080")
  -h    Print help for server mode. (This)
  -http string
        HTTP server listen address. (default ":80")
  -https string
        HTTPS server listen address. (default ":443")
  -https-redirect
        Redirect all http requests to https.
  -remote string
        Remote address for forwarded ports. (default "127.0.0.1")
  -server-cert string
        Server cert location.
  -server-key string
        Server key location.
  -server-name string
        Server name check
```

### Examples
-  Start with defaults
```bash
docker run -it --rm --network host ghcr.io/ahmetozer/more-ports server
```

-  Redirect http request to https
```bash
docker run -it --rm --network host ghcr.io/ahmetozer/more-ports server -https-redirect
```

-  With Cloudflare's 'Authenticated Origin Pulls’, server name control and https redirect
```bash
docker run -it --rm --network host ghcr.io/ahmetozer/more-ports server --server-name myserver.example.com --client-cert /config/client-cloudflare.pem --https-redirect
```

-  Custom certificate
```bash
docker run -it --rm --network host \
--mount type=bind,source="/data/certs/example.com/example.com.cert",target="/cert/my.cert",readonly \
 --mount type=bind,source="/data/keys/example.com/example.com.key",target="/cert/my.key",readonly \
 ghcr.io/ahmetozer/more-ports server --server-name myserver.example.com \
--client-cert /config/client-cloudflare.pem --https-redirect --server-cert /cert/my.cert --server-cert /cert/my.key
```

## Client Side

If the clients only access regular ports like port 80 for HTTP or port 443 for HTTPS, they are not required to use this service but if they want to access different ports, they have to use the `more-ports` application on their system at the client mode.  
You can access the precompiled application on the GitHub releases.

### Configuration

In general, you don't need to make any changes for configurations on the client-side of the application, just set proxy configuration URL on your operating system from the program given URL.

- Start the application

```bash
2021/05/07 01:25:48 More Ports Service
2021/05/07 01:25:48 Client mode
2021/05/07 01:25:48 Remote ports for http :80, https :443
2021/05/07 01:25:48 Client proxy server started at 127.0.0.1:8080
2021/05/07 01:25:48 Client proxy configuration located at http://127.0.0.1:8080/proxy.pac
```

- Change settings on clients os.
In this example for Windows.  
![Client Windows configure](/docs/client-windows-configuration.jpg)


