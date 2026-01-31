{
  mkShell,
  callPackage,

  go,
  gopls,
  gofumpt,
  goreleaser,
  nodejs,
  pnpm,
}:

let
  defaultPackage = callPackage ./default.nix { };
in
mkShell {
  inputsFrom = [ defaultPackage ];

  packages = [
    go
    gopls
    gofumpt
    goreleaser
    nodejs
    pnpm
  ];
}
