# Generate a Sybil DHT node with a given privateKey

Modification of the simple echo application [ldej.nl/post/building-an-echo-application-with-libp2p](https://ldej.nl/post/building-an-echo-application-with-libp2p/).

As the input, supply a private key output from the generateClosest.go. Then, the application bootstraps a routable host and connects it to the DHT.

```shell

$ go build .
$ ./SybilNode --privKey CAESQFAOFXLJF2EQl9kkd63+COLXMSpskmGN8IdbOJa1mAUND2tuv2w17HgfUoXWhiYjgUlgq3dfZKPOXyTsnbf67Tc=
Decoding key: CAESQFAOFXLJF2EQl9kkd63+COLXMSpskmGN8IdbOJa1mAUND2tuv2w17HgfUoXWhiYjgUlgq3dfZKPOXyTsnbf67Tc=
Unmarshalling key: "\b\x01\x12@P\x0e\x15r\xc9\x17a\x10\x97\xd9$w\xad\xfe\b\xe2\xd71*l\x92a\x8d\xf0\x87[8\x96\xb5\x98\x05\r\x0fkn\xbfl5\xecx\x1fR\x85ֆ&#\x81I`\xabw_d\xa3\xce_$읷\xfa\xed7"
Done
2022/10/10 15:59:52 Host ID: 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL
2022/10/10 15:59:52 Connect to me on:
2022/10/10 15:59:52   /ip4/192.168.1.65/tcp/58967/p2p/12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL
2022/10/10 15:59:52   /ip4/127.0.0.1/tcp/58967/p2p/12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL
2022/10/10 15:59:52 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL bootstrapping to QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd
2022/10/10 15:59:52 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL bootstrapping to QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM
2022/10/10 15:59:52 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL bootstrapping to QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu
2022/10/10 15:59:52 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL bootstrapping to QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ
2022/10/10 15:59:52 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL bootstrapping to QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd
2022/10/10 15:59:52 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL bootstrapping to QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu
2022/10/10 15:59:52 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL bootstrapping to QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64
2022/10/10 15:59:52 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL bootstrapping to QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM
2022/10/10 15:59:52 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL bootstrapping to QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64
2022/10/10 15:59:52 context.Background.WithCancel bootstrapDial 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM
2022/10/10 15:59:52 context.Background.WithCancel bootstrapDial 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64
2022/10/10 15:59:52 context.Background.WithCancel bootstrapDial 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64
2022/10/10 15:59:53 context.Background.WithCancel bootstrapDialSuccess QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ
2022/10/10 15:59:53 bootstrapped with QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ
2022/10/10 15:59:53 context.Background.WithCancel bootstrapDial 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ
2022/10/10 15:59:57 context.Background.WithCancel bootstrapDial 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu
2022/10/10 15:59:57 context.Background.WithCancel bootstrapDial 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu
2022/10/10 15:59:57 context.Background.WithCancel bootstrapDial 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd
2022/10/10 15:59:57 context.Background.WithCancel bootstrapDial 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd
2022/10/10 15:59:57 context.Background.WithCancel bootstrapDial 12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM
2022/10/10 15:59:57 I can be reached at:
2022/10/10 15:59:57 /ip4/192.168.1.65/tcp/58967/p2p/12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL
2022/10/10 15:59:57 /ip4/127.0.0.1/tcp/58967/p2p/12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL
2022/10/10 15:59:57 /ip4/81.156.202.148/tcp/58967/p2p/12D3KooWArZJFU1ALZ8WVfEYTSnnjoR9bju6tFN32nDsgtDphXPL
