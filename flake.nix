{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in with pkgs; {
        devShell = mkShell { buildinputs = [ go vgo2nix ]; };
        defaultPackage = buildGoPackage {
          pname = "mpdcord";
          version = "0.1";
          src = ./.;
          goDeps = ./deps.nix;
          goPackagePath = "github.com/bandithedoge/mpdcord";
        };
      });
}
