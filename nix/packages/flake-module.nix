{
  imports = [ ];

  perSystem =
    { config, pkgs, ... }:
    {
      packages.default = pkgs.hello;
    };
}
