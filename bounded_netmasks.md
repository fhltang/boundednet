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

We consider solutions for smaller versions of the problem: for `N <= N0`, a solution to the problem applied to the "first `N` networks", i.e. `{p[0], ..., p[N-1]}`, is a set of `K < M` networks `{q[0], ..., q[K]}` of minimal footprint size.  Let `Solutions(M, N)` be the set of all minimal solutions applied to the first `N` networks.

The key is computing the function `MinSize(M, N)` which we define as "the minimum footprint size of all solutions in `Solutions(M, N)`".  This function can be expressed recursively in `M` and `N` which will allow us to apply the standard dynamic programming trick to compute an `M * N` table.  With the table, we start from `MinSize(M, N0)` and backtrack to obtain the `M` networks which attain a minimal footprint.

Note that `MinSize` is monotonically decreasing in its first argument and monotonically increasing in its second argument.  That is
   * for any `N`, `M1 <= M2` implies `MinSize(M1, N) >= MinSize(M2, N)`, and
   * for any `M`, `N1 <= N2` implies `MinSize(M, N1) <= MinSize(M, N2)`.

## Solution

### Single Network Helper `LeastNetwork(i, j)`

For `0 <= i <= j < N0`, define `LeastNetwork(i, j)` to be the smallest network containing the union of networks `{p[i], ... p[j-1]}`.

For `i == j`, the union of networks is empty, so the smallest network containing the union of networks is the empty network.

For a given `i` and `j`, `LeastNetwork(i, j)` can be computed in `O(1)` time.

### Expressing `MinSize(M, N)` Recursively

If a solution `x` in `Solutions(M, N-1)` has footprint covering `p[N-1]`, then since `MinSize` is monotonically increasing in its second argument we know that `x` is also in `Solutions(M, N)`.  Therefore `MinSize(M, N) == MinSize(M, N-1)`.

Otherwise, we consider for all `n <= N`, partitionings of the `N` networks into the initial `n` networks and remaining `N - n` networks.  For a given `n`, we construct potential solutions by taking each solution in `Solutions(M-1, n)` which covers the first `n` networks and adding `LeastNetwork(n, N)` to cover the remaining `N - n` networks; together we have at most `M` networks.

For a given `n`, it is possible that there is a network `q[i]` in a solution `y` in `Solutions(M-1, n)` which overlaps with `LeastNetwork(n, N)`.  In this case, we consider separately the scenarios where `q[i]` is a subset of `LeastNetwork(n, N)` and vice-versa.

Assuming `q[i]` is a subset of `LeastNetwork(n, N)`, then necessarily there is some `l < n` for which no solution in `Solutions(M-1, l)` overlaps with `LeastNetwork(l, N)`.  To see why this is so, let `j` be any index for which `p[j]` is a subset of `q[i]`.  Thus `p[j]` is a subset of `LeastNetwork(n, N)` (since `q[i]` is a subset of `LeastNetwork(n, N)`).  If we pick `l` to be the least such `j`, then we know that `LeastNetwork(n, N)` is a superset of all networks in `{p[l], ..., p[n-1]}` but none of the networks in `{p[0], ..., p[l-1]}` since we know that the networks are in ascending order and do not overlap.  This means that although the footprint size of `Solutions(M-1, n) union {LeastNetwork(n, N)}` is less than the sum of the footprint sizes of `Solutions(M-1, n)` and `{LeastNetwork(n, N)}`, the expression `min( MinSize(M-1, n) + size(LeastNetwor(n, N)) for n<=N )` correctly computes the minimum size of potential solutions constructed in this way.  Informally, we are saying that we can obtain a better solution by omitting `q[i]` and that we come across this better solution anyway when we consider `Solutions(M-1, j)` and `LeastNetwork(j, N)`.

Assuming `LeastNetwork(n, N)` is a subset of `q[i]`, ... SHIT, I PAINTED MYSELF INTO A CORNER

The value `MinSize(M, N)` can be expressed recursively as
* if there is a solution in `Solutions(M, N-1)` which is also in `Solutions(M, N)`, then `MinSize(M, N) == MinSize(M, N-1)`
* otherwise `MinSize(M, N) == min( MinSize(M-1, n) + size(LeastNetwork(n, N)) for n<=N )`

To help us determine whether a solution in `Solutions(M, N-1)` is also in `Solutions(M, N)`, we compute another function `RightBound(M, N)` defined to be "the right-most address covered by some solution in `Solutions(M, N)`".

Therefore, `max(p[N-1]) <= RightBound(M, N-1)` iff there is a solution in `Solutions(M, N-1)` that is also in `Solutions(M, N)`.

The value `RightBound(M, N)` can also be expressed recursively as
* if `max(p[N-1]) <= RightBound(M, N-1)`, then `RightBound(M, N) == RightBound(M, N-1)`
* otherwise, `RightBound(M, N) == max{ LeastNetwork(n, N) }` for those `n` giving the minimal values in the "otherwise" clause of the recursive expression of `MinSize(M, N)`

In practice, we compute each entry `MinSize(M, N)` and `RightBound(M, N)` at the same time so it is not quite as ugly as the recursive form above.

### Dynamic Programming

We compute the tables `MinSize` and `RightBound` row by row and cell by cell.

### Backtracking

We start at cell `(M, N0)` and apply the following:

   * if `M==1`, then emit network `LeastNetwork(0, N)` and stop
   * if the value of `MinSize(M, N) == MinSize(M, N-1)`, then do not emit any network and continue with `(M, N-1)`
   * if the value of `MinSize(M, N) == min( MinSize(M-1, n) + size(LeastNetwork(n, N)) for n<=N )`, then we pick the least `n` for which the minimum is attained and
      * emit `LeastNetwork(n, N)`
      * continue with `(M-1, n)`
      
To see why this works, if `MinSize(M, N) == MinSize(M, N-1)` then we know that `max(p[N-1]) <= RightBound(M, N-1)` which means that the solution in `Solutions(M, N-1)` with the largest right bound is also in `Solutions(M, N)`; we want to find that solution with the largest right bound.  To ensure we find the solution with the largest right bound, we always choose the least `n` for which we attain the minimal value; this least value of `n` gives the network `LeastNetwork(n, N)` with the largest right bound.

### Asymptotic Complexity

  1. Sort input networks: `O(N * log(N))`
  1. Remove overlaps: `O(N)`
  1. Precompute `LeastNetwork`: `O(N^2)`
  1. Compute tables `MinSize` and `RightBound`: `O(N^2 * M)`
  1. Backtracking: `O(M + N)`
  
Overall: `O(N^2 * M)`