{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    devenv.url = "github:cachix/devenv";
  };

  outputs = { self, nixpkgs, devenv, systems, ... } @ inputs:
    let
      forEachSystem = nixpkgs.lib.genAttrs (import systems);
    in
    {
      devShells = forEachSystem
        (system:
          let
            pkgs = nixpkgs.legacyPackages.${system};
          in
          {
            default = devenv.lib.mkShell {
              inherit inputs pkgs;
              modules = [
                {
                  # https://devenv.sh/reference/options/
                  packages = [
                    pkgs.bash
                    pkgs.coreutils
                    pkgs.findutils
                    pkgs.gci
                    pkgs.ginkgo
                    pkgs.git
                    pkgs.gnused
                    pkgs.gnugrep
                    pkgs.gnumake
                    pkgs.go
                    pkgs.gofumpt
                    pkgs.gojq
                    pkgs.golangci-lint
                    pkgs.golines
                    pkgs.goreleaser
                    pkgs.gotestsum
                    pkgs.gotools
                    pkgs.kubernetes-helm
                    pkgs.pre-commit
                    pkgs.shfmt
                  ];
                }
              ];
            };
          });
    };
}
