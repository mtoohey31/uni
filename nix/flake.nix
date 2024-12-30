{
  description = "uni";

  inputs = {
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }: {
    overlays.default = final: _: {
      uni = final.callPackage
        ({ buildGoModule, python3Packages }: buildGoModule {
          pname = "uni";
          version = "0.1.0";
          src = builtins.path { path = ./..; name = "uni-src"; };
          vendorHash = null;
          nativeBuildInputs = [ python3Packages.cram ];
          checkPhase = ''
            runHook preCheck

            UNI=$GOPATH/bin/uni cram test.t

            runHook postCheck
          '';
        })
        { };
    };
  } // flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs {
        overlays = [ self.overlays.default ];
        inherit system;
      };
      inherit (pkgs) gopls mkShell uni;
    in
    {
      packages.default = uni;

      devShells.default = mkShell {
        inputsFrom = [ uni ];
        packages = [ gopls ];
        shellHook = ''
          export UNI=$PWD/uni
        '';
      };
    });
}
