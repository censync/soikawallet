## Dependencies

zbar  https://zbar.sourceforge.net/

### Ubuntu

```shell
sudo apt-get -y install zbar-tools
```
### RHEL
```shell
sudo yum -y install zbar.
```


### MacOS:

brew

```shell
brew install zbar
```

port

```shell
sudo port install zbar 
```

## Environments

| Variable                        | Type   | Description     |
|---------------------------------|--------|-----------------|
| SOIKAWALLET_MNEMONIC            | string | mnemonic phrase |
| SOIKAWALLET_MNEMONIC_PASSPHRASE | string | passphrase salt |
| SOIKAWALLET_AGREEMENT_ACCEPTED  | string | "true" or empty |