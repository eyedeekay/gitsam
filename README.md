gitsam
======

Easy git-over-i2p repositories using SAM. Run it in the root of an
existing git repository or in a directory containing many git repositories and
automatically create a mirror on I2P.

This is still pretty much a toy. I made it as a goof after I got tired this
Saturday. I do intend to refine it as I get more chances to do so but for now
it's just for fun.

The versions don't mean anything. I have to increment when I update go mod and I
broke it at first.

This simple application combines an I2P Tunnel created by sam-forwarder and a
git repository created by gitkit in order to simply host a group of Git
repositories over the I2P network. It creates two tunnels to serve groups of
repositories, them is read-only and served with plain HTTP-over-I2P, and one of
them is read-write and served with SSH-over-I2P. Optionally, the root folder
may also be a git repository.

Optionally, if one is not already present, it will set up a post-recieve hook
so that the read-only HTTP mirror can double as a web page. Since it uses
eephttpd, in the absence of an index.html file an HTML version of the README.md
file will be shown, and it can be scripted with tengo(A pure-go fast scripting
language).

Eventually, it will have extensive options for optimizing your git configuration.


        Usage of ./gitsam/gitsam:
          -a string
                hostname to serve on (default "127.0.0.1")
          -c	Use an encrypted leaseset(true or false)
          -d string
                the directory of static files to host(default ./www) (default "./")
          -f string
                Use an ini file for configuration (default "none")
          -g	Uze gzip(true or false) (default true)
          -i	save i2p keys(and thus destinations) across reboots (default true)
          -ib int
                Set inbound tunnel backup quantity(0 to 5) (default 1)
          -il int
                Set inbound tunnel length(0 to 7) (default 3)
          -iq int
                Set inbound tunnel quantity(0 to 15) (default 2)
          -iv int
                Set inbound tunnel length variance(-7 to 7)
          -l string
                Type of access list to use, can be "whitelist" "blacklist" or "none". (default "none")
          -lp
                launch a "Page" on read-only HTTP branch (default true)
          -n string
                name to give the tunnel(default gitsam) (default "gitsam")
          -ob int
                Set outbound tunnel backup quantity(0 to 5) (default 1)
          -ol int
                Set outbound tunnel length(0 to 7) (default 3)
          -oq int
                Set outbound tunnel quantity(0 to 15) (default 2)
          -ov int
                Set outbound tunnel length variance(-7 to 7)
          -p string
                port to serve locally on (default "7882")
          -pk string
                Path to the authorized users SSH public key (default "./id_rsa.pub")
          -pp string
                port to serve read-only pages on (default "7883")
          -r	Reduce tunnel quantity when idle(true or false)
          -rc int
                Reduce idle tunnel quantity to X (0 to 5) (default 3)
          -rt int
                Reduce tunnel quantity after X (milliseconds) (default 600000)
          -sh string
                sam host to connect to (default "127.0.0.1")
          -sk string
                Path to the directory with the private keys used to authenticate the server
          -sp string
                sam port to connect to (default "7656")
          -z	Allow zero-hop, non-anonymous tunnels(true or false)
