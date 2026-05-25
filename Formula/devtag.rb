class Devtag < Formula
  desc "Play an audio clip whenever you git push"
  homepage "https://github.com/kelvinxiao01/devtag"
  url "https://github.com/kelvinxiao01/devtag/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "722f976985db8ebdb97aee1d34924b576cbd2404c523de3ae446bb6def8f6cfc"
  license "MIT"
  head "https://github.com/kelvinxiao01/devtag.git", branch: "main"

  depends_on "go" => :build
  depends_on :macos

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  def caveats
    <<~EOS
      Finish setup with:
        devtag set /path/to/your/audio.mp3
        devtag install

      `devtag install` sets git's global core.hooksPath to ~/.devtag/hooks.
      To remove later: `devtag uninstall`.
    EOS
  end

  test do
    assert_match "Usage", shell_output("#{bin}/devtag --help")
  end
end
