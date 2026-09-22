# typed: false
# frozen_string_literal: true

# This file is a reference formula for local development/testing.
# For actual installation, tap the repository:
#   brew tap cortezramos/tuipr
#   brew install tuipr
#
# GoReleaser manages the formula published to the tap repository.

class Tuipr < Formula
  desc "Keyboard-driven Pull Request Lifecycle Manager for your terminal"
  homepage "https://github.com/cortezramos/tuipr"
  version "0.1.0"
  license "MIT"

  depends_on "gh"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/cortezramos/tuipr/releases/download/v#{version}/tuipr_#{version}_Darwin_x86_64.tar.gz"
      sha256 "8be4e7b362ac3c2dbf8c8b9dfe4274e45c6f2486da59fb7bfb4d797abddb0ccc"
    end
    if Hardware::CPU.arm?
      url "https://github.com/cortezramos/tuipr/releases/download/v#{version}/tuipr_#{version}_Darwin_arm64.tar.gz"
      sha256 "138851cd46547127891ad18afa0e415fd917b3bf5fdd28b7a27f6b029c915284"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/cortezramos/tuipr/releases/download/v#{version}/tuipr_#{version}_Linux_x86_64.tar.gz"
      sha256 "976c84c85f7dcbfbc3f3d8e2188f39d84d94991bfd795f3db82d36d1a551823b"
    end
  end

  def install
    bin.install "tuipr"
  end

  test do
    assert_path_exists bin/"tuipr"
    assert_predicate bin/"tuipr", :executable?
  end
end
