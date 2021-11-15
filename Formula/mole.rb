class Mole < Formula
  desc "Tools to deeply process the frontend projects"
  homepage "https://github.com/hsiaosiyuan0/mole"
  url "https://github.com/hsiaosiyuan0/mole/releases/download/v0.0.2/mole_v0.0.2_darwin_amd64.tgz"
  sha256 "d31b7c176866a0dae4e79c582297b251b4205db8ca8ca9a653e4ae6354a8dff2"
  version "0.0.3"

  def install
    bin.install "mole"
  end
end
