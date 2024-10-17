#!/usr/bin/env sh
set -eu

REPO="https://github.com/api/v1/repos/mizuchilabs/ssh-nexus/releases"

# Downloads the latest release and moves it into ~/.local/bin
main() {
	binary="nexus"
	if [ "${1:-}" = "agent" ]; then
		binary="nexus-agent"
	fi
	platform="$(uname -s)"
	arch="$(uname -m)"
	temp="$(mktemp -t "$binary-XXXXXX")"
	latest=$(curl -s "${REPO}"/latest | grep -o '"tag_name":.*' | cut -d '"' -f 4)

	if [ "$platform" = "Darwin" ]; then
		platform="macos"
	elif [ "$platform" = "Linux" ]; then
		platform="linux"
	else
		echo "Unsupported platform $platform"
		exit 1
	fi

	url="https://github.com/mizuchilabs/ssh-nexus/releases/download/${latest}/${binary}_${platform}_${arch}"

	case "$platform-$arch" in
	macos-arm64* | linux-arm64* | linux-armhf | linux-aarch64)
		arch="arm64"
		;;
	macos-x86* | linux-x86* | linux-i686*)
		arch="amd64"
		;;
	*)
		echo "Unsupported platform or architecture"
		exit 1
		;;
	esac

	echo "Downloading $binary from $url"
	if which curl >/dev/null 2>&1; then
		curl() {
			command curl -fL "$@" -o "$temp"
		}
	elif which wget >/dev/null 2>&1; then
		curl() {
			command wget -O- "$@" >"$temp"
		}
	else
		echo "Could not find 'curl' or 'wget' in your path"
		exit 1
	fi

	"$platform" "$@"

	if [ "$(which $binary)" = "$HOME/.local/bin/$binary" ]; then
		echo "$binary has been installed. Run with '$binary'."
	else
		echo "To run $binary from your terminal, you must add ~/.local/bin to your PATH"
		echo "Run:"

		case "$SHELL" in
		*zsh)
			echo "   echo 'export PATH=\$HOME/.local/bin:\$PATH' >> ~/.zshrc"
			echo "   source ~/.zshrc"
			;;
		*fish)
			echo "   fish_add_path -U $HOME/.local/bin"
			;;
		*)
			echo "   echo 'export PATH=\$HOME/.local/bin:\$PATH' >> ~/.bashrc"
			echo "   source ~/.bashrc"
			;;
		esac

		echo "To run $binary now, '~/.local/bin/$binary'"
	fi
}

linux() {
	# Setup ~/.local directories
	mkdir -p "$HOME/.local/bin"

	# Move the binary
	mv "$temp" "$HOME/.local/bin/$binary"

	# Make it executable
	chmod +x "$HOME/.local/bin/$binary"
}

macos() {
	# Setup ~/.local directories
	mkdir -p "$HOME/.local/bin"

	# Move the binary
	mv "$temp" "$HOME/.local/bin/$binary"

	# Make it executable
	chmod +x "$HOME/.local/bin/$binary"
}

main "$@"
