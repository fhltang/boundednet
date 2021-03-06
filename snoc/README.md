# Bounded Networks Problem: A Snoc-Recursive Solution

## Problem Statement

You are given `N` networks and you need to return `M` networks where the `M` networks contain the IPs from the `N` but minimise the number of additional IPs that are not in the `M` networks.

## Definitions

   * An *address* is an integer in set `A = [0, 2^32)`.
   * A *network* is a subset of `A` (i.e. a set of addresses) of the form `[ a * 2^k, (a + 1) * 2^k )` for some `a` and `k`.
   * The *footprint* of a sequence of networks is the union of the networks.
   * The *footprint size* of a sequence of networks (or just *size* of a set of networks) is the cardinality of the footprint, i.e. the number of addresses in the union of the networks.
   
## Initial Observations

*Network intersections*: if two networks have a non-empty intersection, then necessarily one network is a subset of the other or vice-versa.

*Minimal footprint size*: since the footprint of the `M` networks includes that of the `N` networks, then minimising the footprint size of the `M` networks also minimises the additional addresses that are not in the `M` networks.  We say that a footprint is minimal if it has minimal size.

## Solution Overview

Throughout, we assume a fixed number `N0` of networks `{p[0], p[1], ... p[N0-1]}`.

We may assume that these networks:

   1. are sorted by their max address, and
   2. no two networks overlap.
   
If the networks do not satisfy these assumptions, we can sort them and remove overlaps.  This pre-processing step would take `O(N * log(N))` for the sort and `O(N)` to remove the overlaps.

We consider solutions for smaller versions of the problem.  For `N <= N0`, we 

   * define `Presolutions(M, N)` to be the set of *presolutions* where a presolution is a sequence of at most `M` networks whose footprint is a superset of the footprint of the first `N` input networks `{p[0], ..., p[N-1]}`
   * define `MinSize(M, N)` to be the smallest size footprint size of all presolutions in `PreSolutions(M, N)` which is `min{ size(union(y)) for y in Presolutions(M, N) }`
   * define `Solutions(M, N)` to be the subset of `Presolutions(M, N)` whose elements all have footprint size `MinSize(M, N)`.

The solution is based on the observation that we will encounter a solution in `Solutions(M, N)` by considering presolutions of the form

    snoc(x, q)

which is the sequence formed by appending network `q` to sequence `x`.

This observation gives a recursive formulation of `MinSize(M, N)` which allows us to apply the standard dynamic programming trick using an `M * N` table.  The table can be used to backtrack and obtain the `M` networks which attain a minimal footprint.

## Solution

### Monotonicity

Observation: `MinSize` is monotonically decreasing in its first argument and monotonically increasing in its second argument.  That is
   * for any `N`, `M1 <= M2` implies `MinSize(M1, N) >= MinSize(M2, N)`, and
   * for any `M`, `N1 <= N2` implies `MinSize(M, N1) <= MinSize(M, N2)`.
   
Corollary of monotonicity observation: If `x` is a solution in `Solutions(M, N-1)` and `x` covers `p[N-1]` (the `N`th network), then by definition `x` is in `Presolutions(M, N)`.  Since `MinSize` is monotonically increasing in its second argument, we know that in fact `x` is in `Solutions(M, N)`.

### Single Network Helper `LeastNetwork(i, j)`

For `0 <= i < j <= N0`, define `LeastNetwork(i, j)` to be the smallest network containing the union of networks `{p[i], ... p[j-1]}`.

For `i == j`, the union of networks is empty and we define `LeastNetwork(j, j)` to be `{}`, the empty set.

For given `i` and `j`, `LeastNetwork(i, j)` can be computed in `O(1)` time.

### Expressing `MinSize(M, N)` Recursively

We consider for all `n <= N`, partitionings of the `N` networks into the initial `n` networks and remaining `N - n` networks.  For a given `n`, we construct presolutions of the form `snoc(x, LeastNetwork(n, N))` by taking each solution `x` in `Solutions(M-1, n)` together with the network `LeastNetwork(n, N)`.  The former covers the first `n` networks and the latter covers the remaining `N - n` networks.  We know that `snoc(x, LeastNetwork(n, N))` is in `Presolutions(M, N)` since it has at most `M` networks.

We will now show that one such presolution is in fact a solution.  That is, for `M > 1`,

    MinSize(M, N) == min{ MinSize(M-1, n) + size(LeastNetwork(n, N)) for n<=N }
    
Suppose for all `n`, there is no overlap of any solution in `Solutions(M-1, n)` and `LeastNetwork(n, N)`.  In this case, we know that for any `x` in `Solutions(M-1, n)`

    size(union( snoc(x, LeastNetwork(n, N)) )) == size(union(x)) + size(LeastNetwork(n, N))
                                               == MinSize(M-1, n) + size(LeastNetwork(n, N))
    
and so the formula for `MinSize(M, N)` is clearly correct.

Now suppose for some `n`, there is a network `q[i]` in some solution `x` in `Solutions(M-1, n)` which overlaps with `LeastNetwork(n, N)`.  Therefore the sum of the footprint sizes is greater than the footprint size of `snoc(x LeastNetwork(n, N))`.  We will show that since we take the minimum over all `n <= N` that these "overestimates" do not affect the overall answer.  Using one of our initial observations about overlapping networks, it suffices to consider two scenarios: one where `q[i]` is a subset of `LeastNetwork(n, N)` and vice-versa.

Assuming `q[i]` is a subset of `LeastNetwork(n, N)`, then necessarily there is some `L < n` for which `LeastNetwork(n, N)` is a superset of all networks in `{p[L], ..., p[n-1]}` but none of the networks in `{p[0], ..., p[L-1]}`.  This means that `LeastNetwork(n, N) == LeastNetwork(L, N)` since they cover the same networks.  For any `x` in `Solutions(M-1, L)`, since `x` cannot overlap with `LeastNetwork(L, N)`

    size(union( snoc(x, LeastNetwork(L, N)) )) == MinSize(M-1, L) + size(LeastNetwork(L, N))
                                               == MinSize(M-1, L) + size(LeastNetwork(n, N))
                                               <= MinSize(M-1, n) + size(LeastNetwork(n, N))

by monotonicity of `MinSize`.  Therefore the true minimum would be attained for some `n <= L`.
                                      
To see why there is such `L`, let `j` be any index for which `p[j]` is a subset of `q[i]`.  Thus `p[j]` is a subset of `LeastNetwork(n, N)` (since `q[i]` is a subset of `LeastNetwork(n, N)`).  If we pick `L` to be the least such `j`, then we know that `LeastNetwork(n, N)` is a superset of all networks in `{p[L], ..., p[n-1]}` but none of the networks in `{p[0], ..., p[L-1]}` since we know that the networks are in ascending order and do not overlap.

Now suppose that `q[i]` is a network in solution `x` in `Solutions(M-1, n)` and assume `LeastNetwork(n, N)` is a subset of `q[i]`.  This means

    size(union( snoc(x, LeastNetwork(n, N)) )) == size(union( x ))
                                               == MinSize(M-1, n)
 
and so the expression `MinSize(M-1, n) + size(LeastNetwork(n, N)` overestimates when computing the minimum.  However, since `LeastNetwork(n, N)` is a superset of `p[N-1]` and `q[i]` is a superset of `LeastNetwork(n, N)`, this means that `x` is also in `Solutions(M-1, N)`.  By monotonicity of `MinSize`, we know that `MinSize(M-1, n) == MinSize(M-1, N)`.  Thus we have shown that the overestimate does not change the value of the minimum.

### Dynamic Programming

We compute the tables `MinSize` and `RightBound` for each row `M` starting from `1` and for increasing values of `N` from `1` to `N0`. 

### Backtracking

We start at cell `(M, N0)` and apply the following:

   * if `M==1`, then emit network `LeastNetwork(0, N)` and stop
   * if the value of `MinSize(M, N) == MinSize(M, N-1)`, then do not emit any network and continue with `(M, N-1)`
   * if the value of `MinSize(M, N) == MinSize(M-1, n) + size(LeastNetwork(n, N))` for some `n<=N`,
      * emit `LeastNetwork(n, N)` if it is not `{}`
      * continue with cell `(M-1, n)`
      
### Asymptotic Complexity

  1. Sort input networks: `O(N * log(N))`
  1. Remove overlaps: `O(N)`
  1. Precompute `LeastNetwork`: `O(N^2)`
  1. Compute tables `MinSize` and `RightBound`: `O(N^2 * M)`
  1. Backtracking: `O(M + N)`
  
Overall: `O(N^2 * M)`
