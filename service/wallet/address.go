package wallet

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/types"
)

type address struct {
	path *types.DerivationPath
	key  *types.ProtectedKey
	pub  *ecdsa.PublicKey
	addr string

	accountLabelIndex uint32
	addressLabelIndex uint32
	nodeIndex         uint32

	staticKey bool
	//
	// lastSync uint64
}

func (a *address) Address() string {
	return a.addr
}

func (a *address) Path() *types.DerivationPath {
	return a.path
}

func (a *address) IsExternal() bool {
	return a.path.Charge() == types.ChargeExternal
}

func (a *address) AddressIndex() types.AddressIndex {
	return a.path.AddressIndex()
}
func (a *address) IsHardenedAddress() bool {
	return a.path.IsHardenedAddress()
}

func (a *address) CoinType() types.CoinType {
	return a.path.Coin()
}

func (a *address) Account() types.AccountIndex {
	return a.path.Account()
}

func (s *Wallet) addAddress(path *types.DerivationPath) (addr *address, err error) {
	if s.bip44Key == nil {
		return nil, errors.New("BIP-44 key is not set")
	}

	if !types.IsCoinExists(path.Coin()) {
		return nil, errors.New("coin is not exists in SLIP-44 list")
	}

	if _, ok := s.addresses[path.String()]; ok {
		return nil, errors.New("addr already exists")
	}

	// Create addr

	deriveKey, err := s.addressKey(path)
	var (
		key *hdkeychain.ExtendedKey
	)

	if path.IsHardenedAddress() {
		key, err = deriveKey.Derive(hardenedKeyStart + path.AddressIndex().Index)
	} else {
		key, err = deriveKey.Derive(path.AddressIndex().Index)
	}

	if err != nil {
		return nil, errors.New("cannot create addr key")
	}

	ecKey, err := key.ECPrivKey()

	if err != nil {
		// cannot convert to btcec ECPrivKey
		return nil, errors.New("cannot create addr key")
	}

	pubKey := ecKey.ToECDSA().Public().(*ecdsa.PublicKey)

	ctx := types.NewRPCContext(path.Coin(), 0)
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}

	addr = &address{
		path: path,
		key:  types.NewProtectedKey(ecKey.ToECDSA()),
		pub:  pubKey,
		addr: provider.Address(pubKey), // TODO: Move addr marshaller from provider
	}

	s.addresses[path.String()] = addr

	return addr, nil
}

func (s *Wallet) address(path *types.DerivationPath) (*address, error) {
	address, ok := s.addresses[path.String()]
	if !ok {
		return nil, errors.New("addr is not found")
	}
	return address, nil
}

func (s *Wallet) addressKey(path *types.DerivationPath) (*hdkeychain.ExtendedKey, error) {
	if s.bip44Key == nil {
		return nil, errors.New("BIP-44 key is not set")
	}

	// m/44'/60'
	coinKey, err := s.bip44Key.Derive(hardenedKeyStart + uint32(path.Coin()))
	if err != nil {
		return nil, errors.New("cannot initialize coin key")
	}

	// m/44'/60'/0'
	accountKey, err := coinKey.Derive(hardenedKeyStart + uint32(path.Account()))
	if err != nil {
		return nil, errors.New("cannot initialize account key")
	}
	// m/44'/60'/0'/0
	chargeKey, err := accountKey.Derive(uint32(path.Charge()))
	if err != nil {
		return nil, errors.New("cannot initialize charge key")
	}
	return chargeKey, nil
}

func (s *Wallet) isAccountExists(coinType types.CoinType, accountIndex types.AccountIndex) bool {
	for _, address := range s.addresses {
		if address.CoinType() == coinType && address.Account() == accountIndex {
			return true
		}
	}
	return false
}

func (s *Wallet) AddAddresses(dto *dto.AddAddressesDTO) (addresses []*resp.AddressResponse, err error) {
	if len(dto.DerivationPaths) == 0 {
		return nil, errors.New("derivation paths is not set")
	}
	for i := range dto.DerivationPaths {
		dPath, err := types.ParsePath(dto.DerivationPaths[i])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("cannot parse derivation path: %s", err))
		}
		addr, err := s.addAddress(dPath)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("cannot create account: %s", err))
		}
		addresses = append(addresses, &resp.AddressResponse{
			Address:           addr.Address(),
			Path:              addr.Path().String(),
			IsExternal:        addr.IsExternal(),
			AddressIndex:      addr.AddressIndex(),
			IsHardenedAddress: addr.IsHardenedAddress(),
			CoinType:          addr.CoinType(),
			Account:           addr.Account(),
		})
	}
	return addresses, nil
}

func (s *Wallet) GetAddressesByAccount(dto *dto.GetAddressesByAccountDTO) []*resp.AddressResponse {
	var addresses []*resp.AddressResponse

	for _, addr := range s.addresses {
		if addr.Path().Coin() == types.CoinType(dto.CoinType) &&
			addr.Path().Account() == types.AccountIndex(dto.AccountIndex) {
			addresses = append(addresses, &resp.AddressResponse{
				Address:           addr.Address(),
				Path:              addr.Path().String(),
				IsExternal:        addr.IsExternal(),
				AddressIndex:      addr.AddressIndex(),
				IsHardenedAddress: addr.IsHardenedAddress(),
				CoinType:          addr.CoinType(),
				Account:           addr.Account(),
			})
		}
	}

	return addresses
}

func (s *Wallet) GetAllAddresses() []*resp.AddressResponse {
	var addresses []*resp.AddressResponse
	for _, addr := range s.addresses {
		addresses = append(addresses, &resp.AddressResponse{
			Address:           addr.Address(),
			Path:              addr.Path().String(),
			IsExternal:        addr.IsExternal(),
			AddressIndex:      addr.AddressIndex(),
			IsHardenedAddress: addr.IsHardenedAddress(),
			CoinType:          addr.CoinType(),
			Account:           addr.Account(),
		})
	}
	return addresses
}

func (s *Wallet) GetAddressTokensByPath(dto *dto.GetAddressTokensByPathDTO) (tokens map[string]float64, err error) {
	addrPath, err := types.ParsePath(dto.DerivationPath)
	if err != nil {
		return nil, err
	}

	addr, err := s.address(addrPath)

	if err != nil {
		return nil, err
	}

	ctx := types.NewRPCContext(addr.CoinType(), addr.nodeIndex, addr.Address())
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}

	balance, err := provider.GetBalance(ctx)

	if err != nil {
		return nil, err
	}
	return map[string]float64{
		s.getRPCProvider(addr.CoinType()).Currency(): balance,
	}, nil
}
