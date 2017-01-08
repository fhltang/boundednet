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

We consider solutions for smaller versions of the problem: for `N <= N0`, a solution to the problem applied to the "first `N` networks", i.e. `{p[0], ..., p[N-1]}`, is a set of networks `{q[0], ..., q[M-1]}` of minimal footprint size.  Let `Solutions(M, N)` be the set of all minimal solutions applied to the first `N` networks.

The key is computing the function `MinSize(M, N)` which we define as "the minimal footprint size of solutions in `Solutions(M, N)`".  This function can be expressed recursively in `M` and `N` which will allow us to apply the standard dynamic programming trick to compute an `M * N` table.  With the table, we start from `MinSize(M, N0)` and backtrack to obtain the `M` networks which attain a minimal footprint.

## Solution

### Single Network Helper `LeastNetwork(i, j)`

For `0 <= i <= j < N0`, define `LeastNetwork(i, j)` to be the smallest network containing the union of networks `{p[i], ... p[j-1]}`.

For a given `i` and `j`, `LeastNetwork(i, j)` can be computed in `O(1)` time.

### Expressing `MinSize(M, N)` Recursively

If some solution in `Solutions(M, N-1)` is also in `Solutions(M, N)`, then `MinSize(M, N) == MinSize(M, N-1)`.

Otherwise, we consider for all `n <= N`, partitionings of the `N` networks into the initial `n` networks and remaining `N - n` networks.  Each solution in `Solutions(M-1, n)` covers the first `n` networks and `LeastNetwork(n, N)` covers the remaining `N-n` networks; together we have at most `M` networks.

For a given `n`, it is possible that `LeastNetwork(n, N)` overlaps with a network in a solution in `Solutions(M-1, n)` which means that solution has some network that is a subset of `LeastNetwork(n, N)`.  (Note that the alternative of `LeastNetwork(n, N)` is a subset of some network in the solution is not possible because that would mean that solution is also a solution of `Solutions(M, N)`.)  In this case, the footprint size of `Solutions(M-1, n) union {LeastNetwork(n, N)}` is less than the sum of the footprint sizes of `Solutions(M-1, n)` and `{LeastNetwork(n, N)}`.

However, if this happens then necessarily there is some `l < n` for which all solutions in `Solutions(M-1, l)` do not overlap with `LeastNetwork(l, N)`.

The value `MinSize(M, N)` can be expressed recursively as
* if there is a solution in `Solutions(M, N-1)` which is also in `Solutions(M, N)`, then `MinSize(M, N) == MinSize(M, N-1)`
* otherwise `MinSize(M, N) == min( MinSize(M-1, n) + size(LeastNetwork(n, N)) for n<=N )`

To help us determine whether a solution in `Solutions(M, N-1)` is also in `Solutions(M, N)`, we compute another function `RightBound(M, N)` defined to be "the right-most address covered by some solution in `Solutions(M, N)`".

Therefore, `max(p[N-1]) <= RightBound(M, N-1)` iff there is a solution in `Solutions(M, N-1)` that is also in `Solutions(M, N)`.

The value `RightBound(M, N)` can also be expressed recursively as
* if `RightBound(N-1, M)` > `max(p[N-1])`, then `RightBound(N, M) == RightBound(N-1, M)`
* otherwise, `RightBound(M, N) == max{ LeastNetwork(n, N) }` for those `n` giving the minimal values in the "otherwise" clause of the recursive expression of `MinSize(M, N)`

In practice, we compute `MinSize(M, N)` and `RightBound(M, N)` concurrently so it is not quite as ugly as the recursive form above.

### Asymptotic Complexity

  1. Sort input networks: `O(N * log(N))`
  1. Remove overlaps: `O(N)`
  1. Precompute `LeastNetwork`: `O(N^2)`
  1. Compute tables `MinSize` and `RightBound`: `O(N^2 * M)`
  1. Backtracking: `O(M + N)`
  
Overall: `O(N^2 * M)`