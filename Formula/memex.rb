class Memex < Formula
  desc "Memex CLI tool"
  homepage "https://github.com/jai0651/memex"
  head "https://github.com/jai0651/memex.git", branch: "main"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
    pkgshare.install "shell_integration.sh"
  end

  test do
    system "#{bin}/memex", "--help"
  end
end
