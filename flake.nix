{
  description = "PND backend dev environment";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      goVersion = 21; # Change this to update the whole stack
      overlays = [ (final: prev: { go = prev."go_1_${toString goVersion}"; }) ];
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs { inherit overlays system; };
      });
    in
    {
      devShells = forEachSupportedSystem ({ pkgs }: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            # go (specified by overlay)
            go_1_21

            # goimports, godoc, etc.
            gotools

            # https://github.com/mvdan/gofumpt
            gofumpt

            # https://github.com/golangci/golangci-lint
            golangci-lint

            # https://github.com/golang-migrate/migrate
            go-migrate
          ];
        };
      });
    };
}
