# Trieste
Trieste is an Italian coastal and historical frontier city partially enclosed by Slovenia. Why not [visit](https://www.turismofvg.it/Locality/Trieste)?

Trieste is also a sparse prefix tree library for Go based on ideas from [crit-bit](https://cr.yp.to/critbit.html) and [qp tries](https://dotat.at/prog/qp/README.html).

## Status
This is currently a prototype.

## Building
You can get the dependencies necessary for building this library by installing [dep](https://github.com/golang/dep) and running:

```shell
dep ensure -v
```

## Purpose
This work is motivated by the need for prefix-addressable and Merkle tree storage in [Hyperledger Burrow](https://github.com/hyperledger/burrow), but if it ever becomes useful might have general applicability. It aims to provide a reasonably compressed and fast sorted, range-iterable, prefix tree with a full-byte branching factor of 256 (technically 257 including a distinguished terminal). It compresses paths in the tree (from root to leaf) by only introducing nodes on the 'critical byte' on which keys differ and it keeps down storage space by implementing a sparse map of child nodes for a branch using a bitset.

## Future
Burrow has some uses for a basic in-memory version of this tree to provide sorted caching layers, but further extensions may be of interest:

- Providing a lazy Merkle tree hash where hashes are stable over prefixes
- Providing a persistent version where nodes are written immutably to a storage backend
- Writing to disc or a memory-mapped file to provide a prefix tree database
- Providing various kinds of proof to facilitate (among other things) blockchain light clients

## Related work
See the excellent and more mature [Tendermint IAVL](https://github.com/tendermint/iavl) library with similar aims.

See also the [go-ethereum trie package](https://github.com/ethereum/go-ethereum/tree/master/trie) based also on a condensed trie (sometimes called PATRICIA) with support for proofs and Merkle hashes, with a 16 + 1 branching factor, but no sparse array.

## Why bother?
Surprisingly few standalone trie libraries exist for Go. IAVL is good but not so efficient for in-memory storage and because of the way balacning happens AVL is order-dependent and does not maintain stable hashes for prefixes.

A common data structure in-memory and on-disc may be possible with some advantages for serialisation. Also it might be nice to have the tree structure as the native database structure rather than implementing an prefix/AVL tree of a B+ tree (or whatever the database layer is).
