## SoikaWallet - MultiChain Hierarchical Deterministic Wallet

This is non-custodial, multi-chain secured opensource wallet, with AirGap support.

### Dependencies

zbar  https://zbar.sourceforge.net/

#### Ubuntu

```shell
sudo apt-get -y install zbar-tools
```

#### RHEL

```shell
sudo yum -y install zbar
```

#### MacOS:

brew

```shell
brew install zbar
```

port

```shell
sudo port install zbar 
```


#### Build dependencies

```
sudo apt-get -y install libzbar-dev
```

### Components

* core - core sdk
* tui - terminal UI

## License

The soikawallet library (i.e. all code outside of the `cmd` and `tui` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html),
also included in our repository in the `LICENSE.LGPLv3` file.

The soikawallet binaries and ui libraries (i.e. all code inside of the `cmd` and `tui` directory) are licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also
included in our repository in the `LICENSE.GPL3` file.