// Copyright 2023 The soikawallet Authors
// This file is part of the soikawallet library.
//
// The soikawallet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The soikawallet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

const testPassphrase = `TREZOR`

type testAddress struct {
	// NetworkType   string
	Path       string
	Address    string
	PublicKey  string
	PrivateKey string
}

type testVector struct {
	Entropy   string
	Mnemonic  string
	Seed      string
	Key       string
	Addresses []testAddress
}

var testVectors = []testVector{
	{
		"00000000000000000000000000000000",
		"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		"c55257c360c07c72029aebc1b53c05ed0362ada38ead3e3e9efa3708e53495531f09a6987599d18264c1e1c92f2cf141630c7a3c4ab7c81b2f001698e7463b04",
		"xprv9s21ZrQH143K3h3fDYiay8mocZ3afhfULfb5GX8kCBdno77K4HiA15Tg23wpbeF1pLfs1c5SPmYHrEpTuuRhxMwvKDwqdKiGJS9XFKzUsAF",
		[]testAddress{
			// ETH
			{
				//	"",
				"m/44'/60'/0'/0/0",
				"0x9c32F71D4DB8Fb9e1A58B0a80dF79935e7256FA6",
				"0x03986dee3b8afe24cb8ccb2ac23dac3f8c43d22850d14b809b26d6b8aa5a1f4778",
				"0x62f1d86b246c81bdd8f6c166d56896a4a5e1eddbcaebe06480e5c0bc74c28224",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x7AF7283bd1462C3b957e8FAc28Dc19cBbF2FAdfe",
				"0x03462e7b95dab24fe8a57ac897d9026545ec4327c9c5e4a772e5d14cc5422f9489",
				"0x49ee230b1605382ac1c40079191bca937fc30e8c2fa845b7de27a96ffcc4ddbf",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TAyDUYP5rcf56xFwrg8cU1qQwvnWpkeapM",
				"02111c4267dff570eaf2d2be12e0bdb9e3a7dbdb8a765ba86bf4879aa2a5595964",
				"554d613c6ae7cfe1f7cc0814f48e8eab176ca316fd7d1153fcd7a45b73fee11e",
			},
			{
				"m/44'/195'/0'/0/1",
				"TFYaFtDHPwqLuo64eh4CYJiwHT7UivRmst",
				"03f7bc5cef7611e532c2f8390a6c5621e84f3d8d78e9b1ce2ef2db65b7dcc97696",
				"ca62831c3acc50f5f37b9f03b718145b84b583cfdba3a47c5d76e95451366f70",
			},
		},
	},
	{
		"7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
		"legal winner thank year wave sausage worth useful legal winner thank yellow",
		"2e8905819b8723fe2c1d161860e5ee1830318dbf49a83bd451cfb8440c28bd6fa457fe1296106559a3c80937a1c1069be3a3a5bd381ee6260e8d9739fce1f607",
		"xprv9s21ZrQH143K2gA81bYFHqU68xz1cX2APaSq5tt6MFSLeXnCKV1RVUJt9FWNTbrrryem4ZckN8k4Ls1H6nwdvDTvnV7zEXs2HgPezuVccsq",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0x6006ef1944FB519A746d00cDAf715Cbd27a5a008",
				"0x035422486d29f5189ce7e606252d96d81fa446dc8bb5a6221c307e061c20e3089a",
				"0x9f20bfeef91877e3c5f50fc0557a80d25f77a650c83d47601a8193bccb0e678a",
			},
			{
				"m/44'/60'/0'/0/1",
				"0xA5Cd5A7a353DAD72643512f25F395eC1857Bc1CA",
				"0x02e5ee8903ddfc552c92ea484378fc860b39917639ab2af4e901e3e8e24471eca6",
				"0x84ca2976ffb4a13ace73127577d82aba99ee58c8b5eb59acea74d4a05beef97c",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TLKvhyqakrxrN5SYXFpLETSyHvy4QMS84L",
				"0245328a67ebb51164cf72fc9e17783ec09666bf5d41f68f086154bc5f881f6772",
				"eb4a7b73b297017001e6dea04b3761cf3375c891c6c41f3805c376952d9ce92c",
			},
			{
				"m/44'/195'/0'/0/1",
				"TQeMnrxMgCJmgkzdH4NXA3GPJkLiWCooNB",
				"0261ed8a578afa71c7d9c0555ec45dcaea6020383c92101a4522756a64c6ef698b",
				"11689ced259cc68a4e584d12dc67fd4320452685b99ec37a76f2aa8a2f33a274",
			},
		},
	},
	{
		"80808080808080808080808080808080",
		"letter advice cage absurd amount doctor acoustic avoid letter advice cage above",
		"d71de856f81a8acc65e6fc851a38d4d7ec216fd0796d0a6827a3ad6ed5511a30fa280f12eb2e47ed2ac03b5c462a0358d18d69fe4f985ec81778c1b370b652a8",
		"xprv9s21ZrQH143K2shfP28KM3nr5Ap1SXjz8gc2rAqqMEynmjt6o1qboCDpxckqXavCwdnYds6yBHZGKHv7ef2eTXy461PXUjBFQg6PrwY4Gzq",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0x97aa6F4c3e3120E25Ad2Ad3b88E6C13EF21ACE4a",
				"0x03f5eaa5038c7b94282f6d7d292cd12e8e7432f60447f8615982284140fe690f06",
				"0xfb5f0fd04b28b1b8ff80857019b6b1445133f15aad015461fae895bad687ae0e",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x29553c266cc50DBB07bA3E4Faea9405C2C8c731d",
				"0x03f06d830ce059286fcdcd3239b20c2b1cc776b8532795897922da1145215996f7",
				"0x92c795ab3b73c2f9394dcaec49852b22a6946f12f325839d6cd137d02832ba92",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TPGHUwerX55Wbjz5DasnqGHUZBXjD1jmfN",
				"026eb6009682a5114219b3de2a58069b22660697447159c3b8fa243fddf2489982",
				"85b84717ad87341308b83af21383b71d782cd6b2aff86e637d2e0f2fe5a2d923",
			},
			{
				"m/44'/195'/0'/0/1",
				"TWBkEi2CvkZJKUAnNv5Jr2ygPiTF1uyLfh",
				"02818d061a42acbbed66b7a498b14d10a34e63c61344317860aef50850b0c6148b",
				"42940db4e80c949f7eb9f55b2d3db1a627bd361a024a60fb700738fe310f2414",
			},
		},
	},
	{
		"ffffffffffffffffffffffffffffffff",
		"zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo wrong",
		"ac27495480225222079d7be181583751e86f571027b0497b5b5d11218e0a8a13332572917f0f8e5a589620c6f15b11c61dee327651a14c34e18231052e48c069",
		"xprv9s21ZrQH143K2V4oox4M8Zmhi2Fjx5XK4Lf7GKRvPSgydU3mjZuKGCTg7UPiBUD7ydVPvSLtg9hjp7MQTYsW67rZHAXeccqYqrsx8LcXnyd",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0x31a2E8fD7fC06AED00565458566993511c3F2d37",
				"0x03be2c1471c764fef1899a951ec4ad1b00c17d61a84e777aed4323f03dc8bab322",
				"0x4f54c1357bcd855f1ed4315e66dd0771c83318c6726d3c1847510825a72a38e0",
			},
			{
				"m/44'/60'/0'/0/1",
				"0xdcb237346f3Fb099a49d68d9eF11F7bAE2f23052",
				"0x038897ef9101407c48990352e7dbfb4fd142dc79a1a280de5ea12a17fe7a64556c",
				"0x6efa8fa293c1ff8a0245c9bb19ab2dbc800a107ddb78d53fb6ddd1531d12a634",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TH8iQevWUUfewzrJ1XH6khH1Gudzx3rP4D",
				"03e25ea095184f5453f98a71faa257c5e57b0ce207ade18be5d4e2e577b1c800b4",
				"e5ec48c956929c290e1d01adec06f1212fcf0064bf2a9e8eef92dfeea97f6a12",
			},
			{
				"m/44'/195'/0'/0/1",
				"TErYfyr5SqGfBSppidJYBFyqUU44tDVjEn",
				"032d4d7f202964e923dd69f0cb04fe26163d9a3ee522ee46898b3a2db6832b97f9",
				"1daebf7b773a0cf5e56fbb2d63d69f3f6be6c30889392c396db0662899812d0b",
			},
		},
	},
	{
		"000000000000000000000000000000000000000000000000",
		"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon agent",
		"035895f2f481b1b0f01fcf8c289c794660b289981a78f8106447707fdd9666ca06da5a9a565181599b79f53b844d8a71dd9f439c52a3d7b3e8a79c906ac845fa",
		"xprv9s21ZrQH143K3mEDrypcZ2usWqFgzKB6jBBx9B6GfC7fu26X6hPRzVjzkqkPvDqp6g5eypdk6cyhGnBngbjeHTe4LsuLG1cCmKJka5SMkmU",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0x8e5713dC3Fdf4812957924Bd7976907DC455FC42",
				"0x02ec813ca13f55a1cddceadd8b3a01029f5894017ad8fd464761ccba7a9007e089",
				"0xb5cf2b09ca6fa2b8fa3b281348b25b68d224981840b93c6cc15920e205734d6e",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x4Aa575D7318EA222FB8F9ee71BFBe921bA15c38b",
				"0x0307ecf37a7834691fe5929b1e3cebad15d63b3fe83daa8a38244ec00f0196dee2",
				"0x3b7c2a71b187de3f6f9e9cfc3fddf0a9ed0b4938ee527c9b96496628f22c3590",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TLB5ZBNnFcshompAS5EkmJanobwoyC6Gui",
				"03fbfe485b44fbe67bc4c83a0821ef0045ded6740855fab11c6dbca383a50a3c03",
				"a91fe77b341feb58d58026dc381ccb407b5ec9c1872a7e0e4cdb9c33d36f2ae4",
			},
			{
				"m/44'/195'/0'/0/1",
				"TUBaCmMZejWSCQ6eTMJfq4pD6a61phDGrZ",
				"03cc30ff40107708adb10a9ae69b31f34117d5a1f68ae47984822cbf23a56b75dd",
				"3c56ec373680217a4c5320612cd57c0dda0d21a356f9433d8196c42a778bfb64",
			},
		},
	},
	{
		"7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
		"legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal will",
		"f2b94508732bcbacbcc020faefecfc89feafa6649a5491b8c952cede496c214a0c7b3c392d168748f2d4a612bada0753b52a1c7ac53c1e93abd5c6320b9e95dd",
		"xprv9s21ZrQH143K3Lv9MZLj16np5GzLe7tDKQfVusBni7toqJGcnKRtHSxUwbKUyUWiwpK55g1DUSsw76TF1T93VT4gz4wt5RM23pkaQLnvBh7",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0x2DDC23D8400729462b61b2e99eA87DefBB678f58",
				"0x023e0873be61174c75a90ab0c3b00eb1d52053a300e2cd4aa5d3e9d407cf5ebd8a",
				"0x3e89861b314309a79a1d1e8ffad279fc3eefaf2eab45e5215ab7127b0a18e4c0",
			},
			{
				"m/44'/60'/0'/0/1",
				"0xDa566925c7FD47Fe3B4306103B6285FB15A63558",
				"0x03a612da73193598ac4e6b43fba8c1a61c9f5cfb8a6ef7faedbfd20deba35f367f",
				"0xe8efc1c95336bea28f92d873ea055582acb8ca8278035a8ba034d67e906f9a28",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TJFYseFBcvuRThTHJYeJXzF4zEicgx4Cov",
				"02ab5616e11b1eb802b3f75aabf1d972b5ace73bf8854a50977caa8adb1a1a9a65",
				"1c35ef1e974cc1031a86399c94d4743cea4f10fda8b2b83a89aebbf0dd71224b",
			},
			{
				"m/44'/195'/0'/0/1",
				"TNK168A3QxcX57cN9rXEVxUJhELN7bjoqe",
				"02c79aceaf05ec4d6a265b3d0e694c121f86afeb90587835d434586b43d56e7e8c",
				"2e2a72074253f76890691cb16219381ffd2670cbba91b39fa9caaa64ead5256f",
			},
		},
	},
	{
		"808080808080808080808080808080808080808080808080",
		"letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter always",
		"107d7c02a5aa6f38c58083ff74f04c607c2d2c0ecc55501dadd72d025b751bc27fe913ffb796f841c49b1d33b610cf0e91d3aa239027f5e99fe4ce9e5088cd65",
		"xprv9s21ZrQH143K3VPCbxbUtpkh9pRG371UCLDz3BjceqP1jz7XZsQ5EnNkYAEkfeZp62cDNj13ZTEVG1TEro9sZ9grfRmcYWLBhCocViKEJae",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0xa65d92F3c3a91B852b913c4539868ceafF8d6B38",
				"0x02bb90faf5e2276091b8c2a6a38c8696fb98c44a729e573f8832dd177f6dc79d6b",
				"0x17803dd9c8a75b0a4c5824d0e42102d58fbc5e8db0e1c9f4b267c25c6a1116c1",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x6FfD98144BC577215767FcE799d67bC4502cb126",
				"0x02a43785c7db41ba64d8a2570f44bf4322df98b24ea4119b6d2375881af3fc8b49",
				"0x00eae8070dc5738770ae690d8a68ca74ee13158ff62479b91cbfeb776e20ea38",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TGfm254jWMdadDriEbmEN3QMznDAD7p5Ao",
				"03212d74ef1422f238ce843d259e1c6a6aa212d0c29731f16e2884ce565c997c5b",
				"b8bf05013d2d7948039a11566d522ae6904ce36621e3f9ffd6efc13d57b21b89",
			},
			{
				"m/44'/195'/0'/0/1",
				"TLEMisutmNNoFLN7Lqp3GXtTyKSJbAhnKW",
				"0224a389d6eaee1e0fae0925a98ff58da0b69837d266bd85953444c17379bdbf5c",
				"93948201aab25ae922d524164618632d0a80e17f4ff0747f8e951ac2481b40a9",
			},
		},
	},
	{
		"ffffffffffffffffffffffffffffffffffffffffffffffff",
		"zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo when",
		"0cd6e5d827bb62eb8fc1e262254223817fd068a74b5b449cc2f667c3f1f985a76379b43348d952e2265b4cd129090758b3e3c2c49103b5051aac2eaeb890a528",
		"xprv9s21ZrQH143K36Ao5jHRVhFGDbLP6FCx8BEEmpru77ef3bmA928BxsqvVM27WnvvyfWywiFN8K6yToqMaGYfzS6Db1EHAXT5TuyCLBXUfdm",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0xbeB3313714A8afc54460678B2CBf4F1C43520f03",
				"0x02b51c611cb347b685daa79125458f63955ef19a6963e412f97a0b693498d0093d",
				"0x9495bfb328cdf12965e6fc11758c7a91688e57b922d1d4e2b8945944eb203bd9",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x177975E97f12e932a2F969bF86d49B16588b6a42",
				"0x03b5c90a7256e7f75e59fac7ff9084516a4a2dccfe620c10a8a725f9986118f3b5",
				"0x849bf5d4c149a4931bd6767c6d637b22fb209bb242aa59ffe902b279430daa06",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TDhbo6L48LnXVGzeHgzNEvK7wnfkJLLHrF",
				"02ccd99f9c9b91ec7bcfcb375637b0b33829ca3baa499e04febd56b5c9b9e1d608",
				"f276c34ee85d31f1be27305ec3a43fb97e2f2755b6d08e4e795959b6821174ea",
			},
			{
				"m/44'/195'/0'/0/1",
				"TRorGhwNgiscTa7Jk7x1AMG53vjrjFzCZj",
				"02d16dcd3fba5e02cd145cdb0ab36d3ae0e3e1e626c576847c4786a65ba382b294",
				"5af0248eb09f0795795407233807ece83defa28721b10a6e49a6ab2c1197f35c",
			},
		},
	},
	{
		"8080808080808080808080808080808080808080808080808080808080808080",
		"letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic bless",
		"c0c519bd0e91a2ed54357d9d1ebef6f5af218a153624cf4f2da911a0ed8f7a09e2ef61af0aca007096df430022f7a2b6fb91661a9589097069720d015e4e982f",
		"xprv9s21ZrQH143K3CSnQNYC3MqAAqHwxeTLhDbhF43A4ss4ciWNmCY9zQGvAKUSqVUf2vPHBTSE1rB2pg4avopqSiLVzXEU8KziNnVPauTqLRo",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0xd74F28D86e9Bf32cc2e51D18e26A5ed6446a22bb",
				"0x02547e34dc79171f3b7d4a4bef4de92b69816e0bfa3009cd94819958397354928b",
				"0xef8e951e809b4396060fe4eb4b317f9f02bbb54ba32821b62d0b4aada85c03ec",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x5fA1d76C2C093389bb6187Df4D10D92A4FF9E9D7",
				"0x034c091f71f6230b15756295cba58316eed2c28308b1e839fedf8d54882b57abf6",
				"0xe08e8e9a7fddb3a8ff15549757f0754bd3d70169030aeafd574fd3af551eed7f",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TDSHEi5BcZCHrjYpUJ1gYK6QwYTT5BzDzS",
				"0225ca784f53ac258c4b333ee4eebbfb13b3dcd0ca1bb2e9df702cf7dac6246660",
				"a5a700705c9f75c905f8c8740ab1b87c1aa9e76dcb9dffb5fa09066621b5108b",
			},
			{
				"m/44'/195'/0'/0/1",
				"TDeaqC2BJ8YUAoWsyFjAcKq2mGykDh4Vw9",
				"02648d1cd3bd1cb8537106953705e0e9e277e0bd50795ca41fe5dc3abb233b9d5e",
				"a519a4db26ac1e6cf1a39f90cd00fddfae14e8da93fea121ee2c9fcec771f638",
			},
		},
	},
	{
		"77c2b00716cec7213839159e404db50d",
		"jelly better achieve collect unaware mountain thought cargo oxygen act hood bridge",
		"b5b6d0127db1a9d2226af0c3346031d77af31e918dba64287a1b44b8ebf63cdd52676f672a290aae502472cf2d602c051f3e6f18055e84e4c43897fc4e51a6ff",
		"xprv9s21ZrQH143K3xC5SRKnxV4R829AcnKE7XjZu2PixyZh3CexnsvmkBsi5rzqXMhxTkfLJFB6FuHJPWxxvcH5eYvCDvWcYAMXpbpGGiVUDfH",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0x8C9c6BE78C21DbbC26C645413380AEc2d500613A",
				"0x02493304320df5d3eeae065fe42d99d10a17c6dd2fa04770d325384faf4a8a83b5",
				"0x65fc30b50feaa94d72e130efd2bfde4d257ac6c38b8e098b8357f13db47a2bd6",
			},
			{
				"m/44'/60'/0'/0/1",
				"0xdC06F56E9fC3DeB743BCB9255Da6497E50Ec8a8d",
				"0x027e3d215eb08640d0e2c834784de9c1d8cb0bfc87c2192dc861e030afa88bde5b",
				"0x68ec263a6f9b71d18917d397718879547f704962052239756fe0d4e7500b4243",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TUxxKf3nCXBqtGCRh1gx3nonvh2Act4GzE",
				"0258bfde4e7c3aeed7724f7567b07e52396527e01e0725113b38171fab02d61e08",
				"7b135b8cb0e169d01f050e7581b6187c9b3435f1ea40ff8a1ceba05292498215",
			},
			{
				"m/44'/195'/0'/0/1",
				"TBL4LWeTqe63Xm5LBvZJrse1FnTiULBx2G",
				"032de5375013e26878e9aa431b05cac2139b37f3bd0ce43e9d801f537b7d5fe3d5",
				"4eac8cd8d2df0c9a72da6fb1202c89174aee29f3748f63db3ad0e6e2ab917893",
			},
		},
	},
	{
		"b63a9c59a6e641f288ebc103017f1da9f8290b3da6bdef7b",
		"renew stay biology evidence goat welcome casual join adapt armor shuffle fault little machine walk stumble urge swap",
		"9248d83e06f4cd98debf5b6f010542760df925ce46cf38a1bdb4e4de7d21f5c39366941c69e1bdbf2966e0f6e6dbece898a0e2f0a4c2b3e640953dfe8b7bbdc5",
		"xprv9s21ZrQH143K4YsWLquHbdGRh1mRrj5DTdRaj1cUrhftfXx4YJ3Zy41H52GR8nywKkpRSTdM71uZTRztscUdAAPL2Z6JQdW4xVPyzxh5zCG",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0xF1f1112E5a60db67B0e2c14f087cD996Dc3417FB",
				"0x0257a233f54c7d77e585dfc576dfb61a7ead8f28fb73492a101da23efc299222cf",
				"0x25473dc191c5b5b27238bab0f6c4eda434e41914edd419ff6d1e7c115cc49284",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x1BcBee484FBef93570f5ea72c306488CD2E434e0",
				"0x031401b567ba855e5cc54b7adadb9cd473d738b1748b96b90f5c391c47af52cbe0",
				"0xfa9eff4470e07d90959eae8289626ff0d0517cd22339aac5b0fddf5dc40130ab",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TJANikR2HKhkD7u6z7vUMbGfQTNFa9eHKh",
				"03da7c9309c84ee18de685ef1fa9082da5fe057dd84414d0e9c56e30daa1176d4a",
				"0ce3e356eac4a0e5335feaffded8697eb92e281236da0774443d687edc45f988",
			},
			{
				"m/44'/195'/0'/0/1",
				"TUdES57Nr1uMwPBno1zwe3vjH8coaEvF1v",
				"02f7c358c7698a535facde6ac7f964e7859d688de83817ed49998551b151c0a9a7",
				"0dd45692f7a75c8495bb83ade50f6525765f05f4d07e10d39354750e3ba6f411",
			},
		},
	},
	{
		"3e141609b97933b66a060dcddc71fad1d91677db872031e85f4c015c5e7e8982",
		"dignity pass list indicate nasty swamp pool script soccer toe leaf photo multiply desk host tomato cradle drill spread actor shine dismiss champion exotic",
		"ff7f3184df8696d8bef94b6c03114dbee0ef89ff938712301d27ed8336ca89ef9635da20af07d4175f2bf5f3de130f39c9d9e8dd0472489c19b1a020a940da67",
		"xprv9s21ZrQH143K44JSkE9N3huFVGqK5YUroYxjd5eBotHvyBcDXNvjF3uxSiGDuGo7ub2GJgLc3HtvQbQkzies2qNjeM1p8nyTWEzNHuyVqss",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0xcE2F6278C149984B3Fa258892527Aa6FBFA04863",
				"0x027ce6133b2d639929a9029771fd2a051e9cb758b353d519ac945b4b26742eb9b0",
				"0x1ce53139bd144f1cb8ef220d76e1c8f45061d8844fef02b5deb330f47f4ad2ff",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x41c853BE861aBd6e67fD63A11131198E29e3b0B6",
				"0x02bb5f8962bcf384f13871b5d42ea3127d125d0668d404174fcc6a417d9b3e2de3",
				"0x409da55a2ae1f4cf15e3b98c819ff0328df3350b69ccf7115488871dc03198ba",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TEgYRkAipsumtNNqtL1zxVHNwAtwbuL7Qx",
				"02310f6d6bf0f8b7183b9774bbee36872d9681a94a117d027f4bca5aecad7f117c",
				"22db64b8a3f7527ce312840b53cf0e50be175a2cd4bfc2421c5d6772507c04f8",
			},
			{
				"m/44'/195'/0'/0/1",
				"TGXxpSysxBXk8PQQiacju1nFAgVp1Bp585",
				"0272b2075b91191d6d318ea92cd8ad19e731c0fe899b1e2de68c0913876759bd6b",
				"f26ede338a1fe2a17fce2bd701e0e3c79db74d39305da979d1021644b01187f8",
			},
		},
	},
	{
		"0460ef47585604c5660618db2e6a7e7f",
		"afford alter spike radar gate glance object seek swamp infant panel yellow",
		"65f93a9f36b6c85cbe634ffc1f99f2b82cbb10b31edc7f087b4f6cb9e976e9faf76ff41f8f27c99afdf38f7a303ba1136ee48a4c1e7fcd3dba7aa876113a36e4",
		"xprv9s21ZrQH143K2fzHWz7Z7PQj54R9Acrra9W28nnMLzgHonTebXnRD35dmvyaB41A1U1o59duUJ7dF9227Hr84AFY8aAeGNhnetXuecd6t67",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0xf6044F898DE5e816db2Eb9dC1Ba8D2cB58EE5e20",
				"0x029f679a8b7a853c41d2fbc4c5de82a3b967eca52f2d383f1e48f3e7f1a8ed019e",
				"0x37cd606ba3780f6c12eca043486370d9b71cee2536845757f8614e0fe1a2a86b",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x92bafb3ec2aDf822c2Dc1D442690a7AC80f89088",
				"0x030d6f1a6abdf363ad8990daafca2fe728aaecb2c46d087f180def0c048fe36065",
				"0x339149c7ea6df799157a42f8971d1f0bb35d228bbe234bd1ee4cefe5d506b470",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TAnyEbYYV1igmT5ciKNrEEvSVVniGCsGqm",
				"03c8fef46bf6381769820796a24499313dede05fd9831e4479263326160d3bd71d",
				"4c859d233483ac45c3faceca38a87b05ddebbc42971bcb0fd6c1477a2bc66c65",
			},
			{
				"m/44'/195'/0'/0/1",
				"TXNrHe6gK9uzfBrAg9VkN62B46m37kubVF",
				"03d4f8e57cdeb59892329fe939f3da3cdcc6bbbab96c1b2a9930d2c39dbaeb48e0",
				"11a090e5a8c8a32a768cf35bbade8f8711cb78fbab1139ce644ef698ef8eb398",
			},
		},
	},
	{
		"72f60ebac5dd8add8d2a25a797102c3ce21bc029c200076f",
		"indicate race push merry suffer human cruise dwarf pole review arch keep canvas theme poem divorce alter left",
		"3bbf9daa0dfad8229786ace5ddb4e00fa98a044ae4c4975ffd5e094dba9e0bb289349dbe2091761f30f382d4e35c4a670ee8ab50758d2c55881be69e327117ba",
		"xprv9s21ZrQH143K2gts9Sq6Aq67GTVeWXuJM1Eieknp95mWujAcuD2VixUsqaRuU9Hm3Z7Rh9JzukebGqwfbu6gJv42KRBvK4f4K9Cc84r7jaB",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0x408bDc5404e3262a610900DE18df5aF685aEb33B",
				"0x026c1a507fbead220e55aaf2c79b5f6f6c6b41a5748ccbd8f8b759ec3d55b1190d",
				"0x58c5b59121e818bcf1e0b5f5477ede01376fffc846c711fbf7c71252a428b6fd",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x0E1fA3bCD673720D07da7D5d041895e8707AdEc3",
				"0x0218184e75be7d9284746db584da628a303ac0af87f626f09cfce4fb2d4e2c8f95",
				"0x6377c7ffcdb2e9fe77bc66c055c7759fda21b4a397f521782bbc0a510bf63e9c",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TYe8vC776NCwQtHFnsborWfeZG9idmf4f6",
				"020d7442c7348f31789f71ecf9733613c5464ddaac7485aa913adcdb1cbb7bf703",
				"d125feaed57e7a1f0c781d628da6f3bb4cde554fa2d71489b3ae8ffe04925bb1",
			},
			{
				"m/44'/195'/0'/0/1",
				"TB2pdaszLPfHnetBqxWX6pcZfj1ad6oV1S",
				"02fb0bd59d8b36127b2d46b886bf0d4f2bd93945d8a9cb65d7d51d8add1949bbae",
				"cffbe0d9004b0d491d5c5382e3c8b4b6190819120fd1b91d910d7ee1900b70fb",
			},
		},
	},
	{
		"2c85efc7f24ee4573d2b81a6ec66cee209b2dcbd09d8eddc51e0215b0b68e416",
		"clutch control vehicle tonight unusual clog visa ice plunge glimpse recipe series open hour vintage deposit universe tip job dress radar refuse motion taste",
		"fe908f96f46668b2d5b37d82f558c77ed0d69dd0e7e043a5b0511c48c2f1064694a956f86360c93dd04052a8899497ce9e985ebe0c8c52b955e6ae86d4ff4449",
		"xprv9s21ZrQH143K39y7KHx56XraMbqrS7VBxVqSSCUhFvE8MsaBCr9T7zsZwNH7jvdcii9ToB91qvgeacds6ubaNU3TDxvY2bhZMmESAAssoYD",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0x080F2A8a667348d7dD52f50D6D92E93E2fA4A175",
				"0x02cca430c12832f30a3addd7554ff90bde7d8238bdf6e8870996f9079b4f656931",
				"0x9563b492d682cd06f2c18929dd502455c6853fc55e53032345f188dd76431113",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x950847D9B84ebF64C80B7bc9cc2C5aB24407AcEE",
				"0x02f162a17d85d1fe7cadd1305241c6507f498d0c67ebade6cea3506730325a7e70",
				"0xac5a3fc208526385aaff52d4f1dddadc4f15ada3da96c43e42619efd5e52fd2b",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TWZEdf97EjfeMiAuNkBR2vMZLcjw9Bqfr3",
				"02c169f5c6a4ebed7d659d5836ac5756b4775e9805a3b6b041c3eff978f40b547f",
				"722b15bda556b87a21e16657c645884b7343f7618582b093df7bb60a05f0ff6d",
			},
			{
				"m/44'/195'/0'/0/1",
				"TYG1shR1bmYmBseVYhRGry9zdboVC7ifNY",
				"0267de6b2588170c79d2571fbcf11c9d7acb02621e233a2c8f970cda3e55e8fedb",
				"b159510a1398ae5b2a7a8d1545a0005282a4d3cff10acfccc5c2285169863165",
			},
		},
	},
	{
		"eaebabb2383351fd31d703840b32e9e2",
		"turtle front uncle idea crush write shrug there lottery flower risk shell",
		"bdfb76a0759f301b0b899a1e3985227e53b3f51e67e3f2a65363caedf3e32fde42a66c404f18d7b05818c95ef3ca1e5146646856c461c073169467511680876c",
		"xprv9s21ZrQH143K2mweKbPaebAU2b8poVVeqRgi1UBPybm9pLoCRKGgFgD2LbLHvHNsXDk3n1zjT7RujoLyb9huymgMXZLtL2UWqBHgKxdTjFk",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0x01A402776891e9A1C90676B6B36fd5C65886Dec4",
				"0x02a883e2e12a8ae6228ec0ff929863994599b22aea3a2bc8dbadcd9f4e20caf2bc",
				"0xe8fe9b07820aef9e9f0864031059f6b9ae47bdcfd2a9e61ea360296bdc855a11",
			},
			{
				"m/44'/60'/0'/0/1",
				"0xdAC02EF342fC315dddAFB4BCB17A24EABfD87867",
				"0x0262d3e56db390e1ed1e33972d492e623dec99f22a2cadb6eb4b6525cab7dfe36f",
				"0x894f02288ed0061813a79efc7c2cfbde58c1f95d96873f50dfc65a059e56c83d",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TJvC5sn92tRJfrXk1g3BDf7toceZ2tvASH",
				"03ad7e93d6025e08d02c42ac5d6817a4135a64472634a7ddcdb7bbd6336bb32403",
				"29676c5ed7058ac88816926113e6fe5f99eb7521db1ffa2ec0e71935f6c007f4",
			},
			{
				"m/44'/195'/0'/0/1",
				"TQQfmW98kP3TBjHXZby84iDP4r9PbqmPwR",
				"03d3c6fc8813487fb134c1dbca0773d12a45a04e1eab2f81630324914ad08ca350",
				"8dfc2014b96315c591d9194ccbddc901fcf5ffe94f35094b3c295beaaed3d400",
			},
		},
	},
	{
		"7ac45cfe7722ee6c7ba84fbc2d5bd61b45cb2fe5eb65aa78",
		"kiss carry display unusual confirm curtain upgrade antique rotate hello void custom frequent obey nut hole price segment",
		"ed56ff6c833c07982eb7119a8f48fd363c4a9b1601cd2de736b01045c5eb8ab4f57b079403485d1c4924f0790dc10a971763337cb9f9c62226f64fff26397c79",
		"xprv9s21ZrQH143K4M1N4f2Ma5YRADyBqU7wtb18qiZwWTk1rpx49XTsRCUa2iaPhDRBEVAMdGqDCn5iJTvsAUrPQ8NhVYdwZSf5mekdqwcRUS9",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0xebb75c5f4bC6e89B0863627878Ab1bA1F7bCD1Bb",
				"0x026ed83c5748119c6da5e6d3fd004e22d22f9648bf42b78857295eca30a1472f59",
				"0x62f900d427f794a1ef1b79e248bc94e13b237dd231eb28492d69f416621dd5c3",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x41fd23d823FA5a3381BC249931134a07d6CC9b56",
				"0x02e87004df0bba2b88fb1691663fb885dbb3c148a29d73af4f98efb9f89df54a52",
				"0x216ea65a406b0bc2a7f4454be98537abf8317f925332921bb4ad4215e6c25bf7",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TQEaw8wT82G4QXBmAGxQWCFsWUGGj9q3FM",
				"03dd33a368f45bbe1ac4578bace720f129cf54a915585d2f62d8787b804afe4591",
				"27fedc5ea8494e6d50d515e6c799940b4e857ddb850b77989a45cac0c6f1a0cf",
			},
			{
				"m/44'/195'/0'/0/1",
				"TEJB7gYDrDfzyBM9revQgpKm1TdFh1ULDq",
				"0351e2c91c36daccbdf8d866983b148ac1bce62c0745174321c900e875e1bd213d",
				"4a67b7bbcee72c2a77f2afb4c03b4c38994039c8f6beebcff30294036b1c9ac1",
			},
		},
	},
	{
		"4fa1a8bc3e6d80ee1316050e862c1812031493212b7ec3f3bb1b08f168cabeef",
		"exile ask congress lamp submit jacket era scheme attend cousin alcohol catch course end lucky hurt sentence oven short ball bird grab wing top",
		"095ee6f817b4c2cb30a5a797360a81a40ab0f9a4e25ecd672a3f58a0b5ba0687c096a6b14d2c0deb3bdefce4f61d01ae07417d502429352e27695163f7447a8c",
		"xprv9s21ZrQH143K3BDzEvudRjun23x1nqxchPCmyTsRBNmUZwFP6Hsim6UnwpcEA6De2kVpC6UDoVKUFFh9h47cY4DL5363KwwvGQ3jVzU6rXP",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0xBe9F8038EDA2B3dB25C9fEA205365b0430cdF1e2",
				"0x036468766e3ca23aad5030b35838e588ed6760212c8423794ca95f2a773e668e98",
				"0x145c4b827cc4aad071c5a23ba5986016dbeac6ef1b0c0c455330fa99fb01e134",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x18194C3991867A35298EC7f67BdCE83BDB04815f",
				"0x034bb6a626b28d7c646a2ba95424eb3a58de10c5e39db3a8026bf44d897c76dcb9",
				"0x0972234bb7312d5a3c476c65376c1968caf0f93a1f98dc845498b6fd7ec567ce",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TBQEv4FTiLDHkuJ5gYovm9w1qD3v67hS1i",
				"034195c16e12ce54096af57cb8faf1dc83aed3eb098b5b86b6846b104c43bb3c3b",
				"41a833a9ac762b594b9a26e413f61276f8099843d5bbddf4148da9ed95d04fd4",
			},
			{
				"m/44'/195'/0'/0/1",
				"TMhQodynqQP347k8P2BkScJ4Pie7zpDJ9k",
				"02b6b3d6ec174dafeb5a078dc4029f1c868f81e123a78f6bdad1462e32c6904845",
				"abde89ecb308ca28065d325a941dcbd203fc6f80e968b6184e4e6e51c56a31f6",
			},
		},
	},
	{
		"18ab19a9f54a9274f03e5209a2ac8a91",
		"board flee heavy tunnel powder denial science ski answer betray cargo cat",
		"6eff1bb21562918509c73cb990260db07c0ce34ff0e3cc4a8cb3276129fbcb300bddfe005831350efd633909f476c45c88253276d9fd0df6ef48609e8bb7dca8",
		"xprv9s21ZrQH143K2fopRUQMvgrFpXJHmAbGYfwdpKcRh9cp9E2aHDbQA5V9mXRwCRj2nzjwpAXH4sdhGV8xJxpv2BEZxEJrSDsdqwAYBXcF3eu",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0x2e70adaA04d3567497e3a905e535A3aB5259f6f6",
				"0x021cf8352406e2afd2851151bd6293f0ab66a27e65d4002ffd415b67ab8ae0d14b",
				"0x0dcf3a32ef096cb252f9e4b4e29ca535f7326b18b6979dc8d59b2ad747d0323c",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x03b4293852b3f727Acf7a22Ac91EFCD27515Bc11",
				"0x030a66aa0ec74692fbc50936d181c6d4dffbb7745d0b2fc187d12aada33b151469",
				"0xeeaa718b24a2e8fa7736582bd5150f3922dea0e0f0d92d820a6086ea4aa6935a",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TLvc53mFxQ7uVYWVEWgHdm6uqAea2gEPGA",
				"0321d491a6b0b5f90ae2f47c75b997e69dfea1c7959622cedfd0433cc44e4d877d",
				"2c0c9e17e066406c302ac10711d0c0ed8afb1518b6dd1c61ba97e821599b32fe",
			},
			{
				"m/44'/195'/0'/0/1",
				"TPX5h1gYFevwh4HSh6fCPZi9aBZd8FPQYK",
				"0227fd722534550c331d1688775ed7a9eb0fd8b1ea615b402d858477586c9089d8",
				"d9e78927d348a512e9e34a01f31cd8b9165ec9699ddaf8f141897d214e6783c2",
			},
		},
	},
	{
		"18a2e1d81b8ecfb2a333adcb0c17a5b9eb76cc5d05db91a4",
		"board blade invite damage undo sun mimic interest slam gaze truly inherit resist great inject rocket museum chief",
		"f84521c777a13b61564234bf8f8b62b3afce27fc4062b51bb5e62bdfecb23864ee6ecf07c1d5a97c0834307c5c852d8ceb88e7c97923c0a3b496bedd4e5f88a9",
		"xprv9s21ZrQH143K38yVDKj2uhq6e6jXBtMQATbysyZUGyG14JHvFRsHnvEDsW1xMedAm56UYZzwTDLL33ntWgTLkrynyvE4FLDP4DZpJRbMhMn",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0xBEA3547547E3a6088028bce40A1A4526CDE794F2",
				"0x03d369c876bdef9422bcc416584a1c7dc07e9111920d9417815d96c219af2e690c",
				"0x03d369c876bdef9422bcc416584a1c7dc07e9111920d9417815d96c219af2e690c",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x439DD8c43e01F7f66DC7cB369fcc3F4b2272ffE4",
				"0x02788b5dd79d4cf3e080988880efda28fb11c387777984ec590fd6b20ea2f5debb",
				"0x4bb384bd691a726c1af4ac7b435574caf96a00c3f99f4d7862249dccb3fd5ed1",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TWUo9nHU7XMHjssyHd2p1PvcfrWbbvyZdV",
				"03b79c2dce203a7e1cab49084274f43b08d5faa68416e94fde19d19f9d54408a9d",
				"70796779c1c34ad8baf5a2ce9d56196dcb2087f551fef1c10f0044eac9bf8cc1",
			},
			{
				"m/44'/195'/0'/0/1",
				"TToqe6hACLvyxgyzaRiiYYQBzwaVvB3eAT",
				"03bdce05f07aa0cac826a4add37186e738193d73ee7feaf80e1875ab6d0322b3bf",
				"a9755069397e21d83ebfaa4ad0bf6447fd43bca0f0bff6348e3e5158941afb51",
			},
		},
	},
	{
		"15da872c95a13dd738fbf50e427583ad61f18fd99f628c417a61cf8343c90419",
		"beyond stage sleep clip because twist token leaf atom beauty genius food business side grid unable middle armed observe pair crouch tonight away coconut",
		"b15509eaa2d09d3efd3e006ef42151b30367dc6e3aa5e44caba3fe4d3e352e65101fbdb86a96776b91946ff06f8eac594dc6ee1d3e82a42dfe1b40fef6bcc3fd",
		"xprv9s21ZrQH143K47KSAu4o7EV43wqj2sxVHSHGLmY4ZThKffiHSBN2CNb8RtY6sdaNKZKq7mxa9WS3Kv2iBKtGkmD3L9iDBq1x959Uq3hKM32",
		[]testAddress{
			// ETH
			{
				"m/44'/60'/0'/0/0",
				"0x0F90786115eEea62D53e1FBCb364835e854277ec",
				"0x036499fef0b04c3fb37ae8ab959ef6e2281efd6b65fc5ad2b6007fcf9d1f2348f5",
				"0xe7b16158ab014e3c0c6186c63d6a822e4d99e092ab1997f51cc4013ab28a4d05",
			},
			{
				"m/44'/60'/0'/0/1",
				"0x428dC6AF08c075CcB708a57CFE411a96CBdd102E",
				"0x028dd43bece8df4be9c8e4565298bf46c3d923d3d34028286ca5081bada3b6909b",
				"0x02223ae0f866c4073d08d0e651f32b575d381085f6e671b3e39b3203c6f7e1c9",
			},

			// Tron
			{
				"m/44'/195'/0'/0/0",
				"TLHbWZ43QMN9aTvQmG5gRp9SeyL1pLxXy1",
				"028c58efbcaea44a29f4652b3da874f8322cb9f48fe913343a4efbca8dc844cec5",
				"798380d0db33b6f1b9b5595c6a1b5671613680b0e9d3e1c316c170103c83d8cf",
			},
			{
				"m/44'/195'/0'/0/1",
				"TJr76nU96sea53fmb8BLxG68C9SdSBWe3W",
				"02196d44f9606edc6d349418b0a8d99fcba2747c22879c33eacac1d3d1cc9e7641",
				"9396c879e50b28b422cd695aa35191d9809ca27058a381d68f4a2143e6c071aa",
			},
		},
	},
}
