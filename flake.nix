{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };
  outputs = {
    self,
    nixpkgs,
    gomod2nix,
    ...
  }: let
    supportedSystems = [
      "x86_64-linux"
      "aarch64-linux"
      "x86_64-darwin"
      "aarch64-darwin"
    ];
    forEachSupportedSystem = f:
      nixpkgs.lib.genAttrs supportedSystems (system:
        f {
          pkgs = import nixpkgs {
            inherit system;
            overlays = [self.overlays.default];
          };
        });
  in {
    overlays.default = final: prev: {
      buildGoApplication =
        gomod2nix
        .legacyPackages
        .${
          prev.stdenv.system
        }
        .buildGoApplication;
      gomod2nixPkg = gomod2nix.packages.${prev.stdenv.system}.default;
    };
    devShells = forEachSupportedSystem ({pkgs}: {
      default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gotools
          golangci-lint
          gomod2nixPkg

          sqlite

          alejandra
        ];
      };
    });
    packages = forEachSupportedSystem ({pkgs, ...}: rec {
      jimaku-tg-notify = pkgs.buildGoApplication {
        pname = "jimaku-tg-notify";
        version = "0.0.7";
        src = self;
        modules = ./gomod2nix.toml;
      };
      default = jimaku-tg-notify;
    });
  };
}
