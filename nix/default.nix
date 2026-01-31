{ lib, buildGoModule, version ? "unstable" }:

buildGoModule {
  pname = "garbo";
  inherit version;

  src = lib.cleanSource ./.;
  
  vendorHash = null;

  ldflags = [
    "-s"
    "-w"
    "-X main.version=${version}"
  ];

  meta = {
    description = "";
    homepage = "https://github.com/";
    license = lib.licenses.gpl3Plus;
    maintainers = with lib.maintainers; [ pixel-87 ];
    mainProgram = "garbo";
  };
}
