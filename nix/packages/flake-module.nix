{
  imports = [
    ./docker.nix
    ./immich-frame-sync.nix
  ];

  perSystem =
    { config, ... }:
    {
      packages.default = config.packages.immich-frame-sync;
    };
}
