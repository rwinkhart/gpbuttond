# Maintainer: Randall Winkhart <idgr@tutanota.com>
pkgname=gpbuttond
pkgver=0.1.0
pkgrel=0
pkgdesc="Maps GPIO pins to keyboard keys"
url="https://github.com/rwinkhart/gpbuttond"
arch="armhf x86_64"
license="GPL-3"
options="!check"
depends=""
makedepends="go"
source="$pkgname-$pkgver.tar.gz::https://github.com/rwinkhart/gpbuttond/archive/refs/tags/0.1.0.tar.gz"

build() {
        cd "$srcdir"/gpbuttond-0.1.0
        rm go.mod
        go mod init main
        go mod tidy
        go build gpbuttond.go
}

package() {
        cd "$srcdir"/gpbuttond-0.1.0
        install -Dm755 gpbuttond ${pkgdir}/usr/bin/gpbuttond
}

sha512sums="405ddfcbe2d6e827bfc78d044926156012b853b96f557636acc0aeaae17ec1165dee68c43237ad89cbff198109b1918bca4f764033fd786d02d03f0f29127379  gpbuttond-0.1.0.tar.gz"
