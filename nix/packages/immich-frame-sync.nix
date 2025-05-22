{ self, ... }:
{
  perSystem =
    {
      lib,
      pkgs,
      ...
    }:
    {
      packages.immich-frame-sync =
        let
          shortRev = self.shortRev or self.dirtyShortRev;
          rev = self.rev or self.dirtyRev;
          tag = builtins.getEnv "RELEASE_VERSION";

          version = if tag != "" then tag else rev;
        in
        pkgs.buildGoModule {
          name = "immich-frame-sync-${shortRev}";

          src = lib.fileset.toSource {
            fileset = lib.fileset.unions [
              ../../cmd
              ../../db/migrations
              ../../go.mod
              ../../go.sum
              ../../main.go
              ../../pkg
            ];
            root = ../..;
          };

          ldflags = [
            "-X github.com/kalbasit/immich-frame-sync/cmd.Version=${version}"
          ];

          subPackages = [ "." ];

          vendorHash = "sha256-XXMxim2LMyR2qFsXGXWUd5yhSAPYxod/KPwDacEhmWc=";

          doCheck = true;

          nativeBuildInputs = [
            pkgs.dbmate # used for testing
          ];

          postInstall = ''
            mkdir -p $out/share/immich-frame-sync
            cp -r db $out/share/immich-frame-sync/db
          '';

          meta = {
            description = "Synchronize your digital picture frame with random, filtered assets from your Immich library. Supports FTP uploads (extensible) and syncing deletions back to Immich trash";
            homepage = "https://github.com/kalbasit/immich-frame-sync";
            license = lib.licenses.mit;
            mainProgram = "immich-frame-sync";
            maintainers = [ lib.maintainers.kalbasit ];
          };
        };

    };
}
