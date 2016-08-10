with import <nixpkgs> { };
with goPackages;

let 
  myDeps =  (import ./myDeps.nix);
in

buildGoPackage rec {
  name = "cfp-${version}";
  version = "0.0.1";

  shellHook = ''
  '';

  goPackagePath = "github.com/cfp/cfp";

  buildInputs = with myDeps; [ letsencrypt jq gomailv2 captcha go gocraft-web CompileDaemon.bin postgresql gorm inflection pq crypto net];
 # ++ [gocode.bin];
# ++ [ gorilla-sessions gorilla-context gorilla-securecookie];
}

