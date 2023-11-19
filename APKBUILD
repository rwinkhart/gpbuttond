# Maintainer: Randall Winkhart <idgr@tutanota.com>
pkgname=gpbuttond
pkgver=CHANGE-ME
pkgrel=0
pkgdesc="Maps GPIO pins to keyboard keys"
url="https://github.com/rwinkhart/gpbuttond"
arch="armhf x86_64"
license="GPL-3"
options="!check"
depends=""
makedepends="go"
source="$pkgname-v$pkgver.tar.gz::https://github.com/rwinkhart/gpbuttond/archive/refs/tags/v$pkgver.tar.gz"

build() {
        cd "$srcdir"/gpbuttond-v"$pkgver"
        rm go.mod
        go mod init main
        go mod tidy
        go build gpbuttond.go
}

package() {
        cd "$srcdir"/gpbuttond-v"$pkgver"
        install -Dm755 gpbuttond ${pkgdir}/usr/bin/gpbuttond
}

sha512sums="CHANGE-ME  gpbuttond-v$pkgver.tar.gz"
