# what

this is the call for papers form used for the 2016 tuebix.org signup

![A screenshot of the cfp form](https://raw.githubusercontent.com/nixcloud/cfp/master/screenshots/cfp.jpg)

# how

uses:

* go
* gocraft
* gomailv2
* javascript
* materialize

## cfp.sh 
```
#!/bin/sh
source /etc/profile
cd /home/joachim/cfp
nix-shell --command "while true; do go run server.go ; done"
```

## systemd service
```
  systemd.services.cfp = {
    wantedBy = [ "multi-user.target" ];
     after = [ "network.target" ];
     serviceConfig = {
       #Type = "forking";
       User = "joachim";
       ExecStart = ''/home/joachim/cfp.sh'';
       ExecStop = ''
       '';
     };
  };
```

## apache reverse proxy config
```
      ... 
      # nixcloud.io (https)
      {
        hostName = "nixcloud.io";
        serverAliases = [ "nixcloud.io" "www.nixcloud.io" ];

        documentRoot = "/www/nixcloud.io/";
        enableSSL = true;
        sslServerCert = "/ssl/nixcloud.io-2015.crt";
        sslServerKey = "/ssl/nixcloud.io-2015.key";
        sslServerChain = "/ssl/nixcloud.io-2015-intermediata.der";

        extraConfig = ''
          ...
          RewriteRule ^/cfp$ /cfp/ [R]
          ProxyPass /cfp/ http://127.0.0.1:3000/ retry=0
          ProxyPassReverse /cfp/ http://127.0.0.1:3000/
          ...
        '';
      ...
```
# license

public domain except the used libraries. relicense it the way you need it.

author: joachim schiele <js@lastlog.de
