{
  perSystem =
    {
      config,
      pkgs,
      ...
    }:
    {
      packages.docker = pkgs.dockerTools.buildLayeredImage {
        name = "kalbasit/immich-frame-sync";
        contents =
          let
            etc-passwd = pkgs.writeTextFile {
              name = "passwd";
              text = ''
                root:x:0:0:Super User:/homeless-shelter:/dev/null
              '';
              destination = "/etc/passwd";
            };

            etc-group = pkgs.writeTextFile {
              name = "group";
              text = ''
                root:x:0:
              '';
              destination = "/etc/group";
            };
          in
          [
            # required for Open-Telemetry auto-detection of process information
            etc-passwd
            etc-group

            # required for TLS certificate validation
            pkgs.cacert

            # required for working with timezones
            pkgs.tzdata

            # required for migrating the database
            pkgs.dbmate

            # the immich-frame-sync package
            config.packages.immich-frame-sync
          ];
        config = {
          Cmd = [ "/bin/immich-frame-sync" ];
          Env = [
            "DBMATE_MIGRATIONS_DIR=/share/immich-frame-sync/db/migrations"
            "DBMATE_SCHEMA_FILE=/share/immich-frame-sync/db/schema.sql"
            "DBMATE_NO_DUMP_SCHEMA=true"
          ];
          ExposedPorts = {
            "8356/tcp" = { };
          };
          Labels = {
            "org.opencontainers.image.description" = "Drive a picture frame assets using immich";
            "org.opencontainers.image.licenses" = "MIT";
            "org.opencontainers.image.source" = "https://github.com/kalbasit/immich-frame-sync";
            "org.opencontainers.image.title" = "immich-frame-sync";
            "org.opencontainers.image.url" = "https://github.com/kalbasit/immich-frame-sync";
          };
        };
      };

      packages.push-docker-image = pkgs.writeShellScript "push-docker-image" ''
        set -euo pipefail

        if [[ ! -v DOCKER_IMAGE_TAGS ]]; then
          echo "DOCKER_IMAGE_TAGS is not set but is required." >&2
          exit 1
        fi

        for tag in $DOCKER_IMAGE_TAGS; do
          echo "Pushing the image tag $tag for system ${pkgs.hostPlatform.system}. final tag: $tag-${pkgs.hostPlatform.system}"
          ${pkgs.skopeo}/bin/skopeo --insecure-policy copy \
            "docker-archive:${config.packages.docker}" docker://$tag-${pkgs.hostPlatform.system}
        done
      '';
    };
}
