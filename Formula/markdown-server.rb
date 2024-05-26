require 'net/http'
require 'json'
require 'uri'

class MarkdownServer < Formula
  desc "Markdown server written in Go"
  homepage "https://github.com/christianhellsten/markdown-server"
  url "https://github.com/christianhellsten/markdown-server/archive/refs/tags/v5.tar.gz"
  sha256 "6f4556a179f997981647c503096bf35f041322130704a6cfa0def31eda4ff8f2"
  license "MIT"

  def install
    # Fetch the latest release information from GitHub
    releases_url = URI("https://api.github.com/repos/christianhellsten/markdown-server/releases/latest")
    release_data = Net::HTTP.get(releases_url)
    release_json = JSON.parse(release_data)

    # Determine the download URL for the correct platform binary
    asset = release_json["assets"].find { |a| a["name"].include?(os_binary_name) }
    if asset.nil?
      odie "No binary found for the current platform."
    end

    binary_url = URI(asset["browser_download_url"])

    # Download the binary using Net::HTTP
    bin_dir = "#{buildpath}/bin"
    mkdir_p bin_dir
    File.open("#{bin_dir}/markdown-server", "wb") do |file|
      Net::HTTP.get_response(binary_url) do |response|
        response.read_body do |chunk|
          file.write(chunk)
        end
      end
    end

    # Make the binary executable
    chmod "+x", "#{bin_dir}/markdown-server"

    # Install the binary
    bin.install "#{bin_dir}/markdown-server"
  end

  def os_binary_name
    if OS.mac?
      "markdown-server-darwin-amd64"
    elsif OS.linux?
      "markdown-server-linux-amd64"
    elsif OS.windows?
      "markdown-server-windows-amd64.exe"
    else
      raise "Unsupported platform"
    end
  end

  test do
    assert_match "markdown-server version", shell_output("#{bin}/markdown-server --version")
  end
end
