class MarkdownServer < Formula
  desc "Markdown server software"
  homepage "https://github.com/christianhellsten/markdown-server"
  url "https://github.com/christianhellsten/markdown-server/releases/latest/download/markdown-server-darwin-amd64"
  version "latest"
  sha256 "f709326b75ebc9a96b9a767c517323dcb056a5e53646ef5f59814ef2bc64d10a"
  # TODO: sha256 :no_check

  def install
    bin.install "markdown-server-darwin-amd64" => "markdown-server"
  end

  test do
    system "#{bin}/markdown-server", "--version"
  end
end
