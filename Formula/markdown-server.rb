# Homebrew installation formula
class MarkdownServer < Formula
  desc "Markdown server software"
  homepage "https://github.com/christianhellsten/markdown-server"
  version "latest"
  url begin
    os = if OS.mac?
           "darwin"
         elsif OS.linux?
           "linux"
         else
           "windows"
         end

    arch = if Hardware::CPU.intel?
             "amd64"
           elsif Hardware::CPU.arm?
             "arm64"
           end

    ext = OS.mac? || OS.linux? ? "" : ".exe"
    "https://github.com/christianhellsten/markdown-server/releases/latest/download/markdown-server-#{os}-#{arch}#{ext}"
  end

  def install
    bin_name = OS.mac? || OS.linux? ? "markdown-server" : "markdown-server.exe"
    bin.install Dir["markdown-server-*"].first => bin_name
  end

  test do
    bin_name = OS.mac? || OS.linux? ? "markdown-server" : "markdown-server.exe"
    system "#{bin}/#{bin_name}", "--version"
  end
end
