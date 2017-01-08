# Bounded Netmasks Problem

## Problem Statement

You are given `N` networks and you need to return `M` networks where the `M` networks contain the IPs from the `N` but minimise the number of additional IPs that are not in the `M` networks.

## Definitions

* An *address* is an integer in set `A = [0, 2^32)`.
* A *network* is a subset of `A` (i.e. a set of addresses) of the form `[ a * 2^k, (a + 1) * 2^k )` for some `a` and `k`.
* The *footprint* of a set of networks is the union of the networks
* The *footprint size* of a set of networks (or just *size* of a set of networks) is the cardinality of the footprint, i.e. the number of addresses in the union of the networks.
   
## Initial Observations

*Network intersections*: if two networks have a non-empty intersection, then necessarily one network is a subset of the other or vice-versa.

*Minimal footprint size*: since the footprint of the `M` networks includes that of the `N` networks, then minimising the footprint size of the `M` networks also minimises the additional addresses that are not in the `M` networks.  We say that a footprint is minimal if it has minimal size.

## Solution Overview

Throughout, we assume a fixed number `N0` of networks `{p[0], p[1], ... p[N0-1]}`.

We may assume that these networks:
1. are sorted by their max address, and
2. no two networks overlap.
If the networks do not satisfy these assumptions, we can sort them and remove overlaps.  This pre-processing step would take `O(N * log(N))` for the sort and `O(N)` to remove the overlaps.

We consider solutions for smaller versions of the problem: for `N <= N0`, a solution to the problem applied to the "first `N` networks", i.e. `{p[0], ..., p[N-1]}`, is a set of networks `{q[0], ..., q[M-1]}` of minimal footprint size.  Let `P(M, N)` be the set of all minimal solutions applied to the first `N` networks.

The key is computing the function `S(M, N)` which we define as "the minimal footprint size of solutions in `P(M, N)`".  This function can be expressed recursively in `M` and `N` which will allow us to apply the standard dynamic programming trick to compute an `M * N` table.  With the table, we start from `S(M, N0)` and backtrack to obtain the `M` networks which attain a minimal footprint.

## Solution

### Single Network Helper `F(i, j)`

For `0 <= i <= j < N0`, define `F(i, j)` to be the smallest network containing the union of networks `{p[i], ... p[j-1]}`.

For a given `i` and `j`, `F(i, j)` can be computed in `O(1)` time.

### Expressing `S(M, N)` Recursively

If some solution in `P(M, N-1)` is also in `P(M, N)`, then `S(M, N) == S(M, N-1)`.

Otherwise, we consider for all `n <= N`, partitionings of the `N` networks into the initial `n` networks and remaining `N - n` networks.  Each solution in `P(M-1, n)` covers the first `n` networks and `F(n, N)` covers the remaining `N-n` networks; together we have at most `M` networks.

For a given `n`, it is possible that `F(n, N)` overlaps with a network in a solution in `P(M-1, n)` which means that solution has some network that is a subset of `F(n, N)`.  (Note that the alternative of `F(n, N)` is a subset of some network in the solution is not possible because that would mean that solution is also a solution of `P(M, N)`.)  In this case, the footprint size of `P(M-1, n) union {F(n, N)}` is less than the sum of the footprint sizes of `P(M-1, n)` and `{F(n, N)}`.

However, if this happens then necessarily there is some `l < n` for which all solutions in `P(M-1, l)` do not overlap with `F(l, N)`.

The value `S(M, N)` can be expressed recursively as
* if there is a solution in `P(M, N-1)` which is also in `P(M, N)`, then `S(M, N) == S(M, N-1)`
* otherwise `S(M, N) == min( S(M-1, n) + size(F(n, N)) for n<=N )`

To help us determine whether a solution in `P(M, N-1)` is also in `P(M, N)`, we compute another function `R(M, N)` defined to be "the right-most address covered by some solution in `P(M, N)`".

We know that `max(p[N-1]) <= R(M, N-1)` iff there is a solution in `P(M, N-1)` that is also in `P(M, N)`.

The value `R(M, N)` can also be expressed recursively as
* if `R(N-1, M)` > `max(p[n-1])`, then `R(N, M) == R(N-1, M)`
* otherwise, `R(M, N) == max{ F(n, N) }` for those `n` giving the minimal values in the "otherwise" clause of the recursive expression of `S(M, N)`

In practice, we compute `S(M, N)` and `R(M, N)` concurrently so it is not quite as ugly as the recursive form above.