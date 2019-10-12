gitsam
======

Single-serve git-over-i2p repositories using SAM.

This is still pretty much a toy. I made it as a goof after I got tired this
Saturday. I do intend to refine it as I get more chances to do so but for now
it's just for fun.

This simple application combines an I2P Tunnel created by sam-forwarder and a
git repository created by gitkit in order to simply hote a Git repository over
the I2P network. It creates two tunnels that serve the same repository, one of
them is read-only and served with plain HTTP-over-I2P, and one of them is
read-write and served with SSH.

Eventually, it will have extensive options for optimizing your git
configuration.
