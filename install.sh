#!/usr/bin/env sh
set -eu

REPO="https://api.github.com/repos/mizuchilabs/ssh-nexus/releases"

# Downloads the latest release and moves it into ~/.local/bin
main() {
	binary="nexus"
	if [ "${1:-}" = "agent" ]; then
		binary="nexus-agent"
	elif [ "${1:-}" = "uninstall" ]; then
		uninstall
		exit 0
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

	url="https://github.com/mizuchilabs/ssh-nexus/releases/download/${latest}/${binary}_${platform}_${arch}"

	echo "Downloading $binary from $url"
	if which curl >/dev/null 2>&1; then
		command curl -sfL "$url" -o "$temp"
	elif which wget >/dev/null 2>&1; then
		command wget -O "$temp" "$url"
	else
		echo "Could not find 'curl' or 'wget' in your path"
		exit 1
	fi

	# Ensure the file is not empty
	if [ ! -s "$temp" ]; then
		echo "Failed to download $binary, file is empty"
		exit 1
	fi

	# Ensure the file is executable
	"$platform" "$@"

	if echo "$PATH" | grep -q "$HOME/.local/bin"; then
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

uninstall() {
	binary="nexus"
	[ "$(which nexus-agent 2>/dev/null)" ] && binary="nexus-agent"
	bin_path="$HOME/.local/bin/$binary"
	if [ -f "$bin_path" ]; then
		echo "Uninstalling $binary..."
		rm -f "$bin_path"
		echo "$binary has been removed from $HOME/.local/bin."
	else
		echo "$binary is not installed."
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
