{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoPackage {
  pname = "mpdcord";
  version = "0.1";
  src = ./.;
  goDeps = ./deps.nix;
  goPackagePath = "github.com/bandithedoge/mpdcord";
  buildInputs = with pkgs; [ go vgo2nix ];
  preConfigure = ''
    ${pkgs.vgo2nix}/bin/vgo2nix
  '';
}
