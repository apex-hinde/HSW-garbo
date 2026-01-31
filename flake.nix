{
  description = "";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = 
    { self, nixpkgs, ... }:
    let
      forAllSystems = 
        f: nixpkgs.lib.genAttrs nixpkgs.lib.systems.flakeExposed (
          system: f nixpkgs.legacyPackages.${system}
        );
    in 
    {
      packages = forAllSystems (pkgs: {
        garbo = pkgs.callPackage ./nix/default.nix { version = self.shortRev or "unstable"; };
        default = self.packages.${pkgs.stdenv.hostPlatform.system}.garbo;
        });

      devShells = forAllSystems (pkgs: {
        default = pkgs.callPackage ./nix/shell.nix { };
        });

      overlays.default = final: _: {
        garbo = final.callPackage ./nix/default.nix { version = self.shortRev or "unstable"; };
      };
    };
}
